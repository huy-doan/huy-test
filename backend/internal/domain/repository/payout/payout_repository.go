package repository

import (
	"context"

	model "github.com/huydq/test/internal/domain/model/payout"
)

type PayoutRepository interface {
	// List lists all payouts with filtering and pagination
	List(ctx context.Context, filter *model.PayoutFilter) ([]*model.Payout, int, int64, error)
}
