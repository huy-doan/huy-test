package mapper

import (
	"strconv"

	"github.com/huydq/test/internal/datastructure/inputdata"
	"github.com/huydq/test/internal/datastructure/outputdata"
	"github.com/labstack/echo/v4"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token       string   `json:"token"`
	User        UserData `json:"user"`
	RequiresMFA bool     `json:"requires_mfa,omitempty"`
	MFAInfo     *MFAInfo `json:"mfa_info,omitempty"`
}

type UserData struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
}

type MFAInfo struct {
	Type      string `json:"type"`
	ExpiresIn int64  `json:"expires_in"`
}

type AuthMapper struct {
	ctx echo.Context
}

func NewAuthMapper(ctx echo.Context) *AuthMapper {
	return &AuthMapper{
		ctx: ctx,
	}
}

func (m *AuthMapper) ToLoginInputData(req LoginRequest) *inputdata.LoginInputData {
	return &inputdata.LoginInputData{
		Email:     req.Email,
		Password:  req.Password,
		IPAddress: m.ctx.RealIP(),
		UserAgent: m.ctx.Request().UserAgent(),
	}
}

func (m *AuthMapper) ToLoginResponse(output *outputdata.LoginOutputData) LoginResponse {
	response := LoginResponse{
		Token: output.Token,
		User: UserData{
			ID:       strconv.Itoa(output.User.ID),
			Email:    output.User.Email,
			FullName: output.User.FullName,
		},
		RequiresMFA: output.RequiresMFA,
	}

	if output.RequiresMFA && output.MFAInfo != nil {
		response.MFAInfo = &MFAInfo{
			Type:      output.MFAInfo.Type,
			ExpiresIn: int64(output.MFAInfo.ExpiresIn),
		}
	}

	return response
}
