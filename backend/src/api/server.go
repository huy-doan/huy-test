package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vnlab/makeshop-payment/src/api/http/router"
	"github.com/vnlab/makeshop-payment/src/domain/repositories"
	"github.com/vnlab/makeshop-payment/src/infrastructure/auth"
	"github.com/vnlab/makeshop-payment/src/infrastructure/email"
	"github.com/vnlab/makeshop-payment/src/infrastructure/logger"
	"github.com/vnlab/makeshop-payment/src/lib/validator"
	"github.com/vnlab/makeshop-payment/src/usecase"
)

// Server represents the API server
type Server struct {
	httpServer      *http.Server
	jwtService      *auth.JWTService
	userUsecase     *usecase.UserUsecase
	auditLogUsecase *usecase.AuditLogUsecase
}

// NewServer creates a new API server
func NewServer(
	userRepo repositories.UserRepository,
	roleRepo repositories.RoleRepository,
	permissionRepo repositories.PermissionRepository,
	payoutRepo repositories.PayoutRepository,
	payoutRecordRepo repositories.PayoutRecordRepository,
	auditLogRepo repositories.AuditLogRepository,
	auditLogTypeRepo repositories.AuditLogTypeRepository,
	appLogger logger.Logger,
	twoFactorTokenRepo repositories.TwoFactorTokenRepository,
	merchantRepo repositories.MerchantRepository,
) *Server {
	// Set up validator
	validator.Setup()

	// Initialize services
	jwtService := auth.NewJWTService()
	auditLogUsecase := usecase.NewAuditLogUsecase(auditLogRepo, auditLogTypeRepo, userRepo)
	userUsecase := usecase.NewUserUseCase(userRepo, roleRepo, jwtService)
	mailService, err := email.NewMailService()
	if err != nil {
		log.Fatalf("Failed to create mail service: %v", err)
	}

	// Set up HTTP routes with the logger
	handler := router.SetupRouter(
		userRepo,
		roleRepo,
		permissionRepo,
		payoutRepo,
		payoutRecordRepo,
		merchantRepo,
		jwtService,
		auditLogUsecase,
		appLogger,
		twoFactorTokenRepo,
		mailService,
	)

	// Create HTTP server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: handler,
	}

	return &Server{
		httpServer:      httpServer,
		jwtService:      jwtService,
		userUsecase:     userUsecase,
		auditLogUsecase: auditLogUsecase,
	}
}

// Start starts the API server
func (s *Server) Start() error {
	appLogger := logger.GetLogger()

	// Set up graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("Shutting down server...")

		// Log shutdown with the application logger
		appLogger.Info("Shutting down server gracefully", nil)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := s.httpServer.Shutdown(ctx); err != nil {
			appLogger.Error("Server forced to shutdown", map[string]any{
				"error": err.Error(),
			})
		}
	}()

	appLogger.Info("Server starting", map[string]any{
		"address": s.httpServer.Addr,
	})

	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		appLogger.Error("Server failed to start", map[string]any{
			"error": err.Error(),
		})
		return err
	}

	return nil
}
