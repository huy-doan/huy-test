package validator

import "github.com/huydq/test/src/lib/validator"

// CreateUserRequest represents the request body for creating a new user
type CreateUserRequest struct {
	Email         string `json:"email" validate:"required,email"`
	Password      string `json:"password" validate:"required,password_policy"`
	FirstName     string `json:"first_name" validate:"required"`
	LastName      string `json:"last_name" validate:"required"`
	FirstNameKana string `json:"first_name_kana" validate:"required,kana"`
	LastNameKana  string `json:"last_name_kana" validate:"required,kana"`
	RoleID        int    `json:"role_id,omitempty" validate:"omitempty,min=1"`
	EnabledMFA    bool   `json:"enabled_mfa"`
}

// Validate performs validation on the CreateUserRequest
func (r *CreateUserRequest) Validate() error {
	return validator.GetValidate().Struct(r)
}
