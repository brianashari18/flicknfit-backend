package main

import (
	"log"

	"flicknfit_backend/config"
	"flicknfit_backend/database"
	"flicknfit_backend/routes"
	"flicknfit_backend/utils"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

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
	}
	// Automatically migrate the database schema.
	database.Migrate(db, appLogger)

	// Create a new Fiber app instance with default configurations.
	app := fiber.New()

	// Use the built-in Fiber logger middleware to log incoming HTTP requests.
	app.Use(logger.New())

	// Set up the API routes. This function will be defined in the routes package.
	// We pass the Fiber app and the database instance to the route setup function.
	routes.SetupRoutes(app, db)

	// Start the server and listen on the specified port.
	log.Printf("Server starting on port %s", cfg.AppPort)
	if err := app.Listen(":" + cfg.AppPort); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
