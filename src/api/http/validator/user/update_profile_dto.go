package validator

import "github.com/vnlab/makeshop-payment/src/lib/validator"

type UpdateProfileRequest struct {
	FirstName     string `json:"first_name" binding:"required"`
	LastName      string `json:"last_name" binding:"required"`
	FirstNameKana string `json:"first_name_kana" binding:"required"`
	LastNameKana  string `json:"last_name_kana" binding:"required"`
}

// Validate kiểm tra dữ liệu của UpdateProfileRequest
func (r *UpdateProfileRequest) Validate() error {
	v := validator.GetValidate()
	return v.Struct(r)
}
