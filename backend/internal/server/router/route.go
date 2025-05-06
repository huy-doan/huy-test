package router

import (
	"net/http"
	"os"

	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/vnlab/makeshop-payment/internal/middleware"
	"github.com/vnlab/makeshop-payment/internal/pkg/logger"
	"github.com/vnlab/makeshop-payment/src/api/http/handlers"
	"github.com/vnlab/makeshop-payment/src/domain/models"
	"github.com/vnlab/makeshop-payment/src/domain/repositories"
	"github.com/vnlab/makeshop-payment/src/infrastructure/auth"
	"github.com/vnlab/makeshop-payment/src/infrastructure/email"
	"github.com/vnlab/makeshop-payment/src/usecase"
)

// SetupRouter sets up the HTTP router with all routes and middleware
func SetupRouter(
	userRepo repositories.UserRepository,
	roleRepo repositories.RoleRepository,
	jwtService *auth.JWTService,
	auditLogUsecase *usecase.AuditLogUsecase,
	appLogger logger.Logger,
	twoFactorRepo repositories.TwoFactorTokenRepository,
	mailService *email.MailService,
) http.Handler {
	// Create main router
	mux := http.NewServeMux()

	// Create usecases
	userUsecase := usecase.NewUserUseCase(userRepo, roleRepo, jwtService)
	twoFAUsecase := usecase.NewTwoFAUsecase(userRepo, twoFactorRepo, jwtService, mailService)

	// Create handlers
	authHandler := handlers.NewAuthHandler(userUsecase, jwtService, auditLogUsecase, twoFAUsecase)

	// Set up middleware
	errorMiddleware := middleware.ErrorHandler
	corsMiddleware := middleware.CORSMiddleware
	requestLoggerMiddleware := middleware.RequestLoggerMiddleware(appLogger)
	authMiddleware := middleware.NewAuthMiddleware(jwtService)
	languageMiddleware := middleware.LanguageMiddleware

	// Create audit logger builder
	auditLogger := middleware.NewAuditLogger(auditLogUsecase)

	// Health check endpoint
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	// API v1 routes - Authentication
	// Using the new Build pattern for middleware chaining

	// Login - Handler with audit log response middleware
	mux.HandleFunc("POST /api/v1/auth/login", middleware.Build().With(
		auditLogger.WithType(models.AuditLogTypeLogin).AsResponseMiddleware(),
	).Handle(authHandler.Login))

	// Logout - Auth middleware then handler and finally audit log
	mux.HandleFunc("POST /api/v1/auth/logout", middleware.Build(
		authMiddleware,
		auditLogger.WithType(models.AuditLogTypeLogout).AsResponseMiddleware(),
	).Handle(authHandler.Logout))

	mux.HandleFunc("GET /api/v1/auth/me", middleware.Build(
		authMiddleware,
	).Handle(authHandler.Me))

	mux.HandleFunc("POST /api/v1/auth/verify", authHandler.VerifyMFA)
	mux.HandleFunc("POST /api/v1/auth/resend-code", authHandler.ResendCode)

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

	// Apply global middleware using the new Build pattern
	var handler http.Handler = mux
	handler = middleware.Build(
		corsMiddleware,
		middleware.PerformanceMonitor(appLogger),
		errorMiddleware,
		requestLoggerMiddleware,
		middleware.WithLogger(appLogger),
		languageMiddleware,
	).Handle(func(w http.ResponseWriter, r *http.Request) {
		mux.ServeHTTP(w, r)
	})

	return handler
}
