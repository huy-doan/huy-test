package errors

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	ms "github.com/huydq/test/src/lib/messages"
)

// ValidationErrorDetail represents a single validation error
type ValidationErrorDetail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// FormatValidationError formats validation errors to our custom format
func FormatValidationError(err error) *Error {
	if err == nil {
		return nil
	}

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return ValidationError(err.Error(), nil)
	}

	var details []ValidationErrorDetail
	for _, fieldError := range validationErrors {
		detail := ValidationErrorDetail{
			Field:   fieldError.Field(),
			Message: getErrorMessage(fieldError),
		}
		details = append(details, detail)
	}

	message := "Validation failed"
	if len(details) > 0 {
		message = fmt.Sprintf("%s", details[0].Message)
	}

	return ValidationError(message, details)
}

func getErrorMessage(fieldError validator.FieldError) string {
	switch fieldError.Tag() {
	case "required":
		return fmt.Sprintf(ms.ValidateRequired, strings.ToLower(fieldError.Field()))
	case "email":
		return fmt.Sprintf(ms.ValidateEmail, strings.ToLower(fieldError.Field()))
	case "min":
		return fmt.Sprintf(ms.ValidateMin, strings.ToLower(fieldError.Field()), fieldError.Param())
	case "max":
		return fmt.Sprintf(ms.ValidateMax, strings.ToLower(fieldError.Field()), fieldError.Param())
	case "kana":
		return fmt.Sprintf(ms.ValidateKana, strings.ToLower(fieldError.Field()))
	case "password_policy":
		return fmt.Sprintf(ms.ValidatePasswordPolicy, strings.ToLower(fieldError.Field()))
	default:
		return fmt.Sprintf(ms.ValidateField, strings.ToLower(fieldError.Field()), fieldError.Tag())
	}
}
