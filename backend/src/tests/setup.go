package tests

import (
	"sync"

	"github.com/huydq/test/src/infrastructure/auth"
	"github.com/huydq/test/src/infrastructure/logger"
	"github.com/huydq/test/src/infrastructure/persistence/mysql"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

var (
	dbMutex sync.Mutex
)

type TestSuite struct {
	suite.Suite
	JWTService *auth.JWTService
	DB         *gorm.DB
}

var TestDB *gorm.DB
var dbInitOnce sync.Once

func GetTestDB() *gorm.DB {
	dbInitOnce.Do(func() {
		dbMutex.Lock()
		defer dbMutex.Unlock()

		if TestDB == nil {
			appLogger := logger.GetLogger()
			db, err := mysql.NewConnection(appLogger)
			if err != nil {
				panic("failed to connect to test database: " + err.Error())
			}
			TestDB = db
		}
	})
	return TestDB
}

func (s *TestSuite) SetupSuite() {
	s.JWTService = auth.NewJWTService()
	if s.JWTService == nil {
		panic("failed to initialize JWTService")
	}

	// Connect to the test database which should already have migrations and seeds applied
	s.DB = GetTestDB()
}

// TearDownSuite is kept minimal - no truncating
func (s *TestSuite) TearDownSuite() {
	// Do nothing - no longer truncating tables to preserve seed data
}
