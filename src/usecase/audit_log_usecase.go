package usecase

import (
	"context"

	"github.com/huydq/test/src/domain/models"
	"github.com/huydq/test/src/domain/repositories"
)

// AuditLogUsecase handles audit logging business logic
type AuditLogUsecase struct {
	auditLogRepo     repositories.AuditLogRepository
	auditLogTypeRepo repositories.AuditLogTypeRepository
}

// NewAuditLogUsecase creates a new AuditLogUsecase
func NewAuditLogUsecase(
	auditLogRepo repositories.AuditLogRepository,
	auditLogTypeRepo repositories.AuditLogTypeRepository,
) *AuditLogUsecase {
	return &AuditLogUsecase{
		auditLogRepo:     auditLogRepo,
		auditLogTypeRepo: auditLogTypeRepo,
	}
}

func (uc *AuditLogUsecase) logEventByType(ctx context.Context, userID *int, ipAddress *string, userAgent *string, auditTypeID int) error {
	auditLogType, err := uc.auditLogTypeRepo.GetByID(ctx, auditTypeID)
	description := auditLogType.Name

	if err != nil {
		description = ""
	}

	auditLog := &models.AuditLog{
		UserID:      userID,
		Description: &description,
		AuditTypeID: auditTypeID,
		IPAddress:   ipAddress,
		UserAgent:   userAgent,
	}

	return uc.auditLogRepo.Create(ctx, auditLog)
}

func (uc *AuditLogUsecase) LogLoginEvent(ctx context.Context, userID *int, ipAddress *string, userAgent *string) error {
	return uc.logEventByType(ctx, userID, ipAddress, userAgent, models.AuditTypeUserLogin)
}

func (uc *AuditLogUsecase) LogLogoutEvent(ctx context.Context, userID *int, ipAddress *string, userAgent *string) error {
	return uc.logEventByType(ctx, userID, ipAddress, userAgent, models.AuditTypeUserLogout)
}
