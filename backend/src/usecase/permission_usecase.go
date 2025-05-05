package usecase

import (
	"context"

	"github.com/vnlab/makeshop-payment/src/domain/models"
	"github.com/vnlab/makeshop-payment/src/domain/repositories"
)

type PermissionUsecase struct {
	permissionRepo repositories.PermissionRepository
}

func NewPermissionUseCase(permissionRepo repositories.PermissionRepository) *PermissionUsecase {
	return &PermissionUsecase{
		permissionRepo: permissionRepo,
	}
}

func (u *PermissionUsecase) ListPermission(ctx context.Context) ([]*models.Permission, error) {
	permissions, err := u.permissionRepo.List(ctx)
	if err != nil {
		return nil, err
	}
	return permissions, nil
}
