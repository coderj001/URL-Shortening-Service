package database

import (
	"database/sql"
	"fmt"
	"time"

	apitypes "github.com/coderj001/URL-shortener/api"
	"github.com/coderj001/URL-shortener/config"
	_ "github.com/go-sql-driver/mysql"
)

type MySQLStore struct {
	db *sql.DB
}

func (s *MySQLStore) DB() *sql.DB {
	return s.db
}

func (s *MySQLStore) Exec(sql_query string) error {
	if _, err := s.db.Exec(sql_query); err != nil {
		return err
	}
	return nil
}

func (s *MySQLStore) Close() {
	s.db.Close()
}

func NewMySQLStore() (*MySQLStore, error) {
	dsn := config.GetConfig().DB.GetDSN()
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	if err = db.Ping(); err != nil {
		return nil, err
	}
	// Ensure tables exist
	urlsTableQuery := `
		CREATE TABLE IF NOT EXISTS urls (
			id INT AUTO_INCREMENT PRIMARY KEY,
			short_id VARCHAR(255) UNIQUE NOT NULL,
			original_url TEXT NOT NULL,
			expires_at TIMESTAMP NOT NULL,
			INDEX idx_expires_at(expires_at),
			UNIQUE INDEX idx_short_id(short_id)
		);`

	if _, err := db.Exec(urlsTableQuery); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	rateLimitTableQuery := `
		CREATE TABLE IF NOT EXISTS rate_limits (
			id INT AUTO_INCREMENT PRIMARY KEY,
			client_ip VARCHAR(45) UNIQUE NOT NULL,
			remaining INT NOT NULL,
			reset_at TIMESTAMP NOT NULL
		);`

	if _, err := db.Exec(rateLimitTableQuery); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	analyticsTableQuery := `
		CREATE TABLE IF NOT EXISTS analytics (
			id INT AUTO_INCREMENT PRIMARY KEY,
			short_id VARCHAR(255) UNIQUE NOT NULL,
			clicks INTEGER DEFAULT 1,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		);`

	if _, err := db.Exec(analyticsTableQuery); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	usersTableQuery := `
	CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
	 auth_level INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := db.Exec(usersTableQuery); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return &MySQLStore{db: db}, nil
}

func (s *MySQLStore) GetURLAnalytics(shortID string) (*apitypes.URLAnalytics, error) {
	analytics := &apitypes.URLAnalytics{ShortID: shortID}
	err := s.db.QueryRow(
		"SELECT clicks, created_at, updated_at FROM analytics where short_id = ?",
		shortID,
	).Scan(
		&analytics.Clicks,
		&analytics.CreatedAt,
		&analytics.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Ananlytics db yet to be created.")
		}
	}
	return analytics, nil
}

func (s *MySQLStore) UpdateAnalytics(shortID string) error {
	_, err := s.db.Exec(`
		INSERT INTO analytics (short_id) VALUES (?) 
		ON DUPLICATE KEY UPDATE clicks = clicks + 1	`, shortID)
	if err != nil {
		return err
	}
	return nil
}

func (s *MySQLStore) SaveURL(short, original string, expiry time.Duration) error {
	expiresAt := time.Now().Add(expiry * time.Hour)
	_, err := s.db.Exec(
		"INSERT INTO urls (short_id, original_url, expires_at) VALUES (?, ?, ?)",
		short, original, expiresAt,
	)
	return err
}

func (s *MySQLStore) GetURL(short string) (string, error) {
	var original string
	err := s.db.QueryRow(
		"SELECT original_url FROM urls WHERE short_id = ? AND expires_at > NOW()",
		short,
	).Scan(&original)

	if err == sql.ErrNoRows {
		return "", nil
	}
	return original, err
}

func (s *MySQLStore) CheckRateLimit(ip string) (int, time.Time, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return 0, time.Time{}, err
	}
	defer tx.Rollback()

	var remaining int
	var resetAt time.Time

	// Check existing rate limit
	err = tx.QueryRow(
		"SELECT remaining, reset_at FROM rate_limits WHERE client_ip = ? FOR UPDATE",
		ip,
	).Scan(&remaining, &resetAt)

	if err == sql.ErrNoRows {
		// Create new rate limit entry
		resetAt = time.Now().Add(30 * time.Minute)
		remaining = 10
		_, err = tx.Exec(
			"INSERT INTO rate_limits (client_ip, remaining, reset_at) VALUES (?, ?, ?)",
			ip, remaining, resetAt,
		)
	} else if err != nil {
		return 0, time.Time{}, err
	} else {
		if time.Now().After(resetAt) {
			// Reset the counter
			remaining = 10
			resetAt = time.Now().Add(30 * time.Minute)
			_, err = tx.Exec(
				"UPDATE rate_limits SET remaining = ?, reset_at = ? WHERE client_ip = ?",
				remaining, resetAt, ip,
			)
		} else {
			// Decrement remaining
			remaining--
			_, err = tx.Exec(
				"UPDATE rate_limits SET remaining = ? WHERE client_ip = ?",
				remaining, ip,
			)
		}
	}

	if err != nil {
		return 0, time.Time{}, err
	}

	if err = tx.Commit(); err != nil {
		return 0, time.Time{}, err
	}

	return remaining, resetAt, nil
}

func (s *MySQLStore) SaveUser(username, password_hash string) error {
	_, err := s.db.Exec("INSERT INTO users (username, password_hash) VALUES (?, ?)", username, password_hash)
	return err
}

func (s *MySQLStore) GetHashPassward(username string) (string, error) {
	var hashedPassword string
	err := s.db.QueryRow("SELECT password_hash FROM users WHERE username = ?", username).Scan(&hashedPassword)
	return hashedPassword, err
}

func (s *MySQLStore) GetUser(username string) (apitypes.User, error) {
	var user apitypes.User
	err := s.db.QueryRow("SELECT id, username, auth_level, created_at FROM users WHERE username = ?", username).Scan(
		&user.ID,
		&user.Username,
		&user.AuthLevel,
		&user.CreatedAt,
	)
	return user, err
}
