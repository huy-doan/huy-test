package usecase

import (
	"context"

	"github.com/huydq/test/internal/datastructure/converter"
	"github.com/huydq/test/internal/datastructure/inputdata"
	"github.com/huydq/test/internal/datastructure/outputdata"
	"github.com/huydq/test/internal/domain/service"
)

type PermissionUsecase interface {
	GetPermissionsByIDs(ctx context.Context, input *inputdata.GetPermissionsByIDsInput) ([]*outputdata.PermissionOutput, error)
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
) ([]*outputdata.PermissionOutput, error) {
	permissions, err := u.permissionService.GetPermissionsByIDs(ctx, input.IDs)
	if err != nil {
		return nil, err
	}

	return converter.PermissionModelsToOutputs(permissions), nil
}

func (u *permissionUsecaseImpl) ListPermissions(
	ctx context.Context,
) (*outputdata.PermissionListOutput, error) {
	permissions, err := u.permissionService.ListPermissions(ctx)
	if err != nil {
		return nil, err
	}

	return &outputdata.PermissionListOutput{
		Permissions: converter.PermissionModelsToOutputs(permissions),
		Total:       int64(len(permissions)),
	}, nil
}
