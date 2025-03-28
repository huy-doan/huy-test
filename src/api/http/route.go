package http

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/huydq/ddd-project/src/api/http/handlers"
	"github.com/huydq/ddd-project/src/api/http/middleware"
	"github.com/huydq/ddd-project/src/domain/models"
	"github.com/huydq/ddd-project/src/domain/repositories"
	"github.com/huydq/ddd-project/src/infrastructure/auth"
	"github.com/huydq/ddd-project/src/infrastructure/config"
	"github.com/huydq/ddd-project/src/usecase"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRouter sets up the Gin router with all routes and middleware
func SetupRouter(
	router *gin.Engine,
	userRepo repositories.UserRepository,
	roleRepo repositories.RoleRepository,
	jwtService *auth.JWTService,
) *gin.Engine {
	appConfig := config.GetConfig()
	allowOrigin := "*"
	allowHeader := []string{"Origin", "Content-Type", "Accept", "Authorization"}
	if appConfig.FrontUrl != "" {
		allowOrigin = appConfig.FrontUrl
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

	// Set middleware
	authMiddleware := middleware.AuthMiddleware(jwtService)
	adminMiddleware := middleware.RoleMiddleware(string(models.RoleCodeAdmin))

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})
	router.Use(middleware.ErrorHandler())

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

		// Protected routes
		protected := v1.Group("")
		protected.Use(authMiddleware)
		{
			// User routes
			users := protected.Group("/users")
			{
				users.GET("/profile", userHandler.GetProfile)
				users.PUT("/profile", userHandler.UpdateProfile)
				users.POST("/change-password", userHandler.ChangePassword)

				// Admin-only routes
				admin := users.Group("")
				admin.Use(adminMiddleware)
				{
					admin.GET("", userHandler.ListUsers)
					admin.GET("/:id", userHandler.GetUserByID)
				}
			}
		}
	}

	if gin.Mode() != gin.ReleaseMode {
		// Setup Swagger
		router.GET("/swaggers/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
	return router
}
