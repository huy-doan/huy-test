package usecase_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"slices"

	"github.com/huydq/test/src/domain/models"
	"github.com/huydq/test/src/domain/repositories"
	"github.com/huydq/test/src/domain/repositories/filter"
	repoImpl "github.com/huydq/test/src/infrastructure/persistence/repositories"
	"github.com/huydq/test/src/tests"
	"github.com/huydq/test/src/tests/fixtures"
	"github.com/huydq/test/src/usecase"
	"github.com/stretchr/testify/suite"
)

type RoleUsecaseTestSuite struct {
	tests.TestSuite
	roleUsecase     *usecase.RoleUsecase
	roleRepo        repositories.RoleRepository
	permissionRepo  repositories.PermissionRepository
	ctx             context.Context
	adminRole       *models.Role       // Expected to be seeded as ID 1
	normalUserRole  *models.Role       // Expected to be seeded as ID 2
	adminPermission *models.Permission // Expected to be seeded as ID 1
	testRoles       []*models.Role     // Roles created for testing that need cleanup
}

func TestRoleUsecaseSuite(t *testing.T) {
	suite.Run(t, new(RoleUsecaseTestSuite))
}

func (s *RoleUsecaseTestSuite) SetupSuite() {
	// Call the parent SetupSuite method
	s.TestSuite.SetupSuite()
}

func (s *RoleUsecaseTestSuite) SetupTest() {
	// Initialize context
	s.ctx = context.Background()

	// Initialize repositories using the test DB
	s.roleRepo = repoImpl.NewRoleRepository(s.DB)
	s.permissionRepo = repoImpl.NewPermissionRepository(s.DB)

	// Create roleUsecase with real repositories
	s.roleUsecase = usecase.NewRoleUsecase(s.roleRepo, s.permissionRepo)

	// Initialize testRoles slice for tracking created roles
	s.testRoles = make([]*models.Role, 0)

	// Load seeded roles using fixtures
	s.adminRole = fixtures.GetAdminRole(s.ctx, s.roleRepo)
	s.Require().NotNil(s.adminRole, "Admin role not found in database")

	s.normalUserRole = fixtures.GetNormalUserRole(s.ctx, s.roleRepo)
	s.Require().NotNil(s.normalUserRole, "Normal user role not found in database")

	// Load seeded permission using fixtures
	s.adminPermission = fixtures.GetAdminPermission(s.ctx, s.permissionRepo)
	s.Require().NotNil(s.adminPermission, "Admin permission not found in database")
}

func (s *RoleUsecaseTestSuite) createTestRole(nameSuffix string) *models.Role {
	// Generate a unique name and code for the role using timestamp
	timestamp := time.Now().UnixNano()
	roleName := fmt.Sprintf("Test Role %s %d", nameSuffix, timestamp)
	roleCode := fmt.Sprintf("TEST_ROLE_%s_%d", nameSuffix, timestamp)

	testRole := &models.Role{
		Name: roleName,
		Code: roleCode,
	}

	err := s.roleRepo.Create(s.ctx, testRole)
	s.Require().NoError(err)
	s.Require().NotZero(testRole.ID, "Test role ID should not be zero after creation")

	// Add to testRoles for cleanup
	s.testRoles = append(s.testRoles, testRole)

	return testRole
}

func (s *RoleUsecaseTestSuite) TearDownTest() {
	// Clean up all test roles created during the test
	for _, role := range s.testRoles {
		if role != nil && role.ID != 0 {
			err := s.roleRepo.Delete(s.ctx, role.ID)
			if err != nil {
				// Log but don't fail the test if cleanup fails
				fmt.Printf("Warning: Failed to delete test role ID %d: %v\n", role.ID, err)
			}
		}
	}
	// Clear the testRoles slice
	s.testRoles = nil
}

func (s *RoleUsecaseTestSuite) TearDownSuite() {
}

