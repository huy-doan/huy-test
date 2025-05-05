package serializers

import (
	"time"

	"github.com/vnlab/makeshop-payment/src/domain/models"
)

// PermissionResponse represents the serialized permission for API responses and documentation
type PermissionResponse struct {
	ID        int             `json:"id"`
	Name      string          `json:"name"`
	Code      string          `json:"code"`
	Screen    *ScreenResponse `json:"screen,omitempty"`
	CreatedAt string          `json:"created_at"`
	UpdatedAt string          `json:"updated_at"`
}

// ScreenResponse represents the serialized screen for API responses and documentation
type ScreenResponse struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	ScreenCode string `json:"screen_code"`
	ScreenPath string `json:"screen_path"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

type PermissionSerializer struct {
	Permission *models.Permission
}

func NewPermissionSerializer(permission *models.Permission) *PermissionSerializer {
	return &PermissionSerializer{Permission: permission}
}

func (s *PermissionSerializer) Serialize() any {
	if s.Permission == nil {
		return nil
	}

	result := map[string]any{
		"id":         s.Permission.ID,
		"name":       s.Permission.Name,
		"code":       s.Permission.Code,
		"created_at": s.Permission.CreatedAt.Format(time.RFC3339),
		"updated_at": s.Permission.UpdatedAt.Format(time.RFC3339),
	}

	if s.Permission.Screen != nil {
		result["screen"] = map[string]any{
			"id":          s.Permission.Screen.ID,
			"name":        s.Permission.Screen.Name,
			"screen_code": s.Permission.Screen.ScreenCode,
			"screen_path": s.Permission.Screen.ScreenPath,
			"created_at":  s.Permission.Screen.CreatedAt.Format(time.RFC3339),
			"updated_at":  s.Permission.Screen.UpdatedAt.Format(time.RFC3339),
		}
	}

	return result
}

func SerializePermissionCollection(permissions []*models.Permission) []any {
	result := make([]any, len(permissions))

	for i, permission := range permissions {
		result[i] = NewPermissionSerializer(permission).Serialize()
	}

	return result
}
