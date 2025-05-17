package mapper

import (
	"github.com/huydq/test/internal/datastructure/inputdata"
	object "github.com/huydq/test/internal/domain/object/mfa_type"
	generated "github.com/huydq/test/internal/pkg/api/generated"
	utils "github.com/huydq/test/internal/pkg/utils"
	"github.com/labstack/echo/v4"
)

type ResendCodeMapper struct {
	ctx echo.Context
}

func NewResendCodeMapper(ctx echo.Context) *ResendCodeMapper {
	return &ResendCodeMapper{
		ctx: ctx,
	}
}

func (m *ResendCodeMapper) ToResendCodeInputData(req generated.ResendCodeRequest) *inputdata.ResendCodeInputData {
	mfaType := int(object.MFA_TYPE_EMAIL)
	if req.MfaType != nil {
		mfaType = *req.MfaType
	}

	return &inputdata.ResendCodeInputData{
		Email:   string(req.Email),
		MfaType: mfaType,
	}
}

func (m *ResendCodeMapper) ToResendCodeData(canResend bool, remainingTime int, expiresIn int) *generated.ResendCodeResponse {
	if !canResend {
		return &generated.ResendCodeResponse{
			CanResend:     utils.ToPtr(canResend),
			RemainingTime: utils.ToPtr(remainingTime),
		}
	}

	return &generated.ResendCodeResponse{
		CanResend: utils.ToPtr(true),
		ExpiresIn: utils.ToPtr(expiresIn),
	}
}
