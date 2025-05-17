package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"fmt"

	"github.com/huydq/test/src/api/http/errors"
	"github.com/huydq/test/src/api/http/middleware"
	"github.com/huydq/test/src/api/http/response"
	"github.com/huydq/test/src/api/http/serializers"
	validator "github.com/huydq/test/src/api/http/validator/user"
	"github.com/huydq/test/src/domain/models"
	"github.com/huydq/test/src/infrastructure/auth"
	"github.com/huydq/test/src/lib/i18n"
	"github.com/huydq/test/src/lib/utils"
	"github.com/huydq/test/src/usecase"
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
	if !ok || userID == 0 {
		response.Unauthorized(w, "Unauthorized")
		return
	}

	err := h.userUsecase.ChangePassword(r.Context(), userID, req.CurrentPassword, req.NewPassword)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	response.Success(w, nil, i18n.T(r.Context(), "password.reset_success"))
}

// GetUserByID handles getting a user by ID
// @Summary Get a user by ID
// @Description Get a user's details by their ID
// @Tags admin users management
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
// @Router /admin/users/{id} [get]
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		response.BadRequest(w, "Invalid URL path", nil)
		return
	}
	idStr := pathParts[len(pathParts)-1]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(w, i18n.T(r.Context(), "account.not_found"), nil)
		return
	}

	user, err := h.userUsecase.GetUserByID(r.Context(), id)
	if err != nil {
		response.Error(w, errors.InternalError(i18n.T(r.Context(), "common.error")))
		return
	}

	if user == nil {
		response.BadRequest(w, i18n.T(r.Context(), "account.not_found"), nil)
		return
	}

	successMsg := fmt.Sprintf(i18n.T(r.Context(), "common.success"), "ユーザ")
	response.Success(w, serializers.NewUserSerializer(user).Serialize(), i18n.T(r.Context(), successMsg))
}

// ListUsers handles listing users with pagination
// @Summary List users
// @Description Get a list of users with pagination
// @Tags admin users management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query integer false "Page number (default: 1)"
// @Param page_size query integer false "Page size (default: 10, max: 100)"
// @Success 200 {object} response.Response "Success"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 403 {object} response.Response "Forbidden"
// @Failure 500 {object} response.Response "Internal Server Error"
// @Router /admin/users [get]
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

// UpdateUserProfile handles admin user profile updates
// @Summary Update user profile by admin
// @Description Update a user's profile information (admin only)
// @Tags admin users management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param request body validator.UpdateUserRequest true "User update details"
// @Success 200 {object} response.Response{data=models.User}
// @Failure 400 {object} response.Response "Bad Request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 403 {object} response.Response "Forbidden"
// @Failure 404 {object} response.Response "Not Found"
// @Failure 500 {object} response.Response "Internal Server Error"
// @Router /admin/users/{id} [put]
func (h *UserHandler) UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.ExtractParamFromPath(r, "id")
	if err != nil {
		response.BadRequest(w, err.Error(), nil)
		return
	}

	var req validator.UpdateUserRequest
	if err := utils.ParseJSONBody(r, &req); err != nil {
		response.ValidationError(w, err)
		return
	}

	if err := req.Validate(); err != nil {
		response.ValidationError(w, err)
		return
	}

	updateReq := validator.UpdateUserRequest{
		LastName:      req.LastName,
		FirstName:     req.FirstName,
		LastNameKana:  req.LastNameKana,
		FirstNameKana: req.FirstNameKana,
		RoleID:        req.RoleID,
		Email:         req.Email,
		EnabledMFA:    req.EnabledMFA,
	}

	user, err := h.userUsecase.UpdateUser(r.Context(), userID, updateReq)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	successMsg := fmt.Sprintf(i18n.T(r.Context(), "common.success"), "ユーザ")
	response.Success(w, serializers.NewUserSerializer(user).Serialize(), successMsg)
}

