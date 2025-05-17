package persistence

import (
	"context"
	"math"

	model "github.com/huydq/test/internal/domain/model/payout"
	repository "github.com/huydq/test/internal/domain/repository/payout"
	"github.com/huydq/test/internal/infrastructure/persistence/payout/convert"
	"github.com/huydq/test/internal/infrastructure/persistence/payout/dto"
	persistence "github.com/huydq/test/internal/infrastructure/persistence/util"
	"github.com/huydq/test/internal/pkg/database"
	"gorm.io/gorm"
)

type PayoutRepositoryImpl struct {
	db            *gorm.DB
	filterBuilder *persistence.GormFilterBuilder
}

func NewPayoutRepository(db *gorm.DB) repository.PayoutRepository {
	return &PayoutRepositoryImpl{
		db:            db,
		filterBuilder: persistence.NewGormFilterBuilder(),
	}
}

func (r *PayoutRepositoryImpl) List(ctx context.Context, filter *model.PayoutFilter) ([]*model.Payout, int, int64, error) {
	db, err := database.GetTxOrDB(ctx)
	if err != nil {
		return nil, 0, 0, err
	}

	var payoutDTOs []*dto.Payout
	var count int64

	if filter != nil {
		filter.ApplyFilters()
	} else {
		filter = model.NewPayoutFilter()
	}

	query := db.WithContext(ctx).Model(&dto.Payout{})

	query = r.filterBuilder.ApplyBaseFilter(query, &filter.BaseFilter)

	if err := query.Count(&count).Error; err != nil {
		return nil, 0, 0, err
	}

	query = r.filterBuilder.ApplyPagination(query, filter.Pagination)

	query = query.Preload("User")

	if err := query.Find(&payoutDTOs).Error; err != nil {
		return nil, 0, 0, err
	}

	totalPages := int(math.Ceil(float64(count) / float64(filter.Pagination.PageSize)))

	payouts := convert.ToPayoutModels(payoutDTOs)

	return payouts, totalPages, int64(count), nil
}
