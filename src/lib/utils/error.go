// src/lib/utils/error.go
package utils

import (
	"errors"
	"fmt"
	"runtime"
	"strings"

	apiErrors "github.com/vnlab/makeshop-payment/src/api/http/errors"
)

// WrapError wraps an error with a message and source information
func WrapError(err error, message string) error {
	if err == nil {
		return nil
	}

	// Capture caller information
	_, file, line, ok := runtime.Caller(1)
	source := "unknown source"
	if ok {
		// Extract just the filename from the path
		parts := strings.Split(file, "/")
		file = parts[len(parts)-1]
		source = fmt.Sprintf("%s:%d", file, line)
	}

	// Check if it's already our custom error type
	var apiErr *apiErrors.Error
	if errors.As(err, &apiErr) {
		// Create a copy with updated message and source
		return &apiErrors.Error{
			Code:       apiErr.Code,
			Message:    fmt.Sprintf("%s: %s", message, apiErr.Message),
			Type:       apiErr.Type,
			StatusCode: apiErr.StatusCode,
			Details:    apiErr.Details,
			Cause:      apiErr.Cause,
			StackTrace: apiErr.StackTrace,
			Source:     source,
		}
	}

	// Create a new internal error with the original as cause
	return apiErrors.InternalErrorWithCause(message, err)
}

// WrapDatabaseError wraps a database error
func WrapDatabaseError(err error, message string) error {
	if err == nil {
		return nil
	}
	return apiErrors.DatabaseError(message, err)
}

// WrapValidationError wraps a validation error
func WrapValidationError(err error, message string, details interface{}) error {
	if err == nil {
		return nil
	}
	
	return apiErrors.ValidationError(fmt.Sprintf("%s: %v", message, err), details)
}

// WrapExternalServiceError wraps an error from an external service
func WrapExternalServiceError(err error, service, message string) error {
	if err == nil {
		return nil
	}
	
	return apiErrors.ExternalServiceError(message, service, err)
}

// IsNotFoundError checks if an error is a "not found" error
func IsNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	
	var apiErr *apiErrors.Error
	return errors.As(err, &apiErr) && apiErr.Code == "NOT_FOUND"
}

// IsValidationError checks if an error is a validation error
func IsValidationError(err error) bool {
	if err == nil {
		return false
	}
	
	var apiErr *apiErrors.Error
	return errors.As(err, &apiErr) && apiErr.Type == "VALIDATION"
}

// IsAuthError checks if an error is an authorization error
func IsAuthError(err error) bool {
	if err == nil {
		return false
	}
	
	var apiErr *apiErrors.Error
	return errors.As(err, &apiErr) && apiErr.Type == "AUTHORIZATION"
}
