package database

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type MySQLStore struct {
	db *sql.DB
}

func (s *MySQLStore) Close() {
	panic("unimplemented")
}

func NewMySQLStore() (*MySQLStore, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

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
			short_url VARCHAR(255) UNIQUE NOT NULL,
			original_url TEXT NOT NULL,
			clicks INTEGER DEFAULT 0,
			expires_at TIMESTAMP NOT NULL
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

	return &MySQLStore{db: db}, nil
}

func (s *MySQLStore) GetClickCount(shortID string) (int, error) {
	var clicks int
	err := s.db.QueryRow(
		"SELECT clicks FROM urls WHERE short_url = ? AND expires_at > NOW()",
		shortID,
	).Scan(&clicks)
	if err != nil {
		return 0, err
	}
	return clicks, nil
}

func (s *MySQLStore) ClickCount(short string) error {
	var click int
	err := s.db.QueryRow(
		"SELECT clicks FROM urls WHERE short_url = ? AND expires_at > NOW()",
		short,
	).Scan(&click)
	if err != nil {
		return err
	}
	_, err = s.db.Exec("UPDATE urls SET clicks = ? WHERE short_url = ?", click+1, short)
	if err != nil {
		return err
	}
	return nil
}

func (s *MySQLStore) SaveURL(short, original string, expiry time.Duration) error {
	expiresAt := time.Now().Add(expiry * time.Hour)
	_, err := s.db.Exec(
		"INSERT INTO urls (short_url, original_url, expires_at) VALUES (?, ?, ?)",
		short, original, expiresAt,
	)
	return err
}

func (s *MySQLStore) GetURL(short string) (string, error) {
	var original string
	err := s.db.QueryRow(
		"SELECT original_url FROM urls WHERE short_url = ? AND expires_at > NOW()",
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
