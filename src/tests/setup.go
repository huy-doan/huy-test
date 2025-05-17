package tests

import (
	"fmt"

	"github.com/huydq/test/src/infrastructure/auth"
	"github.com/huydq/test/src/infrastructure/config"
	"github.com/huydq/test/src/infrastructure/logger"
	"github.com/huydq/test/src/infrastructure/persistence/mysql"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type TestSuite struct {
	suite.Suite
	JWTService *auth.JWTService
	DB         *gorm.DB
}

// SetupSuite initializes the environment for the test suite
func (s *TestSuite) SetupSuite() {
	// Initialize JWTService
	s.JWTService = auth.NewJWTService()
	if s.JWTService == nil {
		panic("failed to initialize JWTService")
	}

	// Initialize test database connection using pure Go SQLite driver
	appLogger := logger.GetLogger()
	db, err := mysql.NewConnection(appLogger)
	if err != nil {
		panic("failed to connect to test database: " + err.Error())
	}
	s.DB = db

	// Clean up any existing tables first
	s.TruncateDatabase()

	sqlDb, err := db.DB()
	if err != nil {
		panic("failed to get database instance: " + err.Error())
	}

	// Set dialect for migrations
	if err := goose.SetDialect("mysql"); err != nil {
		panic("failed to set goose dialect: " + err.Error())
	}

	// Migrate the schema for the test database
	if err := goose.Up(sqlDb, "/app/config/migrations"); err != nil {
		panic("failed to migrate test database: " + err.Error())
	}

	// Apply seeds with no versioning
	if err := goose.Up(sqlDb, "/app/config/seeds/master", goose.WithNoVersioning()); err != nil {
		panic("failed to seed test database: " + err.Error())
	}
}

// TearDownSuite cleans up resources after all tests in the suite have run
func (s *TestSuite) TearDownSuite() {
	if s.DB != nil {
		// Truncate the database to clean it up for other test suites
		s.TruncateDatabase()

		// Close the database connection
		sqlDb, err := s.DB.DB()
		if err != nil {
			return // Just return if we can't get the DB instance
		}
		_ = sqlDb.Close() // Best effort to close the connection
	}
}

// TruncateDatabase drops all tables in the database to ensure a clean state
func (s *TestSuite) TruncateDatabase() {
	if s.DB == nil {
		return
	}

	appConfig := config.GetConfig()
	if appConfig == nil {
		panic("failed to get app config")
	}

	testDbName := appConfig.DBName

	// Disable foreign key checks to allow cascading deletes
	if err := s.DB.Exec("SET FOREIGN_KEY_CHECKS = 0;").Error; err != nil {
		panic("failed to disable foreign keys: " + err.Error())
	}

	// Get a list of all tables
	var tables []string
	if err := s.DB.Raw(
		fmt.Sprintf(
			"SELECT TABLE_NAME FROM information_schema.tables WHERE table_schema = '%s' AND TABLE_NAME != 'goose_db_version';",
			testDbName,
		),
	).Scan(&tables).Error; err != nil {
		panic("failed to get tables: " + err.Error())
	}

	// Truncate each table
	for _, table := range tables {
		if err := s.DB.Exec(fmt.Sprintf("TRUNCATE TABLE `%s`;", table)).Error; err != nil {
			panic(fmt.Sprintf("failed to truncate table %s: %s", table, err.Error()))
		}
	}

	// Re-enable foreign key checks
	if err := s.DB.Exec("SET FOREIGN_KEY_CHECKS = 1;").Error; err != nil {
		panic("failed to enable foreign keys: " + err.Error())
	}
}
