// src/api/http/handlers/user_handler.go
package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
// @Success 200 {object} models.User
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /users/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, _ := c.Get("userId")
	id, ok := userID.(int)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, err := h.userUsecase.GetUserByID(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// @Summary Update user profile
// @Description Update the profile of the authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body validator.UpdateProfileRequest true "Profile update details"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /users/profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	var req validator.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userId")
	id, ok := userID.(int)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// @Summary Change user password
// @Description Change the password of the authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body validator.ChangePasswordRequest true "Password change details"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /users/change-password [post]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	var req validator.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userId")
	id, ok := userID.(int)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	err := h.userUsecase.ChangePassword(c, id, req.CurrentPassword, req.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

// @Summary Get a user by ID
// @Description Get a user's details by their ID
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 404 {object} map[string]string "Not Found"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /users/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	user, err := h.userUsecase.GetUserByID(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// @Summary List users
// @Description Get a list of users with pagination
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query integer false "Page number (default: 1)"
// @Param page_size query integer false "Page size (default: 10, max: 100)"
// @Success 200 {object} object
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	var req commonValidator.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	// Kiểm tra vai trò admin
	roleCode, exists := c.Get("roleCode")
	if !exists || roleCode != string(models.RoleCodeAdmin) {
		c.JSON(http.StatusForbidden, gin.H{"error": "permission denied"})
		return
	}

	users, totalPages, err := h.userUsecase.ListUsers(c, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users":       users,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": totalPages,
	})
}
