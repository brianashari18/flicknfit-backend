package middlewares

import (
	"flicknfit_backend/errors"
	"flicknfit_backend/utils"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// ErrorHandler is a global error handling middleware for Fiber
func ErrorHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Try to continue with the next handler
		err := c.Next()

		if err != nil {
			return handleError(c, err)
		}

		return nil
	}
}

// handleError processes different types of errors and returns appropriate responses
func handleError(c *fiber.Ctx, err error) error {
	// Check if it's our custom AppError
	if appErr, ok := err.(*errors.AppError); ok {
		// Log internal errors for debugging
		if appErr.InternalErr != nil {
			log.Printf("Internal error: %v", appErr.InternalErr)
		}

		return utils.SendResponse(c, appErr.Code, appErr.Message, fiber.Map{
			"type":    appErr.Type,
			"details": appErr.Details,
		})
	}

	// Check if it's a Fiber error
	if fiberErr, ok := err.(*fiber.Error); ok {
		return utils.SendResponse(c, fiberErr.Code, fiberErr.Message, nil)
	}

	// Log unexpected errors
	log.Printf("Unexpected error: %v", err)

	// Default to internal server error
	return utils.SendResponse(c, http.StatusInternalServerError, "Internal server error", nil)
}

// RecoverHandler recovers from panics and converts them to errors
func RecoverHandler() fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Panic recovered: %v", r)
				err = errors.NewInternalError("Server panic occurred", nil)
			}
		}()

		return c.Next()
	}
}
