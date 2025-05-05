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
	"github.com/vnlab/makeshop-payment/src/infrastructure/email"
	"github.com/vnlab/makeshop-payment/src/infrastructure/logger"
	"github.com/vnlab/makeshop-payment/src/usecase"
)

// SetupRouter sets up the HTTP router with all routes and middleware
func SetupRouter(
	userRepo repositories.UserRepository,
	roleRepo repositories.RoleRepository,
	permissionRepo repositories.PermissionRepository,
	payoutRepo repositories.PayoutRepository,
	payoutRecordRepo repositories.PayoutRecordRepository,
	merchantRepo repositories.MerchantRepository,
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
	roleUsecase := usecase.NewRoleUsecase(roleRepo, permissionRepo)
	permissionUsecase := usecase.NewPermissionUseCase(permissionRepo)
	twoFAUsecase := usecase.NewTwoFAUsecase(userRepo, twoFactorRepo, jwtService, mailService)
	payoutUsecase := usecase.NewPayoutUsecase(payoutRepo, payoutRecordRepo)
	merchantUsecase := usecase.NewMerchantUsecase(merchantRepo)

	// Create handlers
	authHandler := handlers.NewAuthHandler(userUsecase, jwtService, auditLogUsecase, twoFAUsecase)
	userHandler := handlers.NewUserHandler(userUsecase, jwtService)
	roleHandler := handlers.NewRoleHandler(roleUsecase)
	permissionHandler := handlers.NewPermissionHandler(permissionUsecase)
	payoutHandler := handlers.NewPayoutHandler(payoutUsecase)
	auditLogHandler := handlers.NewAuditLogHandler(auditLogUsecase)
	merchantHandler := handlers.NewMerchantHandler(merchantUsecase)

	// Set up middleware
	errorMiddleware := middleware.ErrorHandler
	corsMiddleware := middleware.CORSMiddleware
	requestLoggerMiddleware := middleware.RequestLoggerMiddleware(appLogger)
	authMiddleware := middleware.NewAuthMiddleware(jwtService)
	adminMiddleware := middleware.NewRoleMiddleware(string(models.RoleCodeAdmin))
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

	mux.HandleFunc("POST /api/v1/auth/register", authHandler.Register)

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

	// API v1 routes - User Self-Management
	mux.HandleFunc("GET /api/v1/users/profile", middleware.Build(
		authMiddleware,
	).Handle(userHandler.GetProfile))

	mux.HandleFunc("PUT /api/v1/users/profile", middleware.Build(
		authMiddleware,
	).Handle(userHandler.UpdateProfile))

	mux.HandleFunc("POST /api/v1/users/change-password", middleware.Build(
		authMiddleware,
		auditLogger.WithType(models.AuditLogTypePasswordChange).AsResponseMiddleware(),
	).Handle(userHandler.ChangePassword))

	// Admin-only routes - User Management
	mux.HandleFunc("GET /api/v1/admin/users", middleware.Build(
		authMiddleware,
		adminMiddleware,
	).Handle(userHandler.ListUsers))

	mux.HandleFunc("GET /api/v1/admin/users/{id}", middleware.Build(
		authMiddleware,
		adminMiddleware,
	).Handle(userHandler.GetUserByID))

	mux.HandleFunc("PUT /api/v1/admin/users/{id}", middleware.Build(
		authMiddleware,
		adminMiddleware,
		auditLogger.WithType(models.AuditLogTypeUserUpdate).AsResponseMiddleware(),
	).Handle(userHandler.UpdateUserProfile))

	mux.HandleFunc("POST /api/v1/admin/users", middleware.Build(
		authMiddleware,
		adminMiddleware,
		auditLogger.WithType(models.AuditLogTypeUserCreate).AsResponseMiddleware(),
	).Handle(userHandler.CreateUser))

	mux.HandleFunc("POST /api/v1/admin/users/reset-password", middleware.Build(
		authMiddleware,
		adminMiddleware,
		auditLogger.WithType(models.AuditLogTypePasswordReset).AsResponseMiddleware(),
	).Handle(userHandler.ResetPasswordByAdmin))

	mux.HandleFunc("DELETE /api/v1/admin/users/{id}", middleware.Build(
		authMiddleware,
		adminMiddleware,
		auditLogger.WithType(models.AuditLogTypeUserDelete).AsResponseMiddleware(),
	).Handle(userHandler.DeleteUser))

	// Role CRUD routes - Admin only
	mux.HandleFunc("GET /api/v1/admin/roles", middleware.Build(
		authMiddleware,
		adminMiddleware,
	).Handle(roleHandler.ListRoles))

	mux.HandleFunc("POST /api/v1/admin/roles", middleware.Build(
		authMiddleware,
		adminMiddleware,
	).Handle(roleHandler.CreateRole))

	mux.HandleFunc("GET /api/v1/admin/roles/{id}", middleware.Build(
		authMiddleware,
		adminMiddleware,
	).Handle(roleHandler.GetRoleByID))

	mux.HandleFunc("PUT /api/v1/admin/roles/{id}", middleware.Build(
		authMiddleware,
		adminMiddleware,
	).Handle(roleHandler.UpdateRole))

	mux.HandleFunc("DELETE /api/v1/admin/roles/{id}", middleware.Build(
		authMiddleware,
		adminMiddleware,
	).Handle(roleHandler.DeleteRole))

	mux.HandleFunc("POST /api/v1/admin/roles/permissions/batch", middleware.Build(
		authMiddleware,
		adminMiddleware,
	).Handle(roleHandler.BatchUpdateRolePermissions))

	// Permission CRUD routes - Admin only
	mux.HandleFunc("GET /api/v1/admin/permissions", middleware.Build(
		authMiddleware,
		adminMiddleware,
	).Handle(permissionHandler.ListPermission))

	// Payout routes - Admin only
	mux.HandleFunc("GET /api/v1/admin/payouts", middleware.Build(
		authMiddleware,
		adminMiddleware,
	).Handle(payoutHandler.ListPayouts))

	// AuditLog routes - Admin only
	mux.HandleFunc("GET /api/v1/admin/audit-logs", middleware.Build(
		authMiddleware,
		adminMiddleware,
	).Handle(auditLogHandler.ListAuditLogs))

	mux.HandleFunc("GET /api/v1/admin/audit-logs/users", middleware.Build(
		authMiddleware,
		adminMiddleware,
	).Handle(auditLogHandler.GetAuditLogUsers))

	// Merchant routes - Admin only
	mux.HandleFunc("GET /api/v1/admin/merchants", authMiddleware(adminMiddleware(merchantHandler.ListMerchants)))

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
