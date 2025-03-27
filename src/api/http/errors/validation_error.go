// src/api/http/errors/validation_error.go
package errors

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
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

    // Kiểm tra xem lỗi có phải là ValidationErrors không
    validationErrors, ok := err.(validator.ValidationErrors)
    if !ok {
        // Nếu không phải lỗi validation, trả về lỗi thông thường
        return ValidationError(err.Error(), nil)
    }

    // Tạo mảng chi tiết lỗi
    var details []ValidationErrorDetail
    for _, fieldError := range validationErrors {
        detail := ValidationErrorDetail{
            Field:   fieldError.Field(),
            Message: getErrorMessage(fieldError),
        }
        details = append(details, detail)
    }

    // Tạo thông báo tổng quát
    message := "Validation failed"
    if len(details) > 0 {
        message = fmt.Sprintf("Validation failed: %s", details[0].Message)
    }

    return ValidationError(message, details)
}

// getErrorMessage trả về thông báo lỗi dựa trên loại lỗi validation
func getErrorMessage(fieldError validator.FieldError) string {
    switch fieldError.Tag() {
    case "required":
        return fmt.Sprintf("The %s field is required", strings.ToLower(fieldError.Field()))
    case "email":
        return fmt.Sprintf("The %s field must be a valid email address", strings.ToLower(fieldError.Field()))
    case "min":
        return fmt.Sprintf("The %s field must be at least %s characters", strings.ToLower(fieldError.Field()), fieldError.Param())
    case "max":
        return fmt.Sprintf("The %s field must not be longer than %s characters", strings.ToLower(fieldError.Field()), fieldError.Param())
    case "kana":
        return fmt.Sprintf("The %s field must contain only Katakana characters", strings.ToLower(fieldError.Field()))
    default:
        return fmt.Sprintf("The %s field failed validation: %s", strings.ToLower(fieldError.Field()), fieldError.Tag())
    }
}
