package inputdata

// RegisterInputData represents user registration request data
type RegisterInputData struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
}

// UpdateProfileInputData represents a profile update request
type UpdateProfileInputData struct {
	FullName string `json:"full_name" binding:"required"`
}

// ChangePasswordInputData represents a password change request
type ChangePasswordInputData struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=6"`
}

// UserListInputData represents parameters for listing users
type UserListInputData struct {
	Page      int    `json:"page" default:"1" validate:"min=1"`
	PageSize  int    `json:"page_size" default:"10" validate:"min=1"`
	Search    string `json:"search" validate:"omitempty,max=255"`
	RoleID    *int   `json:"role_id" validate:"omitempty,min=1"`
	SortField string `json:"sort_field" validate:"omitempty"`
	SortOrder string `json:"sort_order" validate:"omitempty,oneof=asc desc"`
}

// CreateUserInputData represents admin user creation request
type CreateUserInputData struct {
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required,min=6"`
	FullName   string `json:"full_name" binding:"required"`
	RoleID     int    `json:"role_id" binding:"required"`
	EnabledMFA bool   `json:"enabled_mfa" default:"false"`
}

// UpdateUserInputData represents admin user update request
type UpdateUserInputData struct {
	Email      *string `json:"email,omitempty"`
	Password   *string `json:"password,omitempty"`
	FullName   *string `json:"full_name,omitempty"`
	RoleID     *int    `json:"role_id,omitempty"`
	EnabledMFA *bool   `json:"enabled_mfa,omitempty"`
}

// ResetPasswordInputData represents password reset request
type ResetPasswordInputData struct {
	NewPassword string `json:"new_password" binding:"required,min=6"`
}
