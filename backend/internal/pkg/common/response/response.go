package response

import (
	"net/http"

	"github.com/huydq/test/internal/pkg/errors"
	"github.com/huydq/test/internal/pkg/utils/messages"
	"github.com/labstack/echo/v4"
)

type BaseDataResponse[T any] struct {
	Success bool    `json:"success"`
	Message *string `json:"message,omitempty"`
	Data    T       `json:"data,omitempty"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   any    `json:"error,omitempty"`
}

type ValidationErrorInfo struct {
	Type    string            `json:"type"`
	Details map[string]string `json:"details"`
}

type StandardErrorInfo struct {
	Code string `json:"code"`
	Type string `json:"type"`
}

func NewSuccessDataResponse[T any](message string, data T) BaseDataResponse[T] {
	return BaseDataResponse[T]{
		Success: true,
		Message: &message,
		Data:    data,
	}
}

func NewSuccessDataResponseWithoutMessage[T any](data T) BaseDataResponse[T] {
	return BaseDataResponse[T]{
		Success: true,
		Data:    data,
	}
}

func SendResponse(ctx echo.Context, status int, message string, data any) error {
	resp := BaseDataResponse[any]{
		Success: true,
		Data:    data,
	}

	if message != "" {
		msg := message
		resp.Message = &msg
	}

	return ctx.JSON(status, resp)
}

func SendOK(ctx echo.Context, message string, data any) error {
	return SendResponse(ctx, http.StatusOK, message, data)
}

func SendCreated(ctx echo.Context, message string, data any) error {
	return SendResponse(ctx, http.StatusCreated, message, data)
}

func SendError(ctx echo.Context, err error) error {
	if appErr, ok := err.(*errors.Error); ok {
		resp := ErrorResponse{
			Success: false,
			Message: appErr.Message,
		}

		if appErr.Type == "VALIDATION" {
			resp.Error = appErr.Details
		} else if appErr.Details != nil {
			resp.Error = appErr.Details
		} else {
			resp.Error = StandardErrorInfo{
				Code: appErr.Code,
				Type: appErr.Type,
			}
		}

		return ctx.JSON(appErr.StatusCode, resp)
	}

	// For unknown errors, create an internal server error
	internalErr := errors.InternalErrorWithCause(messages.MsgInternalError, err)
	resp := ErrorResponse{
		Success: false,
		Message: internalErr.Message,
		Error: StandardErrorInfo{
			Code: internalErr.Code,
			Type: internalErr.Type,
		},
	}

	return ctx.JSON(http.StatusInternalServerError, resp)
}
