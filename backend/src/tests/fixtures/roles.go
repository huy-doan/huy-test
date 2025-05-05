package fixtures

import (
	"context"
	"fmt"
	"time"

	"github.com/vnlab/makeshop-payment/src/domain/models"
	"github.com/vnlab/makeshop-payment/src/domain/repositories"
)

// GetMockRole returns a mock role with the specified ID and code
func GetMockRole(id int, code string) *models.Role {
	return &models.Role{
		ID:   id,
		Code: code,
		Name: getNameForRoleCode(code),
	}
}

// GetMockRoleWithPermissions returns a mock role with the specified permissions
func GetMockRoleWithPermissions(id int, code string, permissions []*models.Permission) *models.Role {
	role := GetMockRole(id, code)
	role.Permissions = permissions
	return role
}

// GetAdminRole returns the admin role from the database
// If it doesn't exist, it returns a mock instance
func GetAdminRole(ctx context.Context, repo repositories.RoleRepository) *models.Role {
	role, err := repo.FindByCode(ctx, string(models.RoleCodeAdmin))
	if err != nil || role == nil {
		// If there's an error or the role doesn't exist, return a mock
		adminPermission := GetMockPermission(1, "アドミン", "ADMIN")
		return GetMockRoleWithPermissions(
			1,
			string(models.RoleCodeAdmin),
			[]*models.Permission{adminPermission},
		)
	}
	return role
}

// GetNormalUserRole returns the normal user role from the database
// If it doesn't exist, it returns a mock instance
func GetNormalUserRole(ctx context.Context, repo repositories.RoleRepository) *models.Role {
	role, err := repo.FindByCode(ctx, string(models.RoleCodeNormalUser))
	if err != nil || role == nil {
		// If there's an error or the role doesn't exist, return a mock with appropriate permissions
		// Get permissions that would typically be assigned to a normal user
		permissions := []*models.Permission{
			GetMockPermission(5, "自分の個人データ変更", "EDIT_OWN_PROFILE"),
			GetMockPermission(6, "自分の行動ログ確認", "VIEW_OWN_LOG"),
		}
		return GetMockRoleWithPermissions(
			2,
			string(models.RoleCodeNormalUser),
			permissions,
		)
	}
	return role
}

// GetBusinessUserRole returns the business user role from the database
// If it doesn't exist, it returns a mock instance
func GetBusinessUserRole(ctx context.Context, repo repositories.RoleRepository) *models.Role {
	role, err := repo.FindByCode(ctx, string(models.RoleCodeBusinessUser))
	if err != nil || role == nil {
		// If there's an error or the role doesn't exist, return a mock with appropriate permissions
		permissions := []*models.Permission{
			GetMockPermission(5, "自分の個人データ変更", "EDIT_OWN_PROFILE"),
			GetMockPermission(8, "振込み承認（事業）", "TRANSFER_APPROVE_BUSINESS"),
			GetMockPermission(6, "自分の行動ログ確認", "VIEW_OWN_LOG"),
		}
		return GetMockRoleWithPermissions(
			3,
			string(models.RoleCodeBusinessUser),
			permissions,
		)
	}
	return role
}

// GetAccountingUserRole returns the accounting user role from the database
// If it doesn't exist, it returns a mock instance
func GetAccountingUserRole(ctx context.Context, repo repositories.RoleRepository) *models.Role {
	role, err := repo.FindByCode(ctx, string(models.RoleCodeAccoutingUser))
	if err != nil || role == nil {
		// If there's an error or the role doesn't exist, return a mock with appropriate permissions
		permissions := []*models.Permission{
			GetMockPermission(5, "自分の個人データ変更", "EDIT_OWN_PROFILE"),
			GetMockPermission(10, "手動振込機能", "MANUAL_TRANSFER"),
			GetMockPermission(9, "振込み承認（経理）", "TRANSFER_APPROVE_ACCOUNTANT"),
			GetMockPermission(6, "自分の行動ログ確認", "VIEW_OWN_LOG"),
		}
		return GetMockRoleWithPermissions(
			4,
			string(models.RoleCodeAccoutingUser),
			permissions,
		)
	}
	return role
}

// GetMockRoles returns all seeded roles as mock objects
func GetMockRoles() []*models.Role {
	return []*models.Role{
		GetMockRole(1, string(models.RoleCodeAdmin)),
		GetMockRole(2, string(models.RoleCodeNormalUser)),
		GetMockRole(3, string(models.RoleCodeBusinessUser)),
		GetMockRole(4, string(models.RoleCodeAccoutingUser)),
	}
}

// FindRoleByID retrieves a role from the database by ID
// If it doesn't exist, it returns a mock instance
func FindRoleByID(ctx context.Context, repo repositories.RoleRepository, id int) *models.Role {
	role, err := repo.FindByID(ctx, id)
	if err != nil || role == nil {
		// If there's an error or the role doesn't exist, return nil
		return nil
	}
	return role
}

// CreateUniqueTestRole creates a unique role for testing with a timestamp suffix
// to ensure it doesn't conflict with other tests
func CreateUniqueTestRole(ctx context.Context, repo repositories.RoleRepository, namePrefix string) (*models.Role, error) {
	// Create a unique name and code using current timestamp to avoid conflicts
	timestamp := time.Now().UnixNano()
	roleName := fmt.Sprintf("Test Role %s %d", namePrefix, timestamp)
	roleCode := fmt.Sprintf("TEST_ROLE_%s_%d", namePrefix, timestamp)

	role := &models.Role{
		Name: roleName,
		Code: roleCode,
	}

	err := repo.Create(ctx, role)
	if err != nil {
		return nil, err
	}

	return role, nil
}

// CreateUniqueTestRoleWithPermissions creates a unique role with the specified permissions
func CreateUniqueTestRoleWithPermissions(ctx context.Context, repo repositories.RoleRepository,
	namePrefix string, permissions []*models.Permission) (*models.Role, error) {

	// Create a unique name and code using current timestamp to avoid conflicts
	timestamp := time.Now().UnixNano()
	roleName := fmt.Sprintf("Test Role With Perms %s %d", namePrefix, timestamp)
	roleCode := fmt.Sprintf("TEST_ROLE_PERMS_%s_%d", namePrefix, timestamp)

	role := &models.Role{
		Name:        roleName,
		Code:        roleCode,
		Permissions: permissions,
	}

	err := repo.Create(ctx, role)
	if err != nil {
		return nil, err
	}

	return role, nil
}

// getNameForRoleCode returns the appropriate name for a given role code
// This should match the names in the seed files
func getNameForRoleCode(code string) string {
	switch code {
	case string(models.RoleCodeAdmin):
		return "システム管理者"
	case string(models.RoleCodeNormalUser):
		return "一般ユーザー"
	case string(models.RoleCodeBusinessUser):
		return "事業担当者"
	case string(models.RoleCodeAccoutingUser):
		return "経理担当者"
	default:
		return "Mock Role"
	}
}
