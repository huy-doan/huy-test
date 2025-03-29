package router

import (
	"net/http"
	"os"

	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/vnlab/makeshop-payment/src/api/http/handlers"
	"github.com/vnlab/makeshop-payment/src/api/http/middleware"
	"github.com/vnlab/makeshop-payment/src/domain/models"
	"github.com/vnlab/makeshop-payment/src/domain/repositories"
	"github.com/vnlab/makeshop-payment/src/infrastructure/auth"
	"github.com/vnlab/makeshop-payment/src/infrastructure/logger"
	"github.com/vnlab/makeshop-payment/src/usecase"
)

// SetupRouter sets up the HTTP router with all routes and middleware
func SetupRouter(
	userRepo repositories.UserRepository,
	roleRepo repositories.RoleRepository,
	jwtService *auth.JWTService,
	appLogger logger.Logger,
) http.Handler {
	// Create main router
	mux := http.NewServeMux()

	// Create usecases
	userUsecase := usecase.NewUserUseCase(userRepo, roleRepo, jwtService)

	// Create handlers
	authHandler := handlers.NewAuthHandler(userUsecase, jwtService)
	userHandler := handlers.NewUserHandler(userUsecase, jwtService)

	// Set up middleware
	errorMiddleware := middleware.ErrorHandler
	corsMiddleware := middleware.CORSMiddleware
	requestLoggerMiddleware := middleware.RequestLoggerMiddleware(appLogger)
	authMiddleware := middleware.NewAuthMiddleware(jwtService)
	adminMiddleware := middleware.NewRoleMiddleware(string(models.RoleCodeAdmin))

	// Health check endpoint
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// API v1 routes - Auth
	mux.HandleFunc("POST /api/v1/auth/login", authHandler.Login)
	mux.HandleFunc("POST /api/v1/auth/register", authHandler.Register)
	mux.HandleFunc("POST /api/v1/auth/logout", authMiddleware(authHandler.Logout))

	// API v1 routes - User
	mux.HandleFunc("GET /api/v1/users/profile", authMiddleware(userHandler.GetProfile))
	mux.HandleFunc("PUT /api/v1/users/profile", authMiddleware(userHandler.UpdateProfile))
	mux.HandleFunc("POST /api/v1/users/change-password", authMiddleware(userHandler.ChangePassword))

	// Admin-only routes
	mux.HandleFunc("GET /api/v1/users", authMiddleware(adminMiddleware(userHandler.ListUsers)))
	mux.HandleFunc("GET /api/v1/users/{id}", authMiddleware(adminMiddleware(userHandler.GetUserByID)))

	// Swagger documentation (in development mode only)
	if os.Getenv("API_ENV") != "production" {
		mux.HandleFunc("GET /swagger/", func(w http.ResponseWriter, r *http.Request) {
			httpSwagger.Handler(
				httpSwagger.URL("/swagger/doc.json"),
			).ServeHTTP(w, r)
		})
		mux.HandleFunc("GET /swagger", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/swagger/", http.StatusMovedPermanently)
		})
	}

	// Apply global middleware
	// The order is important:
	// 1. Request Logger - logs all requests and adds trace ID, creates response writer
	// 2. Error Handler - captures and formats errors, uses response writer from request logger
	// 3. Performance Monitor - tracks request performance (only for slow requests)
	// 4. CORS - handles preflight requests and sets headers
	// 5. WithLogger - adds logger to context for use by other middleware
	var handler http.Handler = mux
	handler = corsMiddleware(handler)
	handler = middleware.PerformanceMonitor(appLogger)(handler)
	handler = errorMiddleware(handler)
	handler = requestLoggerMiddleware(handler)
	handler = middleware.WithLogger(appLogger)(handler)

	return handler
}