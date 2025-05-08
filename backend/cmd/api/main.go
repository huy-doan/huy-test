package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/vnlab/makeshop-payment/internal/controller/auth"
	"github.com/vnlab/makeshop-payment/internal/middleware"
	"github.com/vnlab/makeshop-payment/internal/pkg/config"
	"github.com/vnlab/makeshop-payment/internal/pkg/dbconn"
	"github.com/vnlab/makeshop-payment/internal/pkg/logger"
	"github.com/vnlab/makeshop-payment/internal/pkg/validator"
	"github.com/vnlab/makeshop-payment/internal/server/http"
	server "github.com/vnlab/makeshop-payment/internal/server/router"
	authService "github.com/vnlab/makeshop-payment/src/infrastructure/auth"
	"github.com/vnlab/makeshop-payment/src/infrastructure/email"
	"github.com/vnlab/makeshop-payment/src/infrastructure/persistence/repositories"
	"github.com/vnlab/makeshop-payment/src/usecase"
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

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	roleRepo := repositories.NewRoleRepository(db)
	// permissionRepo := repositories.NewPermissionRepository(db)
	auditLogRepo := repositories.NewAuditLogRepository(db)
	auditLogTypeRepo := repositories.NewAuditLogTypeRepository(db)
	twoFactorTokenRepo := repositories.NewTwoFactorTokenRepository(db)
	// payoutRepo := repositories.NewPayoutRepository(db)
	// payoutRecordRepo := repositories.NewPayoutRecordRepository(db)
	// merchantRepo := repositories.NewMerchantRepository(db)

	// Initialize services
	jwtService := authService.NewJWTService()
	auditLogUsecase := usecase.NewAuditLogUsecase(auditLogRepo, auditLogTypeRepo, userRepo)
	userUsecase := usecase.NewUserUseCase(userRepo, roleRepo, jwtService)

	// Initialize mail service
	mailService, err := email.NewMailService()
	if err != nil {
		log.Fatalf("Failed to create mail service: %v", err)
	}

	twoFAUsecase := usecase.NewTwoFAUsecase(userRepo, twoFactorTokenRepo, jwtService, mailService)
	// roleUsecase := usecase.NewRoleUsecase(roleRepo, permissionRepo)
	// permissionUsecase := usecase.NewPermissionUseCase(permissionRepo)
	// payoutUsecase := usecase.NewPayoutUsecase(payoutRepo, payoutRecordRepo)
	// merchantUsecase := usecase.NewMerchantUsecase(merchantRepo)

	// Create Echo server
	srv := http.NewServer(appLogger)

	// Create validator for Echo
	srv.Echo().Validator = validator.NewValidator()

	// Setup middleware manager with all middleware functions
	middlewareManager := middleware.NewMiddlewareManager(
		appLogger,
		jwtService,
		auditLogUsecase,
		db,
	)

	// Apply global middleware
	srv.Echo().Use(middlewareManager.CORS)
	srv.Echo().Use(middlewareManager.RequestLogger)
	srv.Echo().Use(middlewareManager.Performance)
	srv.Echo().Use(middlewareManager.ErrorHandler)
	srv.Echo().Use(middlewareManager.RequestLogger)

	// Create controllers
	authController := auth.NewAuthController(userUsecase, jwtService, auditLogUsecase, twoFAUsecase)

	// TODO: Implement the rest of the controllers
	// userController := auth.NewUserController(userUsecase, jwtService, auditLogUsecase, appLogger)
	// roleController := auth.NewRoleController(roleUsecase, appLogger)
	// permissionController := auth.NewPermissionController(permissionUsecase, appLogger)
	// payoutController := auth.NewPayoutController(payoutUsecase, appLogger)
	// auditLogController := auth.NewAuditLogController(auditLogUsecase, appLogger)
	// merchantController := auth.NewMerchantController(merchantUsecase, appLogger)

	// Setup routes with the middleware manager
	server.SetupRoutes(
		srv.Echo(),
		authController,
		// userController,
		// roleController,
		// permissionController,
		// payoutController,
		// auditLogController,
		// merchantController,
		middlewareManager,
	)

	// Start server
	appLogger.Info("Starting server", map[string]any{"port": appConfig.ServerPort})
	if err := srv.Start(); err != nil {
		appLogger.Error("Server failed to start", map[string]any{"error": err.Error()})
	}
}
