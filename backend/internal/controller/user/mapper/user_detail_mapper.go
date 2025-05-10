package mapper

import (
	"github.com/huydq/test/internal/domain/model/user"
	"github.com/labstack/echo/v4"
)

// UserDetailResponse represents the standardized response for a single user's details
type UserDetailResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Data    *DetailedUserData `json:"data"`
}

// UserDetailMapper handles mapping for user detail responses
type UserDetailMapper struct {
	ctx echo.Context
}

// NewUserDetailMapper creates a new mapper for user detail responses
func NewUserDetailMapper(ctx echo.Context) *UserDetailMapper {
	return &UserDetailMapper{
		ctx: ctx,
	}
}

// ToUserDetailResponse creates a user detail response
func (m *UserDetailMapper) ToUserDetailResponse(user *user.User) *UserDetailResponse {
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

	return &UserDetailResponse{
		Success: true,
		Message: "ユーザー情報を正常に取得しました",
		Data: &DetailedUserData{
			ID:         user.ID,
			Email:      user.Email,
			FullName:   user.FullName,
			EnabledMFA: user.EnabledMFA,
			MFAType:    mfaTypeData,
			Role:       roleData,
		},
	}
}
