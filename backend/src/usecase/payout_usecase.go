package usecase

import (
	"context"
	"time"

	"github.com/huydq/test/src/domain/models"
	"github.com/huydq/test/src/domain/repositories"
	"github.com/huydq/test/src/domain/repositories/filter"
)

type PayoutUsecase struct {
	payoutRepo       repositories.PayoutRepository
	payoutRecordRepo repositories.PayoutRecordRepository
}

func NewPayoutUsecase(
	payoutRepo repositories.PayoutRepository,
	payoutRecordRepo repositories.PayoutRecordRepository,
) *PayoutUsecase {
	return &PayoutUsecase{
		payoutRepo:       payoutRepo,
		payoutRecordRepo: payoutRecordRepo,
	}
}

// PayoutFilterParams defines filter parameters for payouts at the usecase layer
// This acts as a DTO for receiving filter parameters from handlers
type PayoutFilterParams struct {
	CreatedAt    *time.Time
	SendingDate  *time.Time
	SentDate     *time.Time
	PayoutStatus *int
}

// GetAll retrieves all payouts based on filter criteria
func (uc *PayoutUsecase) GetAll(ctx context.Context, filter *filter.PayoutFilter) ([]*models.Payout, int, int64, error) {
	payouts, totalPages, total, err := uc.payoutRepo.List(ctx, filter)
	if err != nil {
		return nil, 0, 0, err
	}

	// Get payout record counts
	if len(payouts) > 0 {
		var payoutIDs []int
		for _, payout := range payouts {
			payoutIDs = append(payoutIDs, payout.ID)
		}

		counts, err := uc.payoutRecordRepo.CountByPayoutIDs(ctx, payoutIDs)
		if err != nil {
			return nil, 0, 0, err
		}

		for _, payout := range payouts {
			if count, ok := counts[payout.ID]; ok {
				payout.PayoutRecordCount = count
			}
		}
	}

	// Get payout record sum amounts
	if len(payouts) > 0 {
		var payoutIDs []int
		for _, payout := range payouts {
			payoutIDs = append(payoutIDs, payout.ID)
		}
		counts, err := uc.payoutRecordRepo.SumAmountByPayoutIDs(ctx, payoutIDs)
		if err != nil {
			return nil, 0, 0, err
		}

		for _, payout := range payouts {
			if sum, ok := counts[payout.ID]; ok {
				payout.PayoutRecordSumAmount = sum
			}
		}
	}

	return payouts, totalPages, int64(total), nil
}
