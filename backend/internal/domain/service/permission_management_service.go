package service

import (
	"context"

	modelPermission "github.com/huydq/test/internal/domain/model/permission"
	repositoryPermission "github.com/huydq/test/internal/domain/repository/permission"
)

type PermissionService interface {
	GetPermissionsByIDs(ctx context.Context, ids []int) ([]*modelPermission.Permission, error)
	ListPermissions(ctx context.Context) ([]*modelPermission.Permission, error)
}

type permissionServiceImpl struct {
	permissionRepository repositoryPermission.PermissionRepository
}

func NewPermissionService(
	permissionRepository repositoryPermission.PermissionRepository,
) PermissionService {
	return &permissionServiceImpl{
		permissionRepository: permissionRepository,
	}
}

func (s *permissionServiceImpl) GetPermissionsByIDs(
	ctx context.Context,
	ids []int,
) ([]*modelPermission.Permission, error) {
	return s.permissionRepository.FindByIDs(ctx, ids)
}

func (s *permissionServiceImpl) ListPermissions(
	ctx context.Context,
) ([]*modelPermission.Permission, error) {
	return s.permissionRepository.List(ctx)
}
