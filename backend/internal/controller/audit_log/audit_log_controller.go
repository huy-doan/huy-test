package controller

import (
	"github.com/huydq/test/internal/controller/audit_log/mapper"
	usecase "github.com/huydq/test/internal/usecase/audit_log"
	"github.com/labstack/echo/v4"
)

type AuditLogController struct {
	auditLogUsecase usecase.AuditLogUsecase
}

func NewAuditLogController(auditLogUsecase usecase.AuditLogUsecase) *AuditLogController {
	return &AuditLogController{
		auditLogUsecase: auditLogUsecase,
	}
}

func (c *AuditLogController) ListAuditLogs(ctx echo.Context) error {
	listInput, err := mapper.ToListAuditLogInput(ctx)
	if err != nil {
		return ctx.JSON(400, map[string]interface{}{
			"success": false,
			"message": "フィルターのパラメータが無効です",
		})
	}

	output, err := c.auditLogUsecase.List(ctx.Request().Context(), listInput)
	if err != nil {
		return ctx.JSON(500, map[string]interface{}{
			"success": false,
			"message": "監査ログの一覧取得に失敗しました",
		})
	}

	responseData := mapper.MapAuditLogListOutputToResponse(output)
	return ctx.JSON(200, map[string]interface{}{
		"success": true,
		"message": "監査ログの一覧を取得しました",
		"data":    responseData,
	})
}

// func (c *AuditLogController) GetAuditLogUsers(ctx echo.Context) error {
// 	users, err := c.auditLogUsecase.GetAuditLogUsers(ctx.Request().Context())
// 	if err != nil {
// 		return ctx.JSON(500, map[string]interface{}{
// 			"success": false,
// 			"message": "ユーザー一覧の取得に失敗しました",
// 		})
// 	}

// 	responseData := mapper.MapUsersToResponse(users)
// 	return ctx.JSON(200, map[string]interface{}{
// 		"success": true,
// 		"message": "ユーザー一覧を取得しました",
// 		"data":    responseData,
// 	})
// }
