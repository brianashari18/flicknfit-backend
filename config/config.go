package config

import (
	"errors"
	"os"
)

// Config holds all the application-wide configuration read from environment variables.
type Config struct {
	AppPort      string
	DBUser       string
	DBPassword   string
	DBHost       string
	DBPort       string
	DBName       string
	DBSSLMode    string
	DBCertPath   string
	DBKeyPath    string
	DBRootCert   string
	DBEnableSSL  bool
	JwtSecretKey string
	SmtpHost     string
	SmtpUser     string
	SmtpPassword string
	SmtpPort     string
}

// LoadConfig reads the configuration from environment variables.
// It returns a Config struct and an error if any required variable is missing.
func LoadConfig() (*Config, error) {
	// Initialize a new Config struct.
	cfg := &Config{
		AppPort:      os.Getenv("APP_PORT"),
		DBUser:       os.Getenv("DB_USER"),
		DBPassword:   os.Getenv("DB_PASSWORD"),
		DBHost:       os.Getenv("DB_HOST"),
		DBPort:       os.Getenv("DB_PORT"),
		DBName:       os.Getenv("DB_NAME"),
		DBSSLMode:    os.Getenv("DB_SSL_MODE"),
		DBCertPath:   os.Getenv("DB_CERT_PATH"),
		DBKeyPath:    os.Getenv("DB_KEY_PATH"),
		DBRootCert:   os.Getenv("DB_ROOT_CERT"),
		JwtSecretKey: os.Getenv("JWT_SECRET_KEY"),
		SmtpHost:     os.Getenv("SMTP_HOST"),
		SmtpUser:     os.Getenv("SMTP_USER"),
		SmtpPassword: os.Getenv("SMTP_PASSWORD"),
		SmtpPort:     os.Getenv("SMTP_PORT"),
	}

	// Simple validation to ensure critical variables are set.
	if cfg.AppPort == "" {
		return nil, errors.New("APP_PORT environment variable is not set")
	}
	if cfg.DBUser == "" || cfg.DBPassword == "" || cfg.DBHost == "" || cfg.DBPort == "" || cfg.DBName == "" {
		return nil, errors.New("database environment variables are not fully set")
	}
	if cfg.JwtSecretKey == "" {
		return nil, errors.New("JWT_SECRET_KEY environment variable is not set")
	}

	return cfg, nil
}
