package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	apiErrors "github.com/huydq/test/src/api/http/errors"
	"github.com/huydq/test/src/api/http/response"
	"github.com/huydq/test/src/infrastructure/logger"
)

// Logger context key for storing logger in request context
const LoggerContextKey contextKey = "logger"

// MaxBodyLogSize is the maximum size of request body to log (in bytes)
const MaxBodyLogSize = 4096 // 4KB

// ErrorHandler middleware catches and standardizes errors
func ErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get logger from context if available, otherwise use default
		var logInstance logger.Logger
		if ctxLogger, ok := r.Context().Value(LoggerContextKey).(logger.Logger); ok {
			logInstance = ctxLogger
		} else {
			logInstance = logger.GetLogger()
		}

		// Create a copy of request body for logging if needed
		var requestBody []byte
		var err error

		// Only capture request body on POST, PUT, PATCH requests
		if (r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH") && r.Body != nil {
			requestBody, r.Body, err = drainAndReplaceBody(r.Body, MaxBodyLogSize)
			if err != nil {
				logInstance.Error("Failed to read request body for logging", map[string]interface{}{
					"error": err.Error(),
					"path":  r.URL.Path,
				})
			}
		}

		// Use our custom ResponseWriter to capture errors
		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
			err:            nil,
			logger:         logInstance,
			request:        r,
			requestBody:    requestBody,
		}

		// Recover from panics
		defer func() {
			if rec := recover(); rec != nil {
				// Get stack trace
				stackTrace := debug.Stack()

				// Create error details
				errorDetails := map[string]interface{}{
					"error":       fmt.Sprintf("%v", rec),
					"stack_trace": string(stackTrace),
					"path":        r.URL.Path,
					"method":      r.Method,
					"headers":     sanitizeHeaders(r.Header),
					"query":       r.URL.Query(),
				}

				// Add request body if available
				if len(rw.requestBody) > 0 {
					errorDetails["request_body"] = sanitizeRequestBody(rw.requestBody)
				}

				// Log the panic with all details
				rw.logger.Error("Panic in HTTP handler", errorDetails)

				// Return 500 response
				internalErr := apiErrors.InternalError("An unexpected error occurred")
				response.Error(w, internalErr)
			}
		}()

		// Process request with our custom ResponseWriter
		next.ServeHTTP(rw, r)

		// If there's an error, handle it
		if rw.err != nil {
			// Log the error with detailed info
			logErrorWithDetails(rw, rw.err)

			// Handle and format the error for the response
			handleError(w, rw.err)
		}
	})
}

// responseWriter is a custom ResponseWriter that keeps track of errors
type responseWriter struct {
	http.ResponseWriter
	statusCode  int
	err         error
	logger      logger.Logger
	request     *http.Request
	respWriter  *responseStatusWriter // Reference to the request logger's response writer
	requestBody []byte                // Copy of request body for logging
}

// WriteHeader overrides the WriteHeader to keep track of response status code
func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)

	// Log non-success status codes
	if statusCode >= 400 {
		logErrorStatusWithDetails(rw, statusCode, nil)
	}
}

// Error sets the error for this response
func (rw *responseWriter) Error(err error) {
	rw.err = err
}

