package router

import (
	"os"

	"github.com/huydq/test/internal/controller/auth"
	merchantController "github.com/huydq/test/internal/controller/merchant"
	"github.com/huydq/test/internal/controller/payout"
	auditLogObject "github.com/huydq/test/internal/domain/object/audit_log"
	permissionObject "github.com/huydq/test/internal/domain/object/permission"
	"github.com/labstack/echo/v4"

	auditLogController "github.com/huydq/test/internal/controller/audit_log"
	permissionController "github.com/huydq/test/internal/controller/permission"
	roleController "github.com/huydq/test/internal/controller/role"
	"github.com/huydq/test/internal/controller/user"
	"github.com/huydq/test/internal/middleware"
)

func SetupRoutes(
	e *echo.Echo,
	authController *auth.AuthController,
	userController *user.UserController,
	merchantController *merchantController.MerchantController,
	payoutController *payout.PayoutController,
	roleController *roleController.RoleController,
	permissionController *permissionController.PermissionController,
	auditLogController *auditLogController.AuditLogController,
	middlewareManager *middleware.MiddlewareManager,
) {
	if os.Getenv("API_ENV") != "production" {
		SetupSwaggerUI(e)
	}

	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok"})
	})

	// API v1 routes
	api := e.Group("/api/v1")
	adminGroup := api.Group("/admin")
	{
		adminGroup.Use(middlewareManager.JWT)

		// Auth routes
		authGroup := api.Group("/auth")
		authGroup.POST("/login", authController.Login, middlewareManager.AuditLogger.WithType(auditLogObject.AuditLogTypeLogin).AsResponseMiddleware())
		authGroup.POST("/logout", authController.Logout, middlewareManager.JWT, middlewareManager.AuditLogger.WithType(auditLogObject.AuditLogTypeLogout).AsResponseMiddleware())
		authGroup.GET("/me", authController.Me, middlewareManager.JWT)
		authGroup.POST("/verify", authController.VerifyMFA)
		authGroup.POST("/resend-code", authController.ResendCode)

		// User management routes
		userGroup := api.Group("/admin", middlewareManager.JWT, middlewareManager.RoutePermissions(permissionObject.PermissionCodeUserManage))
		userGroup.GET("/users", userController.ListUsers)
		userGroup.POST("/users", userController.CreateUser)
		userGroup.PUT("/users/:id", userController.UpdateUser)
		userGroup.GET("/users/:id", userController.GetUserByID)
		userGroup.DELETE("/users/:id", userController.DeleteUser)
		// Merchant management routes
		merchantGroup := adminGroup.Group("/merchants", middlewareManager.RoutePermissions(permissionObject.PermissionCodeUserManage))
		merchantGroup.GET("", merchantController.ListMerchants)

		// Payout management routes
		payoutGroup := adminGroup.Group("/payouts")
		payoutGroup.GET("", payoutController.ListPayouts)

		// Role routes
		roleGroup := adminGroup.Group("/roles", middlewareManager.RoutePermissions(permissionObject.PermissionCodeUserManage))
		{
			roleGroup.GET("", roleController.ListRoles)
			roleGroup.GET("/:id", roleController.GetRoleByID)
			roleGroup.POST("", roleController.CreateRole)
			roleGroup.PUT("/:id", roleController.UpdateRole)
			roleGroup.DELETE("/:id", roleController.DeleteRole)
			roleGroup.POST("/permissions/batch", roleController.BatchUpdateRolePermissions)
		}

		// Permission routes
		permissionGroup := adminGroup.Group("/permissions")
		{
			permissionGroup.GET("", permissionController.ListPermissions)
		}

		// Audit log routes
		auditLogGroup := adminGroup.Group("/audit-logs", middlewareManager.RoutePermissions(permissionObject.PermissionCodeUserManage))
		{
			auditLogGroup.GET("", auditLogController.ListAuditLogs)
			auditLogGroup.GET("/users", auditLogController.GetAuditLogUsers)
		}
	}
}
