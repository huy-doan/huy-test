package validator

import "github.com/huydq/test/src/lib/validator"

type UpdateProfileRequest struct {
	FirstName     string `json:"first_name" validate:"required"`
	LastName      string `json:"last_name" validate:"required"`
	FirstNameKana string `json:"first_name_kana" validate:"required"`
	LastNameKana  string `json:"last_name_kana" validate:"required"`
}

func (r *UpdateProfileRequest) Validate() error {
	v := validator.GetValidate()
	return v.Struct(r)
}
