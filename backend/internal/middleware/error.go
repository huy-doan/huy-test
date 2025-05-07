package middleware

import (
	"errors"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type ErrorResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Error   interface{} `json:"error,omitempty"`
}

func (m *MiddlewareManager) ErrorMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer func() {
				if r := recover(); r != nil {
					stack := debug.Stack()
					err, ok := r.(error)
					if !ok {
						err = errors.New("unknown panic")
					}

					m.logger.Error("Panic in request handler", map[string]any{
						"error":       err.Error(),
						"stack_trace": string(stack),
						"path":        c.Request().URL.Path,
						"method":      c.Request().Method,
					})

					c.JSON(http.StatusInternalServerError, ErrorResponse{
						Success: false,
						Message: "An internal server error occurred",
					})
				}
			}()

			err := next(c)
			if err != nil {
				return m.handleError(err, c)
			}

			return nil
		}
	}
}

func (m *MiddlewareManager) handleError(err error, c echo.Context) error {
	var httpError *echo.HTTPError
	var validationErrors validator.ValidationErrors

	statusCode := http.StatusInternalServerError
	message := "An unexpected error occurred"
	errorDetails := map[string]interface{}{}

	// Handle different error types
	if errors.As(err, &httpError) {
		statusCode = httpError.Code
		message = getErrorMessage(httpError)
	} else if errors.As(err, &validationErrors) {
		statusCode = http.StatusBadRequest
		message = "Validation failed"
		errorDetails = formatValidationError(validationErrors)
	} else {
		// Log unexpected errors
		m.logger.Error("Unhandled error", map[string]any{
			"error": err.Error(),
			"path":  c.Request().URL.Path,
		})
	}

	return c.JSON(statusCode, ErrorResponse{
		Success: false,
		Message: message,
		Error:   errorDetails,
	})
}

func getErrorMessage(err *echo.HTTPError) string {
	switch msg := err.Message.(type) {
	case string:
		return msg
	case map[string]interface{}:
		if message, ok := msg["message"].(string); ok {
			return message
		}
		return "Error occurred"
	default:
		return "An error occurred"
	}
}

func formatValidationError(validationErrors validator.ValidationErrors) map[string]interface{} {
	// Convert validation errors to a structured format
	errorDetails := make(map[string]string)

	for _, err := range validationErrors {
		fieldName := strings.ToLower(err.Field())
		var message string

		switch err.Tag() {
		case "required":
			message = "This field is required"
		case "email":
			message = "Please enter a valid email address"
		case "min":
			message = "Field value is too short"
		case "max":
			message = "Field value is too long"
		case "password_policy":
			message = "Password must be at least 12 characters and include upper/lowercase letters, numbers, and symbols"
		default:
			message = "Invalid value"
		}

		errorDetails[fieldName] = message
	}

	return map[string]interface{}{
		"type":    "VALIDATION",
		"details": errorDetails,
	}
}
