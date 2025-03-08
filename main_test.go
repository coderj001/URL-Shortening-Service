package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/coderj001/URL-shortener/api/database"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var testDB *database.MySQLStore

func TestMain(m *testing.M) {
	if err := godotenv.Load("./.env"); err != nil {
		panic("Failed to load .env file")
	}

	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestPingPong(t *testing.T) {
	router := gin.Default()
	setupRoutes(router, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	expectedBody := `{"ping":"pong"}`
	assert.Equal(t, expectedBody, w.Body.String())
}

func TestShortenURL(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	testDB, _ = database.NewMySQLStore() // Ensure DB is initialized
	defer func() {
		testDB.Exec("DROP TABLE IF EXISTS rate_limits;")
		testDB.Exec("DROP TABLE IF EXISTS urls;")
		testDB.Close()
	}()
	setupRoutes(router, testDB)

	body := `{"url": "https://example.com", "short": "test123", "expiry": 24}`
	req, _ := http.NewRequest("POST", "/api/v1", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code, "Expected status 200 OK")

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err, "Response should be valid JSON")

	assert.Equal(t, "https://example.com", response["url"], "Expected URL to match")
	assert.Equal(t, "localhost:3000/test123", response["short"], "Expected short code to match")
}

func TestResolveURL(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	testDB, _ = database.NewMySQLStore() // Ensure DB is initialized
	defer func() {
		testDB.Exec("DROP TABLE IF EXISTS rate_limits;")
		testDB.Exec("DROP TABLE IF EXISTS urls;")
		testDB.Close()
	}()
	testDB.SaveURL("test123", "https://example.com", 24)
	setupRoutes(router, testDB)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test123", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusMovedPermanently, w.Code, "Expected status 200 OK")
	assert.Equal(t, "https://example.com", w.Header().Get("Location"))
}
