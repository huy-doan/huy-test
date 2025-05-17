package service

import (
	"context"

	auditLogModel "github.com/huydq/test/internal/domain/model/audit_log"
	userModel "github.com/huydq/test/internal/domain/model/user"
	auditLogRepository "github.com/huydq/test/internal/domain/repository/audit_log"
	userRepository "github.com/huydq/test/internal/domain/repository/user"
)

type AuditLogService interface {
	CreateAuditLog(ctx context.Context, auditLog *auditLogModel.AuditLog) error
	GetAuditLogs(ctx context.Context, filter *auditLogModel.AuditLogFilter) ([]*auditLogModel.AuditLog, int, int64, error)
	GetUsersWithAuditLogs(ctx context.Context) ([]*userModel.User, error)
}

type auditLogServiceImpl struct {
	auditLogRepository auditLogRepository.AuditLogRepository
	userRepository     userRepository.UserRepository
}

func NewAuditLogService(auditLogRepository auditLogRepository.AuditLogRepository, userRepository userRepository.UserRepository) AuditLogService {
	return &auditLogServiceImpl{
		auditLogRepository: auditLogRepository,
		userRepository:     userRepository,
	}
}

func (s *auditLogServiceImpl) CreateAuditLog(ctx context.Context, auditLog *auditLogModel.AuditLog) error {
	return s.auditLogRepository.Create(ctx, auditLog)
}

func (s *auditLogServiceImpl) GetAuditLogs(ctx context.Context, filter *auditLogModel.AuditLogFilter) ([]*auditLogModel.AuditLog, int, int64, error) {
	return s.auditLogRepository.List(ctx, filter)
}

func (s *auditLogServiceImpl) GetUsersWithAuditLogs(ctx context.Context) ([]*userModel.User, error) {
	return s.userRepository.GetUsersWithAuditLogs(ctx)
}
