package outputdata

// AuditLogUserOutput represents the output data for a user with audit logs
type AuditLogUserOutput struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
}

// AuditLogUsersOutput represents the output data for the list of users with audit logs
type AuditLogUsersOutput struct {
	Users []*AuditLogUserOutput `json:"users"`
}
