package models

const (
	MFAStatusInactive = 0
	MFAStatusActive   = 1
)

// MFAType represents a multi-factor authentication type
type MFAType struct {
	ID int `json:"id"`
	BaseColumnTimestamp

	No       int    `json:"no"`
	Title    string `json:"title"`
	IsActive int    `json:"is_active"`
}

// TableName specifies the database table name
func (MFAType) TableName() string {
	return "master_mfa_types"
}

// IsActiveType checks if this MFA type is active
func (m *MFAType) IsActiveType() bool {
	return m.IsActive == MFAStatusActive
}
