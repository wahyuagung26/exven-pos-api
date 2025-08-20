package response

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// ParseValidationErrors converts validator errors to standard error format
func ParseValidationErrors(err error) map[string][]string {
	fieldErrors := make(map[string][]string)
	
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		for _, validationError := range validationErrors {
			field := strings.ToLower(validationError.Field())
			message := getValidationErrorMessage(validationError)
			fieldErrors[field] = append(fieldErrors[field], message)
		}
	}
	
	return fieldErrors
}

// ValidationErrorFromErr response with parsed validation errors
func ValidationErrorFromErr(c echo.Context, err error) error {
	fieldErrors := ParseValidationErrors(err)
	return ValidationError(c, fieldErrors)
}

func getValidationErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fe.Field() + " is required"
	case "email":
		return "Invalid email format"
	case "min":
		return fe.Field() + " must be at least " + fe.Param() + " characters"
	case "max":
		return fe.Field() + " must be at most " + fe.Param() + " characters"
	default:
		return fe.Field() + " is invalid"
	}
}