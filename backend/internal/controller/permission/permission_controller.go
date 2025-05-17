package controller

import (
	"github.com/huydq/test/internal/controller/base"
	"github.com/huydq/test/internal/controller/permission/mapper"
	"github.com/huydq/test/internal/pkg/common/response"
	"github.com/huydq/test/internal/pkg/errors"
	"github.com/huydq/test/internal/pkg/utils/messages"
	usecase "github.com/huydq/test/internal/usecase/permission"
	"github.com/labstack/echo/v4"
)

type PermissionController struct {
	permissionUsecase usecase.PermissionUsecase
	base.BaseController
}

func NewPermissionController(permissionUsecase usecase.PermissionUsecase) *PermissionController {
	return &PermissionController{
		permissionUsecase: permissionUsecase,
		BaseController:    *base.NewBaseController(),
	}
}

// ListPermissions handles the request to list all permissions
func (c *PermissionController) ListPermissions(ctx echo.Context) error {
	// Call usecase to get list of permissions
	output, err := c.permissionUsecase.ListPermissions(ctx.Request().Context())
	if err != nil {
		return response.SendError(ctx, errors.InternalErrorWithCause(messages.MsgListPermissionsError, err))
	}

	// Map usecase output to response format
	responseData := mapper.ToPermissionListResponse(output)
	return response.SendOK(ctx, messages.MsgListPermissionsSuccess, responseData)
}
