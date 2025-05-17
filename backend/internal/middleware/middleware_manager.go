package middleware

import (
	"github.com/huydq/test/internal/domain/service"
	tokenDomainSvc "github.com/huydq/test/internal/domain/service/auth"
	"github.com/huydq/test/internal/infrastructure/adapter/auth"
	"github.com/huydq/test/internal/pkg/logger"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type MiddlewareManager struct {
	logger          logger.Logger
	jwtService      *auth.JWTService
	auditLogService service.AuditLogService
	tokenDomainSvc  tokenDomainSvc.AccessTokenService
	roleService     service.RoleService

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
	tokenDomainSvc tokenDomainSvc.AccessTokenService,
	roleService service.RoleService,
	db *gorm.DB,
) *MiddlewareManager {
	manager := &MiddlewareManager{
		logger:          logger,
		jwtService:      jwtService,
		auditLogService: auditLogService,
		tokenDomainSvc:  tokenDomainSvc,
		roleService:     roleService,
	}

	// Initialize all middleware functions
	manager.JWT = manager.JWTMiddleware(jwtService, tokenDomainSvc)
	manager.CORS = manager.CORSMiddleware()
	manager.ErrorHandler = manager.ErrorMiddleware()
	manager.RequestLogger = manager.RequestLoggerMiddleware()
	manager.Performance = manager.PerformanceMonitor(logger)
	manager.DBContext = manager.DBContextMiddleware(db)
	manager.AuditLogger = manager.NewAuditLogger(auditLogService)

	return manager
}
