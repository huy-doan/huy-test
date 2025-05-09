package repositories

import (
	"context"

	"github.com/huydq/test/src/domain/models"
	"github.com/huydq/test/src/domain/repositories/filter"
)

type PayoutRepository interface {
	List(ctx context.Context, filter *filter.PayoutFilter) ([]*models.Payout, int, int64, error)
}
