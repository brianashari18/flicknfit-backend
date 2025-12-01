package database

import (
	"flicknfit_backend/config"
	"flicknfit_backend/models"
	"fmt"
	"log/slog"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DB is the global variable for the database connection.
var DB *gorm.DB

// ConnectDB connects to the MySQL database. It returns a GORM DB instance and an error.
// The function now accepts the configuration object as a parameter, making it more modular.
func ConnectDB(cfg *config.Config, logger *logrus.Logger) (*gorm.DB, error) {
	var err error

	// Use DSN (Data Source Name) format for MySQL connection string.
	// The DSN format is "user:password@tcp(host:port)/database?charset=utf8mb4&parseTime=True&loc=Local".
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	// Open the database connection.
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Error("Failed to connect to database!", "error", err)
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Ping the database to ensure the connection is active.
	sqlDB, err := DB.DB()
	if err != nil {
		logger.Error("Failed to get database instance!", "error", err)
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Set connection pool settings to prevent connection timeouts and optimize performance.
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	logger.Info("Database connection established!")
	return DB, nil
}

// Migrate performs auto-migration of the database schema based on the defined models.
// It now accepts a logger instance for structured logging.
func Migrate(db *gorm.DB, logger *logrus.Logger) {
	logger.Info("Running database migrations...")
	err := db.AutoMigrate(
		&models.User{},
		&models.Brand{},
		&models.Product{},
		&models.ProductCategory{},
		&models.ProductStyle{},
		&models.ProductItem{},
		&models.ProductVariation{},
		&models.ProductVariationOption{},
		&models.ProductConfiguration{},
		&models.Review{},
		&models.ShoppingCart{},
		&models.ShoppingCartItem{},
		&models.Favorite{},
		&models.UserWardrobe{},
		&models.BodyScanHistory{},
		&models.StyleRecommendation{},
		&models.FaceScanHistory{},
		&models.ColorToneRecommendation{},
	)
	if err != nil {
		logger.Error("Failed to migrate database schema!", slog.Any("error", err))
	}
	logger.Info("Database migrations completed successfully!")
}
