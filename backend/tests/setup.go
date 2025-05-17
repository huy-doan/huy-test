package tests

import (
	"context"
	"log"
	"os"
	"sync"

	"github.com/huydq/test/internal/pkg/database"
	"github.com/huydq/test/internal/pkg/dbconn"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

var (
	dbMutex sync.Mutex
)

type TestSuite struct {
	suite.Suite
	Context context.Context
	DB      *gorm.DB
}

var TestDB *gorm.DB
var ctx context.Context
var dbInitOnce sync.Once

func GetTestDB() (*gorm.DB, context.Context) {
	dbInitOnce.Do(func() {
		dbMutex.Lock()
		defer dbMutex.Unlock()

		if TestDB == nil {
			dbHost := os.Getenv("DB_HOST")
			dbPort := os.Getenv("DB_PORT")
			dbUser := os.Getenv("DB_USER")
			dbPassword := os.Getenv("DB_PASSWORD")
			dbName := os.Getenv("DB_NAME")
			if dbHost == "" || dbPort == "" || dbUser == "" || dbPassword == "" || dbName == "" {
				log.Fatal("Database environment variables are not set")
			}

			db, err := dbconn.NewConnectionWithConfig(dbHost, dbPort, dbUser, dbPassword, dbName)
			if err != nil {
				log.Fatalf("Failed to connect to database: %v", err)
			}
			ctx, err = database.SetDB(context.Background(), db)
			if err != nil {
				panic("failed to connect to test database: " + err.Error())
			}
			TestDB = db
		}
	})
	return TestDB, ctx
}

func (s *TestSuite) SetupSuite() {
	s.DB, s.Context = GetTestDB()
}

func (s *TestSuite) TearDownSuite() {
}
