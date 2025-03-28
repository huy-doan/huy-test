package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vnlab/makeshop-payment/src/api/http/errors"
	"github.com/vnlab/makeshop-payment/src/api/http/response"
	"github.com/vnlab/makeshop-payment/src/api/http/serializers"
	commonValidator "github.com/vnlab/makeshop-payment/src/api/http/validator/common"
	validator "github.com/vnlab/makeshop-payment/src/api/http/validator/user"
	"github.com/vnlab/makeshop-payment/src/domain/models"
	"github.com/vnlab/makeshop-payment/src/infrastructure/auth"
	"github.com/vnlab/makeshop-payment/src/usecase"
)

type UserHandler struct {
	userUsecase *usecase.UserUsecase
	jwtService  *auth.JWTService
}

func NewUserHandler(userUsecase *usecase.UserUsecase, jwtService *auth.JWTService) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
		jwtService:  jwtService,
	}
}

// @Summary Get user profile
// @Description Get the profile of the authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=models.User}
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 500 {object} response.Response "Internal Server Error"
// @Router /users/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, _ := c.Get("userId")
	id, ok := userID.(int)
	if !ok {
		response.Unauthorized(c, "Unauthorized")
		return
	}

	user, err := h.userUsecase.GetUserByID(c, id)
	if err != nil {
		response.Error(c, errors.InternalError("Failed to get user profile"))
		return
	}

	if user == nil {
		response.NotFound(c, "User not found")
		return
	}

	response.Success(c, serializers.NewUserSerializer(user).Serialize(), "User profile retrieved successfully")
}

// @Summary Update user profile
// @Description Update the profile of the authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body validator.UpdateProfileRequest true "Profile update details"
// @Success 200 {object} response.Response{data=models.User}
// @Failure 400 {object} response.Response "Bad Request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 500 {object} response.Response "Internal Server Error"
// @Router /users/profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	var req validator.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	if err := req.Validate(); err != nil {
		response.ValidationError(c, err)
		return
	}

	userID, _ := c.Get("userId")
	id, ok := userID.(int)
	if !ok {
		response.Unauthorized(c, "Unauthorized")
		return
	}

	updateReq := usecase.UpdateProfileRequest{
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		FirstNameKana: req.FirstNameKana,
		LastNameKana:  req.LastNameKana,
	}

	user, err := h.userUsecase.UpdateUserProfile(c, id, updateReq)
	if err != nil {
		response.Error(c, errors.InternalError("Failed to update user profile"))
		return
	}

	response.Success(c, serializers.NewUserSerializer(user).Serialize(), "Profile updated successfully")
}

// @Summary Change user password
// @Description Change the password of the authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body validator.ChangePasswordRequest true "Password change details"
// @Success 200 {object} response.Response "Success"
// @Failure 400 {object} response.Response "Bad Request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 500 {object} response.Response "Internal Server Error"
// @Router /users/change-password [post]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	var req validator.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	if err := req.Validate(); err != nil {
		response.ValidationError(c, err)
		return
	}

	userID, _ := c.Get("userId")
	id, ok := userID.(int)
	if !ok {
		response.Unauthorized(c, "Unauthorized")
		return
	}

	err := h.userUsecase.ChangePassword(c, id, req.CurrentPassword, req.NewPassword)
	if err != nil {
		if err.Error() == "current password is incorrect" {
			response.BadRequest(c, "Current password is incorrect", nil)
			return
		}
		response.Error(c, errors.InternalError("Failed to change password"))
		return
	}

	response.Success(c, nil, "Password changed successfully")
}

// @Summary Get a user by ID
// @Description Get a user's details by their ID
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} response.Response{data=models.User}
// @Failure 400 {object} response.Response "Bad Request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 403 {object} response.Response "Forbidden"
// @Failure 404 {object} response.Response "Not Found"
// @Failure 500 {object} response.Response "Internal Server Error"
// @Router /users/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid user ID", nil)
		return
	}

	user, err := h.userUsecase.GetUserByID(c, id)
	if err != nil {
		response.Error(c, errors.InternalError("Failed to get user"))
		return
	}

	if user == nil {
		response.NotFound(c, "User not found")
		return
	}

	response.Success(c, serializers.NewUserSerializer(user).Serialize(), "User retrieved successfully")
}

// @Summary List users
// @Description Get a list of users with pagination
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query integer false "Page number (default: 1)"
// @Param page_size query integer false "Page size (default: 10, max: 100)"
// @Success 200 {object} response.Response "Success"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 403 {object} response.Response "Forbidden"
// @Failure 500 {object} response.Response "Internal Server Error"
// @Router /users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	var req commonValidator.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	if err := req.Validate(); err != nil {
		response.ValidationError(c, err)
		return
	}

	page := 1
	if req.Page > 0 {
		page = req.Page
	}

	pageSize := 10
	if req.PageSize > 0 {
		pageSize = req.PageSize
	}

	// Check admin role
	roleCode, exists := c.Get("roleCode")
	if !exists || roleCode != string(models.RoleCodeAdmin) {
		response.Forbidden(c, "Permission denied")
		return
	}

	users, totalPages, err := h.userUsecase.ListUsers(c, page, pageSize)
	if err != nil {
		response.Error(c, errors.InternalError("Failed to list users"))
		return
	}

	response.Success(c, gin.H{
		"users":       serializers.SerializeUserCollection(users),
		"page":        page,
		"page_size":   pageSize,
		"total_pages": totalPages,
	}, "Users retrieved successfully")
}
