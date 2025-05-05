package validator

import "github.com/vnlab/makeshop-payment/src/lib/validator"

type UpdateProfileRequest struct {
	FullName string `json:"full_name" validate:"required"`
}

func (r *UpdateProfileRequest) Validate() error {
	v := validator.GetValidate()
	return v.Struct(r)
}
