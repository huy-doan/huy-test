package usecase_test

import (
	"context"
	"testing"

	"github.com/huydq/test/src/domain/models"
	"github.com/huydq/test/src/domain/repositories"
	repoImpl "github.com/huydq/test/src/infrastructure/persistence/repositories"
	"github.com/huydq/test/src/tests"
	"github.com/huydq/test/src/usecase"
	"github.com/stretchr/testify/suite"
)

type PermissionUsecaseTestSuite struct {
	tests.TestSuite
	permissionUsecase *usecase.PermissionUsecase
	permissionRepo    repositories.PermissionRepository
	ctx               context.Context
	adminPermission   *models.Permission // Expected to be seeded as ID 1
}

func TestPermissionUsecaseSuite(t *testing.T) {
	suite.Run(t, new(PermissionUsecaseTestSuite))
}

func (s *PermissionUsecaseTestSuite) SetupSuite() {
	// Call the parent SetupSuite method
	s.TestSuite.SetupSuite()
}

func (s *PermissionUsecaseTestSuite) SetupTest() {
	// Initialize context
	s.ctx = context.Background()

	// Initialize repository using the test DB
	s.permissionRepo = repoImpl.NewPermissionRepository(s.DB)

	// Create permissionUsecase with real repository
	s.permissionUsecase = usecase.NewPermissionUseCase(s.permissionRepo)

	// Load permissions from database (from seeds)
	permissions, err := s.permissionRepo.List(s.ctx)
	s.Require().NoError(err)
	s.Require().NotEmpty(permissions, "No permissions found in the test database")

	// Find admin permission from seeded data - now we're looking for USER_MANAGE permission
	for _, p := range permissions {
		if p.Code == string(models.PermissionCodeUserManage) {
			s.adminPermission = p
			break
		}
	}
	s.Require().NotNil(s.adminPermission, "Admin permission (USER_MANAGE) not found in database")
}

func (s *PermissionUsecaseTestSuite) TearDownTest() {
	// No specific cleanup needed since we're only using existing seeded data
}

func (s *PermissionUsecaseTestSuite) TearDownSuite() {
}

func (s *PermissionUsecaseTestSuite) TestListPermission() {
	// Test listing all permissions
	permissions, err := s.permissionUsecase.ListPermission(s.ctx)

	// Assertions
	s.Assert().NoError(err)
	s.Assert().NotNil(permissions)
	s.Assert().NotEmpty(permissions)
	s.Assert().GreaterOrEqual(len(permissions), 1) // At least one permission should be seeded

	// Verify admin permission (USER_MANAGE) is included
	foundAdminPerm := false
	for _, perm := range permissions {
		if perm.Code == string(models.PermissionCodeUserManage) {
			foundAdminPerm = true
			// Verify the permission properties match what we expect
			s.Assert().Equal(s.adminPermission.ID, perm.ID)
			s.Assert().Equal(s.adminPermission.Name, perm.Name)
			s.Assert().Equal(s.adminPermission.Code, perm.Code)
			break
		}
	}
	s.Assert().True(foundAdminPerm, "Admin permission (USER_MANAGE) should be in the permissions list")
}
