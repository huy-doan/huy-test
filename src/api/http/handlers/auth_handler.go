package handlers

import (
	"net/http"

	"github.com/vnlab/makeshop-payment/src/api/http/errors"
	"github.com/vnlab/makeshop-payment/src/api/http/middleware"
	"github.com/vnlab/makeshop-payment/src/api/http/response"
	"github.com/vnlab/makeshop-payment/src/api/http/serializers"
	validator "github.com/vnlab/makeshop-payment/src/api/http/validator/auth"
	"github.com/vnlab/makeshop-payment/src/infrastructure/auth"
	"github.com/vnlab/makeshop-payment/src/lib/i18n"
	"github.com/vnlab/makeshop-payment/src/lib/utils"
	"github.com/vnlab/makeshop-payment/src/usecase"
)

type AuthHandler struct {
	userUsecase    *usecase.UserUsecase
	jwtService     *auth.JWTService
	turnstileService *auth.TurnstileService
	auditLogUsecase *usecase.AuditLogUsecase
	lockedAccountUsecase *usecase.LockedAccountUsecase
}

func NewAuthHandler(
	userUsecase *usecase.UserUsecase,
	jwtService *auth.JWTService,
	turnstileService *auth.TurnstileService,
	auditLogUsecase *usecase.AuditLogUsecase,
	lockedAccountUsecase *usecase.LockedAccountUsecase,
) *AuthHandler {
	return &AuthHandler{
		userUsecase:     userUsecase,
		jwtService:      jwtService,
		turnstileService: turnstileService,
		auditLogUsecase:  auditLogUsecase,
		lockedAccountUsecase: lockedAccountUsecase,
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
	
	var err = h.lockedAccountUsecase.CheckAccountStatus(r.Context(), req.Email)
	if err != nil {
		response.Unauthorized(w, i18n.T(r.Context(), err.Error()))
		return
	}

	// Verify Turnstile token
	valid, err := h.turnstileService.VerifyToken(req.TurnstileToken)

	if err != nil {
		response.BadRequest(w, i18n.T(r.Context(), "turnstile.verify_failed"), nil)
		return
	}

	if !valid {
		response.BadRequest(w, i18n.T(r.Context(), "turnstile.invalid_token"), nil)
		return
	}

	loginReq := usecase.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
		TurnstileToken: req.TurnstileToken,
	}

	result, err := h.userUsecase.Login(r.Context(), loginReq)
	if err != nil {
		_ = h.lockedAccountUsecase.HandleFailedLogin(r.Context(), req.Email)
		var accountStatusErr = h.lockedAccountUsecase.CheckAccountStatus(r.Context(), req.Email)
		if accountStatusErr != nil {
			response.Unauthorized(w, i18n.T(r.Context(), accountStatusErr.Error()))
			return
		}

		tempLockRemaining, permLockRemaining, countErr := h.lockedAccountUsecase.GetRemainingAttempts(r.Context(), req.Email)
		if countErr == nil {
			if tempLockRemaining >= 0 {
				response.Unauthorized(w, i18n.T(r.Context(), "login.attempts_remaining_temp_lock", tempLockRemaining))
				return
			}
	
			if permLockRemaining >= 0 {
				response.Unauthorized(w, i18n.T(r.Context(), "login.attempts_remaining_perm_lock", permLockRemaining))
				return
			}
		}

		response.Unauthorized(w, i18n.T(r.Context(), "login.failed"))
		return
	}

	// Log successful login event
	ipAddress := r.RemoteAddr
	userAgent := r.UserAgent()
	userID := result.User.ID
	err = h.auditLogUsecase.LogLoginEvent(r.Context(), &userID, &ipAddress, &userAgent)
	h.lockedAccountUsecase.UnlockAccountByEmail(r.Context(), req.Email)
	if err != nil {
		// TODO: log error
	}

	// We transform the login response directly here since it's a special case
	// containing both user data and a token
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
