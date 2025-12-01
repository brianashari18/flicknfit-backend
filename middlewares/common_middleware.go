package middlewares

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

// SetupCommonMiddlewares sets up common security and performance middlewares
func SetupCommonMiddlewares(app *fiber.App) {
	// Request ID middleware - adds unique ID to each request
	app.Use(requestid.New())

	// Logger middleware - logs all requests
	app.Use(logger.New(logger.Config{
		Format:     "${time} [${ip}]:${port} ${status} - ${method} ${path} ${latency}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Local",
	}))

	// CORS middleware - handles Cross-Origin Resource Sharing
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://192.168.2.41:3000,http://192.168.2.41:8000,https://flicknfit.com",
		AllowMethods:     "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-Requested-With",
		ExposeHeaders:    "Content-Length,X-Request-ID",
		AllowCredentials: true,
		MaxAge:           86400, // 24 hours
	}))

	// Helmet middleware - sets security headers
	app.Use(helmet.New())

	// Rate limiting middleware
	app.Use(limiter.New(limiter.Config{
		Max:        100,         // Max requests per window
		Expiration: time.Minute, // Window duration
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP() // Rate limit by IP
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"success": false,
				"message": "Rate limit exceeded. Please try again later.",
			})
		},
	}))
}

// SetupAPIMiddlewares sets up API-specific middlewares
func SetupAPIMiddlewares(api fiber.Router) {
	// Stricter rate limiting for API endpoints
	api.Use(limiter.New(limiter.Config{
		Max:        50,
		Expiration: time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"success": false,
				"message": "API rate limit exceeded. Please try again later.",
			})
		},
	}))
}

// SetupAuthMiddlewares sets up authentication-specific middlewares
func SetupAuthMiddlewares(auth fiber.Router) {
	// Even stricter rate limiting for auth endpoints
	auth.Use(limiter.New(limiter.Config{
		Max:        10,
		Expiration: time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"success": false,
				"message": "Authentication rate limit exceeded. Please try again later.",
			})
		},
	}))
}
