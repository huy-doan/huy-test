package persistence

import (
	"context"
	"math"

	model "github.com/huydq/test/internal/domain/model/audit_log"
	repository "github.com/huydq/test/internal/domain/repository/audit_log"
	"github.com/huydq/test/internal/infrastructure/persistence/audit_log/convert"
	"github.com/huydq/test/internal/infrastructure/persistence/audit_log/dto"
	persistence "github.com/huydq/test/internal/infrastructure/persistence/util"
	"gorm.io/gorm"
)

type AuditLogRepositoryImpl struct {
	db            *gorm.DB
	filterBuilder *persistence.GormFilterBuilder
}

func NewAuditLogRepository(db *gorm.DB) repository.AuditLogRepository {
	return &AuditLogRepositoryImpl{
		db:            db,
		filterBuilder: persistence.NewGormFilterBuilder(),
	}
}

func (r *AuditLogRepositoryImpl) Create(ctx context.Context, auditLog *model.AuditLog) error {
	auditLogDto := convert.ToAuditLogDTO(auditLog)
	return r.db.Create(auditLogDto).Error
}

func (r *AuditLogRepositoryImpl) List(ctx context.Context, filter *model.AuditLogFilter) ([]*model.AuditLog, int, int64, error) {
	var auditLogDtos []*dto.AuditLogDTO
	var count int64

	if filter != nil {
		filter.ApplyFilters()
	} else {
		filter = model.NewAuditLogFilter()
	}

	query := r.db.WithContext(ctx).Model(&dto.AuditLogDTO{})

	query = r.filterBuilder.ApplyBaseFilter(query, &filter.BaseFilter)

	if err := query.Count(&count).Error; err != nil {
		return nil, 0, 0, err
	}

	err := query.Find(&auditLogDtos).Count(&count).Error
	if err != nil {
		return nil, 0, 0, err
	}

	query = r.filterBuilder.ApplyPagination(query, filter.Pagination)

	// query = query.Preload("User")

	if err := query.Find(&auditLogDtos).Error; err != nil {
		return nil, 0, 0, err
	}

	totalPages := int(math.Ceil(float64(count) / float64(filter.Pagination.PageSize)))

	auditLogs := convert.ToAuditLogModels(auditLogDtos)

	return auditLogs, totalPages, int64(count), nil
}
