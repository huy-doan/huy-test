package models

import "time"

// LockedAccount represents a locked user account
type LockedAccount struct {
	ID        int        `json:"id"`
	Email     string     `json:"email"`
	UserID    *int       `json:"user_id"`
	Count     int        `json:"count"`
	LockedAt  *time.Time `json:"locked_at,omitempty"`
	ExpiredAt *time.Time `json:"expired_at,omitempty"`
	BaseColumnTimestamp

	User *User `json:"user"`
}

// Constants for lock thresholds
const (
	TempLockThreshold    = 3 // Number of attempts before temporary lock
	PermLockThreshold    = 5 // Number of attempts before permanent lock
	TempLockDuration     = 5 // Duration in minutes for temporary lock
	MaxFailedMFAAttempts = 5 // Maximum failed 2FA attempts before permanent lock
)

// TableName specifies the table name for LockedAccount
func (la *LockedAccount) TableName() string {
	return "locked_accounts"
}

// IsTemporarilyLocked checks if the account is temporarily locked
func (la *LockedAccount) IsTemporarilyLocked() bool {
	if la.LockedAt != nil && la.ExpiredAt == nil {
		return false
	}

	if la.LockedAt == nil && la.ExpiredAt == nil {
		return false
	}

	return la.ExpiredAt.After(time.Now())
}

// IsPermanentlyLocked checks if the account is permanently locked
func (la *LockedAccount) IsPermanentlyLocked() bool {
	return la.LockedAt != nil && la.ExpiredAt == nil
}

// ShouldTemporarilyLock checks if the account should be temporarily locked
func (la *LockedAccount) ShouldTemporarilyLock() bool {
	return la.Count >= TempLockThreshold && la.Count < PermLockThreshold
}

// ShouldPermanentlyLock checks if the account should be permanently locked
func (la *LockedAccount) ShouldPermanentlyLock() bool {
	return la.Count >= PermLockThreshold || la.UserID == nil
}
