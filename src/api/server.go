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
	httpAPI "github.com/vnlab/makeshop-payment/src/api/http"
	"github.com/vnlab/makeshop-payment/src/domain/repositories"
	"github.com/vnlab/makeshop-payment/src/infrastructure/auth"
	"github.com/vnlab/makeshop-payment/src/lib/validator"
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
) *Server {
	// Set Gin mode
	ginMode := os.Getenv("GIN_MODE")
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
	// Set up graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("Shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := s.httpServer.Shutdown(ctx); err != nil {
			log.Fatalf("Server forced to shutdown: %v", err)
		}
	}()

	log.Printf("Server starting on %s", s.httpServer.Addr)
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}
