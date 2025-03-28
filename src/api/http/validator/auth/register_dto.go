package validator

import "github.com/vnlab/makeshop-payment/src/lib/validator"

type RegisterRequest struct {
	Email         string `json:"email" binding:"required,email"`
	Password      string `json:"password" binding:"required,min=6,max=40"`
	FirstName     string `json:"first_name" binding:"required,min=1,max=191"`
	LastName      string `json:"last_name" binding:"required,min=1,max=191"`
	FirstNameKana string `json:"first_name_kana" binding:"required,min=1,max=191,kana"`
	LastNameKana  string `json:"last_name_kana" binding:"required,min=1,max=191,kana"`
}

func (r *RegisterRequest) Validate() error {
	v := validator.GetValidate()
	return v.Struct(r)
}
