package validator

import "github.com/vnlab/makeshop-payment/src/lib/validator"

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Validate kiểm tra dữ liệu của LoginRequest
func (r *LoginRequest) Validate() error {
	v := validator.GetValidate()
	return v.Struct(r)
}
