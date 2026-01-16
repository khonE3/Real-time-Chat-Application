package config

import (
	"fmt"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerHost string
	ServerPort string

	// PostgreSQL
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// Redis
	RedisURL      string
	RedisPassword string

	// CORS
	CORSOrigins string
}

func Load() *Config {
	// Load .env file if exists
	godotenv.Load()

	return &Config{
		ServerHost: getEnv("SERVER_HOST", "0.0.0.0"),
		ServerPort: getEnv("SERVER_PORT", "8080"),

		// PostgreSQL
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "chatdb"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		// Redis
		RedisURL:      getEnv("REDIS_URL", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),

		// CORS
		CORSOrigins: getEnv("CORS_ORIGINS", "http://localhost:3000"),
	}
}

// GetDatabaseURL builds the PostgreSQL connection string
func (c *Config) GetDatabaseURL() string {
	// URL encode the password to handle special characters like %, @, etc.
	encodedPassword := url.QueryEscape(c.DBPassword)
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.DBUser, encodedPassword, c.DBHost, c.DBPort, c.DBName, c.DBSSLMode)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
