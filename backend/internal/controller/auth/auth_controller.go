package auth

import (
	"net/http"

	"github.com/huydq/test/internal/controller/auth/mapper"
	"github.com/huydq/test/internal/controller/base"
	"github.com/huydq/test/internal/datastructure/inputdata"
	authService "github.com/huydq/test/internal/infrastructure/adapter/auth"
	"github.com/huydq/test/internal/middleware"
	"github.com/huydq/test/internal/pkg/errors"
	authUC "github.com/huydq/test/internal/usecase/auth"
	twoFAUC "github.com/huydq/test/internal/usecase/two_fa"
	"github.com/labstack/echo/v4"
)

type AuthController struct {
	base.BaseController
	authUsecase  *authUC.AuthUsecase
	jwtService   *authService.JWTService
	twoFAUsecase twoFAUC.TwoFAUsecase
}

func NewAuthController(
	authUsecase *authUC.AuthUsecase,
	jwtService *authService.JWTService,
	twoFAUsecase twoFAUC.TwoFAUsecase,
) *AuthController {
	return &AuthController{
		BaseController: *base.NewBaseController(),
		authUsecase:    authUsecase,
		jwtService:     jwtService,
		twoFAUsecase:   twoFAUsecase,
	}
}

// Login handles the login request
func (c *AuthController) Login(ctx echo.Context) error {
	var loginReq mapper.LoginRequest
	if err := c.BindAndValidate(ctx, &loginReq); err != nil {
		return c.SendError(ctx, err)
	}

	authMapper := mapper.NewAuthMapper(ctx)
	loginInput := authMapper.ToLoginInputData(loginReq)
	loginOutput, err := c.authUsecase.Login(ctx.Request().Context(), loginInput)
	if err != nil {
		if err.Error() == "login.failed" {
			return c.SendError(ctx, errors.UnauthorizedError("ログインに失敗しました"))
		}

		return c.SendError(ctx, err)
	}

	if loginOutput.User != nil && loginOutput.User.EnabledMFA {
		twofaInput := &inputdata.GenerateTwoFAInputData{
			UserID:  loginOutput.User.ID,
			MFAType: loginOutput.User.MFAType,
		}

		verificationResp, err := c.twoFAUsecase.Generate2FAToken(ctx.Request().Context(), twofaInput)
		if err != nil {
			return c.SendError(ctx, errors.InternalErrorWithCause("多要素認証トークンの生成に失敗しました", err))
		}

		mfaMapper := mapper.NewMFARequiredMapper(ctx)
		mfaResponse := mfaMapper.ToMFARequiredResponse(
			loginOutput.User.Email,
			"Email",
			verificationResp.ExpiresIn,
			verificationResp.MFAType,
		)

		return c.SendOK(ctx, mfaResponse)
	}

	authSuccessMapper := mapper.NewAuthSuccessMapper(ctx)
	response := authSuccessMapper.ToAuthSuccessResponse(loginOutput.Token, loginOutput.User)
	return c.SendOK(ctx, response)
}

// Logout handles the logout request
func (c *AuthController) Logout(ctx echo.Context) error {
	token := ctx.Get(string(middleware.ContextKey_AuthToken))
	if token == nil {
		return c.SendError(ctx, errors.UnauthorizedError("認証されていません"))
	}

	c.jwtService.BlacklistToken(token.(string))

	return c.SendOK(ctx, map[string]interface{}{
		"success": true,
		"message": "ログアウトしました",
	})
}

// Me handles the me request
func (c *AuthController) Me(ctx echo.Context) error {
	userID := ctx.Get(string(middleware.ContextKey_AuthUserIDKey))
	if userID == nil {
		return c.SendError(ctx, errors.UnauthorizedError("認証されていません"))
	}

	userIDInt, ok := userID.(int)
	if !ok {
		return c.SendError(ctx, errors.InternalError("無効なユーザーIDです"))
	}

	profileOutput, err := c.authUsecase.GetMe(ctx.Request().Context(), userIDInt)
	if err != nil {
		if err.Error() == "user.not_found" {
			return c.SendError(ctx, errors.NotFoundError("ユーザーが見つかりません"))
		}
		return c.SendError(ctx, err)
	}

	profileMapper := mapper.NewUserProfileMapper(ctx)
	response := profileMapper.ToUserProfileResponse(profileOutput)

	return c.SendOK(ctx, response)
}

// VerifyMFA handles the verify mfa request
func (c *AuthController) VerifyMFA(ctx echo.Context) error {
	var verifyReq mapper.VerifyMFARequest
	if err := c.BindAndValidate(ctx, &verifyReq); err != nil {
		return c.SendError(ctx, err)
	}

	verifyMapper := mapper.NewVerifyMFAMapper(ctx)
	verifyInput := verifyMapper.ToVerifyMFAInputData(verifyReq)
	verifyOutput, err := c.twoFAUsecase.Verify2FAToken(ctx.Request().Context(), verifyInput)

	if err != nil {
		switch err.Error() {
		case "user.not_found":
			return c.SendError(ctx, errors.NotFoundError("アカウントが見つかりません"))
		case "mfa.invalid_token":
			return c.SendError(ctx, errors.BadRequestError("無効な認証コードです", nil))
		case "mfa.expired_token":
			return c.SendError(ctx, errors.BadRequestError("認証コードの有効期限が切れています", nil))
		default:
			return c.SendError(ctx, err)
		}
	}

	authSuccessMapper := mapper.NewAuthSuccessMapper(ctx)
	response := authSuccessMapper.ToAuthSuccessResponse(verifyOutput.Token, verifyOutput.User)

	return c.SendOK(ctx, response)
}

// ResendCode handles the resend code request
func (c *AuthController) ResendCode(ctx echo.Context) error {
	var resendReq mapper.ResendCodeRequest
	if err := c.BindAndValidate(ctx, &resendReq); err != nil {
		return c.SendError(ctx, err)
	}

	user, err := c.authUsecase.FindUserByEmail(ctx.Request().Context(), resendReq.Email)
	if err != nil {
		return c.SendError(ctx, err)
	}

	if user == nil {
		return c.SendError(ctx, errors.NotFoundError("アカウントが見つかりません"))
	}

	resendMapper := mapper.NewResendCodeMapper(ctx)
	canResendInput := resendMapper.ToCanResendCodeInputData(resendReq, user.ID)
	canResendOutput, err := c.twoFAUsecase.CanResendCode(ctx.Request().Context(), canResendInput)
	if err != nil {
		return c.SendError(ctx, err)
	}

	if !canResendOutput.CanResend {
		response := resendMapper.ToResendCodeResponse(
			canResendOutput.CanResend,
			canResendOutput.RemainingTime,
			0,
		)
		return ctx.JSON(http.StatusTooManyRequests, response)
	}

	generateInput := resendMapper.ToGenerateTwoFAInputData(user.ID, resendReq.MFAType)
	generateOutput, err := c.twoFAUsecase.Generate2FAToken(ctx.Request().Context(), generateInput)
	if err != nil {
		return c.SendError(ctx, errors.InternalErrorWithCause("認証コードの生成に失敗しました", err))
	}

	response := resendMapper.ToResendCodeResponse(
		true,
		0,
		generateOutput.ExpiresIn,
	)
	return c.SendOK(ctx, response)
}
