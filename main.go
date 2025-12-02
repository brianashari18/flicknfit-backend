package main

import (
	"log"

	"flicknfit_backend/admin"
	"flicknfit_backend/config"
	"flicknfit_backend/container"
	"flicknfit_backend/database"
	"flicknfit_backend/routes"
	"flicknfit_backend/seeders"
	"flicknfit_backend/utils"
	"log/slog"

	_ "flicknfit_backend/docs" // This is required for swag

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

// @title FlickNFit API
// @version 1.0
// @description FlickNFit backend API - Fashion recommendation and wardrobe management system
// @termsOfService http://swagger.io/terms/

// @contact.name FlickNFit API Support
// @contact.url http://www.flicknfit.com/support
// @contact.email support@flicknfit.com

// @license.name MIT
// @license.url https://github.com/flicknfit/backend/blob/main/LICENSE

// @host localhost:8000
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

// init runs before main() to load environment variables from the .env file.
func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables.")
	}
}

func main() {
	appLogger := utils.NewLogger()

	// Initialize and load the application configuration from environment variables.
	cfg, err := config.LoadConfig()
	if err != nil {
		appLogger.Fatalf("Error loading configuration: %v", err)
	}

	// Initialize the database connection.
	db, err := database.ConnectDB(cfg, appLogger)
	if err != nil {
		appLogger.Error("Failed to connect to database", slog.Any("error", err))
		return
	} // Automatically migrate the database schema.
	database.Migrate(db, appLogger)

	// Jalankan semua seeder sekaligus
	if err := seeders.SeedAll(db); err != nil {
		appLogger.Fatalf("Seeder error: %v", err)
	}

	// Initialize the dependency injection container
	appContainer, err := container.NewContainer(db, cfg)
	if err != nil {
		appLogger.Error("Failed to initialize container", slog.Any("error", err))
		return
	}

	// Create a new Fiber app instance with custom configurations.
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// Custom error handler will be handled by our middleware
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		},
		AppName: "FlickNFit API v1.0",
	})
	// Set up the API routes with middlewares
	routes.SetupRoutes(app, db, appContainer)
	// Setup Simple Admin Dashboard
	adminConfig := &admin.AdminConfig{
		DB:     db,
		Config: cfg,
	}

	if err := admin.SetupAdmin(app, adminConfig); err != nil {
		appLogger.Error("Failed to setup admin dashboard", slog.Any("error", err))
		// Continue running without admin - non-critical
		log.Printf("Warning: Admin dashboard not available: %v", err)
	} else {
		log.Printf("✅ Admin dashboard available at http://localhost:%s/admin", cfg.AppPort)
	}

	// Start the server and listen on the specified port.
	log.Printf("Server starting on %s:%s", cfg.AppHost, cfg.AppPort)
	if err := app.Listen(cfg.AppHost + ":" + cfg.AppPort); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
