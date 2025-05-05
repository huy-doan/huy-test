package validator

import "github.com/vnlab/makeshop-payment/src/lib/validator"

type SortingRequest struct {
	SortField string `form:"sort_field" validate:"omitempty"`
	SortOrder string `form:"sort_order" validate:"omitempty,oneof=ASC DESC"`
}

func (r *SortingRequest) Validate() error {
	v := validator.GetValidate()
	return v.Struct(r)
}
