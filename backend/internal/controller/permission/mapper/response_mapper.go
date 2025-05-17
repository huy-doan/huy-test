package mapper

import (
	"github.com/huydq/test/internal/datastructure/outputdata"
	model "github.com/huydq/test/internal/domain/model/permission"
)

func ToPermissionListResponse(output *outputdata.PermissionListOutput) model.PermissionListResponse {
	return model.PermissionListResponse{
		Permissions: output.Permissions,
		Total:       output.Total,
	}
}
