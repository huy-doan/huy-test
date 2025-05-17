package base

import (
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
func (c *BaseController) BindAndValidate(ctx echo.Context, req any) error {
	if err := ctx.Bind(req); err != nil {
		return errors.BadRequestError("リクエストデータが無効です", err.Error())
	}

	if err := ctx.Validate(req); err != nil {
		return errors.FormatValidationError(err)
	}

	return nil
}

// GetIDParam extracts and validates an ID from the URL path
func (c *BaseController) GetIDParam(ctx echo.Context, param string) (int, error) {
	id, err := strconv.Atoi(ctx.Param(param))
	if err != nil {
		paramError := errors.ErrorDetails{
			Field: param,
		}
		return 0, errors.BadRequestError("パラメータが無効です", paramError)
	}
	return id, nil
}
