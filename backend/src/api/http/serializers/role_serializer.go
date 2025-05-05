package serializers

import (
	"time"

	"github.com/vnlab/makeshop-payment/src/domain/models"
)

type RoleSerializer struct {
	Role *models.Role
}

func NewRoleSerializer(role *models.Role) *RoleSerializer {
	return &RoleSerializer{Role: role}
}

func (s *RoleSerializer) Serialize() any {
	if s.Role == nil {
		return nil
	}

	result := map[string]any{
		"id":         s.Role.ID,
		"code":       s.Role.Code,
		"name":       s.Role.Name,
		"created_at": s.Role.CreatedAt.Format(time.RFC3339),
		"updated_at": s.Role.UpdatedAt.Format(time.RFC3339),
	}

	if s.Role.Permissions != nil {
		result["permissions"] = SerializePermissionCollection(s.Role.Permissions)
	}

	return result
}

func SerializeRoleCollection(roles []*models.Role) []any {
	result := make([]any, len(roles))

	for i, role := range roles {
		result[i] = NewRoleSerializer(role).Serialize()
	}

	return result
}
