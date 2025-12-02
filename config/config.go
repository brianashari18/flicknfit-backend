package config

import (
	"errors"
	"os"
)

// Config holds all the application-wide configuration read from environment variables.
type Config struct {
	AppHost                string
	AppPort                string
	DBUser                 string
	DBPassword             string
	DBHost                 string
	DBPort                 string
	DBName                 string
	DBSSLMode              string
	DBCertPath             string
	DBKeyPath              string
	DBRootCert             string
	DBEnableSSL            bool
	JwtSecretKey           string
	SmtpHost               string
	SmtpUser               string
	SmtpPassword           string
	SmtpPort               string
	AIApiURL               string
	FirebaseProjectID      string
	FirebasePrivateKeyPath string
	GroqAPIKey             string
	GroqModel              string
	GeminiAPIKey           string
	GeminiModel            string
	TelkomLLMAPIKey        string
	TelkomModel            string
	TelkomLLMUrl           string
	SupabaseURL            string
	SupabaseKey            string
	SupabaseBucket         string
	EncryptionKey          string
}

// LoadConfig reads the configuration from environment variables.
// It returns a Config struct and an error if any required variable is missing.
func LoadConfig() (*Config, error) {
	// Initialize a new Config struct.
	cfg := &Config{
		AppHost:                os.Getenv("APP_HOST"),
		AppPort:                os.Getenv("APP_PORT"),
		DBUser:                 os.Getenv("DB_USER"),
		DBPassword:             os.Getenv("DB_PASSWORD"),
		DBHost:                 os.Getenv("DB_HOST"),
		DBPort:                 os.Getenv("DB_PORT"),
		DBName:                 os.Getenv("DB_NAME"),
		DBSSLMode:              os.Getenv("DB_SSL_MODE"),
		DBCertPath:             os.Getenv("DB_CERT_PATH"),
		DBKeyPath:              os.Getenv("DB_KEY_PATH"),
		DBRootCert:             os.Getenv("DB_ROOT_CERT"),
		JwtSecretKey:           os.Getenv("JWT_SECRET_KEY"),
		SmtpHost:               os.Getenv("SMTP_HOST"),
		SmtpUser:               os.Getenv("SMTP_USER"),
		SmtpPassword:           os.Getenv("SMTP_PASSWORD"),
		SmtpPort:               os.Getenv("SMTP_PORT"),
		AIApiURL:               os.Getenv("AI_API_URL"),
		FirebaseProjectID:      os.Getenv("FIREBASE_PROJECT_ID"),
		FirebasePrivateKeyPath: os.Getenv("FIREBASE_PRIVATE_KEY_PATH"),
		GroqAPIKey:             os.Getenv("GROQ_API_KEY"),
		GroqModel:              os.Getenv("GROQ_MODEL"),
		GeminiAPIKey:           os.Getenv("GEMINI_API_KEY"),
		GeminiModel:            os.Getenv("GEMINI_MODEL"),
		TelkomLLMAPIKey:        os.Getenv("TELKOM_LLM_API_KEY"),
		TelkomModel:            os.Getenv("TELKOM_MODEL"),
		TelkomLLMUrl:           os.Getenv("TELKOM_LLM_URL"),
		SupabaseURL:            os.Getenv("SUPABASE_URL"),
		SupabaseKey:            os.Getenv("SUPABASE_KEY"),
		SupabaseBucket:         os.Getenv("SUPABASE_BUCKET"),
		EncryptionKey:          os.Getenv("ENCRYPTION_KEY"),
	}

	// Simple validation to ensure critical variables are set.
	if cfg.AppHost == "" {
		return nil, errors.New("APP_HOST environment variable is not set")
	}
	if cfg.AppPort == "" {
		return nil, errors.New("APP_PORT environment variable is not set")
	}
	if cfg.DBUser == "" || cfg.DBPassword == "" || cfg.DBHost == "" || cfg.DBPort == "" || cfg.DBName == "" {
		return nil, errors.New("database environment variables are not fully set")
	}
	if cfg.JwtSecretKey == "" {
		return nil, errors.New("JWT_SECRET_KEY environment variable is not set")
	}
	if cfg.AIApiURL == "" {
		return nil, errors.New("AI_API_URL environment variable is not set")
	}

	return cfg, nil
}
