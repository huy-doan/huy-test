package mysql

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	CONN_MAX_LIFETIME = time.Minute * 10
	MAX_IDLE_CONNS    = 500
	MAX_OPEN_CONNS    = 250
)

// NewConnection creates a new MySQL database connection using GORM
func NewConnection() (*gorm.DB, error) {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	if dbHost == "" {
		dbHost = "localhost"
	}
	if dbPort == "" {
		dbPort = "3306"
	}
	if dbUser == "" {
		dbUser = "apiuser"
	}
	if dbPassword == "" {
		dbPassword = "apipassword"
	}
	if dbName == "" {
		dbName = "msp-db-dev"
	}
	loc := url.QueryEscape("Asia/Tokyo")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=%s",
		dbUser, dbPassword, dbHost, dbPort, dbName, loc)

	// Configure GORM logger
	file, err := os.OpenFile("/app/logs/db-backend.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	logLevel := logger.Warn
	if os.Getenv("GIN_MODE") != gin.ReleaseMode {
		logLevel = logger.Info
	}

	newLogger := logger.New(
		log.New(file, "\r\n", log.LstdFlags), // Thay thế logger.Writer
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logLevel,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
			Colorful:                  false,
		},
	)

	// Open database connection
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

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
