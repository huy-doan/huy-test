package main

import (
	"log"
	"os"
	"path/filepath"

	_ "github.com/huydq/test/docs"
	"github.com/huydq/test/src/api"
	"github.com/huydq/test/src/api/http/middleware"
	"github.com/huydq/test/src/infrastructure/config"
	"github.com/huydq/test/src/infrastructure/logger"
	"github.com/huydq/test/src/infrastructure/persistence/mysql"
	"github.com/huydq/test/src/infrastructure/persistence/repositories"
	"github.com/joho/godotenv"
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
	// Initialize i18n
	if err := middleware.InitI18n(); err != nil {
		log.Fatalf("Failed to initialize i18n: %v", err)
	}

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

	db, err := mysql.NewConnection(appLogger)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get SQL DB: %v", err)
	}
	defer sqlDB.Close()

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	roleRepo := repositories.NewRoleRepository(db)
	auditLogRepo := repositories.NewAuditLogRepository(db)
	auditLogTypeRepo := repositories.NewAuditLogTypeRepository(db)
	twoFactorTokenRepo := repositories.NewTwoFactorTokenRepository(db)

	// Create and start API server
	server := api.NewServer(
		userRepo,
		roleRepo,
		auditLogRepo,
		auditLogTypeRepo,
		appLogger,
		twoFactorTokenRepo,
	)
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
