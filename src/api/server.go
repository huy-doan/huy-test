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

	"github.com/huydq/test/src/api/http/router"
	"github.com/huydq/test/src/domain/repositories"
	"github.com/huydq/test/src/infrastructure/auth"
	"github.com/huydq/test/src/infrastructure/logger"
	"github.com/huydq/test/src/lib/validator"
	"github.com/huydq/test/src/usecase"
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
	auditLogRepo repositories.AuditLogRepository,
	auditLogTypeRepo repositories.AuditLogTypeRepository,
	appLogger logger.Logger,
	twoFactorTokenRepo repositories.TwoFactorTokenRepository,
) *Server {
	// Set up validator
	validator.Setup()

	// Initialize services
	jwtService := auth.NewJWTService()
	auditLogUsecase := usecase.NewAuditLogUsecase(auditLogRepo, auditLogTypeRepo)
	userUsecase := usecase.NewUserUseCase(userRepo, roleRepo, jwtService)

	// Set up HTTP routes with the logger
	handler := router.SetupRouter(
		userRepo,
		roleRepo,
		jwtService,
		auditLogUsecase,
		appLogger,
		twoFactorTokenRepo,
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
			appLogger.Error("Server forced to shutdown", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()

	appLogger.Info("Server starting", map[string]interface{}{
		"address": s.httpServer.Addr,
	})

	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		appLogger.Error("Server failed to start", map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}

	return nil
}
