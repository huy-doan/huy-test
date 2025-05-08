package server

import (
	"github.com/labstack/echo/v4"
	"github.com/vnlab/makeshop-payment/internal/controller/auth"
	"github.com/vnlab/makeshop-payment/internal/middleware"
)

func SetupRoutes(e *echo.Echo, authController *auth.AuthController, middlewareManager *middleware.MiddlewareManager) {
	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok"})
	})

	// API v1 routes
	api := e.Group("/api/v2")

	// Auth routes
	authGroup := api.Group("/auth")
	authGroup.POST("/login", authController.Login)
	authGroup.POST("/register", authController.Register)
	authGroup.POST("/logout", authController.Logout, middlewareManager.JWT)
	authGroup.GET("/me", authController.Me, middlewareManager.JWT)
	authGroup.POST("/verify", authController.VerifyMFA)
	authGroup.POST("/resend-code", authController.ResendCode)
}
