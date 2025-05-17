package auth

import (
	"github.com/huydq/test/internal/controller/auth/mapper"
	"github.com/huydq/test/internal/controller/base"
	userMapper "github.com/huydq/test/internal/controller/user/mapper"
	"github.com/huydq/test/internal/middleware"
	generated "github.com/huydq/test/internal/pkg/api/generated"
	response "github.com/huydq/test/internal/pkg/common/response"
	"github.com/huydq/test/internal/pkg/errors"
	"github.com/huydq/test/internal/pkg/utils/messages"
	authUC "github.com/huydq/test/internal/usecase/auth"
	"github.com/labstack/echo/v4"
)

type AuthController struct {
	base.BaseController
	authUsecase *authUC.AuthUsecase
}

func NewAuthController(
	authUsecase *authUC.AuthUsecase,
) *AuthController {
	return &AuthController{
		BaseController: *base.NewBaseController(),
		authUsecase:    authUsecase,
	}
}

// Login handles the login request
func (c *AuthController) Login(ctx echo.Context) error {
	var loginReq generated.LoginRequest
	if err := c.BindAndValidate(ctx, &loginReq); err != nil {
		return response.SendError(ctx, err)
	}
	loginInput := mapper.NewAuthMapper(ctx).ToLoginInputData(loginReq)

	// Call usecase
	loginOutput, err := c.authUsecase.Login(ctx.Request().Context(), loginInput)
	if err != nil {
		return response.SendError(ctx, errors.UnauthorizedError(messages.MsgLoginFailed))
	}

	ctx.Set(string(middleware.ContextKey_AuditLogTargetUserID), loginOutput.User.ID)

	if loginOutput.RequiresMFA {
		mfaRequiredData := mapper.NewTwoFAMapper(ctx).ToMFARequiredData(
			loginOutput.User.Email,
			loginOutput.MFAInfo.Type,
			loginOutput.MFAInfo.ExpiresIn,
		)

		return response.SendOK(ctx, messages.MsgMFARequired, mfaRequiredData)
	}

	authSuccessData := mapper.NewAuthMapper(ctx).ToLoginSuccessData(loginOutput.Token, loginOutput.User)
	return response.SendOK(ctx, messages.MsgLoginSuccess, authSuccessData)
}

// Logout handles the logout request
func (c *AuthController) Logout(ctx echo.Context) error {
	token := ctx.Get(string(middleware.ContextKey_AuthToken))
	if token == nil {
		return response.SendError(ctx, errors.UnauthorizedError(messages.MsgUnauthenticated))
	}

	err := c.authUsecase.Logout(ctx.Request().Context(), token.(string))

	if err != nil {
		return response.SendError(ctx, errors.InternalError(messages.MsgLogoutFailed))
	}

	return response.SendOK(ctx, messages.MsgLogoutSuccess, nil)
}

// Me handles the me request
func (c *AuthController) Me(ctx echo.Context) error {
	userID := ctx.Get(string(middleware.ContextKey_AuthUserIDKey))
	if userID == nil {
		return response.SendError(ctx, errors.UnauthorizedError(messages.MsgUnauthenticated))
	}

	// Call usecase
	profileOutput, err := c.authUsecase.GetMe(ctx.Request().Context(), userID.(int))
	if err != nil {
		return response.SendError(ctx, errors.NotFoundError(messages.MsgUserNotFound))
	}

	userData := userMapper.ToDetailedUserData(profileOutput.User)

	return response.SendOK(ctx, messages.MsgGetUserSuccess, userData)
}

// VerifyMFA handles the verify mfa request
func (c *AuthController) VerifyMFA(ctx echo.Context) error {
	var verifyReq generated.VerifyMFARequest
	if err := c.BindAndValidate(ctx, &verifyReq); err != nil {
		return response.SendError(ctx, err)
	}
	verifyInput := mapper.NewTwoFAMapper(ctx).ToVerifyMFAInputData(verifyReq)

	// Call usecase
	verifyOutput, err := c.authUsecase.Verify2FAToken(ctx.Request().Context(), verifyInput)
	if err != nil {
		return response.SendError(ctx, errors.BadRequestError(messages.MsgMFAVerifyFailed, nil))
	}

	authSuccessData := mapper.NewAuthMapper(ctx).ToLoginSuccessData(verifyOutput.Token, verifyOutput.User)
	return response.SendOK(ctx, messages.MsgLoginSuccess, authSuccessData)
}

// ResendCode handles the resend code request
func (c *AuthController) ResendCode(ctx echo.Context) error {
	var resendReq generated.ResendCodeRequest
	if err := c.BindAndValidate(ctx, &resendReq); err != nil {
		return response.SendError(ctx, err)
	}
	resendInput := mapper.NewResendCodeMapper(ctx).ToResendCodeInputData(resendReq)

	// Call usecase
	resendOutput, err := c.authUsecase.ResendCode(ctx.Request().Context(), resendInput)
	if err != nil {
		return response.SendError(ctx, errors.BadRequestError(messages.MsgResendCodeFailed, nil))
	}

	if !resendOutput.CanResend {
		responseCodeData := mapper.NewResendCodeMapper(ctx).ToResendCodeData(
			resendOutput.CanResend,
			resendOutput.RemainingTime,
			resendOutput.ExpiresIn,
		)

		return response.SendError(ctx, errors.BadRequestError(messages.MsgResendCodeFailed, responseCodeData))
	}

	responseCodeData := mapper.NewResendCodeMapper(ctx).ToResendCodeData(
		resendOutput.CanResend,
		resendOutput.RemainingTime,
		resendOutput.ExpiresIn,
	)

	return response.SendOK(ctx, messages.MsgAuthCodeSentSuccess, responseCodeData)
}
