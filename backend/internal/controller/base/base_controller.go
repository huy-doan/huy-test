package base

import (
	"net/http"
	"strconv"

	"github.com/huydq/test/internal/pkg/errors"
	"github.com/labstack/echo/v4"
)

// BaseController provides common controller functionality
type BaseController struct{}

// NewBaseController creates a new BaseController
func NewBaseController() *BaseController {
	return &BaseController{}
}

// BindAndValidate binds and validates a request
func (c *BaseController) BindAndValidate(ctx echo.Context, req interface{}) error {
	if err := ctx.Bind(req); err != nil {
		return errors.BadRequestError("リクエストデータが無効です", err.Error())
	}

	if err := ctx.Validate(req); err != nil {
		return errors.ValidationError("入力値の検証に失敗しました", err.Error())
	}

	return nil
}

// SendSuccess sends a successful response
func (c *BaseController) SendSuccess(ctx echo.Context, status int, data interface{}) error {
	return ctx.JSON(status, data)
}

// SendError sends an error response
func (c *BaseController) SendError(ctx echo.Context, err error) error {
	if appErr, ok := err.(*errors.Error); ok {
		return ctx.JSON(appErr.StatusCode, appErr)
	}

	// Generic error if not our custom error type
	internalErr := errors.InternalErrorWithCause("内部サーバーエラーが発生しました", err)
	return ctx.JSON(http.StatusInternalServerError, internalErr)
}

// SendOK sends a 200 OK response
func (c *BaseController) SendOK(ctx echo.Context, data interface{}) error {
	return c.SendSuccess(ctx, http.StatusOK, data)
}

// SendCreated sends a 201 Created response
func (c *BaseController) SendCreated(ctx echo.Context, data interface{}) error {
	return c.SendSuccess(ctx, http.StatusCreated, data)
}

// GetIDParam extracts and validates an ID from the URL path
func (c *BaseController) GetIDParam(ctx echo.Context, param string) (int, error) {
	id, err := strconv.Atoi(ctx.Param(param))
	if err != nil {
		return 0, errors.BadRequestError("パラメータが無効です", "数値である必要があります")
	}
	return id, nil
}
