package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	apiErrors "github.com/huydq/demo/src/api/http/errors"
)

// Response is the standard API response structure
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// Success sends a successful response
func Success(c *gin.Context, data interface{}, message string) {
	serializedData := DefaultSerialize(data)
	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: message,
		Data:    serializedData,
	})
}

// Created sends a 201 Created response
func Created(c *gin.Context, data interface{}, message string) {
	serializedData := DefaultSerialize(data)
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Message: message,
		Data:    serializedData,
	})
}

// Error sends an error response
func Error(c *gin.Context, err error) {
	// Handle different error types
	switch e := err.(type) {
	case *apiErrors.Error:
		// If it's our custom error type, use it directly
		c.JSON(e.StatusCode, Response{
			Success: false,
			Message: e.Message,
			Error: gin.H{
				"code":    e.Code,
				"type":    e.Type,
				"details": e.Details,
			},
		})
	default:
		// For any other error, convert to internal server error
		internalErr := apiErrors.InternalError(err.Error())
		c.JSON(internalErr.StatusCode, Response{
			Success: false,
			Message: internalErr.Message,
			Error: gin.H{
				"code":    internalErr.Code,
				"type":    internalErr.Type,
				"details": internalErr.Details,
			},
		})
	}
}

// ValidationError sends a validation error response
func ValidationError(c *gin.Context, err error) {
	validationErr := apiErrors.FormatValidationError(err)
	c.JSON(validationErr.StatusCode, Response{
		Success: false,
		Message: validationErr.Message,
		Error: gin.H{
			"code":    validationErr.Code,
			"type":    validationErr.Type,
			"details": validationErr.Details,
		},
	})
}

// NotFound sends a not found error response
func NotFound(c *gin.Context, message string) {
	notFoundErr := apiErrors.NotFoundError(message)
	c.JSON(notFoundErr.StatusCode, Response{
		Success: false,
		Message: notFoundErr.Message,
		Error: gin.H{
			"code": notFoundErr.Code,
			"type": notFoundErr.Type,
		},
	})
}

// Unauthorized sends an unauthorized error response
func Unauthorized(c *gin.Context, message string) {
	unauthorizedErr := apiErrors.UnauthorizedError(message)
	c.JSON(unauthorizedErr.StatusCode, Response{
		Success: false,
		Message: unauthorizedErr.Message,
		Error: gin.H{
			"code": unauthorizedErr.Code,
			"type": unauthorizedErr.Type,
		},
	})
}

// Forbidden sends a forbidden error response
func Forbidden(c *gin.Context, message string) {
	forbiddenErr := apiErrors.ForbiddenError(message)
	c.JSON(forbiddenErr.StatusCode, Response{
		Success: false,
		Message: forbiddenErr.Message,
		Error: gin.H{
			"code": forbiddenErr.Code,
			"type": forbiddenErr.Type,
		},
	})
}

// BadRequest sends a bad request error response
func BadRequest(c *gin.Context, message string, details interface{}) {
	badRequestErr := apiErrors.ValidationError(message, details)
	c.JSON(badRequestErr.StatusCode, Response{
		Success: false,
		Message: badRequestErr.Message,
		Error: gin.H{
			"code":    badRequestErr.Code,
			"type":    badRequestErr.Type,
			"details": badRequestErr.Details,
		},
	})
}
