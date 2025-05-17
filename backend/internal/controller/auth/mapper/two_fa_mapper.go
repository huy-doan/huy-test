package mapper

import (
	"github.com/huydq/test/internal/datastructure/inputdata"
	generated "github.com/huydq/test/internal/pkg/api/generated"
	utils "github.com/huydq/test/internal/pkg/utils"
	"github.com/labstack/echo/v4"
)

type TwoFAMapper struct {
	ctx echo.Context
}

func NewTwoFAMapper(ctx echo.Context) *TwoFAMapper {
	return &TwoFAMapper{
		ctx: ctx,
	}
}

func (m *TwoFAMapper) ToVerifyMFAInputData(req generated.VerifyMFARequest) *inputdata.VerifyTwoFAInputData {
	return &inputdata.VerifyTwoFAInputData{
		Email: string(req.Email),
		Token: req.Token,
	}
}

func (m *TwoFAMapper) ToMFARequiredData(
	email string,
	mfaType string,
	expiresIn int,
) *generated.RequiredTwoFaResponse {
	return &generated.RequiredTwoFaResponse{
		User: &generated.User{
			Email: utils.ToPtr(email),
		},
		ExpiresIn:   utils.ToPtr(expiresIn),
		RequiresMfa: utils.ToPtr(true),
		MfaType:     utils.ToPtr(mfaType),
	}
}
