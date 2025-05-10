package mapper

import (
	"github.com/huydq/test/internal/datastructure/outputdata"
	userModel "github.com/huydq/test/internal/domain/model/user"
	"github.com/labstack/echo/v4"
)

type UserProfileResponse struct {
	Success bool             `json:"success"`
	Message string           `json:"message"`
	Data    *UserProfileData `json:"data"`
}

type UserProfileData struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	RoleID   int    `json:"role_id"`
}

type UserProfileMapper struct {
	ctx echo.Context
}

func NewUserProfileMapper(ctx echo.Context) *UserProfileMapper {
	return &UserProfileMapper{
		ctx: ctx,
	}
}

func (m *UserProfileMapper) ToUserProfileResponse(output *outputdata.UserProfileOutputData) *UserProfileResponse {
	return &UserProfileResponse{
		Success: true,
		Message: "ユーザー情報を取得しました",
		Data:    m.mapUserToProfileData(output.User),
	}
}

func (m *UserProfileMapper) mapUserToProfileData(user *userModel.User) *UserProfileData {
	if user == nil {
		return nil
	}

	return &UserProfileData{
		ID:       user.ID,
		Email:    user.Email,
		FullName: user.FullName,
		RoleID:   user.RoleID,
	}
}
