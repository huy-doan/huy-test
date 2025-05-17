package router

import (
	"net/http"
	"os"

	"github.com/huydq/test/src/api/http/handlers"
	"github.com/huydq/test/src/api/http/middleware"
	"github.com/huydq/test/src/domain/models"
	"github.com/huydq/test/src/domain/repositories"
	"github.com/huydq/test/src/infrastructure/auth"
	"github.com/huydq/test/src/infrastructure/logger"
	"github.com/huydq/test/src/usecase"
	httpSwagger "github.com/swaggo/http-swagger"
)

// SetupRouter sets up the HTTP router with all routes and middleware
func SetupRouter(
	userRepo repositories.UserRepository,
	roleRepo repositories.RoleRepository,
	jwtService *auth.JWTService,
	auditLogUsecase *usecase.AuditLogUsecase,
	appLogger logger.Logger,
	twoFactorRepo repositories.TwoFactorTokenRepository,
) http.Handler {
	// Create main router
	mux := http.NewServeMux()

	// Create usecases
	userUsecase := usecase.NewUserUseCase(userRepo, roleRepo, jwtService)
	twoFAUsecase := usecase.NewTwoFAUsecase(userRepo, twoFactorRepo, jwtService)

	// Create handlers
	authHandler := handlers.NewAuthHandler(userUsecase, jwtService, auditLogUsecase, twoFAUsecase)
	userHandler := handlers.NewUserHandler(userUsecase, jwtService)

	// Set up middleware
	errorMiddleware := middleware.ErrorHandler
	corsMiddleware := middleware.CORSMiddleware
	requestLoggerMiddleware := middleware.RequestLoggerMiddleware(appLogger)
	authMiddleware := middleware.NewAuthMiddleware(jwtService)
	adminMiddleware := middleware.NewRoleMiddleware(string(models.RoleCodeAdmin))
	languageMiddleware := middleware.LanguageMiddleware

	// Health check endpoint
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// API v1 routes - Authentication
	mux.HandleFunc("POST /api/v1/auth/login", authHandler.Login)
	mux.HandleFunc("POST /api/v1/auth/register", authHandler.Register)
	mux.HandleFunc("POST /api/v1/auth/logout", authMiddleware(authHandler.Logout))
	mux.HandleFunc("GET /api/v1/auth/me", authMiddleware(authHandler.Me))
	mux.HandleFunc("POST /api/v1/auth/verify", authHandler.VerifyMFA)

	// API v1 routes - User Self-Management
	mux.HandleFunc("GET /api/v1/users/profile", authMiddleware(userHandler.GetProfile))
	mux.HandleFunc("PUT /api/v1/users/profile", authMiddleware(userHandler.UpdateProfile))
	mux.HandleFunc("POST /api/v1/users/change-password", authMiddleware(userHandler.ChangePassword))

	// Admin-only routes - User Management
	mux.HandleFunc("GET /api/v1/admin/users", authMiddleware(adminMiddleware(userHandler.ListUsers)))
	mux.HandleFunc("GET /api/v1/admin/users/{id}", authMiddleware(adminMiddleware(userHandler.GetUserByID)))
	mux.HandleFunc("PUT /api/v1/admin/users/{id}", authMiddleware(adminMiddleware(userHandler.UpdateUserProfile)))
	mux.HandleFunc("POST /api/v1/admin/users", authMiddleware(adminMiddleware(userHandler.CreateUser)))
	mux.HandleFunc("POST /api/v1/admin/users/reset-password", authMiddleware(adminMiddleware(userHandler.ResetPasswordByAdmin)))
	mux.HandleFunc("DELETE /api/v1/admin/users/{id}", authMiddleware(adminMiddleware(userHandler.DeleteUser)))

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
	var handler http.Handler = mux
	handler = corsMiddleware(handler)
	handler = middleware.PerformanceMonitor(appLogger)(handler)
	handler = errorMiddleware(handler)
	handler = requestLoggerMiddleware(handler)
	handler = middleware.WithLogger(appLogger)(handler)
	handler = languageMiddleware(handler)

	return handler
}
