package mapper

import (
	"github.com/huydq/test/internal/domain/model/user"
	"github.com/labstack/echo/v4"
)

type AuthSuccessResponse struct {
	Success bool             `json:"success"`
	Message string           `json:"message"`
	Data    *AuthSuccessData `json:"data"`
}

type AuthSuccessData struct {
	Token string            `json:"token"`
	User  *DetailedUserData `json:"user"`
}

type DetailedUserData struct {
	ID         int          `json:"id"`
	Email      string       `json:"email"`
	FullName   string       `json:"full_name"`
	EnabledMFA bool         `json:"enabled_mfa"`
	MFAType    *MFATypeData `json:"mfa_type,omitempty"`
	Role       *RoleData    `json:"role,omitempty"`
}

type MFATypeData struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	IsActive bool   `json:"is_active"`
}

type RoleData struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

type AuthSuccessMapper struct {
	ctx echo.Context
}

func NewAuthSuccessMapper(ctx echo.Context) *AuthSuccessMapper {
	return &AuthSuccessMapper{
		ctx: ctx,
	}
}

func (m *AuthSuccessMapper) ToAuthSuccessResponse(token string, user *user.User) *AuthSuccessResponse {
	var mfaTypeData *MFATypeData
	if user.EnabledMFA {
		mfaType := "Email"
		mfaTypeData = &MFATypeData{
			ID:       user.MFAType,
			Title:    mfaType,
			IsActive: true,
		}
	}

	var roleData *RoleData
	if user.Role != nil {
		roleData = &RoleData{
			ID:   user.Role.ID,
			Name: user.Role.Name,
			Code: string(user.Role.Code),
		}
	}

	return &AuthSuccessResponse{
		Success: true,
		Message: "ログインに成功しました",
		Data: &AuthSuccessData{
			Token: token,
			User: &DetailedUserData{
				ID:         user.ID,
				Email:      user.Email,
				FullName:   user.FullName,
				EnabledMFA: user.EnabledMFA,
				MFAType:    mfaTypeData,
				Role:       roleData,
			},
		},
	}
}
