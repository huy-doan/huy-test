package repositories

import (
	"context"

	"github.com/vnlab/makeshop-payment/src/domain/models"
	"github.com/vnlab/makeshop-payment/src/domain/repositories/filter"
)

type PayoutRepository interface {
	List(ctx context.Context, filter *filter.PayoutFilter) ([]*models.Payout, int, int64, error)
}
