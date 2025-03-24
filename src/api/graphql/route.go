// src/api/graphql/server.go

package graphql

import (
	"context"
	"errors"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"github.com/vnlab/makeshop-payment/src/api/graphql/generated"
	"github.com/vnlab/makeshop-payment/src/api/graphql/middleware"
	"github.com/vnlab/makeshop-payment/src/api/graphql/resolvers"
	"github.com/vnlab/makeshop-payment/src/infrastructure/auth"
	"github.com/vnlab/makeshop-payment/src/infrastructure/logger"
	"github.com/vnlab/makeshop-payment/src/usecase"
)

// SetupGraphQL configures GraphQL handlers for the given Gin router
func SetupGraphQL(
	router *gin.Engine,
	userUsecase *usecase.UserUsecase,
	jwtService *auth.JWTService,
	appLogger logger.Logger,
) {
	// Set up authentication middleware for GraphQL
	authMiddleware := middleware.GraphQLAuthMiddleware(jwtService)
	loggerMiddleware := middleware.GraphQLLoggerMiddleware(appLogger)

	// Initialize GraphQL handler
	srv := handler.New(generated.NewExecutableSchema(generated.Config{
		Resolvers: resolvers.NewResolver(userUsecase, jwtService),
	}))

	// Configure GraphQL server options
	srv.AddTransport(transport.POST{})
	srv.Use(extension.Introspection{})
	
	// Add custom request logger to the GraphQL handler
	srv.AroundOperations(func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
		ginCtx, exists := ctx.Value("GinContextKey").(*gin.Context)
		if exists {
			loggerInstance, hasLogger := ginCtx.Get("logger")
			if hasLogger {
				log := loggerInstance.(logger.Logger)
				
				// Now we can safely get operation details
				opCtx := graphql.GetOperationContext(ctx)
				if opCtx != nil {
					log.Info("GraphQL operation", map[string]interface{}{
						"operation_name": opCtx.OperationName,
						"operation_type": opCtx.Operation.Operation,
					})
				}
				
				// Add logger to context for resolvers
				ctx = middleware.WithLogger(ctx, log)
			}
		}
		return next(ctx)
	})
	
	// Add custom error presenter to log errors
	srv.SetErrorPresenter(func(ctx context.Context, err error) *gqlerror.Error {
		// get logger from context
		log := middleware.GetLogger(ctx)
		
		// convert error to gqlerror.Error if possible
		var gqlErr *gqlerror.Error
		if errors.As(err, &gqlErr) {
			log.Error("GraphQL error", map[string]interface{}{
				"message": gqlErr.Message,
				"path": gqlErr.Path, 
				"locations": gqlErr.Locations,
			})
			return gqlErr
		}
		
		// If it's not a gqlerror.Error, create a new error
		gqlErr = &gqlerror.Error{
			Message: err.Error(),
		}
		
		// Get information about the field path from context if available
		path := graphql.GetFieldContext(ctx)
		if path != nil {
			gqlErr.Path = path.Path()
		}
		
		log.Error("GraphQL error", map[string]interface{}{
			"message": gqlErr.Message,
			"path": gqlErr.Path,
		})
		
		return gqlErr
	})

	
	// Setup GraphQL endpoint with middleware
	v1 := router.Group("/api/v1")
	{
		graphqlRoute := v1.Group("/graphql")
		graphqlRoute.Use(loggerMiddleware) // Add logger middleware first
		graphqlRoute.Use(authMiddleware)   // Then authentication
		{
			// POST endpoint for GraphQL API
			graphqlRoute.POST("", func(c *gin.Context) {
				// Store Gin context in GraphQL context to access it in the around operations
				ctx := context.WithValue(c.Request.Context(), "GinContextKey", c)
				
				// Add auth info
				ctx = middleware.WithAuth(ctx, c)
				
				// Update request with enhanced context
				c.Request = c.Request.WithContext(ctx)
				
				// Execute GraphQL handler
				srv.ServeHTTP(c.Writer, c.Request)
			})
		}

		// GraphQL Playground (development only)
		if gin.Mode() != gin.ReleaseMode {
			v1.GET("/playground", func(c *gin.Context) {
				playground.Handler("GraphQL Playground", "/api/v1/graphql").ServeHTTP(c.Writer, c.Request)
			})
		}
	}
}
