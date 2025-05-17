package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/huydq/test/internal/controller/auth"
	"github.com/huydq/test/internal/controller/merchant"
	"github.com/huydq/test/internal/controller/user"
	internalAuth "github.com/huydq/test/internal/infrastructure/adapter/auth"
	internalEmail "github.com/huydq/test/internal/infrastructure/adapter/email"
	auditLogPersistence "github.com/huydq/test/internal/infrastructure/persistence/audit_log"
	merchantPersistence "github.com/huydq/test/internal/infrastructure/persistence/merchant"
	payoutPersistence "github.com/huydq/test/internal/infrastructure/persistence/payout"
	payoutRecordPersistence "github.com/huydq/test/internal/infrastructure/persistence/payout_record"
	permissionPersistence "github.com/huydq/test/internal/infrastructure/persistence/permission"
	rolePersistence "github.com/huydq/test/internal/infrastructure/persistence/role"
	twoFactorPersistence "github.com/huydq/test/internal/infrastructure/persistence/two_factor_token"
	userPersistence "github.com/huydq/test/internal/infrastructure/persistence/user"
	"github.com/joho/godotenv"

	auditLogController "github.com/huydq/test/internal/controller/audit_log"
	payoutController "github.com/huydq/test/internal/controller/payout"
	permissionController "github.com/huydq/test/internal/controller/permission"
	roleController "github.com/huydq/test/internal/controller/role"

	"github.com/huydq/test/internal/domain/service"
	twoFactorTokenDomainService "github.com/huydq/test/internal/domain/service/auth"
	"github.com/huydq/test/internal/middleware"
	"github.com/huydq/test/internal/pkg/config"
	"github.com/huydq/test/internal/pkg/dbconn"
	"github.com/huydq/test/internal/pkg/logger"
	"github.com/huydq/test/internal/pkg/validator"
	"github.com/huydq/test/internal/server/http"
	"github.com/huydq/test/internal/server/router"
	auditLogUsecase "github.com/huydq/test/internal/usecase/audit_log"
	authUC "github.com/huydq/test/internal/usecase/auth"
	merchantUC "github.com/huydq/test/internal/usecase/merchant"
	payoutUsecase "github.com/huydq/test/internal/usecase/payout"
	permissionUsecase "github.com/huydq/test/internal/usecase/permission"
	roleUsecase "github.com/huydq/test/internal/usecase/role"
	userUC "github.com/huydq/test/internal/usecase/user"

	authService "github.com/huydq/test/internal/infrastructure/adapter/auth"
)

func init() {
	// Load environment variables from .env file if it exists
	envPath := filepath.Join(".", ".env")
	if _, err := os.Stat(envPath); err == nil {
		if err := godotenv.Load(envPath); err != nil {
			log.Printf("Warning: Error loading .env file: %v", err)
		}
	}
}

func main() {
	// Load configuration
	appConfig := config.LoadConfig()

	// Initialize logger
	logger.InitLogger(&logger.Config{
		LogLevel:         appConfig.LogLevel,
		LogDirectory:     appConfig.LogDirectory,
		EnableConsoleLog: appConfig.EnableConsoleLog,
	})
	appLogger := logger.GetLogger()

	db, err := dbconn.NewConnection(appLogger)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get SQL DB: %v", err)
	}
	defer sqlDB.Close()

	// Initialize internal repositories and services for auth usecase
	internalUserRepo := userPersistence.NewUserRepository(db)
	internalRoleRepo := rolePersistence.NewRoleRepository(db)
	internalPermissionRepo := permissionPersistence.NewPermissionRepository(db)
	internalAuditLogRepo := auditLogPersistence.NewAuditLogRepository(db)
	internalTwoFactorRepo := twoFactorPersistence.NewTwoFactorTokenRepository(db)
	internalMerchantRepo := merchantPersistence.NewMerchantRepository(db)
	internalPayoutRepo := payoutPersistence.NewPayoutRepository(db)
	internalPayoutRecordRepo := payoutRecordPersistence.NewPayoutRecordRepository(db)

	// Initialize services
	jwtService := authService.NewJWTService()
	internalJwtService := internalAuth.NewJWTService()
	mailService, err := internalEmail.NewMailService()
	if err != nil {
		log.Fatalf("Failed to create internal mail service: %v", err)
	}

	auditLogService := service.NewAuditLogService(internalAuditLogRepo, internalUserRepo)
	roleService := service.NewRoleService(internalRoleRepo, internalPermissionRepo)
	permissionService := service.NewPermissionService(internalPermissionRepo)
	payoutService := service.NewPayoutManagementService(internalPayoutRepo, internalPayoutRecordRepo)

	// Initialize usecases
	auditLogUsecase := auditLogUsecase.NewAuditLogUsecase(auditLogService)
	roleUsecase := roleUsecase.NewRoleUsecase(roleService)
	permissionUsecase := permissionUsecase.NewPermissionUsecase(permissionService)

	twoFactorDomainSvc := twoFactorTokenDomainService.NewTwoFactorTokenService(internalUserRepo, internalTwoFactorRepo, mailService)
	userManagementUsecase := userUC.NewManageUsersUsecase(internalUserRepo, internalRoleRepo)
	merchantManagementUsecase := merchantUC.NewManageMerchantsUsecase(internalMerchantRepo)

	authUsecase := authUC.NewAuthUsecase(internalUserRepo, internalTwoFactorRepo, internalJwtService, twoFactorDomainSvc)

	payoutUsecase := payoutUsecase.NewPayoutUsecase(payoutService)

	// Initialize controllers
	authController := auth.NewAuthController(authUsecase)
	userController := user.NewUserController(userManagementUsecase)
	merchantController := merchant.NewMerchantController(merchantManagementUsecase)
	roleController := roleController.NewRoleController(roleUsecase)
	permissionController := permissionController.NewPermissionController(permissionUsecase)
	auditLogController := auditLogController.NewAuditLogController(auditLogUsecase)
	payoutController := payoutController.NewPayoutController(payoutUsecase)

	// Create Echo server
	srv := http.NewServer(appLogger)

	// Create validator for Echo
	srv.Echo().Validator = validator.NewValidator()

	// Setup middleware manager with all middleware functions
	middlewareManager := middleware.NewMiddlewareManager(
		appLogger,
		jwtService,
		auditLogService,
		db,
	)

	// Apply global middleware
	srv.Echo().Use(middlewareManager.CORS)
	srv.Echo().Use(middlewareManager.RequestLogger)
	srv.Echo().Use(middlewareManager.Performance)
	srv.Echo().Use(middlewareManager.ErrorHandler)
	srv.Echo().Use(middlewareManager.RequestLogger)
	srv.Echo().Use(middlewareManager.DBContext)

	// Setup routes with the middleware manager
	router.SetupRoutes(
		srv.Echo(),
		authController,
		userController,
		merchantController,
		payoutController,
		roleController,
		permissionController,
		auditLogController,
		middlewareManager,
	)

	// Start server
	appLogger.Info("Starting server", map[string]any{"port": appConfig.ServerPort})
	if err := srv.Start(); err != nil {
		appLogger.Error("Server failed to start", map[string]any{"error": err.Error()})
	}
}
