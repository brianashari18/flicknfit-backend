package middlewares

import (
	"flicknfit_backend/config"
	"flicknfit_backend/services"
	"flicknfit_backend/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware verifies authentication tokens and sets user context
func AuthMiddleware() fiber.Handler {
	appLogger := utils.NewLogger()

	cfg, err := config.LoadConfig()
	if err != nil {
		appLogger.Fatalf("Error loading configuration: %v", err)
	}

	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return utils.SendError(c, fiber.StatusUnauthorized, "Missing authorization header")
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return utils.SendError(c, fiber.StatusUnauthorized, "Invalid authorization header format")
		}
		tokenString := parts[1]
		claims := &services.CustomClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.JwtSecretKey), nil
		})

		if err != nil || !token.Valid {
			return utils.SendError(c, fiber.StatusUnauthorized, "Invalid or expired access token")
		}

		// Store user information in context for next middlewares/handlers
		c.Locals("userID", claims.UserID)
		c.Locals("role", claims.Role)

		return c.Next()
	}
}
