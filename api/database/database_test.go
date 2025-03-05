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

	// Ensure tables exist
	testDB.db.Exec(`
	CREATE TABLE IF NOT EXISTS urls (
		id INT AUTO_INCREMENT PRIMARY KEY,
		short_url VARCHAR(255) UNIQUE NOT NULL,
		original_url TEXT NOT NULL,
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
	ip := "192.168.1.1"

	remaining, _, err := testDB.CheckRateLimit(ip)
	assert.Nil(t, err, "Expected no error when checking rate limit")
	assert.Equal(t, 10, remaining, "Expected initial remaining limit to be 10")

	// Consume 1 request
	testDB.CheckRateLimit(ip)
	remaining, _, _ = testDB.CheckRateLimit(ip)
	assert.Equal(t, 8, remaining, "Expected remaining limit to decrease")
}
