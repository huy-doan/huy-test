package handlers

import (
	"net/http"

	"github.com/huydq/test/src/api/http/errors"
	"github.com/huydq/test/src/api/http/middleware"
	"github.com/huydq/test/src/api/http/response"
	"github.com/huydq/test/src/api/http/serializers"
	validator "github.com/huydq/test/src/api/http/validator/auth"
	models "github.com/huydq/test/src/domain/models"
	"github.com/huydq/test/src/infrastructure/auth"
	"github.com/huydq/test/src/lib/i18n"
	"github.com/huydq/test/src/lib/utils"
	"github.com/huydq/test/src/usecase"
)

type AuthHandler struct {
	userUsecase     *usecase.UserUsecase
	jwtService      *auth.JWTService
	auditLogUsecase *usecase.AuditLogUsecase
	twoFAUsecase    *usecase.TwoFAUsecase
}

func NewAuthHandler(
	userUsecase *usecase.UserUsecase,
	jwtService *auth.JWTService,
	auditLogUsecase *usecase.AuditLogUsecase,
	twoFAUsecase *usecase.TwoFAUsecase,
) *AuthHandler {
	return &AuthHandler{
		userUsecase:     userUsecase,
		jwtService:      jwtService,
		auditLogUsecase: auditLogUsecase,
		twoFAUsecase:    twoFAUsecase,
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
		response.Unauthorized(w, i18n.T(r.Context(), "login.failed"))
		return
	}

	if result.User.EnabledMFA {
		verificationResp, err := h.twoFAUsecase.Generate2FAToken(
			r.Context(),
			result.User.ID,
			result.User.MFAType,
		)
		if err != nil {
			response.Error(w, errors.InternalError(i18n.T(r.Context(), "common.error")))
			return
		}

		responseData := map[string]interface{}{
			"requires_mfa": true,
			"user": map[string]interface{}{
				"email":    result.User.Email,
				"mfa_type": models.GetMFATypeTitle(verificationResp.MFAType),
			},
			"expires_in":  verificationResp.ExpiresIn,
			"mfa_type_id": verificationResp.MFAType,
		}

		response.Success(w, responseData, i18n.T(r.Context(), "login.mfa_required"))
		return
	}

	ipAddress := r.RemoteAddr
	userAgent := r.UserAgent()
	userID := result.User.ID
	err = h.auditLogUsecase.LogLoginEvent(r.Context(), &userID, &ipAddress, &userAgent)
	if err != nil {
		// TODO: log error
	}

	responseData := map[string]interface{}{
		"token": result.Token,
		"user":  serializers.NewUserSerializer(result.User).Serialize(),
	}
	response.Success(w, responseData, i18n.T(r.Context(), "login.success"))
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

	// Get user ID from context
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		response.Unauthorized(w, "No user ID found")
		return
	}

	// Log logout event
	ipAddress := r.RemoteAddr
	userAgent := r.UserAgent()
	err := h.auditLogUsecase.LogLogoutEvent(r.Context(), &userID, &ipAddress, &userAgent)
	if err != nil {
		// TODO: log error
	}

	h.jwtService.BlacklistToken(token)
	response.Success(w, nil, i18n.T(r.Context(), "logout.success"))
}

// Me handles getting current user information
// @Summary Get current user information
// @Description Get the current authenticated user's information
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=models.User}
// @Failure 401 {object} response.Response "Unauthorized"
// @Router /auth/me [get]
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		response.Unauthorized(w, i18n.T(r.Context(), "account.unauthorized"))
		return
	}

	// Get user information from usecase
	user, err := h.userUsecase.GetUserByID(r.Context(), userID)
	if err != nil {
		response.NotFound(w, i18n.T(r.Context(), "ユーザーが見つかりません."))
		return
	}

	response.Success(w, serializers.NewUserSerializer(user).Serialize(), i18n.T(r.Context(), "ユーザー情報を正常に取得しました"))
}

// VerifyMFA handles MFA verification
// @Summary Verify 2FA token
// @Description Verify a 2FA token to complete the login process
// @Tags auth
// @Accept json
// @Produce json
// @Param request body validator.VerifyRequest true "2FA verification details"
// @Success 200 {object} response.Response{data=usecase.VerifyResponse}
// @Failure 400 {object} response.Response "Bad Request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 500 {object} response.Response "Internal Server Error"
// @Router /auth/verify [post]
func (h *AuthHandler) VerifyMFA(w http.ResponseWriter, r *http.Request) {
	var req validator.VerifyRequest
	if err := utils.ParseJSONBody(r, &req); err != nil {
		response.ValidationError(w, err)
		return
	}

	if err := req.Validate(); err != nil {
		response.ValidationError(w, err)
		return
	}

	result, err := h.twoFAUsecase.Verify2FAToken(r.Context(), req)
	if err != nil {
		switch err.Error() {
		case "user.not_found":
			response.NotFound(w, i18n.T(r.Context(), "account.not_found"))
		case "mfa.invalid_token":
			response.BadRequest(w, i18n.T(r.Context(), "mfa.invalid_token"), nil)
		case "mfa.expired_token":
			response.BadRequest(w, i18n.T(r.Context(), "mfa.expired_token"), nil)
		default:
			response.Error(w, errors.InternalError(i18n.T(r.Context(), "common.error")))
		}
		return
	}

	ipAddress := r.RemoteAddr
	userAgent := r.UserAgent()
	userID := result.User.ID
	err = h.auditLogUsecase.LogLoginEvent(r.Context(), &userID, &ipAddress, &userAgent)
	if err != nil {
		// TODO: log error
	}

	responseData := map[string]interface{}{
		"token": result.Token,
		"user":  serializers.NewUserSerializer(result.User).Serialize(),
	}

	response.Success(w, responseData, i18n.T(r.Context(), "login.success"))
}
