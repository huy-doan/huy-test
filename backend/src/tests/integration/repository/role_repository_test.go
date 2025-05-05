package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/vnlab/makeshop-payment/src/domain/models"
	"github.com/vnlab/makeshop-payment/src/domain/repositories"
	"github.com/vnlab/makeshop-payment/src/domain/repositories/filter"
	repoImpl "github.com/vnlab/makeshop-payment/src/infrastructure/persistence/repositories"
	"github.com/vnlab/makeshop-payment/src/tests"
	"github.com/vnlab/makeshop-payment/src/tests/fixtures"
)

// RoleRepositoryTestSuite defines the test suite for role repository
type RoleRepositoryTestSuite struct {
	tests.TestSuite
	repo            repositories.RoleRepository
	permissionRepo  repositories.PermissionRepository
	testPermissions []*models.Permission
}

// SetupSuite initializes the test suite
func (s *RoleRepositoryTestSuite) SetupSuite() {
	// Call the parent SetupSuite method
	s.TestSuite.SetupSuite()

	// Initialize repositories
	s.repo = repoImpl.NewRoleRepository(s.DB)
	s.permissionRepo = repoImpl.NewPermissionRepository(s.DB)
}

// SetupTest runs before each test
func (s *RoleRepositoryTestSuite) SetupTest() {
	// Fetch permissions to use in our tests
	permissions, err := s.permissionRepo.List(context.Background())
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), permissions, "Test database must have seeded permissions")

	// Store permissions for use in tests
	if len(permissions) >= 2 {
		s.testPermissions = permissions[:2]
	} else {
		s.testPermissions = permissions
	}
}

// TearDownTest runs after each test
func (s *RoleRepositoryTestSuite) TearDownTest() {
}

// TestCreateRole tests the creation of a role
func (s *RoleRepositoryTestSuite) TestCreateRole() {
	ctx := context.Background()

	// Create a unique test role with permissions using fixture helper
	role, err := fixtures.CreateUniqueTestRoleWithPermissions(
		ctx,
		s.repo,
		"Create",
		s.testPermissions,
	)

	// Assertions
	require.NoError(s.T(), err)
	assert.NotZero(s.T(), role.ID)
	assert.NotZero(s.T(), role.CreatedAt)

	// Verify role was created
	savedRole, err := s.repo.FindByID(context.Background(), role.ID)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), savedRole)
	assert.Equal(s.T(), role.Name, savedRole.Name)
	assert.Equal(s.T(), role.Code, savedRole.Code)
	assert.Len(s.T(), savedRole.Permissions, len(s.testPermissions))
}

// TestFindRoleByID tests finding a role by ID
func (s *RoleRepositoryTestSuite) TestFindRoleByID() {
	ctx := context.Background()

	// Create a unique test role for finding by ID
	role, err := fixtures.CreateUniqueTestRole(ctx, s.repo, "FindByID")
	require.NoError(s.T(), err)
	require.NotZero(s.T(), role.ID)

	// Test finding the role
	foundRole, err := s.repo.FindByID(context.Background(), role.ID)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), foundRole)
	assert.Equal(s.T(), role.ID, foundRole.ID)
	assert.Equal(s.T(), role.Name, foundRole.Name)
	assert.Equal(s.T(), role.Code, foundRole.Code)

	// Test with non-existent ID
	nonExistentRole, err := s.repo.FindByID(context.Background(), 999999)
	assert.NoError(s.T(), err)
	assert.Nil(s.T(), nonExistentRole)
}

// TestFindRoleByCode tests finding a role by code
func (s *RoleRepositoryTestSuite) TestFindRoleByCode() {
	ctx := context.Background()

	// Try to fetch an existing seeded role (Admin)
	adminRole, err := s.repo.FindByCode(ctx, string(models.RoleCodeAdmin))
	assert.NoError(s.T(), err)

	// If admin role exists in database, test with it, otherwise create a test role
	var role *models.Role
	if adminRole != nil {
		role = adminRole
	} else {
		role, err = fixtures.CreateUniqueTestRole(ctx, s.repo, "FindByCode")
		require.NoError(s.T(), err)
	}

	// Test finding the role by code
	foundRole, err := s.repo.FindByCode(context.Background(), role.Code)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), foundRole)
	assert.Equal(s.T(), role.ID, foundRole.ID)
	assert.Equal(s.T(), role.Name, foundRole.Name)
	assert.Equal(s.T(), role.Code, foundRole.Code)

	// Test with non-existent code
	nonExistentRole, err := s.repo.FindByCode(context.Background(), "NON_EXISTENT_CODE")
	assert.NoError(s.T(), err)
	assert.Nil(s.T(), nonExistentRole)
}

