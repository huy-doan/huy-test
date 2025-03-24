package models

// Role represents a user role in the system
type Role struct {
	ID int `json:"id"`
	BaseColumnTimestamp

	Name string `json:"name"`
	Code string `json:"code"`
}

// TableName specifies the database table name
func (Role) TableName() string {
	return "roles"
}

// RoleCode defines constants for standard role codes
type RoleCode string

const (
	RoleCodeAdmin         RoleCode = "SYSTEM_ADMIN"
	RoleCodeNormalUser    RoleCode = "GENERAL_USER"
	RoleCodeBusinessUser  RoleCode = "BUSINESS_USER"
	RoleCodeAccoutingUser RoleCode = "ACCOUNTING_USER"
)

// IsAdmin checks if the role is an admin role
func (r *Role) IsAdmin() bool {
	return r.Code == string(RoleCodeAdmin)
}

// IsNormalUser checks if the role is a customer role
func (r *Role) IsNormalUser() bool {
	return r.Code == string(RoleCodeNormalUser)
}

// IsBusinessUser checks if the role is a business user role
func (r *Role) IsBusinessUser() bool {
	return r.Code == string(RoleCodeBusinessUser)
}

// IsAccountingUser checks if the role is an accounting user role
func (r *Role) IsAccountingUser() bool {
	return r.Code == string(RoleCodeAccoutingUser)
}
