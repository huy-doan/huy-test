package repositories

import (
	"context"

	"github.com/vnlab/makeshop-payment/src/domain/models"
)

type PayinDetailRepository interface {
	BulkInsert(ctx context.Context, details []*models.PayPayPayinDetail) error
}