func (s *RoleUsecaseTestSuite) TestListRoles() {
	// Create a test role for listing with unique name
	testRole := s.createTestRole("List")

	// Test with valid pagination
	roleFilter := filter.NewRoleFilter()
	roleFilter.SetPagination(1, 100) // Use large page size to ensure we get all roles
	roles, totalPages, totalCount, err := s.roleUsecase.ListRoles(s.ctx, roleFilter)
	s.Assert().NoError(err)
	s.Assert().NotEmpty(roles)
	s.Assert().GreaterOrEqual(totalPages, 1)        // At least one page should exist
	s.Assert().GreaterOrEqual(totalCount, int64(1)) // At least one record should exist

	// Verify the seeded roles are present
	foundAdmin := false
	foundNormalUser := false
	foundTestRole := false
	for _, role := range roles {
		if role.Code == string(models.RoleCodeAdmin) {
			foundAdmin = true
		} else if role.Code == string(models.RoleCodeNormalUser) {
			foundNormalUser = true
		} else if role.ID == testRole.ID {
			foundTestRole = true
			// Verify properties match what we expect
			s.Assert().Equal(testRole.Name, role.Name)
			s.Assert().Equal(testRole.Code, role.Code)
		}
	}
	s.Assert().True(foundAdmin, "Admin role should be in the list")
	s.Assert().True(foundNormalUser, "Normal user role should be in the list")
	s.Assert().True(foundTestRole, "Test role should be in the list")

	// Test with invalid pagination (should use defaults)
	invalidFilter := filter.NewRoleFilter()
	invalidFilter.SetPagination(0, 0)
	roles, totalPages, totalCount, err = s.roleUsecase.ListRoles(s.ctx, invalidFilter)
	s.Assert().NoError(err)
	s.Assert().NotEmpty(roles)
	s.Assert().GreaterOrEqual(totalPages, 1)
	s.Assert().GreaterOrEqual(totalCount, int64(1))

	// Test with explicit pagination
	smallPageFilter := filter.NewRoleFilter()
	smallPageFilter.SetPagination(1, 2)
	roles, totalPages, totalCount, err = s.roleUsecase.ListRoles(s.ctx, smallPageFilter)
	s.Assert().NoError(err)
	s.Assert().NotEmpty(roles)
	s.Assert().GreaterOrEqual(totalPages, 1)
	s.Assert().GreaterOrEqual(totalCount, int64(1))
	s.Assert().LessOrEqual(len(roles), 2) // Should respect the page size

	// Test with name filter
	nameFilter := filter.NewRoleFilter()
	nameFilter.SetPagination(1, 10)
	testRoleName := "List" // Part of the test role name created above
	nameFilter.Name = &testRoleName
	nameFilter.ApplyFilters() // Apply the filter conditions
	roles, totalPages, totalCount, err = s.roleUsecase.ListRoles(s.ctx, nameFilter)
	s.Assert().NoError(err)
	s.Assert().GreaterOrEqual(totalCount, int64(1), "Total count should be at least 1 when filtering by name")
	s.Assert().GreaterOrEqual(totalPages, 1, "Total pages should be at least 1 when filtering by name")

	// Should find at least our test role with "List" in the name
	foundTestRole = false
	for _, role := range roles {
		if role.ID == testRole.ID {
			foundTestRole = true
			break
		}
	}
	s.Assert().True(foundTestRole, "Should find the test role when filtering by name")
}

func (s *RoleUsecaseTestSuite) TestGetRoleByID() {
	// Create a test role with unique name
	testRole := s.createTestRole("Get")

	// Test with existing role (use admin role from seed)
	role, err := s.roleUsecase.GetRoleByID(s.ctx, s.adminRole.ID)
	s.Assert().NoError(err)
	s.Assert().NotNil(role)
	s.Assert().Equal(s.adminRole.ID, role.ID)
	s.Assert().Equal(s.adminRole.Name, role.Name)
	s.Assert().Equal(s.adminRole.Code, role.Code)

	// Test with existing role (use test role created in setup)
	role, err = s.roleUsecase.GetRoleByID(s.ctx, testRole.ID)
	s.Assert().NoError(err)
	s.Assert().NotNil(role)
	s.Assert().Equal(testRole.ID, role.ID)
	s.Assert().Equal(testRole.Name, role.Name)
	s.Assert().Equal(testRole.Code, role.Code)

	// Test with non-existent role
	role, err = s.roleUsecase.GetRoleByID(s.ctx, 999999)
	s.Assert().NoError(err) // No error, just nil role
	s.Assert().Nil(role)
}

