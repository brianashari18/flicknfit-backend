package utils

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// ErrorResponse is a generic struct for sending a consistent
// JSON error message to the client.
type ErrorResponse struct {
	Message string `json:"message"`
}

// APIResponse represents a standard API response structure
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// AppError represents application-specific errors
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("code: %d, message: %s, error: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("code: %d, message: %s", e.Code, e.Message)
}

// NewAppError creates a new application error
func NewAppError(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// SendError sends a standardized error response
func SendError(c *fiber.Ctx, statusCode int, message string) error {
	return SendResponse(c, statusCode, message, nil)
}

// SendSuccess sends a standardized success response
func SendSuccess(c *fiber.Ctx, message string, data interface{}) error {
	return SendResponse(c, 200, message, data)
}
