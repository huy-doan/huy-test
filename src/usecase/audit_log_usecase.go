package usecase

import (
	"context"
	"github.com/vnlab/makeshop-payment/src/domain/models"
	"github.com/vnlab/makeshop-payment/src/domain/repositories"
)

// AuditLogUsecase handles audit logging business logic
type AuditLogUsecase struct {
	auditLogRepo repositories.AuditLogRepository
	masterAuditLogTypeRepo repositories.MasterAuditLogTypeRepository
}

// NewAuditLogUsecase creates a new AuditLogUsecase
func NewAuditLogUsecase(
	auditLogRepo repositories.AuditLogRepository,
	masterAuditLogTypeRepo repositories.MasterAuditLogTypeRepository,
) *AuditLogUsecase {
	return &AuditLogUsecase{
		auditLogRepo:             auditLogRepo,
		masterAuditLogTypeRepo: masterAuditLogTypeRepo,
	}
}

func (uc *AuditLogUsecase) logEventByType(ctx context.Context, userID *int, ipAddress *string, userAgent *string, auditTypeID int) error {
	masterAuditLogType, err := uc.masterAuditLogTypeRepo.GetByID(ctx, auditTypeID)
	description := masterAuditLogType.Name

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