func (s *RoleUsecaseTestSuite) TestValidatePermissions() {
	// Create a test role with unique name (not needed for this test but adds test coverage)
	s.createTestRole("Validate")

	// Test with valid permission (admin permission from seed)
	err := s.roleUsecase.ValidatePermissions(s.ctx, []int{s.adminPermission.ID})
	s.Assert().NoError(err)

	// Test with empty permissions
	err = s.roleUsecase.ValidatePermissions(s.ctx, []int{})
	s.Assert().NoError(err)

	// Test with non-existent permission
	err = s.roleUsecase.ValidatePermissions(s.ctx, []int{999999})
	s.Assert().Error(err)
	s.Assert().Contains(err.Error(), "one or more permission IDs do not exist")

	// Test with mix of valid and invalid permissions
	err = s.roleUsecase.ValidatePermissions(s.ctx, []int{s.adminPermission.ID, 999999})
	s.Assert().Error(err)
	s.Assert().Contains(err.Error(), "one or more permission IDs do not exist")
}

func (s *RoleUsecaseTestSuite) TestCreateRole() {
	// Test creating role with no permissions - using unique name
	timestamp := time.Now().UnixNano()
	noPermRoleName := fmt.Sprintf("No Perm Role %d", timestamp)
	noPermRoleCode := fmt.Sprintf("NO_PERM_ROLE_%d", timestamp)

	noPermRole := &models.Role{
		Name: noPermRoleName,
		Code: noPermRoleCode,
	}
	err := s.roleUsecase.CreateRole(s.ctx, noPermRole)
	s.Assert().NoError(err)
	s.Assert().NotZero(noPermRole.ID)

	// Add to testRoles for cleanup
	s.testRoles = append(s.testRoles, noPermRole)

	// Verify the role was created correctly
	createdRole, err := s.roleRepo.FindByID(s.ctx, noPermRole.ID)
	s.Assert().NoError(err)
	s.Assert().NotNil(createdRole)
	s.Assert().Equal(noPermRole.Name, createdRole.Name)
	s.Assert().Equal(noPermRole.Code, createdRole.Code)

	// Test creating role with permissions (using admin permission from seed) - using unique name
	timestamp = time.Now().UnixNano()
	withPermRoleName := fmt.Sprintf("With Perm Role %d", timestamp)
	withPermRoleCode := fmt.Sprintf("WITH_PERM_ROLE_%d", timestamp)

	withPermRole := &models.Role{
		Name: withPermRoleName,
		Code: withPermRoleCode,
		Permissions: []*models.Permission{
			{ID: s.adminPermission.ID},
		},
	}
	err = s.roleUsecase.CreateRole(s.ctx, withPermRole)
	s.Assert().NoError(err)
	s.Assert().NotZero(withPermRole.ID)

	// Add to testRoles for cleanup
	s.testRoles = append(s.testRoles, withPermRole)

	// Verify the role was created with permissions
	createdRoleWithPerms, err := s.roleRepo.FindByID(s.ctx, withPermRole.ID)
	s.Assert().NoError(err)
	s.Assert().NotNil(createdRoleWithPerms)
	s.Assert().Equal(withPermRole.Name, createdRoleWithPerms.Name)
	s.Assert().Equal(withPermRole.Code, createdRoleWithPerms.Code)
	s.Assert().NotEmpty(createdRoleWithPerms.Permissions)
	s.Assert().Equal(s.adminPermission.ID, createdRoleWithPerms.Permissions[0].ID)

	// Test creating role with invalid permissions - using unique name
	timestamp = time.Now().UnixNano()
	invalidPermRoleName := fmt.Sprintf("Invalid Perm Role %d", timestamp)
	invalidPermRoleCode := fmt.Sprintf("INVALID_PERM_ROLE_%d", timestamp)

	invalidPermRole := &models.Role{
		Name: invalidPermRoleName,
		Code: invalidPermRoleCode,
		Permissions: []*models.Permission{
			{ID: 999999}, // Non-existent permission
		},
	}
	err = s.roleUsecase.CreateRole(s.ctx, invalidPermRole)
	s.Assert().Error(err)
	s.Assert().Contains(err.Error(), "one or more permission IDs do not exist")

	// Test creating role with duplicate name
	duplicateNameRole := &models.Role{
		Name: noPermRoleName, // Use the same name as the first role
		Code: fmt.Sprintf("UNIQUE_CODE_%d", time.Now().UnixNano()),
	}
	err = s.roleUsecase.CreateRole(s.ctx, duplicateNameRole)
	s.Assert().Error(err)
	s.Assert().Contains(err.Error(), "role name already exists")

	// Test creating role with duplicate code
	duplicateCodeRole := &models.Role{
		Name: fmt.Sprintf("Unique Name %d", time.Now().UnixNano()),
		Code: noPermRoleCode, // Use the same code as the first role
	}
	err = s.roleUsecase.CreateRole(s.ctx, duplicateCodeRole)
	s.Assert().Error(err)
	s.Assert().Contains(err.Error(), "role code already exists")

	// No need to add to testRoles as it wasn't created
}

