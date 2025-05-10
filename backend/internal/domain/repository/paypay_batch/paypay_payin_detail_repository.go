package repositories

import (
	"context"

	models "github.com/huydq/test/internal/infrastructure/persistence/paypay_batch/paypay_payin_detail/dto"
)

type PayinDetailRepository interface {
	BulkInsert(ctx context.Context, details []*models.PayPayPayinDetail) error
}
