package handlers

import (
	"net/http"

	"github.com/vnlab/makeshop-payment/src/api/http/errors"
	"github.com/vnlab/makeshop-payment/src/api/http/middleware"
	"github.com/vnlab/makeshop-payment/src/api/http/response"
	"github.com/vnlab/makeshop-payment/src/api/http/serializers"
	validator "github.com/vnlab/makeshop-payment/src/api/http/validator/auth"
	"github.com/vnlab/makeshop-payment/src/infrastructure/auth"
	"github.com/vnlab/makeshop-payment/src/lib/utils"
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

// Login handles user login
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
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req validator.LoginRequest
	if err := utils.ParseJSONBody(r, &req); err != nil {
		response.ValidationError(w, err)
		return
	}

	if err := req.Validate(); err != nil {
		response.ValidationError(w, err)
		return
	}

	loginReq := usecase.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	result, err := h.userUsecase.Login(r.Context(), loginReq)
	if err != nil {
		response.Unauthorized(w, "Invalid email or password")
		return
	}

	// We transform the login response directly here since it's a special case
	// containing both user data and a token
	responseData := map[string]interface{}{
		"token": result.Token,
		"user":  serializers.NewUserSerializer(result.User).Serialize(),
	}
	response.Success(w, responseData, "Login successful")
}

// Register handles user registration
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
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req validator.RegisterRequest
	if err := utils.ParseJSONBody(r, &req); err != nil {
		response.ValidationError(w, err)
		return
	}

	if err := req.Validate(); err != nil {
		response.ValidationError(w, err)
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

	user, err := h.userUsecase.Register(r.Context(), registerReq)
	if err != nil {
		// Check for specific error cases
		if err.Error() == "email already exists" {
			response.BadRequest(w, "Email already exists", nil)
			return
		}
		response.Error(w, errors.InternalError("Failed to register user"))
		return
	}

	response.Created(w, serializers.NewUserSerializer(user).Serialize(), "User registered successfully")
}

// Logout handles user logout
// @Summary Logout a user
// @Description Logout a user and invalidate their token
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response "Success"
// @Failure 401 {object} response.Response "Unauthorized"
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	token, ok := r.Context().Value(middleware.TokenKey).(string)
	if !ok {
		response.Unauthorized(w, "No token found")
		return
	}

	h.jwtService.BlacklistToken(token)
	response.Success(w, nil, "Logged out successfully")
}
