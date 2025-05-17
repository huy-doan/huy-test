package usecase

import (
	"context"

	"github.com/huydq/test/internal/datastructure/outputdata"
	model "github.com/huydq/test/internal/domain/model/audit_log"
	"github.com/huydq/test/internal/domain/service"
)

type AuditLogUsecase interface {
	List(ctx context.Context, input *model.AuditLogFilter) ([]*model.AuditLog, int, int64, error)
	GetAuditLogUsers(ctx context.Context) (*outputdata.AuditLogUsersOutput, error)
}

type auditLogUsecaseImpl struct {
	auditLogService service.AuditLogService
}

func NewAuditLogUsecase(auditLogService service.AuditLogService) AuditLogUsecase {
	return &auditLogUsecaseImpl{
		auditLogService: auditLogService,
	}
}

func (uc *auditLogUsecaseImpl) List(ctx context.Context, input *model.AuditLogFilter) ([]*model.AuditLog, int, int64, error) {
	if input == nil {
		input = model.NewAuditLogFilter()
	} else {
		input.ApplyFilters()
	}

	return uc.auditLogService.GetAuditLogs(ctx, input)
}

func (uc *auditLogUsecaseImpl) GetAuditLogUsers(ctx context.Context) (*outputdata.AuditLogUsersOutput, error) {
	users, err := uc.auditLogService.GetUsersWithAuditLogs(ctx)
	if err != nil {
		return nil, err
	}

	userOutputs := make([]*outputdata.AuditLogUserOutput, len(users))
	for i, user := range users {
		userOutputs[i] = &outputdata.AuditLogUserOutput{
			ID:       user.ID,
			Email:    user.Email,
			FullName: user.FullName,
		}
	}

	return &outputdata.AuditLogUsersOutput{
		Users: userOutputs,
	}, nil
}
