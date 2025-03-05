package database

import (
	"database/sql"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testDB *MySQLStore

func TestMain(m *testing.M) {
	dsn := "root:rootpassword@tcp(localhost:3306)/url_shortener?parseTime=true"
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		panic("Failed to connect test database, " + err.Error())
	}

	testDB = &MySQLStore{db: db}

	testDB.db.Exec(`
	DROP TABLE IF EXISTS rate_limits;
	DROP TABLE IF EXISTS urls;
	`)

	// Ensure tables exist
	testDB.db.Exec(`
	CREATE TABLE IF NOT EXISTS urls (
		id INT AUTO_INCREMENT PRIMARY KEY,
		short_url VARCHAR(255) UNIQUE NOT NULL,
		original_url TEXT NOT NULL,
		clicks INTEGER DEFAULT 0,
		expires_at TIMESTAMP NOT NULL
	);`)

	testDB.db.Exec(`
	CREATE TABLE IF NOT EXISTS rate_limits (
		id INT AUTO_INCREMENT PRIMARY KEY,
		client_ip VARCHAR(45) UNIQUE NOT NULL,
		remaining INT NOT NULL,
		reset_at TIMESTAMP NOT NULL
	);`)

	exitCode := m.Run()

	db.Close()
	os.Exit(exitCode)
}

func TestSaveURL(t *testing.T) {
	err := testDB.SaveURL("abc123", "https://example.com", 24)
	assert.Nil(t, err, "Expected no error when saving a URL")

	// Verify data is saved
	url, err := testDB.GetURL("abc123")
	assert.Nil(t, err, "Expected no error when retrieving a URL")
	assert.Equal(t, "https://example.com", url, "Expected the retrieved URL to match")
}

func TestGetURL_NotFound(t *testing.T) {
	url, err := testDB.GetURL("nonexistent")
	assert.Nil(t, err, "Expected no error when retrieving a nonexistent URL")
	assert.Empty(t, url, "Expected an empty result for nonexistent URL")
}

func TestRateLimit(t *testing.T) {
	ip := "localhost"

	remaining, _, err := testDB.CheckRateLimit(ip)
	assert.Nil(t, err, "Expected no error when checking rate limit")
	assert.Equal(t, 10, remaining, "Expected initial remaining limit to be 10")

	// Consume 1 request
	testDB.CheckRateLimit(ip)
	remaining, _, _ = testDB.CheckRateLimit(ip)
	assert.Equal(t, 8, remaining, "Expected remaining limit to decrease")
}

func TestClickCount(t *testing.T) {
	// Save a test URL
	testDB.SaveURL("test123", "https://example.com", 24)

	// Increment clicks
	err := testDB.ClickCount("test123")
	assert.Nil(t, err)

	// Verify clicks count
	var clicks int
	err = testDB.db.QueryRow("SELECT clicks FROM urls WHERE short_url = ?", "test123").Scan(&clicks)
	assert.Nil(t, err)
	assert.Equal(t, 1, clicks, "Expected 1 click count")
}
