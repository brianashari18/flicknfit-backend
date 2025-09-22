package middlewares

import (
	"flicknfit_backend/config"
	"flicknfit_backend/services"
	"flicknfit_backend/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware adalah middleware untuk memverifikasi token otentikasi.
func AuthMiddleware() fiber.Handler {
	appLogger := utils.NewLogger()

	cfg, err := config.LoadConfig()
	if err != nil {
		appLogger.Fatalf("Error loading configuration: %v", err)
	}

	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Missing authorization header",
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid authorization header format",
			})
		}

		tokenString := parts[1]
		claims := &services.CustomClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.JwtSecretKey), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid or expired access token",
			})
		}

		c.Locals("userID", claims.UserID)

		return c.Next()
	}
}
