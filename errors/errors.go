package errors

import (
	"fmt"
	"net/http"
)

// AppErrorType represents different types of application errors
type AppErrorType string

const (
	ErrorTypeValidation     AppErrorType = "VALIDATION_ERROR"
	ErrorTypeAuthentication AppErrorType = "AUTHENTICATION_ERROR"
	ErrorTypeAuthorization  AppErrorType = "AUTHORIZATION_ERROR"
	ErrorTypeNotFound       AppErrorType = "NOT_FOUND_ERROR"
	ErrorTypeConflict       AppErrorType = "CONFLICT_ERROR"
	ErrorTypeInternal       AppErrorType = "INTERNAL_ERROR"
	ErrorTypeDatabase       AppErrorType = "DATABASE_ERROR"
	ErrorTypeExternal       AppErrorType = "EXTERNAL_SERVICE_ERROR"
)

// AppError represents a structured application error
type AppError struct {
	Type        AppErrorType `json:"type"`
	Code        int          `json:"code"`
	Message     string       `json:"message"`
	Details     string       `json:"details,omitempty"`
	InternalErr error        `json:"-"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.InternalErr != nil {
		return fmt.Sprintf("%s: %s (internal: %v)", e.Type, e.Message, e.InternalErr)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// New creates a new AppError
func New(errorType AppErrorType, code int, message string) *AppError {
	return &AppError{
		Type:    errorType,
		Code:    code,
		Message: message,
	}
}

// NewWithDetails creates a new AppError with details
func NewWithDetails(errorType AppErrorType, code int, message, details string) *AppError {
	return &AppError{
		Type:    errorType,
		Code:    code,
		Message: message,
		Details: details,
	}
}

// NewWithInternalError creates a new AppError wrapping an internal error
func NewWithInternalError(errorType AppErrorType, code int, message string, err error) *AppError {
	return &AppError{
		Type:        errorType,
		Code:        code,
		Message:     message,
		InternalErr: err,
	}
}

// Predefined common errors
var (
	ErrInvalidInput       = New(ErrorTypeValidation, http.StatusBadRequest, "Invalid input provided")
	ErrUnauthorized       = New(ErrorTypeAuthentication, http.StatusUnauthorized, "Authentication required")
	ErrForbidden          = New(ErrorTypeAuthorization, http.StatusForbidden, "Access denied")
	ErrNotFound           = New(ErrorTypeNotFound, http.StatusNotFound, "Resource not found")
	ErrConflict           = New(ErrorTypeConflict, http.StatusConflict, "Resource already exists")
	ErrInternalServer     = New(ErrorTypeInternal, http.StatusInternalServerError, "Internal server error")
	ErrDatabaseConnection = New(ErrorTypeDatabase, http.StatusInternalServerError, "Database connection failed")
	ErrExternalService    = New(ErrorTypeExternal, http.StatusBadGateway, "External service unavailable")
)

// Factory functions for common errors
func NewValidationError(message string) *AppError {
	return New(ErrorTypeValidation, http.StatusBadRequest, message)
}

func NewAuthenticationError(message string) *AppError {
	return New(ErrorTypeAuthentication, http.StatusUnauthorized, message)
}

func NewAuthorizationError(message string) *AppError {
	return New(ErrorTypeAuthorization, http.StatusForbidden, message)
}

func NewNotFoundError(resource string) *AppError {
	return New(ErrorTypeNotFound, http.StatusNotFound, fmt.Sprintf("%s not found", resource))
}

func NewConflictError(resource string) *AppError {
	return New(ErrorTypeConflict, http.StatusConflict, fmt.Sprintf("%s already exists", resource))
}

func NewInternalError(message string, err error) *AppError {
	return NewWithInternalError(ErrorTypeInternal, http.StatusInternalServerError, message, err)
}

func NewDatabaseError(operation string, err error) *AppError {
	return NewWithInternalError(ErrorTypeDatabase, http.StatusInternalServerError,
		fmt.Sprintf("Database %s failed", operation), err)
}
