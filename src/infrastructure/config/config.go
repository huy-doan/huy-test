package config

import (
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
)

// Config holds all the configuration for the application
type Config struct {
	GinMode string // Gin mode for the server

	// Client configuration
	FrontUrl string

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
	JWTDurationHour int
}

var (
	configInstance *Config
	once     sync.Once
)

// LoadConfig loads the configuration from environment variables
func LoadConfig() *Config {
	once.Do(func() {
		// Load .env file if it exists
		godotenv.Load()
		
		sqlLogLevel := "warn"
		if os.Getenv("API_ENV") != "production" {
			sqlLogLevel = "debug"
		}		

		// Set default values
		configInstance = &Config{
			ServerHost:       "0.0.0.0",
			ServerPort:       "8080",
			LogLevel: 	      "warn",
			LogDirectory:     "/app/logs",
			EnableConsoleLog: true,
			EnableSQLLog:     false,
			SqlLogLevel:	  sqlLogLevel,
			JWTDurationHour:      24, // Hours
		}
		
		// Map of environment variables to configuration fields
		envVars := map[string]*string{
			"SERVER_HOST":    &configInstance.ServerHost,
			"SERVER_PORT":    &configInstance.ServerPort,
			"FRONT_URL":      &configInstance.FrontUrl,
			"GIN_MODE":       &configInstance.GinMode,
			"DB_HOST":        &configInstance.DBHost,
			"DB_PORT":        &configInstance.DBPort,
			"DB_USER":        &configInstance.DBUser,
			"DB_PASSWORD":    &configInstance.DBPassword,
			"DB_NAME":        &configInstance.DBName,
			"LOG_LEVEL":      &configInstance.LogLevel,
			"SQL_LOG_LEVEL":  &configInstance.SqlLogLevel,
			"LOG_DIRECTORY":  &configInstance.LogDirectory,
			"JWT_SECRET":     &configInstance.JWTSecret,
		}

		// Override string fields with environment variables if they exist
		for env, field := range envVars {
			if val := os.Getenv(env); val != "" {
				*field = val
			}
		}

		// Override boolean fields
		boolVars := map[string]*bool{
			"ENABLE_CONSOLE_LOG": &configInstance.EnableConsoleLog,
			"ENABLE_SQL_LOG": &configInstance.EnableSQLLog,
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
		if val := os.Getenv("JWT_EXPIRATION_HOURS"); val != "" {
			if duration, err := strconv.Atoi(val); err == nil {
				configInstance.JWTDurationHour = duration
			}
		}
	})

	return configInstance
}

// GetConfig returns the singleton configInstance of the application configuration
func GetConfig() *Config {
	if configInstance == nil {
		return LoadConfig()
	}
	return configInstance
}
