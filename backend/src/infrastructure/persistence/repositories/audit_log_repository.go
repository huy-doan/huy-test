package repositories

import (
	"context"
	"math"

	"github.com/vnlab/makeshop-payment/src/domain/models"
	"github.com/vnlab/makeshop-payment/src/domain/repositories"
	domainFilter "github.com/vnlab/makeshop-payment/src/domain/repositories/filter"
	infraFilter "github.com/vnlab/makeshop-payment/src/infrastructure/persistence/repositories/filter"
	"gorm.io/gorm"
)

type AuditLogRepositoryImpl struct {
	db            *gorm.DB
	filterBuilder *infraFilter.GormFilterBuilder
}

func NewAuditLogRepository(db *gorm.DB) repositories.AuditLogRepository {
	return &AuditLogRepositoryImpl{
		db:            db,
		filterBuilder: infraFilter.NewGormFilterBuilder(),
	}
}

func (r *AuditLogRepositoryImpl) Create(ctx context.Context, auditLog *models.AuditLog) error {
	return r.db.Create(auditLog).Error
}

func (r *AuditLogRepositoryImpl) List(ctx context.Context, filter *domainFilter.AuditLogFilter) ([]*models.AuditLog, int, int64, error) {
	var auditLogs []*models.AuditLog
	var count int64

	if filter != nil {
		filter.ApplyFilters()
	} else {
		filter = domainFilter.NewAuditLogFilter()
	}

	query := r.db.WithContext(ctx).Model(&models.AuditLog{})

	query = r.filterBuilder.ApplyBaseFilter(query, &filter.BaseFilter)

	if err := query.Count(&count).Error; err != nil {
		return nil, 0, 0, err
	}

	err := query.Find(&auditLogs).Count(&count).Error
	if err != nil {
		return nil, 0, 0, err
	}

	query = r.filterBuilder.ApplyPagination(query, filter.Pagination)

	query = query.Preload("User")

	if err := query.Find(&auditLogs).Error; err != nil {
		return nil, 0, 0, err
	}

	totalPages := int(math.Ceil(float64(count) / float64(filter.Pagination.PageSize)))

	return auditLogs, totalPages, int64(count), nil
}
