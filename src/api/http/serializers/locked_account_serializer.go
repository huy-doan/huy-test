package serializers

import (
	"time"

	"github.com/vnlab/makeshop-payment/src/domain/models"
)

// LockedAccountSerializer transforms LockedAccount models into client-friendly representations
type LockedAccountSerializer struct {
	LockedAccount *models.LockedAccount
}

// NewLockedAccountSerializer creates a new LockedAccountSerializer
func NewLockedAccountSerializer(lockedAccount *models.LockedAccount) *LockedAccountSerializer {
	return &LockedAccountSerializer{LockedAccount: lockedAccount}
}

// Serialize transforms the LockedAccount model to client representation
func (s *LockedAccountSerializer) Serialize() interface{} {
	if s.LockedAccount == nil {
		return nil
	}

	result := map[string]interface{}{
		"id":         s.LockedAccount.ID,
		"user_id":    s.LockedAccount.UserID,
		"count":      s.LockedAccount.Count,
		"created_at": s.LockedAccount.CreatedAt.Format(time.RFC3339),
		"updated_at": s.LockedAccount.UpdatedAt.Format(time.RFC3339),
	}

	if s.LockedAccount.LockedAt != nil {
		result["locked_at"] = s.LockedAccount.LockedAt.Format(time.RFC3339)
	}
	if s.LockedAccount.ExpiredAt != nil {
		result["expired_at"] = s.LockedAccount.ExpiredAt.Format(time.RFC3339)
	}
	if s.LockedAccount.Email != "" {
		result["email"] = s.LockedAccount.Email
	}

	return result
}
