package repositories

import (
	"context"

	payinFileGroupModel "github.com/huydq/test/internal/infrastructure/persistence/paypay_batch/paypay_payin_filegroup/dto"
)

type PayinFileGroupRepository interface {
	//Create methods for PayinFileGroupRepository
	Create(ctx context.Context, group *payinFileGroupModel.PayinFileGroup) error
}
