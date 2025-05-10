package outputdata

import (
	"time"
)

// PermissionOutput represents the output data for a permission
type PermissionOutput struct {
	ID        int           `json:"id"`
	Name      string        `json:"name"`
	Code      string        `json:"code"`
	ScreenID  int           `json:"screen_id"`
	Screen    *ScreenOutput `json:"screen,omitempty"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

// ScreenOutput represents the output data for a screen
type ScreenOutput struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	ScreenCode string    `json:"screen_code"`
	ScreenPath string    `json:"screen_path"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// PermissionListOutput represents the output for listing permissions
type PermissionListOutput struct {
	Permissions []*PermissionOutput `json:"permissions"`
	Total       int64               `json:"total"`
}
