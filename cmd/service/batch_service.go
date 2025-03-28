package service

import (
	"github.com/huydq/demo/src/infrastructure/logger"
	"github.com/huydq/demo/src/infrastructure/persistence/mysql"
	"gorm.io/gorm"
)

type BatchService struct {
	DB          *gorm.DB
}

// NewBatchService init
func NewBatchService() (*BatchService, error) {
	appLogger := logger.GetLogger()
	// Connect to database
	db, err := mysql.NewConnection(appLogger)
	if err != nil {
		return nil, err
	}
	return &BatchService{
		DB:          db,
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
