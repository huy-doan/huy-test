package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// Config holds all the configuration for the application
type Config struct {
	GinMode string // Gin mode for the server

	// Server configuration
	ServerHost string
	ServerPort string
	
	// Database configuration
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	

	// Logger configuration
	LogLevel      string // Log level for application logs
	LogDirectory  string // Directory where log files will be stored
	EnableConsoleLog bool   // Whether to also log to console
	EnableSQLLog  bool   // Whether to log SQL queries
	SqlLogLevel   string // Log level for SQL queries

	// Authentication configuration
	JWTSecret   string
	JWTDuration int
}

// LoadConfig loads the configuration from environment variables
func LoadConfig() *Config {
	// Load .env file if it exists
	godotenv.Load()
	
	sqlLogLevel := "warn"
	if os.Getenv("GIN_MODE") != gin.ReleaseMode {
		sqlLogLevel = "info"
	}		

	// Set default values
	config := &Config{
		ServerHost:       "0.0.0.0",
		ServerPort:       "8080",
		LogLevel: 	      "warn",
		LogDirectory:     "/app/logs",
		EnableConsoleLog: true,
		EnableSQLLog:     false,
		SqlLogLevel:	  sqlLogLevel,
		JWTDuration:      24, // Hours
	}
	
	// Map of environment variables to configuration fields
	envVars := map[string]*string{
		"SERVER_HOST":    &config.ServerHost,
		"SERVER_PORT":    &config.ServerPort,
		"GIN_MODE":       &config.GinMode,
		"DB_HOST":        &config.DBHost,
		"DB_PORT":        &config.DBPort,
		"DB_USER":        &config.DBUser,
		"DB_PASSWORD":    &config.DBPassword,
		"DB_NAME":        &config.DBName,
		"LOG_LEVEL":      &config.LogLevel,
		"SQL_LOG_LEVEL":  &config.SqlLogLevel,
		"LOG_DIRECTORY":  &config.LogDirectory,
		"JWT_SECRET":     &config.JWTSecret,
	}

	// Override string fields with environment variables if they exist
	for env, field := range envVars {
		if val := os.Getenv(env); val != "" {
			*field = val
		}
	}

	// Override boolean fields
	boolVars := map[string]*bool{
		"ENABLE_CONSOLE_LOG": &config.EnableConsoleLog,
		"ENABLE_SQL_LOG": &config.EnableSQLLog,
	}
	for env, field := range boolVars {
		if val := os.Getenv(env); val != "" {
			parsedVal, err := strconv.ParseBool(val)
			if err == nil {
				*field = parsedVal
			}
		}
	}

	// Override integer fields
	if val := os.Getenv("JWT_DURATION"); val != "" {
		if duration, err := strconv.Atoi(val); err == nil {
			config.JWTDuration = duration
		}
	}
	fmt.Printf("config: %v\n", config)
	return config
}
