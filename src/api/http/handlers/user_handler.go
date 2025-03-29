package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/vnlab/makeshop-payment/src/api/http/errors"
	"github.com/vnlab/makeshop-payment/src/api/http/middleware"
	"github.com/vnlab/makeshop-payment/src/api/http/response"
	"github.com/vnlab/makeshop-payment/src/api/http/serializers"
	validator "github.com/vnlab/makeshop-payment/src/api/http/validator/user"
	"github.com/vnlab/makeshop-payment/src/domain/models"
	"github.com/vnlab/makeshop-payment/src/infrastructure/auth"
	"github.com/vnlab/makeshop-payment/src/lib/utils"
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

// GetProfile handles getting the authenticated user's profile
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
func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		response.Unauthorized(w, "Unauthorized")
		return
	}

	user, err := h.userUsecase.GetUserByID(r.Context(), userID)
	if err != nil {
		response.Error(w, errors.InternalError("Failed to get user profile"))
		return
	}

	if user == nil {
		response.NotFound(w, "User not found")
		return
	}

	response.Success(w, serializers.NewUserSerializer(user).Serialize(), "User profile retrieved successfully")
}

// UpdateProfile handles updating the authenticated user's profile
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
func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	var req validator.UpdateProfileRequest
	if err := utils.ParseJSONBody(r, &req); err != nil {
		response.ValidationError(w, err)
		return
	}

	if err := req.Validate(); err != nil {
		response.ValidationError(w, err)
		return
	}

	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		response.Unauthorized(w, "Unauthorized")
		return
	}

	updateReq := usecase.UpdateProfileRequest{
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		FirstNameKana: req.FirstNameKana,
		LastNameKana:  req.LastNameKana,
	}

	user, err := h.userUsecase.UpdateUserProfile(r.Context(), userID, updateReq)
	if err != nil {
		response.Error(w, errors.InternalError("Failed to update user profile"))
		return
	}

	response.Success(w, serializers.NewUserSerializer(user).Serialize(), "Profile updated successfully")
}

// ChangePassword handles changing the authenticated user's password
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
func (h *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var req validator.ChangePasswordRequest
	if err := utils.ParseJSONBody(r, &req); err != nil {
		response.ValidationError(w, err)
		return
	}

	if err := req.Validate(); err != nil {
		response.ValidationError(w, err)
		return
	}

	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		response.Unauthorized(w, "Unauthorized")
		return
	}

	err := h.userUsecase.ChangePassword(r.Context(), userID, req.CurrentPassword, req.NewPassword)
	if err != nil {
		if err.Error() == "current password is incorrect" {
			response.BadRequest(w, "Current password is incorrect", nil)
			return
		}
		response.Error(w, errors.InternalError("Failed to change password"))
		return
	}

	response.Success(w, nil, "Password changed successfully")
}

// GetUserByID handles getting a user by ID
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
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	// Extract the ID from the URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		response.BadRequest(w, "Invalid URL path", nil)
		return
	}
	idStr := pathParts[len(pathParts)-1]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(w, "Invalid user ID", nil)
		return
	}

	user, err := h.userUsecase.GetUserByID(r.Context(), id)
	if err != nil {
		response.Error(w, errors.InternalError("Failed to get user"))
		return
	}

	if user == nil {
		response.NotFound(w, "User not found")
		return
	}

	response.Success(w, serializers.NewUserSerializer(user).Serialize(), "User retrieved successfully")
}

// ListUsers handles listing users with pagination
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
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("page_size")

	page := 1
	if pageStr != "" {
		if pageVal, err := strconv.Atoi(pageStr); err == nil && pageVal > 0 {
			page = pageVal
		}
	}

	pageSize := 10
	if pageSizeStr != "" {
		if pageSizeVal, err := strconv.Atoi(pageSizeStr); err == nil && pageSizeVal > 0 && pageSizeVal <= 100 {
			pageSize = pageSizeVal
		}
	}

	// Check admin role
	roleCode, ok := r.Context().Value(middleware.RoleCodeKey).(string)
	if !ok || roleCode != string(models.RoleCodeAdmin) {
		response.Forbidden(w, "Permission denied")
		return
	}

	users, totalPages, err := h.userUsecase.ListUsers(r.Context(), page, pageSize)
	if err != nil {
		response.Error(w, errors.InternalError("Failed to list users"))
		return
	}

	responseData := map[string]interface{}{
		"users":       serializers.SerializeUserCollection(users),
		"page":        page,
		"page_size":   pageSize,
		"total_pages": totalPages,
	}

	response.Success(w, responseData, "Users retrieved successfully")
}
