package database

import (
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
)

type DBTestSuite struct {
	suite.Suite
	db *MySQLStore
}

func (s *DBTestSuite) SetupSuite() {
	if err := godotenv.Load("../../.env"); err != nil {
		panic("Failed to load .env file")
	}
	s.db, _ = NewMySQLStore()
}

func (s *DBTestSuite) TearDownSuite() {
	s.db.Exec("DROP TABLE IF EXISTS rate_limits;")
	s.db.Exec("DROP TABLE IF EXISTS urls;")
	s.db.Close()
}

func (s *DBTestSuite) TestSaveURL() {
	err := s.db.SaveURL("abc123", "https://example.com", 24)
	s.Nil(err, "Expected no error when saving a URL")

	// Verify data is saved
	url, err := s.db.GetURL("abc123")
	s.Nil(err, "Expected no error when retrieving a URL")
	s.Equal("https://example.com", url, "Expected the retrieved URL to match")
}

func (s *DBTestSuite) TestGetURL_NotFound() {
	url, err := s.db.GetURL("nonexistent")
	s.Nil(err, "Expected no error when retrieving a nonexistent URL")
	s.Empty(url, "Expected an empty result for nonexistent URL")
}

func (s *DBTestSuite) TestRateLimit() {
	ip := "localhost"

	remaining, _, err := s.db.CheckRateLimit(ip)
	s.Nil(err, "Expected no error when checking rate limit")
	s.Equal(10, remaining, "Expected initial remaining limit to be 10")

	// Consume 1 request
	s.db.CheckRateLimit(ip)
	remaining, _, _ = s.db.CheckRateLimit(ip)
	s.Equal(8, remaining, "Expected remaining limit to decrease")
}

func (s *DBTestSuite) TestClickCount() {
	// Save a test URL
	s.db.SaveURL("test123", "https://example.com", 24)

	// Increment clicks
	err := s.db.ClickCount("test123")
	s.Nil(err)

	// Verify clicks count
	var clicks int
	err = s.db.db.QueryRow("SELECT clicks FROM urls WHERE short_url = ?", "test123").Scan(&clicks)
	s.Nil(err)
	s.Equal(1, clicks, "Expected 1 click count")
}

func TestDB(t *testing.T) {
	suite.Run(t, new(DBTestSuite))
}
