package user

import (
	"github.com/huydq/test/internal/controller/base"
	"github.com/huydq/test/internal/controller/user/mapper"
	"github.com/huydq/test/internal/pkg/errors"
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
	var request mapper.UserListRequest
	if err := c.BindAndValidate(ctx, &request); err != nil {
		return c.SendError(ctx, err)
	}

	if request.Page <= 0 {
		request.Page = 1
	}
	if request.PageSize <= 0 {
		request.PageSize = 10
	}

	inputData := request.ToUserListInputData()
	users, totalPages, totalCount, err := c.userUsecase.ListUsers(ctx.Request().Context(), inputData)
	if err != nil {
		return c.SendError(ctx, err)
	}

	userListSuccessMapper := mapper.NewUserListSuccessMapper(ctx)
	response := userListSuccessMapper.ToUserListSuccessResponse(users, totalPages, totalCount, request.Page, request.PageSize)

	return c.SendOK(ctx, response)
}

// CreateUser handles the request to create a new user
func (c *UserController) CreateUser(ctx echo.Context) error {
	var request mapper.CreateUserRequest
	if err := c.BindAndValidate(ctx, &request); err != nil {
		return c.SendError(ctx, err)
	}

	inputData := request.ToCreateUserInputData()
	newUser, err := c.userUsecase.CreateUser(ctx.Request().Context(), inputData)
	if err != nil {
		return c.SendError(ctx, err)
	}

	createUserSuccessMapper := mapper.NewCreateUserSuccessMapper(ctx)
	response := createUserSuccessMapper.ToCreateUserSuccessResponse(newUser)

	return c.SendCreated(ctx, response)
}

// UpdateUser handles the request to update an existing user
func (c *UserController) UpdateUser(ctx echo.Context) error {
	userID, err := c.GetIDParam(ctx, "id")
	if err != nil {
		return c.SendError(ctx, err)
	}

	var request mapper.UpdateUserRequest
	if err := c.BindAndValidate(ctx, &request); err != nil {
		return c.SendError(ctx, err)
	}

	inputData := request.ToUpdateUserInputData()
	updatedUser, err := c.userUsecase.UpdateUser(ctx.Request().Context(), userID, inputData)
	if err != nil {
		return c.SendError(ctx, err)
	}

	updateUserSuccessMapper := mapper.NewUpdateUserSuccessMapper(ctx)
	response := updateUserSuccessMapper.ToUpdateUserSuccessResponse(updatedUser)

	return c.SendOK(ctx, response)
}

// DeleteUser handles the request to delete a user
func (c *UserController) DeleteUser(ctx echo.Context) error {
	userID, err := c.GetIDParam(ctx, "id")
	if err != nil {
		return c.SendError(ctx, err)
	}

	err = c.userUsecase.DeleteUser(ctx.Request().Context(), userID)
	if err != nil {
		return c.SendError(ctx, err)
	}

	response := mapper.NewDeleteUserResponse()
	return c.SendOK(ctx, response.ToResponseMap())
}

// GetUserByID handles the request to get a user by ID
func (c *UserController) GetUserByID(ctx echo.Context) error {
	userID, err := c.GetIDParam(ctx, "id")
	if err != nil {
		return c.SendError(ctx, err)
	}

	user, err := c.userUsecase.GetUserByID(ctx.Request().Context(), userID)
	if err != nil {
		return c.SendError(ctx, err)
	}

	if user == nil {
		return c.SendError(ctx, errors.NotFoundError("User not found"))
	}

	userDetailMapper := mapper.NewUserDetailMapper(ctx)
	response := userDetailMapper.ToUserDetailResponse(user)

	return c.SendOK(ctx, response)
}
