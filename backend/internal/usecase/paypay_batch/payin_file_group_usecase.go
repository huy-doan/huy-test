package paypay_batch

import (
	"context"

	repositories "github.com/huydq/test/internal/domain/repository/paypay_batch"
	payinFileGroupModel "github.com/huydq/test/internal/infrastructure/persistence/paypay_batch/paypay_payin_filegroup/dto"
)

type PayinFileGroupUsecase struct {
	repo repositories.PayinFileGroupRepository
}

func NewPayinFileGroupUsecase(repo repositories.PayinFileGroupRepository) *PayinFileGroupUsecase {
	return &PayinFileGroupUsecase{repo: repo}
}

func (uc *PayinFileGroupUsecase) CreateGroup(ctx context.Context, group *payinFileGroupModel.PayinFileGroup) error {
	return uc.repo.Create(ctx, group)
}
