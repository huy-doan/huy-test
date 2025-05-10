package persistence

import (
	"context"

	repository "github.com/huydq/test/internal/domain/repository/payout_record"
	persistence "github.com/huydq/test/internal/infrastructure/persistence/util"
	"gorm.io/gorm"
)

type PayoutRecordRepositoryImpl struct {
	db            *gorm.DB
	filterBuilder *persistence.GormFilterBuilder
}

func NewPayoutRecordRepository(db *gorm.DB) repository.PayoutRecordRepository {
	return &PayoutRecordRepositoryImpl{
		db:            db,
		filterBuilder: persistence.NewGormFilterBuilder(),
	}
}

// CountByPayoutID counts the number of payout records associated with a specific payout ID
func (r *PayoutRecordRepositoryImpl) CountByPayoutID(ctx context.Context, payoutID int) (int, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Table("payout_record").
		Where("payout_id = ? AND deleted_at IS NULL", payoutID).
		Count(&count).Error

	if err != nil {
		return 0, err
	}

	return int(count), nil
}

// CountByPayoutIDs counts the number of payout records associated with multiple payout IDs
func (r *PayoutRecordRepositoryImpl) CountByPayoutIDs(ctx context.Context, payoutIDs []int) (map[int]int, error) {
	counts := make(map[int]int)

	if len(payoutIDs) == 0 {
		return counts, nil
	}

	var results []struct {
		PayoutID int
		Count    int64
	}

	err := r.db.WithContext(ctx).
		Table("payout_record").
		Select("payout_id, COUNT(*) AS count").
		Where("payout_id IN (?) AND deleted_at IS NULL", payoutIDs).
		Group("payout_id").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	for _, result := range results {
		counts[result.PayoutID] = int(result.Count)
	}

	return counts, nil
}

// SumAmountByPayoutID calculates the total amount of payout records associated with a specific payout ID
func (r *PayoutRecordRepositoryImpl) SumAmountByPayoutID(ctx context.Context, payoutID int) (float64, error) {
	var totalAmount float64

	err := r.db.WithContext(ctx).
		Table("payout_record").
		Where("payout_id = ? AND deleted_at IS NULL", payoutID).
		Select("SUM(amount)").
		Scan(&totalAmount).Error

	if err != nil {
		return 0, err
	}

	return totalAmount, nil
}

// SumAmountByPayoutIDs calculates the total amount of payout records associated with multiple payout IDs
func (r *PayoutRecordRepositoryImpl) SumAmountByPayoutIDs(ctx context.Context, payoutIDs []int) (map[int]float64, error) {
	totalAmounts := make(map[int]float64)

	if len(payoutIDs) == 0 {
		return totalAmounts, nil
	}

	var results []struct {
		PayoutID int
		Sum      float64
	}

	err := r.db.WithContext(ctx).
		Table("payout_record").
		Select("payout_id, SUM(amount) AS sum").
		Where("payout_id IN (?) AND deleted_at IS NULL", payoutIDs).
		Group("payout_id").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	for _, result := range results {
		totalAmounts[result.PayoutID] = result.Sum
	}

	return totalAmounts, nil
}
