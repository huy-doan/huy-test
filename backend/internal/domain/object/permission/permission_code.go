package object

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
