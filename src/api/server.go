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
	"github.com/vnlab/makeshop-payment/src/api/graphql"
	httpAPI "github.com/vnlab/makeshop-payment/src/api/http"
	"github.com/vnlab/makeshop-payment/src/domain/repositories"
	"github.com/vnlab/makeshop-payment/src/infrastructure/auth"
	"github.com/vnlab/makeshop-payment/src/infrastructure/logger"
	"github.com/vnlab/makeshop-payment/src/lib/validator"
	"github.com/vnlab/makeshop-payment/src/usecase"
)

// Server represents the API server
type Server struct {
	router         *gin.Engine
	httpServer     *http.Server
	jwtService     *auth.JWTService
	userUsecase    *usecase.UserUsecase
}

// NewServer creates a new API server
func NewServer(
	userRepo repositories.UserRepository,
	roleRepo repositories.RoleRepository,
	appLogger logger.Logger,
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
	userUsecase := usecase.NewUserUseCase(userRepo, roleRepo, jwtService)

	// Set up HTTP routes - FIX: Save the router returned from SetupRouter
	router = httpAPI.SetupRouter(
		router,
		jwtService,
	)

	// Set up GraphQL
	graphql.SetupGraphQL(
		router,  // This router instance is created but never assigned to the Server struct
		userUsecase,
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
		Handler: router,
	}

	return &Server{
		router:         router,
		httpServer:     httpServer,
		jwtService:     jwtService,
		userUsecase:    userUsecase,
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