func (s *RoleUsecaseTestSuite) TestUpdateRole() {
	// Create a role to update with unique name
	roleToUpdate := s.createTestRole("Update")

	// Create another role to test duplicate name validation
	anotherRole := s.createTestRole("Another")

	// Test updating role with valid data (using admin permission from seed)
	timestamp := time.Now().UnixNano()
	updatedRoleName := fmt.Sprintf("Updated Role Name %d", timestamp)
	updateData := &usecase.UpdateRoleRequest{
		Name:          updatedRoleName,
		PermissionIDs: []int{s.adminPermission.ID},
	}
	err := s.roleUsecase.UpdateRole(s.ctx, roleToUpdate.ID, updateData)
	s.Assert().NoError(err)

	// Verify role was updated
	updatedRole, err := s.roleRepo.FindByID(s.ctx, roleToUpdate.ID)
	s.Assert().NoError(err)
	s.Assert().NotNil(updatedRole)
	s.Assert().Equal(updateData.Name, updatedRole.Name)
	// Code field should remain unchanged since it's not included in the UpdateRoleRequest
	s.Assert().Equal(roleToUpdate.Code, updatedRole.Code)
	s.Assert().Equal(1, len(updatedRole.Permissions))
	s.Assert().Equal(s.adminPermission.ID, updatedRole.Permissions[0].ID)

	// Test updating with duplicate name (using another role's name)
	duplicateNameUpdateData := &usecase.UpdateRoleRequest{
		Name:          anotherRole.Name, // Try to use another role's name
		PermissionIDs: []int{s.adminPermission.ID},
	}
	err = s.roleUsecase.UpdateRole(s.ctx, roleToUpdate.ID, duplicateNameUpdateData)
	s.Assert().Error(err)
	s.Assert().Contains(err.Error(), "role name already exists")

	// Test updating non-existent role
	nonExistentRoleName := fmt.Sprintf("Non-existent Role %d", time.Now().UnixNano())
	err = s.roleUsecase.UpdateRole(s.ctx, 999999, &usecase.UpdateRoleRequest{
		Name: nonExistentRoleName,
	})
	s.Assert().Error(err)
	s.Assert().Contains(err.Error(), "role not found")

	// Test updating with invalid permissions
	invalidRoleName := fmt.Sprintf("Invalid Perms Role %d", time.Now().UnixNano())
	invalidUpdateData := &usecase.UpdateRoleRequest{
		Name:          invalidRoleName,
		PermissionIDs: []int{999999}, // Non-existent permission
	}
	err = s.roleUsecase.UpdateRole(s.ctx, roleToUpdate.ID, invalidUpdateData)
	s.Assert().Error(err)
	s.Assert().Contains(err.Error(), "one or more permission IDs do not exist")
}

func (s *RoleUsecaseTestSuite) TestDeleteRole() {
	// Create a role to delete with unique name
	roleToDelete := s.createTestRole("Delete")

	// Delete the role
	err := s.roleUsecase.DeleteRole(s.ctx, roleToDelete.ID)
	s.Assert().NoError(err)

	// Verify the role was deleted
	role, err := s.roleRepo.FindByID(s.ctx, roleToDelete.ID)
	s.Assert().NoError(err)
	s.Assert().Nil(role)

	// Remove from testRoles since we've already deleted it
	for i, r := range s.testRoles {
		if r.ID == roleToDelete.ID {
			s.testRoles = slices.Delete(s.testRoles, i, i+1)
			break
		}
	}

	// Test deleting non-existent role
	err = s.roleUsecase.DeleteRole(s.ctx, 999999)
	s.Assert().NoError(err) // No error expected, just nil role

	// Verify the role doesn't appear in list results
	roleFilter := filter.NewRoleFilter()
	roleFilter.SetPagination(1, 100)
	roles, _, _, err := s.roleUsecase.ListRoles(s.ctx, roleFilter)
	s.Assert().NoError(err)

	// Verify the deleted role is not in the results
	for _, r := range roles {
		s.Assert().NotEqual(roleToDelete.ID, r.ID, "Deleted role should not appear in results")
	}

	// Try to get the deleted role by ID - should return nil
	deletedRole, err := s.roleUsecase.GetRoleByID(s.ctx, roleToDelete.ID)
	s.Assert().NoError(err)
	s.Assert().Nil(deletedRole, "GetRoleByID should return nil for deleted role")
}

