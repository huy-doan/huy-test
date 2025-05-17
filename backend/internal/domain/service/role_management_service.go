package service

import (
	"context"

	modelPermission "github.com/huydq/test/internal/domain/model/permission"
	modelRole "github.com/huydq/test/internal/domain/model/role"
	objectPermission "github.com/huydq/test/internal/domain/object/permission"
	repositoryPermission "github.com/huydq/test/internal/domain/repository/permission"
	repositoryRole "github.com/huydq/test/internal/domain/repository/role"
)

type RoleService interface {
	GetRoleByID(ctx context.Context, id int) (*modelRole.Role, error)
	GetRoleByCode(ctx context.Context, code string) (*modelRole.Role, error)
	GetRoleByName(ctx context.Context, name string) (*modelRole.Role, error)
	CreateRole(ctx context.Context, role *modelRole.Role) error
	UpdateRole(ctx context.Context, role *modelRole.Role) error
	DeleteRole(ctx context.Context, id int) error
	ListRoles(ctx context.Context) ([]*modelRole.Role, error)

	UpdateRolePermissions(ctx context.Context, roleID int, permissionIDs []int) error
	GetPermissionsByIDs(ctx context.Context, ids []int) ([]*modelPermission.Permission, error)

	// Newly added methods that were previously in PermissionService
	HasPermission(ctx context.Context, roleID int, permissions ...objectPermission.PermissionCode) (bool, error)
}

type roleServiceImpl struct {
	roleRepository       repositoryRole.RoleRepository
	permissionRepository repositoryPermission.PermissionRepository
}

func NewRoleService(
	roleRepository repositoryRole.RoleRepository,
	permissionRepository repositoryPermission.PermissionRepository,
) RoleService {
	return &roleServiceImpl{
		roleRepository:       roleRepository,
		permissionRepository: permissionRepository,
	}
}

func (s *roleServiceImpl) GetRoleByID(ctx context.Context, id int) (*modelRole.Role, error) {
	return s.roleRepository.FindByID(ctx, id)
}

func (s *roleServiceImpl) GetRoleByCode(ctx context.Context, code string) (*modelRole.Role, error) {
	return s.roleRepository.FindByCode(ctx, code)
}

func (s *roleServiceImpl) GetRoleByName(ctx context.Context, name string) (*modelRole.Role, error) {
	return s.roleRepository.FindByName(ctx, name)
}

func (s *roleServiceImpl) CreateRole(ctx context.Context, role *modelRole.Role) error {
	return s.roleRepository.Create(ctx, role)
}

func (s *roleServiceImpl) UpdateRole(ctx context.Context, role *modelRole.Role) error {
	return s.roleRepository.Update(ctx, role)
}

func (s *roleServiceImpl) DeleteRole(ctx context.Context, id int) error {
	return s.roleRepository.Delete(ctx, id)
}

func (s *roleServiceImpl) ListRoles(ctx context.Context) ([]*modelRole.Role, error) {
	return s.roleRepository.List(ctx)
}

func (s *roleServiceImpl) UpdateRolePermissions(ctx context.Context, roleID int, permissionIDs []int) error {
	role, err := s.roleRepository.FindByID(ctx, roleID)
	if err != nil {
		return err
	}

	permissions, err := s.permissionRepository.FindByIDs(ctx, permissionIDs)
	if err != nil {
		return err
	}

	role.Permissions = permissions

	return s.roleRepository.Update(ctx, role)
}

func (s *roleServiceImpl) GetPermissionsByIDs(ctx context.Context, ids []int) ([]*modelPermission.Permission, error) {
	return s.permissionRepository.FindByIDs(ctx, ids)
}

func (s *roleServiceImpl) HasPermission(ctx context.Context, roleID int, permissions ...objectPermission.PermissionCode) (bool, error) {
	role, err := s.roleRepository.FindByID(ctx, roleID)
	if err != nil {
		return false, err
	}

	if role == nil {
		return false, nil
	}

	return role.HasPermission(permissions...), nil
}
