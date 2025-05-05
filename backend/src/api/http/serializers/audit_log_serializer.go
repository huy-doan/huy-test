package serializers

import (
	"time"

	"github.com/vnlab/makeshop-payment/src/domain/models"
)

// AuditLogSerializer transforms AuditLog models into client-friendly representations
type AuditLogSerializer struct {
	AuditLog *models.AuditLog
}

// NewAuditLogSerializer creates a new AuditLogSerializer
func NewAuditLogSerializer(auditLog *models.AuditLog) *AuditLogSerializer {
	return &AuditLogSerializer{AuditLog: auditLog}
}

// Serialize transforms the AuditLog model to client representation
func (s *AuditLogSerializer) Serialize() any {
	if s.AuditLog == nil {
		return nil
	}

	result := map[string]any{
		"id":             s.AuditLog.ID,
		"audit_log_type": s.AuditLog.AuditLogType,
		"created_at":     s.AuditLog.CreatedAt.Format(time.RFC3339),
		"updated_at":     s.AuditLog.UpdatedAt.Format(time.RFC3339),
	}

	// Add optional fields only if they exist
	if s.AuditLog.UserID != nil {
		result["user_id"] = *s.AuditLog.UserID
	}

	if s.AuditLog.Description != nil {
		result["description"] = *s.AuditLog.Description
	}

	if s.AuditLog.TransactionID != nil {
		result["transaction_id"] = *s.AuditLog.TransactionID
	}

	if s.AuditLog.PayoutID != nil {
		result["payout_id"] = *s.AuditLog.PayoutID
	}

	if s.AuditLog.PayinID != nil {
		result["payin_id"] = *s.AuditLog.PayinID
	}

	if s.AuditLog.UserAgent != nil {
		result["user_agent"] = *s.AuditLog.UserAgent
	}

	if s.AuditLog.IPAddress != nil {
		result["ip_address"] = *s.AuditLog.IPAddress
	}

	// Include user data if available
	if s.AuditLog.User != nil {
		result["user"] = map[string]any{
			"id":        s.AuditLog.User.ID,
			"email":     s.AuditLog.User.Email,
			"full_name": s.AuditLog.User.FullName,
		}
	}

	return result
}

// SerializeAuditLogCollection serializes a collection of audit logs
func SerializeAuditLogCollection(auditLogs []*models.AuditLog) []any {
	result := make([]any, len(auditLogs))

	for i, auditLog := range auditLogs {
		result[i] = NewAuditLogSerializer(auditLog).Serialize()
	}

	return result
}
