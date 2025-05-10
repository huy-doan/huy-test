package controller

import (
	"net/http"

	"github.com/huydq/test/internal/controller/permission/mapper"
	usecase "github.com/huydq/test/internal/usecase/permission"
	"github.com/labstack/echo/v4"
)

type PermissionController struct {
	permissionUsecase usecase.PermissionUsecase
}

func NewPermissionController(permissionUsecase usecase.PermissionUsecase) *PermissionController {
	return &PermissionController{
		permissionUsecase: permissionUsecase,
	}
}

// ListPermissions handles the request to list all permissions
func (c *PermissionController) ListPermissions(ctx echo.Context) error {
	// Call usecase to get list of permissions
	output, err := c.permissionUsecase.ListPermissions(ctx.Request().Context())
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "権限一覧の取得に失敗しました",
		})
	}

	// Map usecase output to response format
	responseData := mapper.MapPermissionListOutputToResponse(output)
	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "権限一覧を取得しました",
		"data":    responseData,
	})
}

// GetPermissionsByIDs handles the request to get specific permissions by IDs
func (c *PermissionController) GetPermissionsByIDs(ctx echo.Context) error {
	// Get permission IDs from request
	input, err := mapper.MapRequestToGetPermissionsByIDsInput(ctx)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "リクエストパラメータが無効です",
		})
	}

	// Call usecase to get permissions
	output, err := c.permissionUsecase.GetPermissionsByIDs(ctx.Request().Context(), input)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "権限の取得に失敗しました",
		})
	}

	// Map usecase output to response format
	responseData := mapper.MapPermissionsToResponse(output)
	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "権限を取得しました",
		"data":    responseData,
	})
}
