package response

import (
	"encoding/json"
	"net/http"

	apiErrors "github.com/huydq/test/src/api/http/errors"
)

// Response is the standard API response structure
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// writeJSON writes a JSON response with the given status code
func writeJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

// Success sends a successful response
func Success(w http.ResponseWriter, data interface{}, message string) {
	serializedData := DefaultSerialize(data)
	writeJSON(w, http.StatusOK, Response{
		Success: true,
		Message: message,
		Data:    serializedData,
	})
}

// Created sends a 201 Created response
func Created(w http.ResponseWriter, data interface{}, message string) {
	serializedData := DefaultSerialize(data)
	writeJSON(w, http.StatusCreated, Response{
		Success: true,
		Message: message,
		Data:    serializedData,
	})
}

// Error sends an error response
func Error(w http.ResponseWriter, err error) {
	// Handle different error types
	switch e := err.(type) {
	case *apiErrors.Error:
		// If it's our custom error type, use it directly
		writeJSON(w, e.StatusCode, Response{
			Success: false,
			Message: e.Message,
			Error: map[string]interface{}{
				"code":    e.Code,
				"type":    e.Type,
				"details": e.Details,
			},
		})
	default:
		// For any other error, convert to internal server error
		internalErr := apiErrors.InternalError(err.Error())
		writeJSON(w, internalErr.StatusCode, Response{
			Success: false,
			Message: internalErr.Message,
			Error: map[string]interface{}{
				"code":    internalErr.Code,
				"type":    internalErr.Type,
				"details": internalErr.Details,
			},
		})
	}
}

// ValidationError sends a validation error response
func ValidationError(w http.ResponseWriter, err error) {
	validationErr := apiErrors.FormatValidationError(err)
	writeJSON(w, validationErr.StatusCode, Response{
		Success: false,
		Message: validationErr.Message,
		Error: map[string]interface{}{
			"code":    validationErr.Code,
			"type":    validationErr.Type,
			"details": validationErr.Details,
		},
	})
}

// NotFound sends a not found error response
func NotFound(w http.ResponseWriter, message string) {
	notFoundErr := apiErrors.NotFoundError(message)
	writeJSON(w, notFoundErr.StatusCode, Response{
		Success: false,
		Message: notFoundErr.Message,
		Error: map[string]interface{}{
			"code": notFoundErr.Code,
			"type": notFoundErr.Type,
		},
	})
}

// Unauthorized sends an unauthorized error response
func Unauthorized(w http.ResponseWriter, message string) {
	unauthorizedErr := apiErrors.UnauthorizedError(message)
	writeJSON(w, unauthorizedErr.StatusCode, Response{
		Success: false,
		Message: unauthorizedErr.Message,
		Error: map[string]interface{}{
			"code": unauthorizedErr.Code,
			"type": unauthorizedErr.Type,
		},
	})
}

// Forbidden sends a forbidden error response
func Forbidden(w http.ResponseWriter, message string) {
	forbiddenErr := apiErrors.ForbiddenError(message)
	writeJSON(w, forbiddenErr.StatusCode, Response{
		Success: false,
		Message: forbiddenErr.Message,
		Error: map[string]interface{}{
			"code": forbiddenErr.Code,
			"type": forbiddenErr.Type,
		},
	})
}

// BadRequest sends a bad request error response
func BadRequest(w http.ResponseWriter, message string, details interface{}) {
	badRequestErr := apiErrors.ValidationError(message, details)
	writeJSON(w, badRequestErr.StatusCode, Response{
		Success: false,
		Message: badRequestErr.Message,
		Error: map[string]interface{}{
			"code":    badRequestErr.Code,
			"type":    badRequestErr.Type,
			"details": badRequestErr.Details,
		},
	})
}
