package usecase

import (
	"context"

	repository "github.com/huydq/test/batch/domain/repository/payin"
	model "github.com/huydq/test/internal/domain/model/payin"
)

type PayinFileGroupUsecase struct {
	repo repository.PayinFileGroupRepository
}

func NewPayinFileGroupUsecase(repo repository.PayinFileGroupRepository) *PayinFileGroupUsecase {
	return &PayinFileGroupUsecase{repo: repo}
}

func (uc *PayinFileGroupUsecase) CreateGroup(ctx context.Context, group *model.PayinFileGroup) error {
	return uc.repo.Create(ctx, group)
}
