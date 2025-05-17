package repository

import (
	"context"
)

type PayoutRecordRepository interface {
	CountByPayoutID(ctx context.Context, payoutID int) (int, error)
	CountByPayoutIDs(ctx context.Context, payoutIDs []int) (map[int]int, error)

	SumAmountByPayoutID(ctx context.Context, payoutID int) (float64, error)
	SumAmountByPayoutIDs(ctx context.Context, payoutIDs []int) (map[int]float64, error)
}
