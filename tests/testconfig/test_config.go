package testconfig

import (
	"flicknfit_backend/config"
	"flicknfit_backend/database"
	"flicknfit_backend/utils"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestConfig holds test-specific configuration
type TestConfig struct {
	DB     *gorm.DB
	Config *config.Config
}

// NewTestConfig creates a new test configuration with in-memory SQLite database
func NewTestConfig() *TestConfig {
	// Create in-memory SQLite database for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to test database:", err)
	}
	// Run migrations
	logger := utils.NewLogger()
	database.Migrate(db, logger)

	// Create test config
	cfg := &config.Config{
		DBHost:       "192.168.2.41",
		DBPort:       "3306",
		DBUser:       "test",
		DBPassword:   "test",
		DBName:       "flicknfit_test",
		JwtSecretKey: "test-jwt-secret-key-for-testing-purposes-only",
		AppPort:      "8080",
	}

	return &TestConfig{
		DB:     db,
		Config: cfg,
	}
}

// Cleanup cleans up test resources
func (tc *TestConfig) Cleanup() {
	// Clear all tables
	tc.DB.Exec("DELETE FROM favorites")
	tc.DB.Exec("DELETE FROM reviews")
	tc.DB.Exec("DELETE FROM user_wardrobes")
	tc.DB.Exec("DELETE FROM shopping_cart_items")
	tc.DB.Exec("DELETE FROM shopping_carts")
	tc.DB.Exec("DELETE FROM product_items")
	tc.DB.Exec("DELETE FROM products")
	tc.DB.Exec("DELETE FROM brands")
	tc.DB.Exec("DELETE FROM users")
}
