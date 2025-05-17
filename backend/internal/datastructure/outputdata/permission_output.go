package outputdata

import model "github.com/huydq/test/internal/domain/model/permission"

// PermissionListOutput represents the output for listing permissions
type PermissionListOutput struct {
	Permissions []*model.Permission `json:"permissions"`
	Total       int64               `json:"total"`
}
