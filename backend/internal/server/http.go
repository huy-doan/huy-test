package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/vnlab/makeshop-payment/internal/pkg/log"
)

// Server represents the HTTP server using Echo framework
type Server struct {
	echo     *echo.Echo
	address  string
	logger   log.Logger
}

// NewServer creates a new Echo server instance
func NewServer(logger log.Logger) *Server {
	// Initialize Echo
	e := echo.New()
	
	// Set address from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	address := fmt.Sprintf(":%s", port)

	return &Server{
		echo:     e,
		address:  address,
		logger:   logger,
	}
}

// Setup configures the Echo server with middleware and routes
func (s *Server) Setup() {
	// Add basic middleware
	s.echo.Use(middleware.Logger())
	s.echo.Use(middleware.Recover())
	s.echo.Use(middleware.CORS())
	
	// Swagger documentation
	s.echo.GET("/swagger/*", echoSwagger.WrapHandler)
}

// GetEcho returns the echo instance for direct manipulation
func (s *Server) GetEcho() *echo.Echo {
	return s.echo
}

// Start starts the HTTP server with graceful shutdown
func (s *Server) Start() error {
	// Start server in a goroutine
	go func() {
		if err := s.echo.Start(s.address); err != nil {
			s.logger.Error("Failed to start server", map[string]any{"error": err.Error()})
		}
	}()

	s.logger.Info("Server started", map[string]any{"address": s.address})

	// Setup graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	s.logger.Info("Shutting down server...", nil)
	if err := s.echo.Shutdown(ctx); err != nil {
		s.logger.Error("Error during server shutdown", map[string]any{"error": err.Error()})
		return err
	}

	return nil
}