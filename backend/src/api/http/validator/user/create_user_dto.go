package validator

import "github.com/vnlab/makeshop-payment/src/lib/validator"

// CreateUserRequest represents the request body for creating a new user
type CreateUserRequest struct {
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password" validate:"required,password_policy"`
	FullName   string `json:"full_name" validate:"required"`
	RoleID     int    `json:"role_id,omitempty" validate:"omitempty,min=1"`
	EnabledMFA bool   `json:"enabled_mfa"`
}

// Validate performs validation on the CreateUserRequest
func (r *CreateUserRequest) Validate() error {
	return validator.GetValidate().Struct(r)
}
