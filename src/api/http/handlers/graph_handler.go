package handlers

import (
	"context"
	"errors"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/gin-gonic/gin"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"github.com/vnlab/makeshop-payment/src/api/graphql/generated"
	"github.com/vnlab/makeshop-payment/src/api/graphql/middleware"
	"github.com/vnlab/makeshop-payment/src/api/graphql/resolvers"
	"github.com/vnlab/makeshop-payment/src/infrastructure/auth"
	"github.com/vnlab/makeshop-payment/src/infrastructure/logger"
	"github.com/vnlab/makeshop-payment/src/usecase"
)

type Graph interface {
	QueryHandler() gin.HandlerFunc
}

// GraphHandler handles GraphQL request processing
type GraphHandler struct {
	UserUsecase *usecase.UserUsecase
	JwtService  *auth.JWTService
}

// NewGraphHandler creates a new GraphHandler
func NewGraphHandler(us *usecase.UserUsecase, js *auth.JWTService) Graph {
	return &GraphHandler{
		UserUsecase: us,
		JwtService: js,
	}
}

// QueryHandler godoc
// @Summary GraphQL query endpoint
// @Description Process GraphQL queries and mutations
// @Tags graphql
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param query body object true "GraphQL query with optional variables and operationName"
// @Success 200 {object} object "GraphQL response"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /graphql [post]
func (h *GraphHandler) QueryHandler() gin.HandlerFunc {
	// TODO: Implement GraphQL loader

	graphHandler := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
		Resolvers: resolvers.NewResolver(h.UserUsecase, h.JwtService),
	}))

	// Configure GraphQL server options
	graphHandler.AddTransport(transport.POST{})
	graphHandler.Use(extension.Introspection{})
	
	// Add custom request logger to the GraphQL handler for logging
	graphHandler.AroundOperations(func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
		ginCtx, exists := ctx.Value("GinContextKey").(*gin.Context)
		if exists {
			loggerInstance, hasLogger := ginCtx.Get("logger")
			if hasLogger {
				log := loggerInstance.(logger.Logger)
				
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
	graphHandler.SetErrorPresenter(func(ctx context.Context, err error) *gqlerror.Error {
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

    return func(c *gin.Context) {
		// Send authentication information from Gin context to GraphQL context
		ctx := context.WithValue(c.Request.Context(), "GinContextKey", c)
		ctx = middleware.WithAuth(ctx, c)
		c.Request = c.Request.WithContext(ctx)

		graphHandler.ServeHTTP(c.Writer, c.Request)
    }
}