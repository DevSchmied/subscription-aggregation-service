package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// App configuration structure
type Config struct {
	AppPort    string
	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string
	DBSSLMode  string
}

// Load and validate configuration
func Load() (*Config, error) {
	// Load local .env file
	_ = godotenv.Load()
	cfg := &Config{
		AppPort:    os.Getenv("APP_PORT"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBName:     os.Getenv("DB_NAME"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBSSLMode:  os.Getenv("DB_SSLMODE"),
	}

	// Validate required fields
	if cfg.AppPort == "" {
		return nil, fmt.Errorf("APP_PORT is required")
	}
	if _, err := strconv.Atoi(cfg.AppPort); err != nil {
		return nil, fmt.Errorf("APP_PORT must be numeric")
	}

	if cfg.DBHost == "" {
		cfg.DBHost = "localhost"
	}

	if cfg.DBPort == "" {
		return nil, fmt.Errorf("DB_PORT is required")
	}
	if _, err := strconv.Atoi(cfg.DBPort); err != nil {
		return nil, fmt.Errorf("DB_PORT must be numeric")
	}

	if cfg.DBName == "" {
		return nil, fmt.Errorf("DB_NAME is required")
	}
	if cfg.DBUser == "" {
		return nil, fmt.Errorf("DB_USER is required")
	}
	if cfg.DBPassword == "" {
		return nil, fmt.Errorf("DB_PASSWORD is required")
	}

	// Apply default values
	if cfg.DBSSLMode == "" {
		cfg.DBSSLMode = "disable"
	}

	return cfg, nil
}
