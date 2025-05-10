package mapper

import (
	"fmt"

	"github.com/huydq/test/internal/datastructure/inputdata"
	"github.com/labstack/echo/v4"
)

// VerifyMFARequest represents the request for verifying MFA token
type VerifyMFARequest struct {
	Email string `json:"email" validate:"required,email"`
	Token string `json:"token" validate:"required"`
}

// VerifyMFAResponse represents the response for MFA verification
type VerifyMFAResponse struct {
	Success bool           `json:"success"`
	Message string         `json:"message"`
	Data    *VerifyMFAData `json:"data"`
}

// VerifyMFAData represents the data in the MFA verification response
type VerifyMFAData struct {
	Token string   `json:"token"`
	User  UserData `json:"user"`
}

// VerifyMFAMapper handles mapping for MFA verification
type VerifyMFAMapper struct {
	ctx echo.Context
}

// NewVerifyMFAMapper creates a new mapper for MFA verification
func NewVerifyMFAMapper(ctx echo.Context) *VerifyMFAMapper {
	return &VerifyMFAMapper{
		ctx: ctx,
	}
}

// ToVerifyMFAInputData maps request to input data
func (m *VerifyMFAMapper) ToVerifyMFAInputData(req VerifyMFARequest) *inputdata.VerifyTwoFAInputData {
	return &inputdata.VerifyTwoFAInputData{
		Email: req.Email,
		Token: req.Token,
	}
}

// ToVerifyMFAResponse maps output data to response
func (m *VerifyMFAMapper) ToVerifyMFAResponse(data map[string]interface{}) *VerifyMFAResponse {
	userData := data["user"].(map[string]interface{})

	return &VerifyMFAResponse{
		Success: true,
		Message: "ログインに成功しました",
		Data: &VerifyMFAData{
			Token: data["token"].(string),
			User: UserData{
				ID:       fmt.Sprintf("%v", userData["id"]),
				Email:    userData["email"].(string),
				FullName: userData["full_name"].(string),
			},
		},
	}
}
