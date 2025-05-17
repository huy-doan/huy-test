package logger

import (
	"errors"
	"fmt"
	"runtime/debug"

	"maps"

	apiErrors "github.com/huydq/test/internal/pkg/errors"
)

// LogErrorWithContext logs an error with detailed contextual information
func (l *loggerImpl) LogErrorWithContext(err error, msg string, contextData map[string]any) {
	if contextData == nil {
		contextData = make(map[string]any)
	}

	// Ensure message is set
	if msg == "" {
		msg = "Error occurred"
	}

	// Extract error details
	errorFields := extractErrorDetails(err)

	// Merge error fields with context data
	maps.Copy(contextData, errorFields)

	// Log the error
	l.Error(msg, contextData)
}

// extractErrorDetails extracts detailed information from various error types
func extractErrorDetails(err error) map[string]any {
	if err == nil {
		return nil
	}

	errorFields := map[string]any{
		"error": err.Error(),
	}

	// Check for our custom API error
	var apiErr *apiErrors.Error
	if errors.As(err, &apiErr) {
		errorFields["error_code"] = apiErr.Code
		errorFields["error_type"] = apiErr.Type

		if apiErr.Details != nil {
			errorFields["error_details"] = apiErr.Details
		}

		if apiErr.Source != "" {
			errorFields["error_source"] = apiErr.Source
		}

		if apiErr.StackTrace != "" {
			errorFields["stack_trace"] = apiErr.StackTrace
		}

		if apiErr.Cause != nil {
			errorFields["error_cause"] = apiErr.Cause.Error()
		}
	} else {
		// For non-API errors, add stack trace for better debugging
		errorFields["stack_trace"] = string(debug.Stack())

		// Also check for wrapped errors
		var wrappedErr error
		if errors.Unwrap(err) != nil {
			wrappedErr = errors.Unwrap(err)
			errorFields["error_cause"] = wrappedErr.Error()
		}
	}

	return errorFields
}

// LogError logs an error with default message
func (l *loggerImpl) LogError(err error, contextData map[string]any) {
	l.LogErrorWithContext(err, "", contextData)
}

// ErrorWithContext creates and logs an error with context
func (l *loggerImpl) ErrorWithContext(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)

	// Create error information
	errorFields := map[string]any{
		"stack_trace": string(debug.Stack()),
	}

	// Log the error
	l.Error(msg, errorFields)
}
