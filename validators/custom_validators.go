package validators

import (
	"flicknfit_backend/constants"
	"regexp"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

// RegisterCustomValidators registers all custom validation rules
func RegisterCustomValidators(v *validator.Validate) {
	v.RegisterValidation("password", validatePassword)
	v.RegisterValidation("username", validateUsername)
	v.RegisterValidation("phone", validatePhone)
	v.RegisterValidation("role", validateRole)
	v.RegisterValidation("product_status", validateProductStatus)
}

// validatePassword validates password complexity
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	// Check length
	if len(password) < constants.MinPasswordLength || len(password) > constants.MaxPasswordLength {
		return false
	}

	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	// Require at least 3 of the 4 character types
	count := 0
	if hasUpper {
		count++
	}
	if hasLower {
		count++
	}
	if hasNumber {
		count++
	}
	if hasSpecial {
		count++
	}

	return count >= 3
}

// validateUsername validates username format and length
func validateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()

	// Check length
	if len(username) < constants.MinUsernameLength || len(username) > constants.MaxUsernameLength {
		return false
	}

	// Check format: alphanumeric, underscore, hyphen only
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, username)
	return matched
}

// validatePhone validates phone number format
func validatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()

	// Indonesian phone number validation
	// Supports formats: +62xxx, 62xxx, 08xxx, 8xxx
	patterns := []string{
		`^\+62[0-9]{8,13}$`, // +62xxxxxxxxx
		`^62[0-9]{8,13}$`,   // 62xxxxxxxxx
		`^08[0-9]{8,12}$`,   // 08xxxxxxxxx
		`^8[0-9]{8,12}$`,    // 8xxxxxxxxx
	}

	for _, pattern := range patterns {
		if matched, _ := regexp.MatchString(pattern, phone); matched {
			return true
		}
	}

	return false
}

// validateRole validates user role
func validateRole(fl validator.FieldLevel) bool {
	role := strings.ToLower(fl.Field().String())
	return role == constants.RoleUser || role == constants.RoleAdmin
}

// validateProductStatus validates product status
func validateProductStatus(fl validator.FieldLevel) bool {
	status := strings.ToLower(fl.Field().String())
	validStatuses := []string{
		constants.ProductStatusActive,
		constants.ProductStatusInactive,
		constants.ProductStatusDraft,
	}

	for _, validStatus := range validStatuses {
		if status == validStatus {
			return true
		}
	}

	return false
}

// GetValidationErrorMessage returns user-friendly error messages for validation failures
func GetValidationErrorMessage(tag string, field string, param string) string {
	switch tag {
	case "required":
		return field + " is required"
	case "email":
		return field + " must be a valid email address"
	case "min":
		return field + " must be at least " + param + " characters long"
	case "max":
		return field + " must be at most " + param + " characters long"
	case "password":
		return field + " must be 8-72 characters long and contain at least 3 of: uppercase, lowercase, number, special character"
	case "username":
		return field + " must be 3-50 characters long and contain only letters, numbers, underscore, or hyphen"
	case "phone":
		return field + " must be a valid Indonesian phone number"
	case "role":
		return field + " must be either 'user' or 'admin'"
	case "product_status":
		return field + " must be 'active', 'inactive', or 'draft'"
	case "numeric":
		return field + " must be a number"
	case "gte":
		return field + " must be greater than or equal to " + param
	case "lte":
		return field + " must be less than or equal to " + param
	default:
		return field + " is invalid"
	}
}