func (s *RoleUsecaseTestSuite) TestBatchUpdateRolePermissions() {
	// Create multiple test roles with unique names
	role1 := s.createTestRole("BatchUpdateA")
	role2 := s.createTestRole("BatchUpdateB")
	role3 := s.createTestRole("BatchUpdateC")

	// Initially, ensure roles have no permissions
	for _, role := range []*models.Role{role1, role2, role3} {
		roleFromDB, err := s.roleRepo.FindByID(s.ctx, role.ID)
		s.Assert().NoError(err)
		s.Assert().Empty(roleFromDB.Permissions, "Role should start with no permissions")
	}

	// Get a couple of permissions to assign
	permissions, err := s.permissionRepo.List(s.ctx)
	s.Require().NoError(err)
	s.Require().GreaterOrEqual(len(permissions), 2, "Need at least 2 permissions for this test")

	permission1ID := permissions[0].ID
	permission2ID := permissions[1].ID

	// Create batch update data - update roles with different permission combinations
	updates := []struct {
		ID            int
		PermissionIDs []int
	}{
		{ID: role1.ID, PermissionIDs: []int{permission1ID}},                // First permission only
		{ID: role2.ID, PermissionIDs: []int{permission2ID}},                // Second permission only
		{ID: role3.ID, PermissionIDs: []int{permission1ID, permission2ID}}, // Both permissions
		{ID: 999999, PermissionIDs: []int{permission1ID}},                  // Non-existent role ID
	}

	// Perform batch update
	updatedRoleIDs, err := s.roleUsecase.BatchUpdateRolePermissions(s.ctx, updates)
	s.Assert().NoError(err)

	// Should get back 3 successfully updated role IDs
	s.Assert().Len(updatedRoleIDs, 3, "Should have updated 3 roles")
	s.Assert().Contains(updatedRoleIDs, role1.ID, "Role 1 should be in successfully updated IDs")
	s.Assert().Contains(updatedRoleIDs, role2.ID, "Role 2 should be in successfully updated IDs")
	s.Assert().Contains(updatedRoleIDs, role3.ID, "Role 3 should be in successfully updated IDs")

	// Verify each role has the correct permissions
	// Role 1 should have permission 1
	role1Updated, err := s.roleRepo.FindByID(s.ctx, role1.ID)
	s.Assert().NoError(err)
	s.Assert().Len(role1Updated.Permissions, 1, "Role 1 should have 1 permission")
	s.Assert().Equal(permission1ID, role1Updated.Permissions[0].ID, "Role 1 should have permission1")

	// Role 2 should have permission 2
	role2Updated, err := s.roleRepo.FindByID(s.ctx, role2.ID)
	s.Assert().NoError(err)
	s.Assert().Len(role2Updated.Permissions, 1, "Role 2 should have 1 permission")
	s.Assert().Equal(permission2ID, role2Updated.Permissions[0].ID, "Role 2 should have permission2")

	// Role 3 should have both permissions
	role3Updated, err := s.roleRepo.FindByID(s.ctx, role3.ID)
	s.Assert().NoError(err)
	s.Assert().Len(role3Updated.Permissions, 2, "Role 3 should have 2 permissions")

	// Check both permissions are present in role3 using a tagged switch
	foundPerm1 := false
	foundPerm2 := false
	for _, perm := range role3Updated.Permissions {
		switch perm.ID {
		case permission1ID:
			foundPerm1 = true
		case permission2ID:
			foundPerm2 = true
		}
	}
	s.Assert().True(foundPerm1, "Role 3 should have permission1")
	s.Assert().True(foundPerm2, "Role 3 should have permission2")

	// Test with invalid permission IDs
	invalidUpdates := []struct {
		ID            int
		PermissionIDs []int
	}{
		{ID: role1.ID, PermissionIDs: []int{999999}}, // Invalid permission ID
	}

	// This should fail validation and return an error
	_, err = s.roleUsecase.BatchUpdateRolePermissions(s.ctx, invalidUpdates)
	s.Assert().Error(err)
	s.Assert().Contains(err.Error(), "one or more permission IDs do not exist")
}
