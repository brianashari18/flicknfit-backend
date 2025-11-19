package utils

import (
	"flicknfit_backend/validators"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validator wraps the validator instance with custom functionality
type Validator struct {
	validator *validator.Validate
}

// NewValidator creates a new validator instance with custom validators registered
func NewValidator() *validator.Validate {
	v := validator.New()

	// Register custom validators
	validators.RegisterCustomValidators(v)

	return v
}

// NewValidatorWrapper creates a new validator wrapper with custom validators registered
func NewValidatorWrapper() *Validator {
	v := validator.New()

	// Register custom validators
	validators.RegisterCustomValidators(v)

	return &Validator{
		validator: v,
	}
}

// Struct validates a struct and returns formatted error messages
func (v *Validator) Struct(s interface{}) error {
	err := v.validator.Struct(s)
	if err != nil {
		return formatValidationErrors(err)
	}
	return nil
}

// Var validates a single variable
func (v *Validator) Var(field interface{}, tag string) error {
	return v.validator.Var(field, tag)
}

// ValidationError represents a formatted validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrors represents multiple validation errors
type ValidationErrors struct {
	Errors []ValidationError `json:"errors"`
}

// Error implements the error interface
func (ve ValidationErrors) Error() string {
	var messages []string
	for _, err := range ve.Errors {
		messages = append(messages, err.Message)
	}
	return strings.Join(messages, "; ")
}

// formatValidationErrors converts validator errors to user-friendly messages
func formatValidationErrors(err error) ValidationErrors {
	var validationErrors []ValidationError

	if validatorErrors, ok := err.(validator.ValidationErrors); ok {
		for _, validatorErr := range validatorErrors {
			field := strings.ToLower(validatorErr.Field())
			message := validators.GetValidationErrorMessage(
				validatorErr.Tag(),
				field,
				validatorErr.Param(),
			)

			validationErrors = append(validationErrors, ValidationError{
				Field:   field,
				Message: message,
			})
		}
	}

	return ValidationErrors{Errors: validationErrors}
}
