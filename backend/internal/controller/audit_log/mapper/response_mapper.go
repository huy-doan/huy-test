package mapper

import (
	"github.com/huydq/test/internal/datastructure/outputdata"
	model "github.com/huydq/test/internal/domain/model/audit_log"
)

type AuditLogUserResponse struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
}

type AuditLogUsersResponse struct {
	Users []*AuditLogUserResponse `json:"users"`
}

type AuditLogResponse struct {
	ID            int                   `json:"id"`
	UserID        *int                  `json:"user_id"`
	User          *AuditLogUserResponse `json:"user"`
	AuditLogType  string                `json:"audit_log_type"`
	Description   *string               `json:"description"`
	TransactionID *int                  `json:"transaction_id"`
	PayoutID      *int                  `json:"payout_id"`
	PayinID       *int                  `json:"payin_id"`
	UserAgent     string                `json:"user_agent"`
	IPAddress     string                `json:"ip_address"`
	CreatedAt     string                `json:"created_at"`
	UpdatedAt     string                `json:"updated_at"`
}

type AuditLogListResponse struct {
	AuditLogs  []*AuditLogResponse `json:"audit_logs"`
	TotalPages int                 `json:"total_pages"`
	Total      int64               `json:"total"`
	Page       int                 `json:"page"`
	PageSize   int                 `json:"page_size"`
}

func ToAuditLogResponse(auditLog *model.AuditLog) *AuditLogResponse {
	if auditLog == nil {
		return nil
	}

	var user *AuditLogUserResponse
	if auditLog.User != nil {
		user = &AuditLogUserResponse{
			ID:       auditLog.User.ID,
			Email:    auditLog.User.Email,
			FullName: auditLog.User.FullName,
		}
	}

	return &AuditLogResponse{
		ID:            auditLog.ID,
		UserID:        auditLog.UserID,
		User:          user,
		AuditLogType:  string(auditLog.AuditLogType),
		Description:   auditLog.Description,
		TransactionID: auditLog.TransactionID,
		PayoutID:      auditLog.PayoutID,
		PayinID:       auditLog.PayinID,
		UserAgent:     auditLog.UserAgent.String(),
		IPAddress:     auditLog.IPAddress.String(),
		CreatedAt:     auditLog.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     auditLog.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func ToAuditLogListResponse(auditLogs []*model.AuditLog, page, pageSize, totalPages int, total int64) *AuditLogListResponse {
	auditLogResponses := make([]*AuditLogResponse, len(auditLogs))
	for i, auditLog := range auditLogs {
		auditLogResponses[i] = ToAuditLogResponse(auditLog)
	}
	return &AuditLogListResponse{
		AuditLogs:  auditLogResponses,
		TotalPages: totalPages,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
	}
}

func ToAuditLogUsersData(users []*outputdata.AuditLogUserOutput) *AuditLogUsersResponse {
	if users == nil {
		return &AuditLogUsersResponse{
			Users: []*AuditLogUserResponse{},
		}
	}

	usersResponse := make([]*AuditLogUserResponse, len(users))
	for i, user := range users {
		usersResponse[i] = ToAuditLogUser(user)
	}

	return &AuditLogUsersResponse{
		Users: usersResponse,
	}
}

func ToAuditLogUser(user *outputdata.AuditLogUserOutput) *AuditLogUserResponse {
	if user == nil {
		return nil
	}

	return &AuditLogUserResponse{
		ID:       user.ID,
		Email:    user.Email,
		FullName: user.FullName,
	}
}
