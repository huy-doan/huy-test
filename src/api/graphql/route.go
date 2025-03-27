package graphql

import (
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/vnlab/makeshop-payment/src/api/graphql/middleware"
	"github.com/vnlab/makeshop-payment/src/api/http/handlers"
	"github.com/vnlab/makeshop-payment/src/infrastructure/auth"
	"github.com/vnlab/makeshop-payment/src/usecase"
)

// SetupGraphQL configures GraphQL handlers for the given Gin router
func SetupGraphQL(
	router *gin.Engine,
	userUsecase *usecase.UserUsecase,
	jwtService *auth.JWTService,
) {
	// Set up authentication middleware for GraphQL
	graphAuthMiddleware := middleware.GraphQLAuthMiddleware(jwtService)

	// Initialize GraphQL handler
	graphHandler := handlers.NewGraphHandler(userUsecase, jwtService)

	// Setup GraphQL endpoint with middleware
	v1 := router.Group("/api/v1")
	{
		graphqlRoute := v1.Group("/graphql")
		graphqlRoute.Use(graphAuthMiddleware)
		{
			// Main endpoint for GraphQL API
			graphqlRoute.POST("", graphHandler.QueryHandler())
		}

		// GraphQL Playground (development only)
		if gin.Mode() != gin.ReleaseMode {
			v1.GET("/playground", func(c *gin.Context) {
				playground.Handler("GraphQL Playground", "/api/v1/graphql").ServeHTTP(c.Writer, c.Request)
			})
		}
	}
}
