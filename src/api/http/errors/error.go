package errors

import "fmt"

// Error represents a domain error
type Error struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	Type       string `json:"type"`
	StatusCode int    `json:"status_code"`
	Details    any    `json:"details,omitempty"`
}

// Error implements the error interface
func (e *Error) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// NewError creates a new error instance
func NewError(code string, message string, errorType string, statusCode int, details any) *Error {
	return &Error{
		Code:       code,
		Message:    message,
		Type:       errorType,
		StatusCode: statusCode,
		Details:    details,
	}
}

// ValidationError creates a new validation error
func ValidationError(message string, details any) *Error {
	return NewError("VALIDATION_ERROR", message, "VALIDATION", 400, details)
}

// NotFoundError creates a new not found error
func NotFoundError(message string) *Error {
	return NewError("NOT_FOUND", message, "NOT_FOUND", 404, nil)
}

// UnauthorizedError creates a new unauthorized error
func UnauthorizedError(message string) *Error {
	return NewError("UNAUTHORIZED", message, "AUTHORIZATION", 401, nil)
}

// ForbiddenError creates a new forbidden error
func ForbiddenError(message string) *Error {
	return NewError("FORBIDDEN", message, "AUTHORIZATION", 403, nil)
}

// InternalError creates a new internal server error
func InternalError(message string) *Error {
	return NewError("INTERNAL_ERROR", message, "SERVER", 500, nil)
}
