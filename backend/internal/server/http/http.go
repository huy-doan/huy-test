package http

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/huydq/test/internal/pkg/common/response"
	"github.com/huydq/test/internal/pkg/errors"
	"github.com/huydq/test/internal/pkg/logger"
	ms "github.com/huydq/test/internal/pkg/utils/messages"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	echo    *echo.Echo
	address string
	logger  logger.Logger
}

func NewServer(logger logger.Logger) *Server {
	e := echo.New()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	address := fmt.Sprintf(":%s", port)

	e.HTTPErrorHandler = customHTTPErrorHandler

	return &Server{
		echo:    e,
		address: address,
		logger:  logger,
	}
}

func customHTTPErrorHandler(err error, c echo.Context) {
	if he, ok := err.(*echo.HTTPError); ok {
		if he.Code == http.StatusNotFound {
			notFoundErr := errors.NotFoundError(ms.MsgRouteNotFound)
			response.SendError(c, notFoundErr)
			return
		}
	}

	c.Echo().DefaultHTTPErrorHandler(err, c)
}

func (s *Server) SetupMiddleware() {
	s.echo.Use(middleware.Logger())
	s.echo.Use(middleware.Recover())
	s.echo.Use(middleware.CORS())
	s.echo.Use(middleware.RequestID())
}

func (s *Server) Echo() *echo.Echo {
	return s.echo
}

func (s *Server) Start() error {
	go func() {
		if err := s.echo.Start(s.address); err != nil {
			s.logger.Error("Failed to start server", map[string]any{"error": err.Error()})
		}
	}()

	s.logger.Info("Server started", map[string]any{"address": s.address})

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
