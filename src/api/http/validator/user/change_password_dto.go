package validator

import "github.com/vnlab/makeshop-payment/src/lib/validator"

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=6"`
}

func (r *ChangePasswordRequest) Validate() error {
	v := validator.GetValidate()
	return v.Struct(r)
}
