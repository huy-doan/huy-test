package models

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User represents a user entity in the system
type User struct {
	ID int `json:"id"`
	BaseColumnTimestamp

	Email        string `json:"email"`
	PasswordHash string `json:"-"`
	RoleID       int    `json:"role_id"`
	Role         *Role  `json:"role"`
	EnabledMFA   bool   `json:"enabled_mfa"`
	MFAType      int    `json:"mfa_type"`
	FullName     string `json:"full_name"`
}

// TableName specifies the database table name
func (User) TableName() string {
	return "user"
}

// GetFullName returns the user's full name
func (u *User) GetFullName() string {
	return u.FullName
}

// NewUser creates a new user with the given details
func NewUser(email, password, fullName string, roleID int) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &User{
		Email:        email,
		PasswordHash: string(hashedPassword),
		RoleID:       roleID,
		EnabledMFA:   true,
		FullName:     fullName,
	}, nil
}

// VerifyPassword verifies the provided password against the stored hash
func (u *User) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

// ChangePassword changes the user's password
func (u *User) ChangePassword(newPassword string) error {
	if newPassword == "" {
		return errors.New("new password cannot be empty")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.PasswordHash = string(hashedPassword)

	return nil
}

// UpdateProfile updates the user's profile information
func (u *User) UpdateProfile(fullName string) error {
	u.FullName = fullName

	return nil
}

// SetMFA configures the MFA settings for a user
func (u *User) SetMFA(enabled bool, mfaType int) {
	u.EnabledMFA = enabled
	u.MFAType = mfaType
	u.UpdatedAt = time.Now()
}

// IsAdmin checks if the user has admin privileges
func (u *User) IsAdmin() bool {
	return u.Role != nil && u.Role.IsAdmin()
}

// IsNormalUser checks if the user is a customer
func (u *User) IsNormalUser() bool {
	return u.Role != nil && u.Role.IsNormalUser()
}
