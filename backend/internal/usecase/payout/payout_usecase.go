package usecase

import (
	"context"
	"errors"

	model "github.com/huydq/test/internal/domain/model/payout"
	"github.com/huydq/test/internal/domain/service"
)

var (
	ErrInvalidInput    = errors.New("入力が無効です")
	ErrPayoutNotFound  = errors.New("支払いが見つかりません")
	ErrOperationFailed = errors.New("操作に失敗しました")
)

type PayoutUsecase interface {
	ListPayouts(ctx context.Context, filter *model.PayoutFilter) ([]*model.Payout, int, int64, error)
}

type payoutUsecaseImpl struct {
	payoutService service.PayoutManagementService
}

func NewPayoutUsecase(payoutService service.PayoutManagementService) PayoutUsecase {
	return &payoutUsecaseImpl{
		payoutService: payoutService,
	}
}

// ListPayouts retrieves a list of payouts based on the provided filter criteria
func (u *payoutUsecaseImpl) ListPayouts(ctx context.Context, filter *model.PayoutFilter) ([]*model.Payout, int, int64, error) {
	if filter == nil {
		filter = model.NewPayoutFilter()
	}

	filter.ApplyFilters()
	return u.payoutService.ListPayouts(ctx, filter)
}
