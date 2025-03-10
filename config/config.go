package config

import (
	"fmt"
	"os"
)

type DBConfig struct {
	Host     string
	Username string
	Password string
	Port     string
	Name     string
}

func (db *DBConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		db.Username,
		db.Password,
		db.Host,
		db.Port,
		db.Name,
	)
}

type Config struct {
	DB    *DBConfig
	Port  string
	Host  string
	Debug bool
}

func (c *Config) GetDomain() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

func GetConfig() *Config {
	return &Config{
		Port:  getEnv("PORT", "3000"),
		Host:  getEnv("HOST", "localhost"),
		Debug: getEnv("DEBUG", "true") == "true",
		DB: &DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			Username: getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", "rootpassword"),
			Name:     getEnv("DB_NAME", "url_shortener"),
		},
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
