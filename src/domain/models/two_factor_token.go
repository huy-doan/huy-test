package models

import (
	"time"
)

type TwoFactorToken struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Token     string    `json:"-"`
	MFAType   int       `json:"mfa_type"`
	User      *User     `json:"-"`
	IsUsed    bool      `json:"is_used"`
	ExpiredAt time.Time `json:"expired_at"`

	BaseColumnTimestamp
}

func (TwoFactorToken) TableName() string {
	return "two_factor_tokens"
}

func (t *TwoFactorToken) IsValid() bool {
	return !t.IsUsed && time.Now().Before(t.ExpiredAt)
}

func (t *TwoFactorToken) MarkAsUsed() {
	t.IsUsed = true
}
