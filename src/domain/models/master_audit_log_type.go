package models

// MasterAuditLogType represents the type of audit log
type MasterAuditLogType struct {
	ID   int    `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
	BaseColumnTimestamp
}

// TableName specifies the table name for MasterAuditLogType
func (MasterAuditLogType) TableName() string {
	return "master_audit_log_types"
}

// Audit log type constants
const (
	AuditTypeUserLogin                  uint8 = 1
	AuditTypeUserLogout                 uint8 = 2
	AuditTypePasswordChange             uint8 = 3
	AuditTypeUserCreate                 uint8 = 4
	AuditTypeUserUpdate                 uint8 = 5
	AuditTypeUserDelete                 uint8 = 6
	AuditTypeRoleAssign                 uint8 = 7
	AuditTypeRoleUpdate                 uint8 = 8
	AuditTypeTransactionBusinessApprove uint8 = 9
	AuditTypeTransactionBusinessReject  uint8 = 10
	AuditTypeTransactionAccountApprove  uint8 = 11
	AuditTypeTransactionAccountReject   uint8 = 12
	AuditTypeTransferBusinessApprove    uint8 = 13
	AuditTypeTransferBusinessReject     uint8 = 14
	AuditTypeTransferAccountApprove     uint8 = 15
	AuditTypeTransferAccountReject      uint8 = 16
	AuditTypeMFAEnabled                 uint8 = 17
	AuditTypeMFADisabled                uint8 = 18
	AuditTypeAccountLocked              uint8 = 19
	AuditTypeAccountUnlocked            uint8 = 20
	AuditTypeAPIAccess                  uint8 = 21
)
