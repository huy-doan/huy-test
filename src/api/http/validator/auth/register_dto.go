package validator

import "github.com/huydq/test/src/lib/validator"

type RegisterRequest struct {
	Email         string `json:"email" validate:"required,email"`
	Password      string `json:"password" validate:"required,min=6,max=40"`
	FirstName     string `json:"first_name" validate:"required,min=1,max=191"`
	LastName      string `json:"last_name" validate:"required,min=1,max=191"`
	FirstNameKana string `json:"first_name_kana" validate:"required,min=1,max=191,kana"`
	LastNameKana  string `json:"last_name_kana" validate:"required,min=1,max=191,kana"`
}

func (r *RegisterRequest) Validate() error {
	v := validator.GetValidate()
	return v.Struct(r)
}
