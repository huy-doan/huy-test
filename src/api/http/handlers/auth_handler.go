package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	validator "github.com/vnlab/makeshop-payment/src/api/http/validator/auth"
	"github.com/vnlab/makeshop-payment/src/infrastructure/auth"
	"github.com/vnlab/makeshop-payment/src/usecase"
)

type AuthHandler struct {
	userUsecase *usecase.UserUsecase
	jwtService  *auth.JWTService
}

func NewAuthHandler(userUsecase *usecase.UserUsecase, jwtService *auth.JWTService) *AuthHandler {
	return &AuthHandler{
		userUsecase: userUsecase,
		jwtService:  jwtService,
	}
}

// @Summary Login a user
// @Description Login with username and password to get an access token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body validator.LoginRequest true "Login credentials"
// @Success 200 {object} usecase.LoginResponse
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req validator.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	loginReq := usecase.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	result, err := h.userUsecase.Login(c, loginReq)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary Register a new user
// @Description Register a new user with the provided details
// @Tags auth
// @Accept json
// @Produce json
// @Param request body validator.RegisterRequest true "User registration details"
// @Success 201 {object} models.User
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req validator.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	registerReq := usecase.RegisterRequest{
		Email:         req.Email,
		Password:      req.Password,
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		FirstNameKana: req.FirstNameKana,
		LastNameKana:  req.LastNameKana,
	}

	user, err := h.userUsecase.Register(c, registerReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// @Summary Logout a user
// @Description Logout a user and invalidate their token
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	token, exists := c.Get("token")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no token found"})
		return
	}

	tokenStr, ok := token.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid token format"})
		return
	}

	h.jwtService.BlacklistToken(tokenStr)
	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}
