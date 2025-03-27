// src/api/http/route.go - cập nhật
package http

import (
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/vnlab/makeshop-payment/src/api/http/handlers"
	"github.com/vnlab/makeshop-payment/src/api/http/middleware"
	"github.com/vnlab/makeshop-payment/src/domain/repositories"
	"github.com/vnlab/makeshop-payment/src/infrastructure/auth"
	"github.com/vnlab/makeshop-payment/src/usecase"
)

// SetupRouter sets up the Gin router with all routes and middleware
func SetupRouter(
	router *gin.Engine,
	userRepo repositories.UserRepository,
	roleRepo repositories.RoleRepository,
	jwtService *auth.JWTService,
) *gin.Engine {
	allowOrigin := "*"
	allowHeader := []string{"Origin", "Content-Type", "Accept", "Authorization"}
	if os.Getenv("API_FRONT_URL") != "" {
		allowOrigin = os.Getenv("API_FRONT_URL")
	}

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{allowOrigin}
	config.AllowCredentials = true
	config.AllowHeaders = allowHeader
	router.Use(cors.New(config))

	// Create usecases
	userUsecase := usecase.NewUserUseCase(userRepo, roleRepo, jwtService)

	// Create handlers
	authHandler := handlers.NewAuthHandler(userUsecase, jwtService)
	userHandler := handlers.NewUserHandler(userUsecase, jwtService)

	// Set up authentication middleware
	authMiddleware := middleware.AuthMiddleware(jwtService)
	adminMiddleware := middleware.RoleMiddleware("SYSTEM_ADMIN")

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// router.Use(middleware.ErrorHandler())

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Auth routes
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/register", authHandler.Register)
			auth.POST("/logout", authMiddleware, authHandler.Logout)
		}

		// User routes (protected)
		users := v1.Group("/users")
		users.Use(authMiddleware)
		{
			users.GET("", adminMiddleware, userHandler.ListUsers)
			users.GET("/:id", userHandler.GetUserByID)
			users.GET("/profile", userHandler.GetProfile)
			users.PUT("/profile", userHandler.UpdateProfile)
			users.POST("/change-password", userHandler.ChangePassword)
		}
	}

	if gin.Mode() != gin.ReleaseMode {
		// Setup Swagger
		router.GET("/swaggers/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	return router
}
