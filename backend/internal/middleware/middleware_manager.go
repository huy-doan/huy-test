package middleware

import (
	"github.com/huydq/test/internal/pkg/logger"
	"github.com/huydq/test/src/domain/models"
	"github.com/huydq/test/src/infrastructure/auth"
	"github.com/huydq/test/src/usecase"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type MiddlewareManager struct {
	logger          logger.Logger
	jwtService      *auth.JWTService
	auditLogUsecase *usecase.AuditLogUsecase

	// Middleware functions
	JWT           echo.MiddlewareFunc
	Admin         echo.MiddlewareFunc
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
	auditLogUsecase *usecase.AuditLogUsecase,
	db *gorm.DB,
) *MiddlewareManager {
	manager := &MiddlewareManager{
		logger:          logger,
		jwtService:      jwtService,
		auditLogUsecase: auditLogUsecase,
	}

	// Initialize all middleware functions
	manager.JWT = manager.JWTMiddleware(jwtService)
	manager.Admin = manager.RoleAuthorization(string(models.RoleCodeAdmin))
	manager.CORS = manager.CORSMiddleware()
	manager.ErrorHandler = manager.ErrorMiddleware()
	manager.RequestLogger = manager.RequestLoggerMiddleware()
	manager.Performance = manager.PerformanceMonitor(logger)
	manager.DBContext = manager.DBContextMiddleware(db)
	manager.AuditLogger = manager.NewAuditLogger(auditLogUsecase)

	return manager
}
