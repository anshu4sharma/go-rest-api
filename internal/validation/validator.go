package validation

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()

	// Register custom validation functions
	_ = validate.RegisterValidation("password", validatePassword)
	_ = validate.RegisterValidation("name", validateName)
}

// ValidateStruct validates a struct using validator tags
func ValidateStruct(s interface{}) error {
	err := validate.Struct(s)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
		}

		var errorMessages []string
		for _, err := range err.(validator.ValidationErrors) {
			errorMessages = append(errorMessages, formatError(err))
		}
		return fmt.Errorf("%s", strings.Join(errorMessages, "; "))
	}
	return nil
}

// Custom validation functions
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	// Password must contain at least:
	// - 8 characters
	// - 1 uppercase letter
	// - 1 lowercase letter
	// - 1 number
	// - 1 special character
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(password)

	return len(password) >= 8 && hasUpper && hasLower && hasNumber && hasSpecial
}

func validateName(fl validator.FieldLevel) bool {
	name := fl.Field().String()

	// Name must:
	// - Be between 2 and 50 characters
	// - Contain only letters, spaces, and hyphens
	// - Not start or end with space or hyphen
	if len(name) < 2 || len(name) > 50 {
		return false
	}

	if strings.HasPrefix(name, " ") || strings.HasPrefix(name, "-") ||
		strings.HasSuffix(name, " ") || strings.HasSuffix(name, "-") {
		return false
	}

	return regexp.MustCompile(`^[a-zA-Z\s-]+$`).MatchString(name)
}

// Helper function to format validation errors
func formatError(err validator.FieldError) string {
	field := strings.ToLower(err.Field())

	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "password":
		return fmt.Sprintf("%s must be at least 8 characters long and contain at least one uppercase letter, one lowercase letter, one number, and one special character", field)
	case "name":
		return "name must be between 2 and 50 characters long and contain only letters, spaces, and hyphens"
	default:
		return fmt.Sprintf("%s failed validation: %s", field, err.Tag())
	}
}
