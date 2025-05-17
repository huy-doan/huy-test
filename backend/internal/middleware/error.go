package middleware

import (
	"errors"
	"net/http"
	"runtime/debug"

	"github.com/go-playground/validator/v10"
	"github.com/huydq/test/internal/pkg/common/response"
	appErrors "github.com/huydq/test/internal/pkg/errors"
	messages "github.com/huydq/test/internal/pkg/utils/messages"
	"github.com/labstack/echo/v4"
)

type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   any    `json:"error,omitempty"`
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

					internalErr := appErrors.InternalError(messages.MsgInternalError)
					response.SendError(c, internalErr)
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

	// Handle echo.HTTPError
	if errors.As(err, &httpError) {
		// Handle standard HTTP errors with our custom format
		switch httpError.Code {
		case http.StatusNotFound:
			notFoundErr := appErrors.NotFoundError(messages.MsgRouteNotFound)
			return response.SendError(c, notFoundErr)
		case http.StatusMethodNotAllowed:
			methodNotAllowedErr := appErrors.NewError(
				messages.CodeMethodNotAllowed,
				messages.MsgMethodNotAllowed,
				messages.TypeClientError,
				http.StatusMethodNotAllowed,
				nil,
			)
			return response.SendError(c, methodNotAllowedErr)
		case http.StatusUnauthorized:
			unauthorizedErr := appErrors.UnauthorizedError(messages.MsgUnauthorized)
			return response.SendError(c, unauthorizedErr)
		case http.StatusForbidden:
			forbiddenErr := appErrors.ForbiddenError(messages.MsgForbidden)
			return response.SendError(c, forbiddenErr)
		case http.StatusBadRequest:
			badRequestErr := appErrors.BadRequestError(messages.MsgBadRequest, httpError.Message)
			return response.SendError(c, badRequestErr)
		default:
			message := getErrorMessage(httpError)
			customErr := appErrors.NewError(
				messages.CodeUnknownError,
				message,
				messages.TypeServerError,
				httpError.Code,
				nil,
			)
			return response.SendError(c, customErr)
		}
	} else if errors.As(err, &validationErrors) {
		formattedErr := appErrors.FormatValidationError(validationErrors)
		return response.SendError(c, formattedErr)
	} else if appErr, ok := err.(*appErrors.Error); ok {
		return response.SendError(c, appErr)
	} else {
		// Log unexpected errors
		m.logger.Error("Unhandled error", map[string]any{
			"error": err.Error(),
			"path":  c.Request().URL.Path,
		})

		internalErr := appErrors.InternalErrorWithCause(messages.MsgInternalError, err)
		return response.SendError(c, internalErr)
	}
}

func getErrorMessage(err *echo.HTTPError) string {
	switch msg := err.Message.(type) {
	case string:
		return msg
	case map[string]any:
		if message, ok := msg["message"].(string); ok {
			return message
		}
		return "Error occurred"
	default:
		return "An error occurred"
	}
}
