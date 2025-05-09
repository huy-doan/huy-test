package repositories

import (
	"context"
	"math"

	"github.com/huydq/test/src/domain/models"
	"github.com/huydq/test/src/domain/repositories"
	domainFilter "github.com/huydq/test/src/domain/repositories/filter"
	infraFilter "github.com/huydq/test/src/infrastructure/persistence/repositories/filter"
	"gorm.io/gorm"
)

// PayoutRepositoryImpl implements the PayoutRepository interface
type PayoutRepositoryImpl struct {
	db            *gorm.DB
	filterBuilder *infraFilter.GormFilterBuilder
}

// NewPayoutRepository creates a new PayoutRepository
func NewPayoutRepository(db *gorm.DB) repositories.PayoutRepository {
	return &PayoutRepositoryImpl{
		db:            db,
		filterBuilder: infraFilter.NewGormFilterBuilder(),
	}
}

// List retrieves payouts with pagination and filtering
func (r *PayoutRepositoryImpl) List(ctx context.Context, filter *domainFilter.PayoutFilter) ([]*models.Payout, int, int64, error) {
	var payouts []*models.Payout
	var count int64

	if filter != nil {
		filter.ApplyFilters()
	} else {
		filter = domainFilter.NewPayoutFilter()
	}

	query := r.db.WithContext(ctx).Model(&models.Payout{})
	query = r.filterBuilder.ApplyBaseFilter(query, &filter.BaseFilter)

	if err := query.Count(&count).Error; err != nil {
		return nil, 0, 0, err
	}

	query = r.filterBuilder.ApplyPagination(query, filter.Pagination)

	query = query.Preload("User")

	if err := query.Find(&payouts).Error; err != nil {
		return nil, 0, 0, err
	}

	totalPages := int(math.Ceil(float64(count) / float64(filter.Pagination.PageSize)))

	return payouts, totalPages, int64(count), nil
}
