package validator

import (
	"time"

	"github.com/vnlab/makeshop-payment/src/lib/validator"
	"errors"
	"github.com/vnlab/makeshop-payment/src/lib/i18n"
)
// UpdateLockedAccountRequest represents an admin locked account update/create request
type UpdateLockedAccountRequest struct {
	LockedAt  *time.Time `json:"locked_at,omitempty"`
	ExpiredAt *time.Time `json:"expired_at,omitempty"`
	Count     int       `json:"count" validate:"required,min=0"`
	Email     string     `json:"email" validate:"required,email"`
}

func (r *UpdateLockedAccountRequest) Validate() error {
	v := validator.GetValidate()

	if r.LockedAt != nil {
		if r.ExpiredAt != nil && r.LockedAt.After(*r.ExpiredAt) {
			return errors.New(i18n.T(nil, "locked_account.locked_at_after_expired_at"))
		}
	} else {
		if r.ExpiredAt != nil {
			return errors.New(i18n.T(nil, "locked_account.expired_at_required"))
		}
	}

	return v.Struct(r)
} 
