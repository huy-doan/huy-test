package repository

import (
	"context"
	"log"
	"testing"

	"slices"

	"github.com/huydq/test/src/domain/models"
	"github.com/huydq/test/src/domain/repositories"
	repoImpl "github.com/huydq/test/src/infrastructure/persistence/repositories"
	"github.com/huydq/test/src/tests"
	"github.com/huydq/test/src/tests/fixtures"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// PermissionRepositoryTestSuite defines the test suite for permission repository
type PermissionRepositoryTestSuite struct {
	tests.TestSuite
	repo        repositories.PermissionRepository
	permissions []*models.Permission
	testAdmin   *models.Permission
}

// SetupSuite initializes the test suite
func (s *PermissionRepositoryTestSuite) SetupSuite() {
	// Call the parent SetupSuite method
	s.TestSuite.SetupSuite()

	// Initialize the repository
	s.repo = repoImpl.NewPermissionRepository(s.DB)
}

// SetupTest runs before each test
func (s *PermissionRepositoryTestSuite) SetupTest() {
	ctx := context.Background()

	// Get all permissions to use in tests
	allPermissions, err := s.repo.List(ctx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), allPermissions, "Test database must have seeded permissions")

	// Store permissions for use in tests
	s.permissions = allPermissions

	// Try to get the admin permission (should exist from seeds)
	s.testAdmin = fixtures.GetAdminPermission(ctx, s.repo)
	require.NotNil(s.T(), s.testAdmin, "Admin permission should exist in the database")
}

// TestFindPermissionsByIDs tests finding permissions by multiple IDs
func (s *PermissionRepositoryTestSuite) TestFindPermissionsByIDs() {
	ctx := context.Background()

	// Skip the test if no permissions are found
	if len(s.permissions) < 2 {
		s.T().Skip("Need at least 2 permissions in database for this test")
	}

	// Collect permission IDs
	var permissionIDs []int
	for _, p := range s.permissions {
		permissionIDs = append(permissionIDs, p.ID)
	}

	// Make sure we have at least 2 valid permission IDs for testing
	validIDs := []int{}
	for _, id := range permissionIDs {
		if len(validIDs) < 2 {
			validIDs = append(validIDs, id)
		} else {
			break
		}
	}

	// Test with multiple IDs
	testCases := []struct {
		name             string
		ids              []int
		expectNil        bool
		minExpectedCount int
	}{
		{
			name:             "Find by single ID",
			ids:              slices.Clone(validIDs[:1]), // Make a copy to avoid slice modification issues
			expectNil:        false,
			minExpectedCount: 1,
		},
		{
			name:             "Find by multiple IDs",
			ids:              slices.Clone(validIDs), // Make a copy of the entire validIDs slice
			expectNil:        false,
			minExpectedCount: 2,
		},
		{
			name:             "Find with non-existent IDs",
			ids:              []int{999999},
			expectNil:        false,
			minExpectedCount: 0,
		},
		{
			name:             "Find with mixed existent and non-existent IDs",
			ids:              append([]int{validIDs[0]}, 999999),
			expectNil:        false,
			minExpectedCount: 1,
		},
		{
			name:             "Find with empty ID array",
			ids:              []int{},
			expectNil:        false,
			minExpectedCount: 0,
		},
	}

	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			permissions, err := s.repo.FindByIDs(ctx, tc.ids)

			log.Printf("TestFindPermissionsByIDs: %s - IDs: %v - Found permissions count: %d",
				tc.name, tc.ids, len(permissions))

			if tc.expectNil {
				assert.Nil(t, permissions)
			} else {
				assert.NoError(t, err)

				// For the non-existent ID case
				if tc.minExpectedCount == 0 && len(tc.ids) > 0 {
					// If we're testing with IDs that don't exist, we should get an empty slice
					assert.Empty(t, permissions)
				} else if tc.minExpectedCount > 0 {
					// If we're expecting something, validate the returned permissions
					assert.NotNil(t, permissions)
					assert.GreaterOrEqual(t, len(permissions), tc.minExpectedCount)

					// Check that returned permissions match the requested IDs
					for _, p := range permissions {
						found := slices.Contains(tc.ids, p.ID)
						assert.True(t, found, "Returned permission ID not found in requested IDs")
					}
				}
			}
		})
	}
}

// TestListPermissions tests listing all permissions
func (s *PermissionRepositoryTestSuite) TestListPermissions() {
	ctx := context.Background()

	// Test the list method
	permissions, err := s.repo.List(ctx)

	// Assertions
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), permissions)
	assert.NotEmpty(s.T(), permissions, "Permissions list should not be empty")

	// Check that permissions have expected properties
	for _, p := range permissions {
		assert.NotZero(s.T(), p.ID, "Permission ID should not be zero")
		assert.NotEmpty(s.T(), p.Name, "Permission name should not be empty")
		assert.NotEmpty(s.T(), p.Code, "Permission code should not be empty")
	}

	// Verify the admin-related permissions exist
	// Instead of checking for IsAdmin, we check for specific admin permission codes
	foundUserManage := false
	for _, p := range permissions {
		if p.Code == string(models.PermissionCodeUserManage) {
			foundUserManage = true
			break
		}
	}
	assert.True(s.T(), foundUserManage, "USER_MANAGE permission should exist in the permissions list")

	// Verify that we get at least the predefined number of permissions from the seed
	assert.GreaterOrEqual(s.T(), len(permissions), 1, "Should have at least the basic permissions from seed")
}

// TestFindPermissionByID tests the helper method from fixtures that finds a permission by ID
func (s *PermissionRepositoryTestSuite) TestFindPermissionByID() {
	ctx := context.Background()

	// Skip if no permissions are found
	if len(s.permissions) == 0 {
		s.T().Skip("No permissions found in the database")
		return
	}

	// Test finding an existing permission (admin)
	adminID := s.testAdmin.ID
	foundPermission := fixtures.FindPermissionByID(ctx, s.repo, adminID)
	assert.NotNil(s.T(), foundPermission)
	assert.Equal(s.T(), adminID, foundPermission.ID)
	assert.Equal(s.T(), string(models.PermissionCodeUserManage), foundPermission.Code)

	// Test finding a non-existent permission
	nonExistentPermission := fixtures.FindPermissionByID(ctx, s.repo, 999999)
	assert.Nil(s.T(), nonExistentPermission)
}

func (s *PermissionRepositoryTestSuite) TearDownTest() {
}

// TestPermissionRepository executes the permission repository test suite
func TestPermissionRepository(t *testing.T) {
	suite.Run(t, new(PermissionRepositoryTestSuite))
}
