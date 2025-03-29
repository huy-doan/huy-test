package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	_ "github.com/vnlab/makeshop-payment/docs"
	"github.com/vnlab/makeshop-payment/src/api"
	"github.com/vnlab/makeshop-payment/src/infrastructure/config"
	"github.com/vnlab/makeshop-payment/src/infrastructure/logger"
	"github.com/vnlab/makeshop-payment/src/infrastructure/persistence/mysql"
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
		LogLevel:        appConfig.LogLevel,
		LogDirectory:    appConfig.LogDirectory,
		EnableConsoleLog: appConfig.EnableConsoleLog,
		EnableSQLLog:    appConfig.EnableSQLLog,
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

	// Create and start API server
	server := api.NewServer(userRepo, roleRepo, appLogger)
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