// logErrorStatusWithDetails logs HTTP error statuses with detailed information
func logErrorStatusWithDetails(rw *responseWriter, statusCode int, err error) {
	// Get request start time for duration calculation if available
	var duration time.Duration
	if startTime, ok := rw.request.Context().Value("request_start_time").(time.Time); ok {
		duration = time.Since(startTime)
	}

	// Build detailed error log
	logFields := map[string]interface{}{
		"status_code":  statusCode,
		"path":         rw.request.URL.Path,
		"method":       rw.request.Method,
		"duration_ms":  duration.Milliseconds(),
		"headers":      sanitizeHeaders(rw.request.Header),
		"query_params": rw.request.URL.Query(),
	}

	// Add error information if available
	if err != nil {
		var apiErr *apiErrors.Error
		if errors.As(err, &apiErr) {
			logFields["error_code"] = apiErr.Code
			logFields["error_type"] = apiErr.Type
			logFields["error_details"] = apiErr.Details
		} else {
			logFields["error"] = err.Error()
		}

		// Include validation errors
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			validationDetails := make([]string, 0, len(validationErrors))
			for _, fieldErr := range validationErrors {
				validationDetails = append(validationDetails, fmt.Sprintf(
					"Field '%s': failed '%s' validation",
					fieldErr.Field(),
					fieldErr.Tag(),
				))
			}
			logFields["validation_errors"] = validationDetails
		}
	}

	// Add sanitized request body if available
	if len(rw.requestBody) > 0 {
		logFields["request_body"] = sanitizeRequestBody(rw.requestBody)
	}

	// Log with appropriate level based on status code
	if statusCode >= 500 {
		// Include stack trace for 5xx errors
		logFields["stack_trace"] = string(debug.Stack())
		rw.logger.Error("Server error response", logFields)
	} else {
		rw.logger.Error("HTTP error response", logFields)
	}
}

// logErrorWithDetails logs errors with detailed context information
func logErrorWithDetails(rw *responseWriter, err error) {
	// Determine status code based on error type
	statusCode := http.StatusInternalServerError
	var apiErr *apiErrors.Error
	if errors.As(err, &apiErr) {
		statusCode = apiErr.StatusCode
	}

	// Log the error with all available details
	logErrorStatusWithDetails(rw, statusCode, err)
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
	// Return as internal server error
	internalErr := apiErrors.InternalError("An unexpected error occurred")
	response.Error(w, internalErr)
}

// WithLogger adds logger to the request context
func WithLogger(logger logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Add logger to context
			ctx := context.WithValue(r.Context(), LoggerContextKey, logger)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// drainAndReplaceBody reads the body and then replaces it with a new reader
func drainAndReplaceBody(body io.ReadCloser, maxSize int64) ([]byte, io.ReadCloser, error) {
	// Read up to maxSize bytes
	limitReader := io.LimitReader(body, maxSize)
	bodyBytes, err := io.ReadAll(limitReader)
	if err != nil {
		return nil, body, err
	}

	// Close original body
	body.Close()

	// Create new reader from captured bytes
	return bodyBytes, io.NopCloser(strings.NewReader(string(bodyBytes))), nil
}

// sanitizeHeaders removes sensitive information from headers
func sanitizeHeaders(headers http.Header) map[string]string {
	result := make(map[string]string)

	for key, values := range headers {
		// Skip sensitive headers
		keyVal := strings.ToLower(key)
		if key == "Authorization" || key == "Cookie" || key == "Set-Cookie" ||
			strings.Contains(keyVal, "token") ||
			strings.Contains(keyVal, "password") ||
			strings.Contains(keyVal, "secret") {
			continue
		}

		if len(values) > 0 {
			result[key] = values[0]
		}
	}

	return result
}

// sanitizeRequestBody removes sensitive fields from request body
func sanitizeRequestBody(bodyBytes []byte) interface{} {
	// Try to parse as JSON
	var bodyData map[string]interface{}
	err := json.Unmarshal(bodyBytes, &bodyData)
	if err != nil {
		// If not valid JSON, return as string (truncated if needed)
		bodyStr := string(bodyBytes)
		if len(bodyStr) > 1000 {
			return bodyStr[:1000] + "... (truncated)"
		}
		return bodyStr
	}

	// Sanitize JSON fields
	for key := range bodyData {
		// Hide sensitive fields
		keyVal := strings.ToLower(key)
		if strings.Contains(keyVal, "password") ||
			strings.Contains(keyVal, "username") ||
			strings.Contains(keyVal, "email") ||
			strings.Contains(keyVal, "token") ||
			strings.Contains(keyVal, "secret") ||
			strings.Contains(keyVal, "key") {
			bodyData[key] = "[REDACTED]"
		}
	}

	return bodyData
}
