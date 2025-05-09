package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"slices"

	"github.com/huydq/test/src/api/http/handlers"
	"github.com/huydq/test/src/api/http/middleware"
	"github.com/huydq/test/src/api/http/response"
	"github.com/huydq/test/src/domain/models"
	"github.com/huydq/test/src/domain/repositories"
	repoImpl "github.com/huydq/test/src/infrastructure/persistence/repositories"
	"github.com/huydq/test/src/tests"
	"github.com/huydq/test/src/usecase"
	"github.com/stretchr/testify/suite"
)

type RoleHandlerTestSuite struct {
	tests.TestSuite
	roleRepo        repositories.RoleRepository
	permissionRepo  repositories.PermissionRepository
	roleUsecase     *usecase.RoleUsecase
	roleHandler     *handlers.RoleHandler
	seedPermissions []*models.Permission
	testRoles       []*models.Role
}

func (s *RoleHandlerTestSuite) SetupSuite() {
	// Call the parent SetupSuite method
	s.TestSuite.SetupSuite()

	// Initialize repositories
	s.roleRepo = repoImpl.NewRoleRepository(s.DB)
	s.permissionRepo = repoImpl.NewPermissionRepository(s.DB)

	// Initialize usecases
	s.roleUsecase = usecase.NewRoleUsecase(s.roleRepo, s.permissionRepo)

	// Initialize handler
	s.roleHandler = handlers.NewRoleHandler(s.roleUsecase)
}

func (s *RoleHandlerTestSuite) SetupTest() {
	// Use seed permissions for tests
	s.seedPermissions = s.getSeedPermissions()
}

func (s *RoleHandlerTestSuite) getSeedPermissions() []*models.Permission {
	permissions, err := s.permissionRepo.List(context.Background())
	s.Require().NoError(err)
	s.Require().NotEmpty(permissions)
	return permissions
}

func (s *RoleHandlerTestSuite) createTestRole(nameSuffix string) *models.Role {
	role := &models.Role{
		Name: "Test Role " + nameSuffix,
		Code: "test_role_" + nameSuffix,
	}

	err := s.roleRepo.Create(context.Background(), role)
	s.Require().NoError(err)
	s.Require().NotZero(role.ID)

	// Add to cleanup list
	s.testRoles = append(s.testRoles, role)

	return role
}

func (s *RoleHandlerTestSuite) TearDownTest() {
	// Clean up roles created during the test
	for _, role := range s.testRoles {
		err := s.roleRepo.Delete(context.Background(), role.ID)
		if err != nil {
			s.T().Logf("Error cleaning up role %d: %v", role.ID, err)
		}
	}
	s.testRoles = nil
}

func (s *RoleHandlerTestSuite) TearDownSuite() {
	// No need to clean up seed permissions
}

// createAuthContext creates a context that mimics an authenticated admin user
func (s *RoleHandlerTestSuite) createAuthContext() context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, middleware.UserIDKey, 1)
	ctx = context.WithValue(ctx, middleware.RoleCodeKey, string(models.RoleCodeAdmin))
	return ctx
}

// createNormalUserContext creates a context for a non-admin user
func (s *RoleHandlerTestSuite) createNormalUserContext() context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, middleware.UserIDKey, 2)
	ctx = context.WithValue(ctx, middleware.RoleCodeKey, string(models.RoleCodeNormalUser))
	return ctx
}

func TestRoleHandlerSuite(t *testing.T) {
	suite.Run(t, new(RoleHandlerTestSuite))
}

