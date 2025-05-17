package repository

import (
	"context"

	model "github.com/huydq/test/internal/domain/model/paypay"
)

type PaypayPayinTransactionRepository interface {
	BulkInsert(ctx context.Context, transactions []*model.PaypayPayinTransaction) error
}
