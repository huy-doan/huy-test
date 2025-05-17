package container

import (
	"fmt"
	"os"
	"path/filepath"

	appConfig "github.com/huydq/test/internal/pkg/config"
	"github.com/huydq/test/internal/pkg/dbconn"
	"github.com/huydq/test/internal/pkg/logger"

	"gorm.io/gorm"
)

type LogBatchConfig struct {
	LogLevel         string
	LogDirectory     string
	EnableConsoleLog bool
}

type BatchService struct {
	Logger    logger.Logger
	DB        *gorm.DB
	AppConfig *appConfig.Config
}

func DefaultLogConfig(appConfig *appConfig.Config) *LogBatchConfig {
	batchLogDir := filepath.Join(appConfig.LogDirectory, "batch")
	os.MkdirAll(batchLogDir, 0755)

	return &LogBatchConfig{
		LogLevel:         appConfig.LogLevel,
		LogDirectory:     batchLogDir,
		EnableConsoleLog: true,
	}
}

func NewBatchContainer() (*BatchService, error) {
	return NewBatchContainerWithConfig()
}

func NewBatchContainerWithConfig() (*BatchService, error) {
	appConfig := appConfig.GetConfig()
	logConfig := DefaultLogConfig(appConfig)
	cliLogger := logger.InitCLILogger(&logger.CLILoggerConfig{
		LogLevel:         logConfig.LogLevel,
		LogDirectory:     logConfig.LogDirectory,
		EnableConsoleLog: logConfig.EnableConsoleLog,
	})

	db, err := dbconn.NewConnection(cliLogger)
	if err != nil {
		cliLogger.Error("Failed to connect to database", map[string]any{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("database connection failed: %w", err)
	}

	return &BatchService{
		DB:        db,
		Logger:    cliLogger,
		AppConfig: appConfig,
	}, nil
}

// Close releases resources when finished
func (s *BatchService) Close() error {
	sqlDB, err := s.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
