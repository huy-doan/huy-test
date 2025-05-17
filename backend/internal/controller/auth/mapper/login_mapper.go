package mapper

import (
	"github.com/huydq/test/internal/datastructure/inputdata"
	"github.com/huydq/test/internal/domain/model/user"
	generated "github.com/huydq/test/internal/pkg/api/generated"
	utils "github.com/huydq/test/internal/pkg/utils"
	"github.com/labstack/echo/v4"
)

type AuthMapper struct {
	ctx echo.Context
}

func NewAuthMapper(ctx echo.Context) *AuthMapper {
	return &AuthMapper{
		ctx: ctx,
	}
}

func (m *AuthMapper) ToLoginSuccessData(token string, user *user.User) *generated.LoginResponse {
	response := &generated.LoginResponse{
		Token: utils.ToPtr(token),
		User: &generated.User{
			Id:         utils.ToPtr(user.ID),
			Email:      utils.ToPtr(user.Email),
			EnabledMfa: utils.ToPtr(user.EnabledMFA),
			FullName:   utils.ToPtr(user.FullName),
			Role: &generated.Role{
				Id:   utils.ToPtr(user.Role.ID),
				Name: utils.ToPtr(user.Role.Name),
			},
			MfaType: &generated.MfaType{
				Id:       utils.ToPtr(user.MFAType),
				IsActive: utils.ToPtr(true),
				Title:    utils.ToPtr("Email"),
			},
		},
	}

	return response
}

func (m *AuthMapper) ToLoginInputData(req generated.LoginRequest) *inputdata.LoginInputData {
	return &inputdata.LoginInputData{
		Email:     string(req.Email),
		Password:  req.Password,
		IPAddress: m.ctx.RealIP(),
		UserAgent: m.ctx.Request().UserAgent(),
	}
}
