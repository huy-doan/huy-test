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

	Email         string  `json:"email"`
	PasswordHash  string  `json:"-"` // Never exposed in JSON
	RoleID        int     `json:"role_id"`
	Role          *Role   `json:"role"`
	EnabledMFA    bool    `json:"enabled_mfa"`
	MFAType       int     `json:"mfa_type"`
	LastName      string  `json:"last_name"`
	FirstName     string  `json:"first_name"`
	LastNameKana  string  `json:"last_name_kana"`
	FirstNameKana string  `json:"first_name_kana"`
	AvatarURL     *string `json:"avatar_url,omitempty"`
}

// TableName specifies the database table name
func (User) TableName() string {
	return "users"
}

// FullName returns the user's full name
func (u *User) FullName() string {
	return u.FirstName + " " + u.LastName
}

// FullNameKana returns the user's full name in kana
func (u *User) FullNameKana() string {
	return u.FirstNameKana + " " + u.LastNameKana
}

// NewUser creates a new user with the given details
func NewUser(email, password, firstName, lastName, firstNameKana, lastNameKana string, roleID int) (*User, error) {
	// Basic validation
	if email == "" {
		return nil, errors.New("email cannot be empty")
	}
	if password == "" {
		return nil, errors.New("password cannot be empty")
	}
	if firstName == "" || lastName == "" {
		return nil, errors.New("first name and last name cannot be empty")
	}
	if firstNameKana == "" || lastNameKana == "" {
		return nil, errors.New("first name kana and last name kana cannot be empty")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &User{
		Email:         email,
		PasswordHash:  string(hashedPassword),
		RoleID:        roleID,
		EnabledMFA:    true, // Default to enabled
		FirstName:     firstName,
		LastName:      lastName,
		FirstNameKana: firstNameKana,
		LastNameKana:  lastNameKana,
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
func (u *User) UpdateProfile(firstName, lastName, firstNameKana, lastNameKana string) error {
	if firstName == "" || lastName == "" {
		return errors.New("first name and last name cannot be empty")
	}

	if firstNameKana == "" || lastNameKana == "" {
		return errors.New("first name kana and last name kana cannot be empty")
	}

	u.FirstName = firstName
	u.LastName = lastName
	u.FirstNameKana = firstNameKana
	u.LastNameKana = lastNameKana
	u.UpdatedAt = time.Now()
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
