package validator

import "github.com/huydq/test/src/lib/validator"

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,password_policy"`
}

func (r *ChangePasswordRequest) Validate() error {
	v := validator.GetValidate()
	return v.Struct(r)
}
