package mapper

import (
	controllerUtil "github.com/huydq/test/internal/controller/util"
	model "github.com/huydq/test/internal/domain/model/audit_log"
	"github.com/huydq/test/internal/domain/model/util"
	generated "github.com/huydq/test/internal/pkg/api/generated"
)

func ToAuditLogFilter(request *generated.AuditLogListRequest) *model.AuditLogFilter {
	filter := model.NewAuditLogFilter()
	filter.SetPagination(request.Page, request.PageSize)
	sortOrder := util.Ascending
	if request.SortOrder == "desc" {
		sortOrder = util.Descending
	}
	filter.SetSort(request.SortField, sortOrder)
	if request.CreatedAt != "" {
		createdAt, err := controllerUtil.ParseDateFromString(request.CreatedAt)
		if err == nil {
			filter.CreatedAt = &createdAt
		}
	}
	if request.UserId != 0 {
		filter.UserID = &request.UserId
	}
	if request.Description != "" {
		filter.Description = &request.Description
	}
	if request.AuditLogType != "" {
		filter.AuditLogType = &request.AuditLogType
	}

	return filter
}
