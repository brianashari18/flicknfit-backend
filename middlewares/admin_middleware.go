package middlewares

import (
	"flicknfit_backend/utils"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// AdminMiddleware checks if the authenticated user has admin role
// This middleware should be used after AuthMiddleware
func AdminMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get user role from previous auth middleware
		role, exists := c.Locals("role").(string)
		if !exists {
			return utils.SendResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		}

		// Check if user has admin role
		if role != "admin" {
			return utils.SendResponse(c, http.StatusForbidden, "Admin access required", nil)
		}

		// Continue to next handler if verification successful
		return c.Next()
	}
}
