package controller

import (
	"github.com/huydq/test/internal/controller/audit_log/mapper"
	"github.com/huydq/test/internal/controller/base"
	generated "github.com/huydq/test/internal/pkg/api/generated"
	response "github.com/huydq/test/internal/pkg/common/response"
	messages "github.com/huydq/test/internal/pkg/utils/messages"
	usecase "github.com/huydq/test/internal/usecase/audit_log"
	"github.com/labstack/echo/v4"
)

type AuditLogController struct {
	base.BaseController
	auditLogUsecase usecase.AuditLogUsecase
}

func NewAuditLogController(auditLogUsecase usecase.AuditLogUsecase) *AuditLogController {
	return &AuditLogController{
		auditLogUsecase: auditLogUsecase,
		BaseController:  *base.NewBaseController(),
	}
}

func (c *AuditLogController) ListAuditLogs(ctx echo.Context) error {
	var request generated.AuditLogListRequest
	if err := c.BindAndValidate(ctx, &request); err != nil {
		return response.SendError(ctx, err)
	}

	listInput := mapper.ToAuditLogFilter(&request)

	auditLogs, totalPages, total, err := c.auditLogUsecase.List(ctx.Request().Context(), listInput)
	if err != nil {
		return response.SendError(ctx, err)
	}

	responseData := mapper.ToAuditLogListResponse(
		auditLogs, request.Page,
		request.PageSize,
		totalPages,
		total,
	)

	return response.SendOK(ctx, messages.MsgListAuditLogsSuccess, responseData)
}

func (c *AuditLogController) GetAuditLogUsers(ctx echo.Context) error {
	output, err := c.auditLogUsecase.GetAuditLogUsers(ctx.Request().Context())
	if err != nil {
		return response.SendError(ctx, err)
	}

	responseData := mapper.ToAuditLogUsersData(output.Users)
	return response.SendOK(ctx, messages.MsgGetAuditLogUsersSuccess, responseData)
}