// CreateUser handles the creation of a new user
// @Summary Create a new user
// @Description Create a new user with the provided information
// @Tags admin users management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user body validator.CreateUserRequest true "User information"
// @Success 201 {object} response.Response{data=models.User}
// @Failure 400 {object} response.Response "Bad Request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 403 {object} response.Response "Forbidden"
// @Failure 409 {object} response.Response "Conflict"
// @Failure 500 {object} response.Response "Internal Server Error"
// @Router /admin/users [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req validator.CreateUserRequest
	if err := utils.ParseJSONBody(r, &req); err != nil {
		response.ValidationError(w, err)
		return
	}

	if err := req.Validate(); err != nil {
		response.ValidationError(w, err)
		return
	}

	newUser, err := h.userUsecase.CreateUser(r.Context(), &req)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	successMsg := fmt.Sprintf(i18n.T(r.Context(), "common.success"), "ユーザ")
	response.Created(w, serializers.NewUserSerializer(newUser).Serialize(), successMsg)
}

// DeleteUser handles the deletion of a user
// @Summary Delete a user
// @Description Delete a user by ID
// @Tags admin users management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} response.Response "Success"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 403 {object} response.Response "Forbidden"
// @Failure 404 {object} response.Response "Not Found"
// @Failure 500 {object} response.Response "Internal Server Error"
// @Router /admin/users/{id} [delete]
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.ExtractParamFromPath(r, "id")
	if err != nil {
		response.NotFound(w, i18n.T(r.Context(), "account.not_found"))
		return
	}

	currentUserID := r.Context().Value(middleware.UserIDKey).(int)
	if currentUserID == userID {
		response.BadRequest(w, i18n.T(r.Context(), "cannot_delete_self"), nil)
		return
	}

	user, err := h.userUsecase.GetUserByID(r.Context(), userID)
	if err != nil {
		response.Error(w, errors.InternalError(i18n.T(r.Context(), "common.error")))
		return
	}

	if user == nil {
		response.NotFound(w, i18n.T(r.Context(), "account.not_found"))
		return
	}

	err = h.userUsecase.DeleteUser(r.Context(), userID)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	successMsg := fmt.Sprintf(i18n.T(r.Context(), "common.success"), "ユーザ")

	response.Success(w, nil, successMsg)
}

// ResetPasswordByAdmin handles resetting a user's password
// @Summary Reset user password
// @Description Reset a user's password by email
// @Tags admin users management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body validator.AdminChangePasswordRequest true "Password reset details"
// @Success 200 {object} response.Response "Success"
// @Failure 400 {object} response.Response "Bad Request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 403 {object} response.Response "Forbidden"
// @Failure 404 {object} response.Response "Not Found"
// @Failure 500 {object} response.Response "Internal Server Error"
// @Router /admin/users/reset-password [post]
func (h *UserHandler) ResetPasswordByAdmin(w http.ResponseWriter, r *http.Request) {
	var req validator.AdminChangePasswordRequest
	if err := utils.ParseJSONBody(r, &req); err != nil {
		response.ValidationError(w, err)
		return
	}

	if err := req.Validate(); err != nil {
		response.ValidationError(w, err)
		return
	}

	err := h.userUsecase.ResetPassword(r.Context(), req.UserID, req.NewPassword)

	if err != nil {
		h.handleError(w, r, err)
		return
	}

	response.Success(w, nil, i18n.T(r.Context(), "password.reset_success"))
}

// handleError centralizes error handling for user-related operations
func (h *UserHandler) handleError(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		return
	}

	switch err.Error() {
	case "account.not_found":
		response.NotFound(w, i18n.T(r.Context(), "account.not_found"))
	case "email.already_exists":
		response.BadRequest(w, i18n.T(r.Context(), "email.already_exists"), nil)
	case "role.not_found":
		response.NotFound(w, i18n.T(r.Context(), "role.not_found"))
	case "cannot_delete_self":
		response.BadRequest(w, i18n.T(r.Context(), "cannot_delete_self"), nil)
	case "password.current_incorrect":
		response.BadRequest(w, i18n.T(r.Context(), "password.current_incorrect"), nil)
	default:
		response.Error(w, errors.InternalError(i18n.T(r.Context(), "common.error")))
	}
}
