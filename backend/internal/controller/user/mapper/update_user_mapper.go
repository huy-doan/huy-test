package mapper

import (
	"github.com/huydq/test/internal/datastructure/inputdata"
)

// UpdateUserRequest represents the HTTP request for updating a user
type UpdateUserRequest struct {
	Email      *string `json:"email,omitempty" validate:"omitempty,email"`
	Password   *string `json:"password,omitempty" validate:"omitempty,min=6"`
	FullName   *string `json:"full_name,omitempty" validate:"omitempty"`
	RoleID     *int    `json:"role_id,omitempty" validate:"omitempty"`
	EnabledMFA *bool   `json:"enabled_mfa,omitempty"`
}

// ToUpdateUserInputData converts the request to UpdateUserInputData
func (r *UpdateUserRequest) ToUpdateUserInputData() *inputdata.UpdateUserInputData {
	return &inputdata.UpdateUserInputData{
		Email:      r.Email,
		Password:   r.Password,
		FullName:   r.FullName,
		RoleID:     r.RoleID,
		EnabledMFA: r.EnabledMFA,
	}
}

// Validate validates the UpdateUserRequest
func (r *UpdateUserRequest) Validate() error {
	// You can add custom validation logic here
	return nil
}

// UpdateUserResponse represents the HTTP response for user update
type UpdateUserResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	User    interface{} `json:"user"`
}

// ToResponseMap converts the UpdateUserResponse to a map for JSON response
func (r *UpdateUserResponse) ToResponseMap() map[string]interface{} {
	return map[string]interface{}{
		"success": r.Success,
		"message": r.Message,
		"user":    r.User,
	}
}

// NewUpdateUserResponse creates a new UpdateUserResponse
func NewUpdateUserResponse(user interface{}) *UpdateUserResponse {
	return &UpdateUserResponse{
		Success: true,
		Message: "User updated successfully",
		User:    user,
	}
}
