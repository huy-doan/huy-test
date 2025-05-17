package user

import (
	"github.com/huydq/test/internal/controller/base"
	"github.com/huydq/test/internal/controller/user/mapper"
	generated "github.com/huydq/test/internal/pkg/api/generated"
	response "github.com/huydq/test/internal/pkg/common/response"
	"github.com/huydq/test/internal/pkg/errors"
	"github.com/huydq/test/internal/pkg/utils/messages"
	"github.com/huydq/test/internal/usecase/user"
	"github.com/labstack/echo/v4"
)

// UserController handles HTTP requests related to user management
type UserController struct {
	base.BaseController
	userUsecase user.UserManagementUsecase
}

// NewUserController creates a new user controller
func NewUserController(userUsecase user.UserManagementUsecase) *UserController {
	return &UserController{
		BaseController: *base.NewBaseController(),
		userUsecase:    userUsecase,
	}
}

// ListUsers handles the request to list users with pagination and filtering
func (c *UserController) ListUsers(ctx echo.Context) error {
	var request generated.UserListRequest
	if err := c.BindAndValidate(ctx, &request); err != nil {
		return response.SendError(ctx, err)
	}

	inputData := mapper.ToUserListInputData(&request)
	users, totalPages, totalCount, err := c.userUsecase.ListUsers(ctx.Request().Context(), inputData)
	if err != nil {
		return response.SendError(ctx, err)
	}

	userListSuccessData := mapper.ToUserListSuccessData(users, totalPages, totalCount, request.Page, request.PageSize)

	return response.SendOK(ctx, messages.MsgListUsersSuccess, userListSuccessData)
}

// CreateUser handles the request to create a new user
func (c *UserController) CreateUser(ctx echo.Context) error {
	var request generated.CreateUserRequest
	if err := c.BindAndValidate(ctx, &request); err != nil {
		return response.SendError(ctx, err)
	}

	inputData := mapper.ToCreateUserInputData(&request)
	newUser, err := c.userUsecase.CreateUser(ctx.Request().Context(), inputData)
	if err != nil {
		return response.SendError(ctx, errors.BadRequestError(messages.MsgCreateUserFailed, err.Error()))
	}

	userData := mapper.ToDetailedUserData(newUser)
	return response.SendCreated(ctx, messages.MsgCreateUserSuccess, userData)
}

// UpdateUser handles the request to update an existing user
func (c *UserController) UpdateUser(ctx echo.Context) error {
	userID, err := c.GetIDParam(ctx, "id")
	if err != nil {
		return response.SendError(ctx, err)
	}

	var request generated.UpdateUserRequest
	if err := c.BindAndValidate(ctx, &request); err != nil {
		return response.SendError(ctx, err)
	}

	inputData := mapper.ToUpdateUserInputData(&request)
	updatedUser, err := c.userUsecase.UpdateUser(ctx.Request().Context(), userID, inputData)
	if err != nil {
		return response.SendError(ctx, errors.BadRequestError(messages.MsgUpdateUserFailed, err.Error()))
	}

	userData := mapper.ToDetailedUserData(updatedUser)
	return response.SendOK(ctx, messages.MsgUpdateUserSuccess, userData)
}

// DeleteUser handles the request to delete a user
func (c *UserController) DeleteUser(ctx echo.Context) error {
	userID, err := c.GetIDParam(ctx, "id")
	if err != nil {
		return response.SendError(ctx, err)
	}

	err = c.userUsecase.DeleteUser(ctx.Request().Context(), userID)
	if err != nil {
		return response.SendError(ctx, errors.BadRequestError(messages.MsgDeleteUserFailed, err.Error()))
	}

	return response.SendOK(ctx, messages.MsgDeleteUserSuccess, nil)
}

// GetUserByID handles the request to get a user by ID
func (c *UserController) GetUserByID(ctx echo.Context) error {
	userID, err := c.GetIDParam(ctx, "id")
	if err != nil {
		return response.SendError(ctx, errors.BadRequestError(messages.MsgGetUserFailed, err.Error()))
	}

	user, err := c.userUsecase.GetUserByID(ctx.Request().Context(), userID)
	if err != nil {
		return response.SendError(ctx, errors.BadRequestError(messages.MsgGetUserFailed, err.Error()))
	}

	if user == nil {
		return response.SendError(ctx, errors.NotFoundError(messages.MsgUserNotFound))
	}

	userData := mapper.ToDetailedUserData(user)

	return response.SendOK(ctx, messages.MsgGetUserSuccess, userData)
}
