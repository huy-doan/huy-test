package models

// Role represents a user role in the system
type Permission struct {
	ID int `json:"id"`
	BaseColumnTimestamp

	Name     string  `json:"name"`
	Code     string  `json:"code"`
	ScreenID int     `json:"screen_id"`
	Screen   *Screen `json:"screen" gorm:"foreignKey:ScreenID"`
}

// TableName specifies the database table name
func (Permission) TableName() string {
	return "permission"
}

// PermissionCode defines constants for standard permission codes
type PermissionCode string

const (
	// Admin-related permissions
	PermissionCodeUserManage     PermissionCode = "USER_MANAGE"
	PermissionCodeUserRoleChange PermissionCode = "USER_ROLE_CHANGE"
	PermissionCodeSystemLogView  PermissionCode = "SYSTEM_LOG_VIEW"

	// User-related permissions
	PermissionCodeEditOwnProfile PermissionCode = "EDIT_OWN_PROFILE"
	PermissionCodeViewOwnLog     PermissionCode = "VIEW_OWN_LOG"
	PermissionCodeViewAdminPanel PermissionCode = "VIEW_ADMIN_PANEL"

	// Transfer-related permissions
	PermissionCodeTransferApproveBusiness   PermissionCode = "TRANSFER_APPROVE_BUSINESS"
	PermissionCodeTransferApproveAccountant PermissionCode = "TRANSFER_APPROVE_ACCOUNTANT"
	PermissionCodeManualTransfer            PermissionCode = "MANUAL_TRANSFER"
)
