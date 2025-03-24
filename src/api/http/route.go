package http

import (
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/vnlab/makeshop-payment/src/infrastructure/auth"
)

// SetupRouter sets up the Gin router with all routes and middleware
func SetupRouter(
	router *gin.Engine,
	jwtService *auth.JWTService,
) *gin.Engine {
	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{os.Getenv("FRONT_URL")}
	config.AllowCredentials = true
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	router.Use(cors.New(config))

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	if gin.Mode() != gin.ReleaseMode {
		// Setup Swagger
		router.GET("/swaggers/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
	return router
}
