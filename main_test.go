package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/coderj001/URL-shortener/database"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	db     *database.MySQLStore
	router *gin.Engine
}

func (s *TestSuite) SetupSuite() {
	if err := godotenv.Load("./.env"); err != nil {
		panic("Failed to load .env file")
	}
	s.db, _ = database.NewMySQLStore()
	s.router = gin.Default()
}

func (s *TestSuite) TearDownSuite() {
	s.db.Exec("DROP TABLE IF EXISTS rate_limits;")
	s.db.Exec("DROP TABLE IF EXISTS urls;")
	s.db.Exec("DROP TABLE IF EXISTS analytics;")
	s.db.Close()
}

func (s *TestSuite) TestPingPong() {
	router := gin.Default()
	setupRoutes(router, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)
	expectedBody := `{"ping":"pong"}`
	s.Equal(expectedBody, w.Body.String())
}

func (s *TestSuite) TestShortenURL() {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	setupRoutes(router, s.db)

	body := `{"url": "https://youtube.com", "short": "test129", "expiry": 24}`
	req, _ := http.NewRequest("POST", "/api/v1", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	s.Equal(http.StatusOK, w.Code, "Expected status 200 OK")

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	s.Nil(err, "Response should be valid JSON")

	s.Equal("https://youtube.com", response["url"], "Expected URL to match")
	s.Equal("test129", response["short"], "Expected short code to match")
}

func (s *TestSuite) TestResolveURL() {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	s.db.SaveURL("test124", "https://google.com", 24)
	setupRoutes(router, s.db)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test124", nil)
	router.ServeHTTP(w, req)

	s.Equal(http.StatusMovedPermanently, w.Code, "Expected status 200 OK")
	s.Equal("https://google.com", w.Header().Get("Location"))
}

func (s *TestSuite) TestAnalyticsURL() {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	s.db.SaveURL("test125", "https://github.com", 24)
	s.db.UpdateAnalytics("test125")
	s.db.UpdateAnalytics("test125")

	setupRoutes(router, s.db)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/analytics/test125", nil)
	router.ServeHTTP(w, req)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	s.Nil(err, "Response should be valid JSON")
	s.Equal(http.StatusOK, w.Code, "Expected status 200 OK")
	s.Equal("test125", response["short_id"], "Expected short_id to match test125")
	s.Equal(2, int(response["clicks"].(float64)), "Expected clicks to be 2")
}

func (s *TestSuite) TestSaveURL() {
	err := s.db.SaveURL("abc1234", "https://gmail.com", 24)
	s.Nil(err, "Expected no error when saving a URL")

	// Verify data is saved
	url, err := s.db.GetURL("abc1234")
	s.Nil(err, "Expected no error when retrieving a URL")
	s.Equal("https://gmail.com", url, "Expected the retrieved URL to match")
}

func (s *TestSuite) TestGetURL_NotFound() {
	url, err := s.db.GetURL("nonexistent")
	s.Nil(err, "Expected no error when retrieving a nonexistent URL")
	s.Empty(url, "Expected an empty result for nonexistent URL")
}

func (s *TestSuite) TestRateLimit() {
	ip := "localhost"

	remaining, _, err := s.db.CheckRateLimit(ip)
	s.Nil(err, "Expected no error when checking rate limit")
	s.Equal(10, remaining, "Expected initial remaining limit to be 10")

	// Consume 1 request
	s.db.CheckRateLimit(ip)
	remaining, _, _ = s.db.CheckRateLimit(ip)
	s.Equal(8, remaining, "Expected remaining limit to decrease")
}

func (s *TestSuite) TestClickCount() {
	// Save a test URL
	s.db.SaveURL("test123", "https://example.com", 24)

	// Increment clicks
	err := s.db.UpdateAnalytics("test123")
	s.Nil(err)

	// Verify clicks count
	var clicks int
	err = s.db.DB().QueryRow("SELECT clicks FROM analytics WHERE short_id = ?", "test123").Scan(&clicks)
	s.Nil(err)
	s.Equal(1, clicks, "Expected 1 click count")
}

func TestAPI(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
