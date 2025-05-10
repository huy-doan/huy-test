package mapper

import (
	"github.com/huydq/test/internal/datastructure/inputdata"
)

type CreateUserRequest struct {
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password" validate:"required,min=6"`
	FullName   string `json:"full_name" validate:"required"`
	RoleID     int    `json:"role_id" validate:"required"`
	EnabledMFA bool   `json:"enabled_mfa"`
}

func (r *CreateUserRequest) ToCreateUserInputData() *inputdata.CreateUserInputData {
	return &inputdata.CreateUserInputData{
		Email:      r.Email,
		Password:   r.Password,
		FullName:   r.FullName,
		RoleID:     r.RoleID,
		EnabledMFA: r.EnabledMFA,
	}
}

type CreateUserResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	User    interface{} `json:"user"`
}

func (r *CreateUserResponse) ToResponseMap() map[string]interface{} {
	return map[string]interface{}{
		"success": r.Success,
		"message": r.Message,
		"user":    r.User,
	}
}

func NewCreateUserResponse(user interface{}) *CreateUserResponse {
	return &CreateUserResponse{
		Success: true,
		Message: "User created successfully",
		User:    user,
	}
}
