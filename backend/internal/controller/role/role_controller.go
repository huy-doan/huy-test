package controller

import (
	"strconv"

	"github.com/huydq/test/internal/controller/base"
	"github.com/huydq/test/internal/controller/role/mapper"
	model "github.com/huydq/test/internal/domain/model/role"
	generated "github.com/huydq/test/internal/pkg/api/generated"
	"github.com/huydq/test/internal/pkg/common/response"
	"github.com/huydq/test/internal/pkg/errors"
	"github.com/huydq/test/internal/pkg/utils/messages"
	usecase "github.com/huydq/test/internal/usecase/role"
	"github.com/labstack/echo/v4"
)

type RoleController struct {
	base.BaseController
	roleUsecase usecase.RoleUsecase
}

func NewRoleController(roleUsecase usecase.RoleUsecase) *RoleController {
	return &RoleController{
		roleUsecase:    roleUsecase,
		BaseController: *base.NewBaseController(),
	}
}

func (c *RoleController) ListRoles(ctx echo.Context) error {
	// Get roles from usecase
	output, err := c.roleUsecase.ListRoles(ctx.Request().Context())
	if err != nil {
		return response.SendError(ctx, errors.InternalErrorWithCause(messages.MsgListRolesError, err))
	}

	responseData := mapper.ToRoleListResponse(output)
	return response.SendOK(ctx, messages.MsgListRolesSuccess, responseData)
}

func (c *RoleController) GetRoleByID(ctx echo.Context) error {
	id, err := c.GetIDParam(ctx, "id")
	if err != nil {
		return response.SendError(ctx, errors.BadRequestError(messages.MsgIDRequiredError, err))
	}

	// Get role from usecase
	output, err := c.roleUsecase.GetRoleByID(ctx.Request().Context(), id)
	if err != nil {
		return response.SendError(ctx, errors.InternalErrorWithCause(messages.MsgGetRoleError, err))
	}

	if output == nil {
		return response.SendError(ctx, errors.NotFoundError(messages.MsgRoleNotFoundError))
	}
	// Map to response format
	responseData := mapper.ToRoleResponse(output)
	return response.SendOK(ctx, messages.MsgGetRoleSuccess, responseData)
}

func (c *RoleController) CreateRole(ctx echo.Context) error {
	var request generated.CreateRoleRequest
	if err := c.BindAndValidate(ctx, &request); err != nil {
		return response.SendError(ctx, errors.FormatValidationError(err))
	}
	role, err := mapper.ToCreateRoleModel(request)
	if err != nil {
		return response.SendError(ctx, errors.BadRequestError(err.Error(), nil))
	}

	output, err := c.roleUsecase.CreateRole(ctx.Request().Context(), role)
	if err != nil {
		return response.SendError(ctx, errors.InternalErrorWithCause(messages.MsgCreateRoleError, err))
	}
	responseData := mapper.ToRoleResponse(output)
	return response.SendOK(ctx, messages.MsgCreateRoleSuccess, responseData)
}

func (c *RoleController) UpdateRole(ctx echo.Context) error {
	id, err := c.GetIDParam(ctx, "id")
	if err != nil {
		return response.SendError(ctx, errors.BadRequestError(messages.MsgIDRequiredError, err))
	}

	var request generated.UpdateRoleRequest
	if err := c.BindAndValidate(ctx, &request); err != nil {
		return response.SendError(ctx, errors.FormatValidationError(err))
	}

	input, err := mapper.ToUpdateRoleInput(request)
	if err != nil {
		return response.SendError(ctx, errors.BadRequestError(err.Error(), nil))
	}

	// Update role using usecase
	output, err := c.roleUsecase.UpdateRole(ctx.Request().Context(), id, input)
	if err != nil {
		return response.SendError(ctx, errors.InternalErrorWithCause(messages.MsgUpdateRoleError, err))
	}
	responseData := mapper.ToRoleResponse(output)
	return response.SendOK(ctx, messages.MsgUpdateRoleSuccess, responseData)
}

func (c *RoleController) DeleteRole(ctx echo.Context) error {
	id := ctx.Param("id")

	if id == "" {
		return response.SendError(ctx, errors.BadRequestError(messages.MsgIDRequiredError, nil))
	}
	idInt, err := strconv.Atoi(id)

	if err != nil {
		return response.SendError(ctx, errors.BadRequestError(messages.MsgIDRequiredError, err))
	}

	err = c.roleUsecase.DeleteRole(ctx.Request().Context(), idInt)
	if err != nil {
		return response.SendError(ctx, errors.InternalErrorWithCause(messages.MsgDeleteRoleError, err))
	}

	return response.SendOK(ctx, messages.MsgDeleteRoleSuccess, nil)
}

func (c *RoleController) BatchUpdateRolePermissions(ctx echo.Context) error {
	var request model.BatchUpdateRolePermissionsRequest
	if err := ctx.Bind(&request); err != nil {
		return response.SendError(ctx, errors.FormatValidationError(err))
	}
	input, err := mapper.ToBatchUpdateRolePermissionsInput(request)
	if err != nil {
		return response.SendError(ctx, errors.BadRequestError(err.Error(), nil))
	}

	output, err := c.roleUsecase.BatchUpdateRolePermissions(ctx.Request().Context(), input)
	if err != nil {
		return response.SendError(ctx, errors.InternalErrorWithCause(messages.MsgBatchUpdateRolePermissionsError, err))
	}
	responseData := mapper.ToBatchUpdateRolePermissionsResponse(output)
	return response.SendOK(ctx, messages.MsgBatchUpdateRolePermissionsSuccess, responseData)
}
