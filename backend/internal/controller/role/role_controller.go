package controller

import (
	"net/http"

	"github.com/huydq/test/internal/controller/role/mapper"
	usecase "github.com/huydq/test/internal/usecase/role"
	"github.com/labstack/echo/v4"
)

type RoleController struct {
	roleUsecase usecase.RoleUsecase
}

func NewRoleController(roleUsecase usecase.RoleUsecase) *RoleController {
	return &RoleController{
		roleUsecase: roleUsecase,
	}
}

// ListRoles handles the request to list all roles
func (c *RoleController) ListRoles(ctx echo.Context) error {
	// Get roles from usecase
	output, err := c.roleUsecase.ListRoles(ctx.Request().Context())
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "ロール一覧の取得に失敗しました",
		})
	}

	// Map to response format
	responseData := mapper.MapRoleListOutputToResponse(output)
	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "ロール一覧を取得しました",
		"data":    responseData,
	})
}

// GetRoleByID handles the request to get a role by ID
func (c *RoleController) GetRoleByID(ctx echo.Context) error {
	// Get role ID from path parameter
	id, err := mapper.ExtractIDFromPath(ctx)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "無効なID形式です",
		})
	}

	// Get role from usecase
	output, err := c.roleUsecase.GetRoleByID(ctx.Request().Context(), id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "ロールの取得に失敗しました",
		})
	}

	if output == nil {
		return ctx.JSON(http.StatusNotFound, map[string]interface{}{
			"success": false,
			"message": "指定されたロールが見つかりませんでした",
		})
	}

	// Map to response format
	responseData := mapper.MapRoleToResponse(output)
	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "ロールを取得しました",
		"data":    responseData,
	})
}

// CreateRole handles the request to create a new role
func (c *RoleController) CreateRole(ctx echo.Context) error {
	// Parse request body into input model
	input, err := mapper.MapRequestToCreateRoleInput(ctx)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "リクエストパラメータが無効です",
		})
	}

	output, err := c.roleUsecase.CreateRole(ctx.Request().Context(), input)
	if err != nil {
		switch err.Error() {
		case "role with name already exists":
			return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"message": "同じ名前のロールが既に存在します",
			})
		case "role with code already exists":
			return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"message": "同じコードのロールが既に存在します",
			})
		case "one or more permission IDs do not exist":
			return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"message": "1つ以上の権限IDが存在しません",
			})
		default:
			return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
				"success": false,
				"message": "ロールの作成に失敗しました",
			})
		}
	}

	// Map to response format
	responseData := mapper.MapRoleToResponse(output)
	return ctx.JSON(http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "ロールが作成されました",
		"data":    responseData,
	})
}

// UpdateRole handles the request to update an existing role
func (c *RoleController) UpdateRole(ctx echo.Context) error {
	// Get role ID from path parameter
	id, err := mapper.ExtractIDFromPath(ctx)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "無効なID形式です",
		})
	}

	// Parse request body into input model
	input, err := mapper.MapRequestToUpdateRoleInput(ctx)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "リクエストパラメータが無効です",
		})
	}

	// Update role using usecase
	output, err := c.roleUsecase.UpdateRole(ctx.Request().Context(), id, input)
	if err != nil {
		// Handle specific errors
		switch err.Error() {
		case "role with name already exists":
			return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"message": "同じ名前のロールが既に存在します",
			})
		case "role with code already exists":
			return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"message": "同じコードのロールが既に存在します",
			})
		case "one or more permission IDs do not exist":
			return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"message": "1つ以上の権限IDが存在しません",
			})
		case "role with ID not found":
			return ctx.JSON(http.StatusNotFound, map[string]interface{}{
				"success": false,
				"message": "指定されたロールが見つかりませんでした",
			})
		default:
			return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
				"success": false,
				"message": "ロールの更新に失敗しました",
			})
		}
	}

	// Map to response format
	responseData := mapper.MapRoleToResponse(output)
	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "ロールが更新されました",
		"data":    responseData,
	})
}

// DeleteRole handles the request to delete a role
func (c *RoleController) DeleteRole(ctx echo.Context) error {
	// Get role ID from path parameter
	id, err := mapper.ExtractIDFromPath(ctx)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "無効なID形式です",
		})
	}

	// Delete role using usecase
	err = c.roleUsecase.DeleteRole(ctx.Request().Context(), id)
	if err != nil {
		// Handle specific errors
		if err.Error() == "role with ID not found" {
			return ctx.JSON(http.StatusNotFound, map[string]interface{}{
				"success": false,
				"message": "指定されたロールが見つかりませんでした",
			})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "ロールの削除に失敗しました",
		})
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "ロールが削除されました",
	})
}

// BatchUpdateRolePermissions handles the request to update permissions for multiple roles at once
func (c *RoleController) BatchUpdateRolePermissions(ctx echo.Context) error {
	// Parse request body into input model
	input, err := mapper.MapRequestToBatchUpdateRolePermissionsInput(ctx)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "リクエストパラメータが無効です",
		})
	}

	// Update role permissions in batch using usecase
	output, err := c.roleUsecase.BatchUpdateRolePermissions(ctx.Request().Context(), input)
	if err != nil {
		if err.Error() == "one or more permission IDs do not exist" {
			return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"message": "1つ以上の権限IDが存在しません",
			})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "ロール権限の一括更新に失敗しました",
		})
	}

	// Map to response format
	responseData := mapper.MapBatchUpdateRolePermissionsOutputToResponse(output)
	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "ロール権限が一括更新されました",
		"data":    responseData,
	})
}
