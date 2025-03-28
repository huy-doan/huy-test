package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/huydq/demo/src/api/http/errors"
	"github.com/huydq/demo/src/api/http/response"
	"github.com/huydq/demo/src/api/http/serializers"
	validator "github.com/huydq/demo/src/api/http/validator/auth"
	"github.com/huydq/demo/src/infrastructure/auth"
	"github.com/huydq/demo/src/usecase"
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
// @Success 200 {object} response.Response{data=usecase.LoginResponse}
// @Failure 400 {object} response.Response "Bad Request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 500 {object} response.Response "Internal Server Error"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req validator.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	if err := req.Validate(); err != nil {
		response.ValidationError(c, err)
		return
	}

	loginReq := usecase.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	result, err := h.userUsecase.Login(c, loginReq)
	if err != nil {
		response.Unauthorized(c, "Invalid email or password")
		return
	}

	// We transform the login response directly here since it's a special case
	// containing both user data and a token
	responseData := map[string]interface{}{
		"token": result.Token,
		"user":  serializers.NewUserSerializer(result.User).Serialize(),
	}
	response.Success(c, responseData, "Login successful")
}

// @Summary Register a new user
// @Description Register a new user with the provided details
// @Tags auth
// @Accept json
// @Produce json
// @Param request body validator.RegisterRequest true "User registration details"
// @Success 201 {object} response.Response{data=models.User}
// @Failure 400 {object} response.Response "Bad Request"
// @Failure 500 {object} response.Response "Internal Server Error"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req validator.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	if err := req.Validate(); err != nil {
		response.ValidationError(c, err)
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
		// Check for specific error cases
		if err.Error() == "email already exists" {
			response.BadRequest(c, "Email already exists", nil)
			return
		}
		response.Error(c, errors.InternalError("Failed to register user"))
		return
	}

	response.Created(c, serializers.NewUserSerializer(user).Serialize(), "User registered successfully")
}

// @Summary Logout a user
// @Description Logout a user and invalidate their token
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response "Success"
// @Failure 401 {object} response.Response "Unauthorized"
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	token, exists := c.Get("token")
	if !exists {
		response.Unauthorized(c, "No token found")
		return
	}

	tokenStr, ok := token.(string)
	if !ok {
		response.Error(c, errors.InternalError("Invalid token format"))
		return
	}

	h.jwtService.BlacklistToken(tokenStr)
	response.Success(c, nil, "Logged out successfully")
}
