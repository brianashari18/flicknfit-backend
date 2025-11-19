package unit

import (
	"flicknfit_backend/middlewares"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	t.Run("should return unauthorized when no authorization header", func(t *testing.T) {
		// Arrange
		app := fiber.New()
		app.Use(middlewares.AuthMiddleware())
		app.Get("/test", func(c *fiber.Ctx) error {
			return c.SendStatus(200)
		})

		req := httptest.NewRequest("GET", "/test", nil)

		// Act
		resp, err := app.Test(req)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 401, resp.StatusCode)
	})

	t.Run("should return unauthorized when invalid token format", func(t *testing.T) {
		// Arrange
		app := fiber.New()
		app.Use(middlewares.AuthMiddleware())
		app.Get("/test", func(c *fiber.Ctx) error {
			return c.SendStatus(200)
		})

		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "InvalidToken")

		// Act
		resp, err := app.Test(req)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 401, resp.StatusCode)
	})

	t.Run("should return unauthorized when token is empty", func(t *testing.T) {
		// Arrange
		app := fiber.New()
		app.Use(middlewares.AuthMiddleware())
		app.Get("/test", func(c *fiber.Ctx) error {
			return c.SendStatus(200)
		})

		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer ")

		// Act
		resp, err := app.Test(req)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 401, resp.StatusCode)
	})
}

func TestErrorHandler(t *testing.T) {
	t.Run("should handle fiber error", func(t *testing.T) {
		// Arrange
		app := fiber.New()
		app.Use(middlewares.ErrorHandler())
		app.Get("/test", func(c *fiber.Ctx) error {
			return fiber.NewError(404, "Not Found")
		})

		req := httptest.NewRequest("GET", "/test", nil)

		// Act
		resp, err := app.Test(req)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 404, resp.StatusCode)
	})

	t.Run("should handle generic error", func(t *testing.T) {
		// Arrange
		app := fiber.New()
		app.Use(middlewares.ErrorHandler())
		app.Get("/test", func(c *fiber.Ctx) error {
			return fiber.NewError(500, "Internal Server Error")
		})

		req := httptest.NewRequest("GET", "/test", nil)

		// Act
		resp, err := app.Test(req)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 500, resp.StatusCode)
	})
}

func TestRecoverHandler(t *testing.T) {
	t.Run("should recover from panic", func(t *testing.T) {
		// Arrange
		app := fiber.New()
		app.Use(middlewares.RecoverHandler())
		app.Get("/test", func(c *fiber.Ctx) error {
			panic("test panic")
		})

		req := httptest.NewRequest("GET", "/test", nil)

		// Act
		resp, err := app.Test(req)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 500, resp.StatusCode)
	})
}
