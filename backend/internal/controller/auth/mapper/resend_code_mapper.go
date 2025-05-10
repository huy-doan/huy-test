package mapper

import (
	"github.com/huydq/test/internal/datastructure/inputdata"
	"github.com/labstack/echo/v4"
)

// ResendCodeRequest represents the request for resending an MFA verification code
type ResendCodeRequest struct {
	Email   string `json:"email" validate:"required,email"`
	MFAType int    `json:"mfa_type" validate:"required"`
}

// ResendCodeResponse represents the response for resending verification code
type ResendCodeResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    *ResendCodeData `json:"data,omitempty"`
}

// ResendCodeData represents data in the resend code response
type ResendCodeData struct {
	CanResend     bool  `json:"can_resend"`
	RemainingTime int   `json:"remaining_time,omitempty"`
	ExpiresIn     int64 `json:"expires_in,omitempty"`
}

// ResendCodeMapper handles mapping for code resend operations
type ResendCodeMapper struct {
	ctx echo.Context
}

// NewResendCodeMapper creates a new mapper for code resend operations
func NewResendCodeMapper(ctx echo.Context) *ResendCodeMapper {
	return &ResendCodeMapper{
		ctx: ctx,
	}
}

// ToCanResendCodeInputData maps request to input data
func (m *ResendCodeMapper) ToCanResendCodeInputData(req ResendCodeRequest, userID int) *inputdata.CanResendCodeInputData {
	return &inputdata.CanResendCodeInputData{
		UserID:  userID,
		MFAType: req.MFAType,
	}
}

// ToGenerateTwoFAInputData maps request to input data for token generation
func (m *ResendCodeMapper) ToGenerateTwoFAInputData(userID int, mfaType int) *inputdata.GenerateTwoFAInputData {
	return &inputdata.GenerateTwoFAInputData{
		UserID:  userID,
		MFAType: mfaType,
	}
}

// ToResendCodeResponse maps output data to response
func (m *ResendCodeMapper) ToResendCodeResponse(canResend bool, remainingTime int, expiresIn int64) *ResendCodeResponse {
	if !canResend {
		return &ResendCodeResponse{
			Success: false,
			Message: "リクエストが多すぎます。しばらく待ってから再度お試しください。",
			Data: &ResendCodeData{
				CanResend:     canResend,
				RemainingTime: remainingTime,
			},
		}
	}

	return &ResendCodeResponse{
		Success: true,
		Message: "認証コードが正常に送信されました",
		Data: &ResendCodeData{
			CanResend: true,
			ExpiresIn: expiresIn,
		},
	}
}
