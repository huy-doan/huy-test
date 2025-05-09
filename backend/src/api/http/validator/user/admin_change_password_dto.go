package validator

import "github.com/huydq/test/src/lib/validator"

type AdminChangePasswordRequest struct {
	UserID      int    `json:"user_id" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,password_policy"`
}

func (r *AdminChangePasswordRequest) Validate() error {
	v := validator.GetValidate()
	return v.Struct(r)
}
