package repository

import (
	"context"

	model "github.com/huydq/test/internal/domain/model/paypay"
)

type PaypayPayinDetailRepository interface {
	BulkInsert(ctx context.Context, details []*model.PaypayPayinDetail) error
}
