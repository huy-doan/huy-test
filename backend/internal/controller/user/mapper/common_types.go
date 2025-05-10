package mapper

import (
	"time"
)

type DetailedUserData struct {
	ID         int          `json:"id"`
	Email      string       `json:"email"`
	FullName   string       `json:"full_name"`
	EnabledMFA bool         `json:"enabled_mfa"`
	MFAType    *MFATypeData `json:"mfa_type,omitempty"`
	Role       *RoleData    `json:"role,omitempty"`
	CreatedAt  string       `json:"created_at"`
	UpdatedAt  string       `json:"updated_at"`
}

type MFATypeData struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	IsActive bool   `json:"is_active"`
}

type RoleData struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Code      string `json:"code"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

func FormatTime(t time.Time) string {
	return t.Format(time.RFC3339)
}
