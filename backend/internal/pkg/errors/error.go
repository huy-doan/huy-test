package errors

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"

	ms "github.com/huydq/test/internal/pkg/utils/messages"
)

// Error represents a domain error
type Error struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	Type       string `json:"type"`
	StatusCode int    `json:"status_code"`
	Details    any    `json:"details,omitempty"`

	// For internal tracking - not exposed in JSON response
	Cause      error  `json:"-"`
	StackTrace string `json:"-"`
	Source     string `json:"-"` // File and line where error occurred
}

type ErrorDetails struct {
	Service    string            `json:"service,omitempty"`     // For external service errors
	Field      string            `json:"field,omitempty"`       // For validation errors
	FieldInfos []ValidationField `json:"field_infos,omitempty"` // For validation errors with multiple fields
}

type ValidationField struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Error implements the error interface
func (e *Error) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap returns the underlying cause
func (e *Error) Unwrap() error {
	return e.Cause
}

// newErrorWithStack creates a new error instance with stack trace
func newErrorWithStack(code string, message string, errorType string, statusCode int, details any, cause error) *Error {
	// Capture caller information
	_, file, line, ok := runtime.Caller(2) // Skip this function and calling function
	source := "unknown source"
	if ok {
		// Extract just the filename from the path
		parts := strings.Split(file, "/")
		file = parts[len(parts)-1]
		source = fmt.Sprintf("%s:%d", file, line)
	}

	err := &Error{
		Code:       code,
		Message:    message,
		Type:       errorType,
		StatusCode: statusCode,
		Details:    details,
		Cause:      cause,
		Source:     source,
	}

	// Capture stack trace only for internal server errors
	if statusCode >= http.StatusInternalServerError {
		stack := make([]byte, 4096)
		n := runtime.Stack(stack, false)
		err.StackTrace = string(stack[:n])
	}

	return err
}

// NewError creates a new error instance
func NewError(code string, message string, errorType string, statusCode int, details any) *Error {
	return newErrorWithStack(code, message, errorType, statusCode, details, nil)
}

func NewErrorWithCause(code string, message string, errorType string, statusCode int, details any, cause error) *Error {
	return newErrorWithStack(code, message, errorType, statusCode, details, cause)
}

func ValidationError(message string, details any) *Error {
	return NewError(ms.CodeValidationError, message, ms.TypeValidationError, http.StatusBadRequest, details)
}

// NotFoundError creates a new not found error
func NotFoundError(message string) *Error {
	return NewError(ms.CodeNotFound, message, ms.TypeNotFoundError, http.StatusNotFound, nil)
}

func UnauthorizedError(message string) *Error {
	return NewError(ms.CodeUnauthorized, message, ms.TypeAuthorizationError, http.StatusUnauthorized, nil)
}

func ForbiddenError(message string) *Error {
	return NewError(ms.CodeForbidden, message, ms.TypeAuthorizationError, http.StatusForbidden, nil)
}

func BadRequestError(message string, details any) *Error {
	return NewError(ms.CodeBadRequest, message, ms.TypeClientError, http.StatusBadRequest, details)
}

func InternalError(message string) *Error {
	return NewError(ms.CodeInternalError, message, ms.TypeServerError, http.StatusInternalServerError, nil)
}

func InternalErrorWithCause(message string, cause error) *Error {
	return NewErrorWithCause(ms.CodeInternalError, message, ms.TypeServerError, http.StatusInternalServerError, nil, cause)
}

func DatabaseError(message string, cause error) *Error {
	return NewErrorWithCause(ms.CodeDatabaseError, message, ms.TypeDatabaseError, http.StatusInternalServerError, nil, cause)
}

func ExternalServiceError(message string, service string, cause error) *Error {
	details := ErrorDetails{
		Service: service,
	}
	return NewErrorWithCause(
		ms.CodeExternalServiceError,
		message,
		ms.TypeExternalServiceError,
		http.StatusInternalServerError,
		details,
		cause,
	)
}
