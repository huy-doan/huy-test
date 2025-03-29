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
	"github.com/vnlab/makeshop-payment/src/infrastructure/logger"
	"github.com/vnlab/makeshop-payment/src/lib/validator"
)

// Server represents the API server
type Server struct {
	httpServer *http.Server
	jwtService *auth.JWTService
}

// NewServer creates a new API server
func NewServer(
	userRepo repositories.UserRepository,
	roleRepo repositories.RoleRepository,
	appLogger logger.Logger,
) *Server {
	// Set up validator
	validator.Setup()

	// Initialize services
	jwtService := auth.NewJWTService()

	// Set up HTTP routes with the logger
	handler := router.SetupRouter(
		userRepo,
		roleRepo,
		jwtService,
		appLogger,
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
		httpServer: httpServer,
		jwtService: jwtService,
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
