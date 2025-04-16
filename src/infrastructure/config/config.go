package config

import (
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
)

// Config holds all the configuration for the application
type Config struct {
	// Environment (development, staging, production)
	ApiEnv string

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
	LogLevel         string // Log level for application logs
	LogDirectory     string // Directory where log files will be stored
	EnableConsoleLog bool   // Whether to also log to console
	EnableSQLLog     bool   // Whether to log SQL queries
	SqlLogLevel      string // Log level for SQL queries

	// Authentication configuration
	JWTSecret       string
	JWTDurationHour int

	// Two Factor Authentication configuration
	MFATokenExpiryMinutes int

	// Email configuration
	SMTPHost         string
	SMTPPort         int
	SMTPUsername     string
	SMTPPassword     string
	SMTPFromEmail    string
	SMTPFromName     string
	SMTPUseAuth      bool
	SMTPUseTLS       bool
	EmailTemplateDir string
}

var (
	configInstance *Config
	once           sync.Once
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
			ServerHost:            "0.0.0.0",
			ServerPort:            "8080",
			LogLevel:              "warn",
			LogDirectory:          "/app/logs",
			EnableConsoleLog:      true,
			EnableSQLLog:          false,
			SqlLogLevel:           sqlLogLevel,
			JWTDurationHour:       24, // Hours
			MFATokenExpiryMinutes: 30, // Minutes
			EmailTemplateDir:      "src/infrastructure/email/templates",
			SMTPFromName:          "Makeshop Payment",
			SMTPUseAuth:           true,
			SMTPUseTLS:            true,
		}

		// Map of environment variables to configuration fields
		envVars := map[string]*string{
			"API_ENV":            &configInstance.ApiEnv,
			"SERVER_HOST":        &configInstance.ServerHost,
			"SERVER_PORT":        &configInstance.ServerPort,
			"FRONT_URL":          &configInstance.FrontUrl,
			"DB_HOST":            &configInstance.DBHost,
			"DB_PORT":            &configInstance.DBPort,
			"DB_USER":            &configInstance.DBUser,
			"DB_PASSWORD":        &configInstance.DBPassword,
			"DB_NAME":            &configInstance.DBName,
			"LOG_LEVEL":          &configInstance.LogLevel,
			"SQL_LOG_LEVEL":      &configInstance.SqlLogLevel,
			"LOG_DIRECTORY":      &configInstance.LogDirectory,
			"JWT_SECRET":         &configInstance.JWTSecret,
			"SMTP_HOST":          &configInstance.SMTPHost,
			"SMTP_USERNAME":      &configInstance.SMTPUsername,
			"SMTP_PASSWORD":      &configInstance.SMTPPassword,
			"SMTP_FROM_EMAIL":    &configInstance.SMTPFromEmail,
			"SMTP_FROM_NAME":     &configInstance.SMTPFromName,
			"EMAIL_TEMPLATE_DIR": &configInstance.EmailTemplateDir,
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
			"ENABLE_SQL_LOG":     &configInstance.EnableSQLLog,
			"SMTP_USE_AUTH":      &configInstance.SMTPUseAuth,
			"SMTP_USE_TLS":       &configInstance.SMTPUseTLS,
		}

		for env, field := range boolVars {
			if val := os.Getenv(env); val != "" {
				parsedVal, err := strconv.ParseBool(val)
				if err == nil {
					*field = parsedVal
				}
			}
		}

		// Override boolean fields
		intVars := map[string]*int{
			"JWT_EXPIRATION_HOURS": &configInstance.JWTDurationHour,
			"SMTP_PORT":            &configInstance.SMTPPort,
		}

		for env, field := range intVars {
			if val := os.Getenv(env); val != "" {
				intVal, err := strconv.Atoi(val)
				if err == nil {
					*field = intVal
				}
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
