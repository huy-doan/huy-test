// src/api/http/middleware/error.go
package middleware

import (
	"errors"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	apiErrors "github.com/vnlab/makeshop-payment/src/api/http/errors"
	"github.com/vnlab/makeshop-payment/src/api/http/response"
)

// ErrorHandler middleware catches and standardizes errors
func ErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Use our custom ResponseWriter to capture errors
		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
			err:            nil,
		}

		// Process request with our custom ResponseWriter
		next.ServeHTTP(rw, r)

		// If there's an error, handle it
		if rw.err != nil {
			handleError(w, rw.err)
		}
	})
}

// responseWriter is a custom ResponseWriter that keeps track of errors
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	err        error
}

// WriteHeader overrides the WriteHeader to keep track of response status code
func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

// Error sets the error for this response
func (rw *responseWriter) Error(err error) {
	rw.err = err
}

// handleError standardizes error responses
func handleError(w http.ResponseWriter, err error) {
	var apiError *apiErrors.Error

	// Check if it's already our custom error type
	if errors.As(err, &apiError) {
		// Already formatted, just return it
		response.Error(w, apiError)
		return
	}

	// Check if it's a validation error
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		formattedErr := apiErrors.FormatValidationError(err)
		response.Error(w, formattedErr)
		return
	}

	// Handle other error types
	// For now, return as internal server error
	log.Printf("Unhandled error: %v", err)
	internalErr := apiErrors.InternalError("An unexpected error occurred")
	response.Error(w, internalErr)
}
