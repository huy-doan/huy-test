package serializers

import (
	"time"

	"github.com/huydq/test/src/domain/models"
)

// UserSerializer transforms User models into client-friendly representations
type UserSerializer struct {
	User *models.User
}

// NewUserSerializer creates a new UserSerializer
func NewUserSerializer(user *models.User) *UserSerializer {
	return &UserSerializer{User: user}
}

// Serialize transforms the User model to client representation
func (s *UserSerializer) Serialize() any {
	if s.User == nil {
		return nil
	}

	result := map[string]any{
		"id":          s.User.ID,
		"email":       s.User.Email,
		"full_name":   s.User.FullName,
		"created_at":  s.User.CreatedAt.Format(time.RFC3339),
		"updated_at":  s.User.UpdatedAt.Format(time.RFC3339),
		"enabled_mfa": s.User.EnabledMFA,
	}

	if s.User.Role != nil {
		result["role"] = NewRoleSerializer(s.User.Role).Serialize()
	}

	if s.User.MFAType != 0 {
		result["mfa_type"] = map[string]any{
			"id":        s.User.MFAType,
			"title":     models.GetMFATypeTitle(s.User.MFAType),
			"is_active": s.User.EnabledMFA,
		}
	}

	return result
}

// SerializeCollection serializes a collection of users
func SerializeUserCollection(users []*models.User) []any {
	result := make([]any, len(users))

	for i, user := range users {
		result[i] = NewUserSerializer(user).Serialize()
	}

	return result
}