// TestUpdateRole tests updating a role
func (s *RoleRepositoryTestSuite) TestUpdateRole() {
	ctx := context.Background()

	// Create a role with one permission initially
	initialPermissions := s.testPermissions[:1]
	role, err := fixtures.CreateUniqueTestRoleWithPermissions(
		ctx,
		s.repo,
		"Update",
		initialPermissions,
	)
	require.NoError(s.T(), err)

	// Update the role
	updatedName := "Updated Role Name"
	role.Name = updatedName
	if role.Permissions != nil && len(s.testPermissions) >= 2 {
		role.Permissions = s.testPermissions // Update to include all test permissions
	}

	err = s.repo.Update(context.Background(), role)
	assert.NoError(s.T(), err)

	// Verify role was updated
	updatedRole, err := s.repo.FindByID(context.Background(), role.ID)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), updatedRole)
	assert.Equal(s.T(), updatedName, updatedRole.Name)
	assert.Equal(s.T(), role.Code, updatedRole.Code)

	// Check permissions were updated if we have multiple permissions
	if len(s.testPermissions) >= 2 {
		assert.Len(s.T(), updatedRole.Permissions, len(s.testPermissions))
	}
}

// TestDeleteRole tests deleting a role
func (s *RoleRepositoryTestSuite) TestDeleteRole() {
	ctx := context.Background()

	// Create a unique test role for deletion
	role, err := fixtures.CreateUniqueTestRole(ctx, s.repo, "Delete")
	require.NoError(s.T(), err)

	// Test deleting the role
	err = s.repo.Delete(context.Background(), role.ID)
	assert.NoError(s.T(), err)

	// Verify role was deleted (soft delete)
	deletedRole, err := s.repo.FindByID(context.Background(), role.ID)
	assert.NoError(s.T(), err)
	assert.Nil(s.T(), deletedRole)
}

// TestListRoles tests listing roles with pagination
func (s *RoleRepositoryTestSuite) TestListRoles() {
	ctx := context.Background()

	// Create multiple roles for testing
	for i := 1; i <= 5; i++ {
		role, err := fixtures.CreateUniqueTestRole(ctx, s.repo, "List")
		require.NoError(s.T(), err)
		require.NotNil(s.T(), role)
	}

	// Test cases for pagination and filtering
	s.T().Run("First page with 2 items", func(t *testing.T) {
		roleFilter := filter.NewRoleFilter()
		roleFilter.SetPagination(1, 2)

		roles, totalPages, total, err := s.repo.List(context.Background(), roleFilter)
		assert.NoError(t, err)
		assert.NotNil(t, roles)
		assert.LessOrEqual(t, len(roles), 2) // Should respect page size
		assert.GreaterOrEqual(t, totalPages, 1)
		assert.GreaterOrEqual(t, total, int64(1)) // At least one role should exist
	})

	s.T().Run("Second page with 2 items", func(t *testing.T) {
		roleFilter := filter.NewRoleFilter()
		roleFilter.SetPagination(2, 2)

		roles, totalPages, total, err := s.repo.List(context.Background(), roleFilter)
		assert.NoError(t, err)
		assert.NotNil(t, roles)
		assert.GreaterOrEqual(t, totalPages, 1)
		assert.GreaterOrEqual(t, total, int64(1)) // At least one role should exist
	})

	s.T().Run("Get all with large page size", func(t *testing.T) {
		roleFilter := filter.NewRoleFilter()
		roleFilter.SetPagination(1, 100)

		roles, totalPages, total, err := s.repo.List(context.Background(), roleFilter)
		assert.NoError(t, err)
		assert.NotNil(t, roles)
		assert.GreaterOrEqual(t, len(roles), 5) // We created 5 roles + any existing ones
		assert.GreaterOrEqual(t, totalPages, 1)
		assert.GreaterOrEqual(t, total, int64(5)) // We created at least 5 roles
	})

	s.T().Run("Filter by name", func(t *testing.T) {
		roleFilter := filter.NewRoleFilter()
		roleFilter.SetPagination(1, 100)
		listName := "List" // Should match our test roles
		roleFilter.Name = &listName
		roleFilter.ApplyFilters() // Important: Apply the filters to convert them to conditions

		roles, totalPages, total, err := s.repo.List(context.Background(), roleFilter)
		assert.NoError(t, err)
		assert.NotNil(t, roles)
		assert.GreaterOrEqual(t, len(roles), 1)
		assert.GreaterOrEqual(t, totalPages, 1)
		assert.GreaterOrEqual(t, total, int64(1)) // At least one role should match filter

		// Verify that filtered results all have "List" in the name
		for _, role := range roles {
			assert.Contains(t, role.Name, "List")
		}
	})
}

// TestRoleRepository executes the role repository test suite
func TestRoleRepository(t *testing.T) {
	suite.Run(t, new(RoleRepositoryTestSuite))
}