func (s *RoleHandlerTestSuite) TestListRoles() {
	// Create test roles for listing
	s.createTestRole("List1")
	s.createTestRole("List2")
	s.createTestRole("List3")

	// Test with admin user - should succeed
	s.Run("Admin_Success", func() {
		req := httptest.NewRequest("GET", "/api/v1/admin/roles?page=1&page_size=10", nil)
		req = req.WithContext(s.createAuthContext())
		w := httptest.NewRecorder()

		s.roleHandler.ListRoles(w, req)

		s.Equal(http.StatusOK, w.Code)

		var resp response.Response
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		s.NoError(err)
		s.True(resp.Success)

		// Verify data contains roles
		data, ok := resp.Data.(map[string]any)
		s.True(ok)
		roles, ok := data["roles"].([]any)
		s.True(ok)
		s.GreaterOrEqual(len(roles), 3)
	})

	// Test with non-admin user - should fail with forbidden
	s.Run("NonAdmin_Forbidden", func() {
		req := httptest.NewRequest("GET", "/api/v1/admin/roles", nil)
		req = req.WithContext(s.createNormalUserContext())
		w := httptest.NewRecorder()

		s.roleHandler.ListRoles(w, req)

		s.Equal(http.StatusForbidden, w.Code)
	})

	// Test with pagination parameters
	s.Run("WithPagination", func() {
		req := httptest.NewRequest("GET", "/api/v1/admin/roles?page=1&page_size=2", nil)
		req = req.WithContext(s.createAuthContext())
		w := httptest.NewRecorder()

		// The handler should parse these parameters from the URL query
		s.roleHandler.ListRoles(w, req)

		s.Equal(http.StatusOK, w.Code)

		var resp response.Response
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		s.NoError(err)

		data, ok := resp.Data.(map[string]any)
		s.True(ok)
		s.Equal(float64(1), data["page"])
		s.Equal(float64(2), data["page_size"])
	})

	// Test with name filter parameter
	s.Run("WithNameFilter", func() {
		// Create a role with a unique name for filtering
		uniqueRole := s.createTestRole("UniqueNameForFiltering")

		// Request with name filter parameter
		req := httptest.NewRequest("GET", "/api/v1/admin/roles?name=UniqueNameForFiltering", nil)
		req = req.WithContext(s.createAuthContext())
		w := httptest.NewRecorder()

		// The handler should parse the name parameter from the URL query
		s.roleHandler.ListRoles(w, req)

		s.Equal(http.StatusOK, w.Code)

		var resp response.Response
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		s.NoError(err)
		s.True(resp.Success)

		// Verify data contains filtered roles
		data, ok := resp.Data.(map[string]any)
		s.True(ok)
		roles, ok := data["roles"].([]any)
		s.True(ok)

		// Should find at least the unique role
		found := false
		for _, roleData := range roles {
			role := roleData.(map[string]any)
			if role["name"].(string) == uniqueRole.Name {
				found = true
				break
			}
		}
		s.True(found, "Should find the role with the filtered name")
	})
}

func (s *RoleHandlerTestSuite) TestGetRoleByID() {
	// Create a test role
	testRole := s.createTestRole("GetTest")

	// Test with admin user - should succeed
	s.Run("Admin_Success", func() {
		req := httptest.NewRequest("GET", "/api/v1/admin/roles/"+s.idToString(testRole.ID), nil)
		req.URL.Path = "/api/v1/admin/roles/" + s.idToString(testRole.ID)
		req = req.WithContext(s.createAuthContext())
		w := httptest.NewRecorder()

		s.roleHandler.GetRoleByID(w, req)

		s.Equal(http.StatusOK, w.Code)

		var resp response.Response
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		s.NoError(err)
		s.True(resp.Success)

		roleData, ok := resp.Data.(map[string]any)
		s.True(ok)
		s.Equal(testRole.Name, roleData["name"])
		s.Equal(testRole.Code, roleData["code"])
	})

	// Test with non-admin user - should fail with forbidden
	s.Run("NonAdmin_Forbidden", func() {
		req := httptest.NewRequest("GET", "/api/v1/admin/roles/"+s.idToString(testRole.ID), nil)
		req.URL.Path = "/api/v1/admin/roles/" + s.idToString(testRole.ID)
		req = req.WithContext(s.createNormalUserContext())
		w := httptest.NewRecorder()

		s.roleHandler.GetRoleByID(w, req)

		s.Equal(http.StatusForbidden, w.Code)
	})

	// Test with invalid ID - should fail with bad request
	s.Run("InvalidID", func() {
		req := httptest.NewRequest("GET", "/api/v1/admin/roles/invalid", nil)
		req.URL.Path = "/api/v1/admin/roles/invalid"
		req = req.WithContext(s.createAuthContext())
		w := httptest.NewRecorder()

		s.roleHandler.GetRoleByID(w, req)

		s.Equal(http.StatusBadRequest, w.Code)
	})

	// Test with non-existent ID - should fail with not found
	s.Run("NonExistentID", func() {
		req := httptest.NewRequest("GET", "/api/v1/admin/roles/9999", nil)
		req.URL.Path = "/api/v1/admin/roles/9999"
		req = req.WithContext(s.createAuthContext())
		w := httptest.NewRecorder()

		s.roleHandler.GetRoleByID(w, req)

		s.Equal(http.StatusNotFound, w.Code)
	})
}

