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
	LogLevel         string
	LogDirectory     string
	EnableConsoleLog bool
	EnableSQLLog     bool
	SqlLogLevel      string

	// Authentication configuration
	JWTSecret       string
	JWTDurationHour int

	// Two Factor Authentication configuration
	MFATokenExpiryMinutes  int
	MFATokenResendInterval int

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

	// S3 configuration
	S3Bucket           string
	S3Region           string
	AwsAccessKeyID     string
	AwsSecretAccessKey string

	// SSH configuration
	SSHUser                                 string
	SSHHost                                 string
	SSHPort                                 string
	SSHPassword                             string
	RemoteDir                               string
	TransactionDetailsNoShippingRelatedPath string
	TransactionDetailsSummaryPath           string
	TransactionDetailsShippingRelatedPath   string
	TopUpDetailsPath                        string
	TopUpSummaryDetailsPath                 string
	TopUpReportPath                         string
	ValidInvoicesPath                       string
	ValidInvoicesDuplicatePath              string
	ValidInvoicesSpreadsheetsPath           string
	LocalDir                                string

	ProviderID int
}

var (
	configInstance *Config
	once           sync.Once
)

// LoadConfig loads the configuration from environment variables
func LoadConfig() *Config {
	once.Do(func() {
		_ = godotenv.Load()

		sqlLogLevel := "warn"
		if os.Getenv("API_ENV") != "production" {
			sqlLogLevel = "debug"
		}

		configInstance = &Config{
			ServerHost:             "0.0.0.0",
			ServerPort:             "8080",
			LogLevel:               "warn",
			LogDirectory:           "/app/logs",
			EnableConsoleLog:       true,
			EnableSQLLog:           false,
			SqlLogLevel:            sqlLogLevel,
			JWTDurationHour:        24,
			MFATokenExpiryMinutes:  30,
			MFATokenResendInterval: 1,
			EmailTemplateDir:       "internal/resource/templates/email",
			SMTPFromName:           "Makeshop Payment",
			SMTPUseAuth:            true,
			SMTPUseTLS:             true,
			ProviderID:             1,
		}

		envVars := map[string]*string{
			"API_ENV":               &configInstance.ApiEnv,
			"SERVER_HOST":           &configInstance.ServerHost,
			"SERVER_PORT":           &configInstance.ServerPort,
			"FRONT_URL":             &configInstance.FrontUrl,
			"DB_HOST":               &configInstance.DBHost,
			"DB_PORT":               &configInstance.DBPort,
			"DB_USER":               &configInstance.DBUser,
			"DB_PASSWORD":           &configInstance.DBPassword,
			"DB_NAME":               &configInstance.DBName,
			"LOG_LEVEL":             &configInstance.LogLevel,
			"SQL_LOG_LEVEL":         &configInstance.SqlLogLevel,
			"LOG_DIRECTORY":         &configInstance.LogDirectory,
			"JWT_SECRET":            &configInstance.JWTSecret,
			"SMTP_HOST":             &configInstance.SMTPHost,
			"SMTP_USERNAME":         &configInstance.SMTPUsername,
			"SMTP_PASSWORD":         &configInstance.SMTPPassword,
			"SMTP_FROM_EMAIL":       &configInstance.SMTPFromEmail,
			"SMTP_FROM_NAME":        &configInstance.SMTPFromName,
			"EMAIL_TEMPLATE_DIR":    &configInstance.EmailTemplateDir,
			"S3_BUCKET":             &configInstance.S3Bucket,
			"S3_REGION":             &configInstance.S3Region,
			"AWS_ACCESS_KEY_ID":     &configInstance.AwsAccessKeyID,
			"AWS_SECRET_ACCESS_KEY": &configInstance.AwsSecretAccessKey,
			"SSH_USER":              &configInstance.SSHUser,
			"SSH_HOST":              &configInstance.SSHHost,
			"SSH_PORT":              &configInstance.SSHPort,
			"SSH_PASSWORD":          &configInstance.SSHPassword,
			"REMOTE_DIR":            &configInstance.RemoteDir,
			"LOCAL_DIR":             &configInstance.LocalDir,
			"TRANSACTION_DETAILS_NO_SHIPPING_RELATED_PATH": &configInstance.TransactionDetailsNoShippingRelatedPath,
			"TRANSACTION_DETAILS_SUMMARY_PATH":             &configInstance.TransactionDetailsSummaryPath,
			"TRANSACTION_DETAILS_SHIPPING_RELATED_PATH":    &configInstance.TransactionDetailsShippingRelatedPath,
			"TOP_UP_DETAILS_PATH":                          &configInstance.TopUpDetailsPath,
			"TOP_UP_SUMMARY_DETAILS_PATH":                  &configInstance.TopUpSummaryDetailsPath,
			"TOP_UP_REPORT_PATH":                           &configInstance.TopUpReportPath,
			"VALID_INVOICES_PATH":                          &configInstance.ValidInvoicesPath,
			"VALID_INVOICES_DUPLICATE_PATH":                &configInstance.ValidInvoicesDuplicatePath,
			"VALID_INVOICES_SPREADSHEETS_PATH":             &configInstance.ValidInvoicesSpreadsheetsPath,
		}

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
