package dbconn

import (
	"fmt"
	"net/url"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/huydq/test/internal/pkg/config"
	"github.com/huydq/test/internal/pkg/logger"
)

const (
	CONN_MAX_LIFETIME = time.Minute * 10
	MAX_IDLE_CONNS    = 500
	MAX_OPEN_CONNS    = 250
)

func NewConnectionWithConfig(dbHost, dbPort, dbUser, dbPassword, dbName string) (*gorm.DB, error) {
	// Configure connection string with Tokyo timezone
	loc := url.QueryEscape("Asia/Tokyo")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=%s",
		dbUser, dbPassword, dbHost, dbPort, dbName, loc)

	// Open database connection
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get SQL DB: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(MAX_IDLE_CONNS)
	sqlDB.SetMaxOpenConns(MAX_OPEN_CONNS)
	sqlDB.SetConnMaxLifetime(CONN_MAX_LIFETIME)

	return db, nil
}

// NewConnection creates a new MySQL database connection using GORM
func NewConnection(appLogger logger.Logger) (*gorm.DB, error) {
	appConfig := config.GetConfig()
	dbHost := appConfig.DBHost
	dbPort := appConfig.DBPort
	dbUser := appConfig.DBUser
	dbPassword := appConfig.DBPassword
	dbName := appConfig.DBName

	db, err := NewConnectionWithConfig(dbHost, dbPort, dbUser, dbPassword, dbName)

	if err != nil {
		appLogger.Error("Failed to connect to database", map[string]any{
			"error": err.Error(),
			"host":  dbHost,
			"port":  dbPort,
			"user":  dbUser,
			"name":  dbName,
		})
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		appLogger.Error("Failed to get SQL DB handle", map[string]any{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to get SQL DB: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(MAX_IDLE_CONNS)
	sqlDB.SetMaxOpenConns(MAX_OPEN_CONNS)
	sqlDB.SetConnMaxLifetime(CONN_MAX_LIFETIME)

	return db, nil
}
