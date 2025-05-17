package usecase

import (
	"context"

	"github.com/huydq/test/internal/datastructure/inputdata"
	"github.com/huydq/test/internal/datastructure/outputdata"
	model "github.com/huydq/test/internal/domain/model/permission"
	"github.com/huydq/test/internal/domain/service"
)

type PermissionUsecase interface {
	GetPermissionsByIDs(ctx context.Context, input *inputdata.GetPermissionsByIDsInput) ([]*model.Permission, error)
	ListPermissions(ctx context.Context) (*outputdata.PermissionListOutput, error)
}

type permissionUsecaseImpl struct {
	permissionService service.PermissionService
}

func NewPermissionUsecase(
	permissionService service.PermissionService,
) PermissionUsecase {
	return &permissionUsecaseImpl{
		permissionService: permissionService,
	}
}

func (u *permissionUsecaseImpl) GetPermissionsByIDs(
	ctx context.Context,
	input *inputdata.GetPermissionsByIDsInput,
) ([]*model.Permission, error) {
	permissions, err := u.permissionService.GetPermissionsByIDs(ctx, input.IDs)
	if err != nil {
		return nil, err
	}

	return permissions, nil
}

func (u *permissionUsecaseImpl) ListPermissions(
	ctx context.Context,
) (*outputdata.PermissionListOutput, error) {
	permissions, err := u.permissionService.ListPermissions(ctx)
	if err != nil {
		return nil, err
	}

	return &outputdata.PermissionListOutput{
		Permissions: permissions,
		Total:       int64(len(permissions)),
	}, nil
}
