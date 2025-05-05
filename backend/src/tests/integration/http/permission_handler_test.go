package http_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/vnlab/makeshop-payment/src/api/http/handlers"
	"github.com/vnlab/makeshop-payment/src/api/http/middleware"
	"github.com/vnlab/makeshop-payment/src/api/http/response"
	"github.com/vnlab/makeshop-payment/src/domain/models"
	"github.com/vnlab/makeshop-payment/src/domain/repositories"
	repoImpl "github.com/vnlab/makeshop-payment/src/infrastructure/persistence/repositories"
	"github.com/vnlab/makeshop-payment/src/tests"
	"github.com/vnlab/makeshop-payment/src/usecase"
)

type PermissionHandlerTestSuite struct {
	tests.TestSuite
	permissionRepo    repositories.PermissionRepository
	permissionUsecase *usecase.PermissionUsecase
	permissionHandler *handlers.PermissionHandler
	seedPermissions   []*models.Permission
}

// SetupSuite initializes the test suite
func (s *PermissionHandlerTestSuite) SetupSuite() {
	// Call the parent SetupSuite method
	s.TestSuite.SetupSuite()

	// Initialize repositories
	s.permissionRepo = repoImpl.NewPermissionRepository(s.DB)

	// Initialize usecase
	s.permissionUsecase = usecase.NewPermissionUseCase(s.permissionRepo)

	// Initialize handler
	s.permissionHandler = handlers.NewPermissionHandler(s.permissionUsecase)
}

// SetupTest runs before each test
func (s *PermissionHandlerTestSuite) SetupTest() {
	// Fetch seed permissions for tests - just reference them, don't modify
	permissions, err := s.permissionRepo.List(context.Background())
	s.Require().NoError(err)
	s.Require().NotEmpty(permissions, "Seed permissions must exist in the database")
	s.seedPermissions = permissions
}

// TearDownTest cleans up after each test
func (s *PermissionHandlerTestSuite) TearDownTest() {
	// Clean up is handled by TearDownSuite
}

// TearDownSuite cleans up after all tests
func (s *PermissionHandlerTestSuite) TearDownSuite() {
	// No cleanup needed for seed permissions
}

// createAdminAuthContext creates a context that mimics an authenticated admin user
func (s *PermissionHandlerTestSuite) createAdminAuthContext() context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, middleware.UserIDKey, 1)
	ctx = context.WithValue(ctx, middleware.RoleCodeKey, string(models.RoleCodeAdmin))
	return ctx
}

// createRegularUserAuthContext creates a context that mimics a non-admin user
func (s *PermissionHandlerTestSuite) createRegularUserAuthContext() context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, middleware.UserIDKey, 2)
	ctx = context.WithValue(ctx, middleware.RoleCodeKey, string(models.RoleCodeNormalUser))
	return ctx
}

// TestListPermission tests the ListPermission handler
func (s *PermissionHandlerTestSuite) TestListPermission() {
	// Test with admin user - should succeed
	s.Run("AdminUser_Success", func() {
		// Create request
		req := httptest.NewRequest("GET", "/api/v1/admin/permissions", nil)
		req = req.WithContext(s.createAdminAuthContext())
		w := httptest.NewRecorder()

		// Execute
		s.permissionHandler.ListPermission(w, req)

		// Assert
		s.Equal(http.StatusOK, w.Code)

		var resp response.Response
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		s.NoError(err)
		s.True(resp.Success)

		// Verify data contains permissions
		data, ok := resp.Data.(map[string]any)
		s.True(ok)
		permissionsData, ok := data["permissions"].([]any)
		s.True(ok)
		s.GreaterOrEqual(len(permissionsData), 5) // We added 5 permissions
	})

	// Test with non-admin user - should fail with forbidden
	s.Run("RegularUser_Forbidden", func() {
		// Create request with non-admin context
		req := httptest.NewRequest("GET", "/api/v1/admin/permissions", nil)
		req = req.WithContext(s.createRegularUserAuthContext())
		w := httptest.NewRecorder()

		// Execute
		s.permissionHandler.ListPermission(w, req)

		// Assert - should be forbidden for non-admin
		s.Equal(http.StatusForbidden, w.Code)
	})

	// Test with no auth context - should fail with forbidden
	s.Run("NoContext_Forbidden", func() {
		// Create request with no auth context
		req := httptest.NewRequest("GET", "/api/v1/admin/permissions", nil)
		w := httptest.NewRecorder()

		// Execute
		s.permissionHandler.ListPermission(w, req)

		// Assert - should be forbidden since no auth context
		s.Equal(http.StatusForbidden, w.Code)
	})
}

// TestPermissionHandlerSuite runs the permission handler test suite
func TestPermissionHandlerSuite(t *testing.T) {
	suite.Run(t, new(PermissionHandlerTestSuite))
}
