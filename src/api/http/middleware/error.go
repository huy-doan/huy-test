package middleware

import (
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	apiErrors "github.com/huydq/demo/src/api/http/errors"
)

// ErrorHandler middleware catches and standardizes errors
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Process request
		c.Next()

		// Check if there are any errors
		if len(c.Errors) > 0 {
			// Get the last error
			err := c.Errors.Last()
			handleError(c, err.Err)
		}
	}
}

// handleError standardizes error responses
func handleError(c *gin.Context, err error) {
	var apiError *apiErrors.Error

	// Check if it's already our custom error type
	if errors.As(err, &apiError) {
		// Already formatted, just return it
		c.JSON(apiError.StatusCode, apiError)
		return
	}

	// Check if it's a validation error
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		formattedErr := apiErrors.FormatValidationError(err)
		c.JSON(formattedErr.StatusCode, formattedErr)
		return
	}

	// Handle other error types
	// For now, return as internal server error
	log.Printf("Unhandled error: %v", err)
	internalErr := apiErrors.InternalError("An unexpected error occurred")
	c.JSON(internalErr.StatusCode, internalErr)
}
