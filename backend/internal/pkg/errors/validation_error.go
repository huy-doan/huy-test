package errors

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	ms "github.com/huydq/test/internal/pkg/utils/messages"
)

type ValidationErrorResponse struct {
	Type    string                 `json:"type"`
	Details map[string]string      `json:"details"`
	Fields  []ValidationErrorField `json:"fields"`
}

type ValidationErrorField struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Tag     string `json:"tag,omitempty"`
	Value   string `json:"value,omitempty"`
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

	fields := make([]ValidationErrorField, 0, len(validationErrors))
	detailsMap := make(map[string]string)

	for _, fieldError := range validationErrors {
		fieldName := strings.ToLower(fieldError.Field())
		errorMessage := getErrorMessage(fieldError)

		field := ValidationErrorField{
			Field:   fieldName,
			Message: errorMessage,
			Tag:     fieldError.Tag(),
		}

		if fieldError.Value() != nil {
			if strValue, ok := fieldError.Value().(string); ok {
				field.Value = strValue
			}
		}

		fields = append(fields, field)

		detailsMap[fieldName] = errorMessage
	}

	validationResponse := ValidationErrorResponse{
		Type:    ms.TypeValidationError,
		Details: detailsMap,
		Fields:  fields,
	}

	message := ms.MsgValidationFailed
	if len(fields) > 0 {
		message = fields[0].Message
	}

	return ValidationError(message, validationResponse)
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
