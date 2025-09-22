package middlewares

import (
	"flicknfit_backend/services"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// AdminMiddleware adalah middleware untuk memverifikasi token otentikasi
// dan memeriksa apakah pengguna memiliki peran 'admin'.
func AdminMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"message": "Missing authorization header",
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid authorization header format",
			})
		}

		tokenString := parts[1]
		claims := &services.CustomClaims{}

		// Parse token tanpa verifikasi, hanya untuk mendapatkan claims
		_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return nil, nil
		})
		if err != nil && claims.Role != "admin" {
			return c.Status(http.StatusForbidden).JSON(fiber.Map{
				"message": "Access denied. Admin role required.",
			})
		}

		// Verifikasi role
		if claims.Role != "admin" {
			return c.Status(http.StatusForbidden).JSON(fiber.Map{
				"message": "Access denied. Admin role required.",
			})
		}

		// Lanjutkan ke handler berikutnya jika verifikasi berhasil
		return c.Next()
	}
}
