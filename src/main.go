package main

import (
	"log"
	"os"
	"path/filepath"

	_ "github.com/huydq/demo/docs"
	"github.com/huydq/demo/src/api"
	"github.com/huydq/demo/src/infrastructure/persistence/mysql"
	"github.com/huydq/demo/src/infrastructure/persistence/repositories"
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
// @host      localhost:3011
// @BasePath  /api/v1
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Connect to database
	db, err := mysql.NewConnection()
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
	server := api.NewServer(userRepo, roleRepo)
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
