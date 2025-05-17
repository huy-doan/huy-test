package repository

import (
	"context"

	model "github.com/huydq/test/internal/domain/model/paypay"
)

type PaypayPayinSummaryRepository interface {
	BulkInsert(ctx context.Context, summaries []*model.PaypayPayinSummary) error
}
