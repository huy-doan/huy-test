package middleware

import (
	"github.com/huydq/test/internal/domain/service"
	"github.com/huydq/test/internal/pkg/logger"
	"github.com/huydq/test/src/infrastructure/auth"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type MiddlewareManager struct {
	logger                      logger.Logger
	jwtService                  *auth.JWTService
	auditLogService             service.AuditLogService
	permissionMiddlewareService service.PermissionMiddlewareService

	// Middleware functions
	JWT           echo.MiddlewareFunc
	Language      echo.MiddlewareFunc
	CORS          echo.MiddlewareFunc
	ErrorHandler  echo.MiddlewareFunc
	RequestLogger echo.MiddlewareFunc
	Performance   echo.MiddlewareFunc
	DBContext     echo.MiddlewareFunc
	AuditLogger   *AuditLogBuilder
}

func NewMiddlewareManager(
	logger logger.Logger,
	jwtService *auth.JWTService,
	auditLogService service.AuditLogService,
	permissionMiddlewareService service.PermissionMiddlewareService,
	db *gorm.DB,
) *MiddlewareManager {
	manager := &MiddlewareManager{
		logger:                      logger,
		jwtService:                  jwtService,
		auditLogService:             auditLogService,
		permissionMiddlewareService: permissionMiddlewareService,
	}

	// Initialize all middleware functions
	manager.JWT = manager.JWTMiddleware(jwtService)
	manager.CORS = manager.CORSMiddleware()
	manager.ErrorHandler = manager.ErrorMiddleware()
	manager.RequestLogger = manager.RequestLoggerMiddleware()
	manager.Performance = manager.PerformanceMonitor(logger)
	manager.DBContext = manager.DBContextMiddleware(db)
	manager.AuditLogger = manager.NewAuditLogger(auditLogService)

	return manager
}
