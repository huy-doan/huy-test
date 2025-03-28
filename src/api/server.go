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

	"github.com/gin-gonic/gin"
	httpAPI "github.com/huydq/ddd-project/src/api/http"
	"github.com/huydq/ddd-project/src/domain/repositories"
	"github.com/huydq/ddd-project/src/infrastructure/auth"
	"github.com/huydq/ddd-project/src/infrastructure/config"
	"github.com/huydq/ddd-project/src/infrastructure/logger"
	"github.com/huydq/ddd-project/src/lib/validator"
)

// Server represents the API server
type Server struct {
	router     *gin.Engine
	httpServer *http.Server
	jwtService *auth.JWTService
}

// NewServer creates a new API server
func NewServer(
	userRepo repositories.UserRepository,
	roleRepo repositories.RoleRepository,
	appLogger logger.Logger,
) *Server {
	appConfig := config.GetConfig()
	// Set Gin mode
	ginMode := appConfig.GinMode
	if ginMode != "" {
		gin.SetMode(ginMode)
	}

	// Set up validator
	validator.Setup()

	// Create router
	router := gin.Default()

	// Initialize services
	jwtService := auth.NewJWTService()

	// Set up HTTP routes
	router = httpAPI.SetupRouter(
		router,
		userRepo,
		roleRepo,
		jwtService,
	)

	// Create HTTP server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: router,
	}

	return &Server{
		router:     router,
		httpServer: httpServer,
		jwtService: jwtService,
	}
}

// Start starts the API server
func (s *Server) Start() error {
	// Get the global logger
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