func (s *RoleHandlerTestSuite) TestCreateRole() {
	// Test successful creation with admin user
	s.Run("Admin_Success", func() {
		requestBody := map[string]any{
			"name":           "New Test Role",
			"code":           "new_test_role",
			"permission_ids": []int{s.seedPermissions[0].ID},
		}
		bodyBytes, _ := json.Marshal(requestBody)

		req := httptest.NewRequest("POST", "/api/v1/admin/roles", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(s.createAuthContext())
		w := httptest.NewRecorder()

		s.roleHandler.CreateRole(w, req)

		s.Equal(http.StatusCreated, w.Code)

		var resp response.Response
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		s.NoError(err)
		s.True(resp.Success)

		roleData := resp.Data.(map[string]any)
		s.Equal("New Test Role", roleData["name"])
		s.Equal("new_test_role", roleData["code"])

		// Add to cleanup
		createdID := int(roleData["id"].(float64))
		s.testRoles = append(s.testRoles, &models.Role{ID: createdID})
	})

	// Test with non-admin user - should fail with forbidden
	s.Run("NonAdmin_Forbidden", func() {
		requestBody := map[string]any{
			"name": "Should Fail Role",
			"code": "should_fail",
		}
		bodyBytes, _ := json.Marshal(requestBody)

		req := httptest.NewRequest("POST", "/api/v1/admin/roles", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(s.createNormalUserContext())
		w := httptest.NewRecorder()

		s.roleHandler.CreateRole(w, req)

		s.Equal(http.StatusForbidden, w.Code)
	})

	// Test with invalid permissions - should fail with bad request
	s.Run("InvalidPermissions", func() {
		requestBody := map[string]any{
			"name":           "Invalid Perm Role",
			"code":           "invalid_perm",
			"permission_ids": []int{9999}, // Non-existent permission
		}
		bodyBytes, _ := json.Marshal(requestBody)

		req := httptest.NewRequest("POST", "/api/v1/admin/roles", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(s.createAuthContext())
		w := httptest.NewRecorder()

		s.roleHandler.CreateRole(w, req)

		s.Equal(http.StatusBadRequest, w.Code)
	})

	// Test with invalid request body - should fail with bad request
	s.Run("InvalidRequestBody", func() {
		requestBody := map[string]any{
			// Missing required fields
		}
		bodyBytes, _ := json.Marshal(requestBody)

		req := httptest.NewRequest("POST", "/api/v1/admin/roles", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(s.createAuthContext())
		w := httptest.NewRecorder()

		s.roleHandler.CreateRole(w, req)

		s.Equal(http.StatusBadRequest, w.Code)
	})
}

func (s *RoleHandlerTestSuite) TestUpdateRole() {
	// Create a role to update
	testRole := s.createTestRole("UpdateTest")

	// Test successful update with admin user
	s.Run("Admin_Success", func() {
		// Generate a unique name for this update
		timestamp := time.Now().UnixNano()
		updatedRoleName := fmt.Sprintf("Updated Role Name %d", timestamp)

		updateBody := map[string]any{
			"name":           updatedRoleName,
			"permission_ids": []int{s.seedPermissions[0].ID},
		}
		bodyBytes, _ := json.Marshal(updateBody)

		req := httptest.NewRequest("PUT", "/api/v1/admin/roles/"+s.idToString(testRole.ID), bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		req.URL.Path = "/api/v1/admin/roles/" + s.idToString(testRole.ID)
		req = req.WithContext(s.createAuthContext())
		w := httptest.NewRecorder()

		s.roleHandler.UpdateRole(w, req)

		s.Equal(http.StatusOK, w.Code)

		var resp response.Response
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		s.NoError(err)
		s.True(resp.Success)

		// Check if Data is not nil before attempting to access it as a map
		s.NotNil(resp.Data, "Response data should not be nil")

		// Safe type assertion with check
		roleData, ok := resp.Data.(map[string]any)
		s.True(ok, "Response data should be a map[string]interface{}")
		s.Equal(updatedRoleName, roleData["name"])
		s.Equal(testRole.Code, roleData["code"]) // Code shouldn't change
	})

	// Test with non-admin user - should fail with forbidden
	s.Run("NonAdmin_Forbidden", func() {
		timestamp := time.Now().UnixNano()
		updateBody := map[string]any{
			"name": fmt.Sprintf("Should Fail Update %d", timestamp),
		}
		bodyBytes, _ := json.Marshal(updateBody)

		req := httptest.NewRequest("PUT", "/api/v1/admin/roles/"+s.idToString(testRole.ID), bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		req.URL.Path = "/api/v1/admin/roles/" + s.idToString(testRole.ID)
		req = req.WithContext(s.createNormalUserContext())
		w := httptest.NewRecorder()

		s.roleHandler.UpdateRole(w, req)

		s.Equal(http.StatusForbidden, w.Code)
	})

	// Test with invalid role ID - should fail with not found
	s.Run("InvalidRoleID", func() {
		timestamp := time.Now().UnixNano()
		updateBody := map[string]any{
			"name": fmt.Sprintf("Invalid Role ID %d", timestamp),
		}
		bodyBytes, _ := json.Marshal(updateBody)

		req := httptest.NewRequest("PUT", "/api/v1/admin/roles/9999", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		req.URL.Path = "/api/v1/admin/roles/9999"
		req = req.WithContext(s.createAuthContext())
		w := httptest.NewRecorder()

		s.roleHandler.UpdateRole(w, req)

		s.Equal(http.StatusNotFound, w.Code)
	})

	// Test with invalid permissions - should fail with bad request
	s.Run("InvalidPermissions", func() {
		timestamp := time.Now().UnixNano()
		updateBody := map[string]any{
			"name":           fmt.Sprintf("Invalid Perms Update %d", timestamp),
			"permission_ids": []int{9999}, // Non-existent permission
		}
		bodyBytes, _ := json.Marshal(updateBody)

		req := httptest.NewRequest("PUT", "/api/v1/admin/roles/"+s.idToString(testRole.ID), bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		req.URL.Path = "/api/v1/admin/roles/" + s.idToString(testRole.ID)
		req = req.WithContext(s.createAuthContext())
		w := httptest.NewRecorder()

		s.roleHandler.UpdateRole(w, req)

		s.Equal(http.StatusBadRequest, w.Code)
	})
}

func (s *RoleHandlerTestSuite) TestDeleteRole() {
	// Create a role to delete
	testRole := s.createTestRole("DeleteTest")

	// Test successful deletion with admin user
	s.Run("Admin_Success", func() {
		req := httptest.NewRequest("DELETE", "/api/v1/admin/roles/"+s.idToString(testRole.ID), nil)
		req.URL.Path = "/api/v1/admin/roles/" + s.idToString(testRole.ID)
		req = req.WithContext(s.createAuthContext())
		w := httptest.NewRecorder()

		s.roleHandler.DeleteRole(w, req)

		s.Equal(http.StatusOK, w.Code)

		// Verify role is actually deleted
		deletedRole, err := s.roleRepo.FindByID(context.Background(), testRole.ID)
		s.NoError(err)
		s.Nil(deletedRole)

		// Remove from cleanup list since it's already deleted
		for i, r := range s.testRoles {
			if r.ID == testRole.ID {
				s.testRoles = slices.Delete(s.testRoles, i, i+1)
				break
			}
		}
	})

	// Create another role for the next tests
	testRole = s.createTestRole("DeleteTest2")

	// Test with non-admin user - should fail with forbidden
	s.Run("NonAdmin_Forbidden", func() {
		req := httptest.NewRequest("DELETE", "/api/v1/admin/roles/"+s.idToString(testRole.ID), nil)
		req.URL.Path = "/api/v1/admin/roles/" + s.idToString(testRole.ID)
		req = req.WithContext(s.createNormalUserContext())
		w := httptest.NewRecorder()

		s.roleHandler.DeleteRole(w, req)

		s.Equal(http.StatusForbidden, w.Code)
	})

	// Test with non-existent ID - should fail with not found
	s.Run("NonExistentID", func() {
		req := httptest.NewRequest("DELETE", "/api/v1/admin/roles/9999", nil)
		req.URL.Path = "/api/v1/admin/roles/9999"
		req = req.WithContext(s.createAuthContext())
		w := httptest.NewRecorder()

		s.roleHandler.DeleteRole(w, req)

		s.Equal(http.StatusNotFound, w.Code)
	})

	// Test with invalid ID - should fail with bad request
	s.Run("InvalidID", func() {
		req := httptest.NewRequest("DELETE", "/api/v1/admin/roles/invalid", nil)
		req.URL.Path = "/api/v1/admin/roles/invalid"
		req = req.WithContext(s.createAuthContext())
		w := httptest.NewRecorder()

		s.roleHandler.DeleteRole(w, req)

		s.Equal(http.StatusBadRequest, w.Code)
	})
}

func (s *RoleHandlerTestSuite) TestDuplicateNameAndCodeHandling() {
	// Create an initial role for testing duplicates
	initialRole := s.createTestRole("Initial")

	// Test duplicate name
	s.Run("DuplicateName_Failure", func() {
		requestBody := map[string]any{
			"name":           initialRole.Name, // Duplicate name
			"code":           "unique_code_" + strconv.Itoa(int(time.Now().UnixNano())),
			"permission_ids": []int{s.seedPermissions[0].ID},
		}
		bodyBytes, _ := json.Marshal(requestBody)

		req := httptest.NewRequest("POST", "/api/v1/admin/roles", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(s.createAuthContext())
		w := httptest.NewRecorder()

		s.roleHandler.CreateRole(w, req)

		// Should return bad request for duplicate name
		s.Equal(http.StatusBadRequest, w.Code)

		var resp response.Response
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		s.NoError(err)
		s.False(resp.Success)
		s.Contains(resp.Message, "duplicate_name")
	})

	// Test duplicate code
	s.Run("DuplicateCode_Failure", func() {
		requestBody := map[string]any{
			"name":           "Unique Name " + strconv.Itoa(int(time.Now().UnixNano())),
			"code":           initialRole.Code, // Duplicate code
			"permission_ids": []int{s.seedPermissions[0].ID},
		}
		bodyBytes, _ := json.Marshal(requestBody)

		req := httptest.NewRequest("POST", "/api/v1/admin/roles", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(s.createAuthContext())
		w := httptest.NewRecorder()

		s.roleHandler.CreateRole(w, req)

		// Should return bad request for duplicate code
		s.Equal(http.StatusBadRequest, w.Code)

		var resp response.Response
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		s.NoError(err)
		s.False(resp.Success)
		s.Contains(resp.Message, "duplicate_code")
	})

	// Test duplicate name during update
	s.Run("UpdateDuplicateName_Failure", func() {
		// Create another role to try to update with duplicate name
		anotherRole := s.createTestRole("Another")

		updateBody := map[string]any{
			"name":           initialRole.Name, // Try to update to the name of the initial role
			"permission_ids": []int{s.seedPermissions[0].ID},
		}
		bodyBytes, _ := json.Marshal(updateBody)

		req := httptest.NewRequest("PUT", "/api/v1/admin/roles/"+s.idToString(anotherRole.ID), bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		req.URL.Path = "/api/v1/admin/roles/" + s.idToString(anotherRole.ID)
		req = req.WithContext(s.createAuthContext())
		w := httptest.NewRecorder()

		s.roleHandler.UpdateRole(w, req)

		// Should return bad request for duplicate name
		s.Equal(http.StatusBadRequest, w.Code)

		var resp response.Response
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		s.NoError(err)
		s.False(resp.Success)
		s.Contains(resp.Message, "duplicate_name")
	})
}

func (s *RoleHandlerTestSuite) TestBatchUpdateRolePermissions() {
	// Create multiple test roles for batch updating
	role1 := s.createTestRole("BatchUpdate1")
	role2 := s.createTestRole("BatchUpdate2")
	role3 := s.createTestRole("BatchUpdate3")

	// Make sure we have at least 2 permissions for testing
	s.Require().GreaterOrEqual(len(s.seedPermissions), 2, "Need at least 2 permissions for this test")
	perm1ID := s.seedPermissions[0].ID
	perm2ID := s.seedPermissions[1].ID

	// Test successful batch update with admin user
	s.Run("Admin_Success", func() {
		requestBody := []map[string]any{
			{
				"id":             role1.ID,
				"permission_ids": []int{perm1ID},
			},
			{
				"id":             role2.ID,
				"permission_ids": []int{perm2ID},
			},
			{
				"id":             role3.ID,
				"permission_ids": []int{perm1ID, perm2ID},
			},
			{
				"id":             9999, // Non-existent role, should be skipped
				"permission_ids": []int{perm1ID},
			},
		}
		bodyBytes, _ := json.Marshal(requestBody)

		req := httptest.NewRequest("POST", "/api/v1/admin/roles/permissions/batch", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(s.createAuthContext())
		w := httptest.NewRecorder()

		s.roleHandler.BatchUpdateRolePermissions(w, req)

		s.Equal(http.StatusOK, w.Code)

		var resp response.Response
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		s.NoError(err)
		s.True(resp.Success)

		// Verify response contains updated role IDs
		data, ok := resp.Data.(map[string]any)
		s.True(ok, "Response data should be a map")
		updatedRoles, ok := data["updated_roles"].([]any)
		s.True(ok, "Updated roles should be an array")
		s.Len(updatedRoles, 3, "Should have updated 3 roles")

		// Verify the total count is returned
		totalUpdated, ok := data["total_updated"].(float64)
		s.True(ok, "Total updated should be a number")
		s.Equal(float64(3), totalUpdated)

		// Verify that the roles were actually updated with the correct permissions
		// Role 1 should have permission 1 only
		role1Updated, err := s.roleRepo.FindByID(context.Background(), role1.ID)
		s.NoError(err)
		s.NotNil(role1Updated)
		s.Len(role1Updated.Permissions, 1, "Role 1 should have 1 permission")
		s.Equal(perm1ID, role1Updated.Permissions[0].ID)

		// Role 2 should have permission 2 only
		role2Updated, err := s.roleRepo.FindByID(context.Background(), role2.ID)
		s.NoError(err)
		s.NotNil(role2Updated)
		s.Len(role2Updated.Permissions, 1, "Role 2 should have 1 permission")
		s.Equal(perm2ID, role2Updated.Permissions[0].ID)

		// Role 3 should have both permissions
		role3Updated, err := s.roleRepo.FindByID(context.Background(), role3.ID)
		s.NoError(err)
		s.NotNil(role3Updated)
		s.Len(role3Updated.Permissions, 2, "Role 3 should have 2 permissions")
	})

	// Test with non-admin user - should fail with forbidden
	s.Run("NonAdmin_Forbidden", func() {
		requestBody := []map[string]any{
			{
				"id":             role1.ID,
				"permission_ids": []int{perm1ID},
			},
		}
		bodyBytes, _ := json.Marshal(requestBody)

		req := httptest.NewRequest("POST", "/api/v1/admin/roles/permissions/batch", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(s.createNormalUserContext())
		w := httptest.NewRecorder()

		s.roleHandler.BatchUpdateRolePermissions(w, req)

		s.Equal(http.StatusForbidden, w.Code)
	})

	// Test with invalid permissions - should fail with bad request
	s.Run("InvalidPermissions", func() {
		requestBody := []map[string]any{
			{
				"id":             role1.ID,
				"permission_ids": []int{9999}, // Non-existent permission ID
			},
		}
		bodyBytes, _ := json.Marshal(requestBody)

		req := httptest.NewRequest("POST", "/api/v1/admin/roles/permissions/batch", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(s.createAuthContext())
		w := httptest.NewRecorder()

		s.roleHandler.BatchUpdateRolePermissions(w, req)

		s.Equal(http.StatusBadRequest, w.Code)
	})

	// Test with invalid request body - should fail with bad request
	s.Run("InvalidRequestBody", func() {
		// Empty array should be valid but do nothing
		requestBody := []map[string]any{}
		bodyBytes, _ := json.Marshal(requestBody)

		req := httptest.NewRequest("POST", "/api/v1/admin/roles/permissions/batch", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(s.createAuthContext())
		w := httptest.NewRecorder()

		s.roleHandler.BatchUpdateRolePermissions(w, req)

		s.Equal(http.StatusOK, w.Code)

		// Invalid JSON body should fail
		req = httptest.NewRequest("POST", "/api/v1/admin/roles/permissions/batch", bytes.NewBufferString("{invalid json"))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(s.createAuthContext())
		w = httptest.NewRecorder()

		s.roleHandler.BatchUpdateRolePermissions(w, req)

		s.Equal(http.StatusBadRequest, w.Code)
	})
}

// Helper function to convert ID to string
func (s *RoleHandlerTestSuite) idToString(id int) string {
	return strconv.Itoa(id)
}
