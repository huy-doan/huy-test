package handlers

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/gin-gonic/gin"
	"github.com/vnlab/makeshop-payment/src/api/graphql/generated"
	"github.com/vnlab/makeshop-payment/src/api/graphql/middleware"
	"github.com/vnlab/makeshop-payment/src/api/graphql/resolvers"
	"github.com/vnlab/makeshop-payment/src/infrastructure/auth"
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

    return func(c *gin.Context) {
		// Send authentication information from Gin context to GraphQL context
		ctx := middleware.WithAuth(c.Request.Context(), c)
		c.Request = c.Request.WithContext(ctx)

		graphHandler.ServeHTTP(c.Writer, c.Request)
    }
}