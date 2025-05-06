package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	_ "github.com/vnlab/makeshop-payment/docs"
	"github.com/vnlab/makeshop-payment/internal/pkg/config"
	"github.com/vnlab/makeshop-payment/internal/pkg/dbconn"
	"github.com/vnlab/makeshop-payment/internal/pkg/logger"
	"github.com/vnlab/makeshop-payment/internal/server/http"
	"github.com/vnlab/makeshop-payment/src/infrastructure/persistence/repositories"
)

func init() {
	// Load environment variables from .env file if it exists
	envPath := filepath.Join(".", ".env")
	if _, err := os.Stat(envPath); err == nil {
		if err := godotenv.Load(envPath); err != nil {
			log.Printf("Warning: Error loading .env file: %v", err)
		}
	}
}

// @title           Makeshop Payment API
// @version         1.0
// @description     Payment API for Makeshop
// @termsOfService  http://swagger.io/terms/
//
// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io
//
// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT
//
// @BasePath  /api/v1
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Connect to database
	appConfig := config.LoadConfig()

	// Initialize logger as a singleton
	logger.InitLogger(&logger.Config{
		LogLevel:         appConfig.LogLevel,
		LogDirectory:     appConfig.LogDirectory,
		EnableConsoleLog: appConfig.EnableConsoleLog,
		EnableSQLLog:     appConfig.EnableSQLLog,
	})

	// Get the global logger instance
	appLogger := logger.GetLogger()

	db, err := dbconn.NewConnection(appLogger)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get SQL DB: %v", err)
	}
	defer func() {
		if err := sqlDB.Close(); err != nil {
			appLogger.Error("Failed to close database connection", map[string]any{
				"error": err.Error(),
			})
		}
	}()

	// TODO: use wire and echo

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	roleRepo := repositories.NewRoleRepository(db)
	auditLogRepo := repositories.NewAuditLogRepository(db)
	auditLogTypeRepo := repositories.NewAuditLogTypeRepository(db)
	twoFactorTokenRepo := repositories.NewTwoFactorTokenRepository(db)

	// Create and start API server
	apiServer := http.NewServer(
		userRepo,
		roleRepo,
		auditLogRepo,
		auditLogTypeRepo,
		appLogger,
		twoFactorTokenRepo,
	)

	if err := apiServer.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
