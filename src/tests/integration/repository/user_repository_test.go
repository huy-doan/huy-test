package repository_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/huydq/test/src/domain/models"
	"github.com/huydq/test/src/domain/repositories"
	repo "github.com/huydq/test/src/infrastructure/persistence/repositories"
	"github.com/huydq/test/src/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
)

type UserRepositoryTestSuite struct {
	tests.TestSuite
	userRepo repositories.UserRepository
	roleRepo repositories.RoleRepository
	testRole *models.Role
}

func (suite *UserRepositoryTestSuite) SetupSuite() {
	// Call the parent SetupSuite method
	suite.TestSuite.SetupSuite()

	// Initialize repositories
	suite.userRepo = repo.NewUserRepository(suite.DB)
	suite.roleRepo = repo.NewRoleRepository(suite.DB)
}

func (suite *UserRepositoryTestSuite) SetupTest() {
	// Get the normal user role from the database (already seeded)
	ctx := context.Background()
	role, err := suite.roleRepo.FindByCode(ctx, string(models.RoleCodeNormalUser))
	if err != nil {
		suite.T().Fatalf("Failed to find normal user role: %v", err)
	}
	if role == nil {
		suite.T().Fatal("Normal user role not found in database, check seeds")
	}
	suite.testRole = role
}

func (suite *UserRepositoryTestSuite) TearDownTest() {
	// Clean up users after each test but keep the roles
	suite.DB.Exec("DELETE FROM users")
}

func (suite *UserRepositoryTestSuite) TestCreate() {
	ctx := context.Background()

	// Create a test user
	password := "testpassword"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := &models.User{
		Email:         "test@example.com",
		PasswordHash:  string(hashedPassword),
		RoleID:        suite.testRole.ID,
		FirstName:     "Test",
		LastName:      "User",
		FirstNameKana: "テスト",
		LastNameKana:  "ユーザー",
	}

	// Test creating a user
	err := suite.userRepo.Create(ctx, user)

	// Assertions
	assert.NoError(suite.T(), err)
	assert.NotZero(suite.T(), user.ID, "User ID should be set after creation")

	// Verify user was created in the database
	var count int64
	suite.DB.Model(&models.User{}).Where("email = ?", user.Email).Count(&count)
	assert.Equal(suite.T(), int64(1), count, "User should exist in database")
}

func (suite *UserRepositoryTestSuite) TestFindByID() {
	ctx := context.Background()

	// Create a test user
	password := "testpassword"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	originalUser := &models.User{
		Email:         "find-by-id@example.com",
		PasswordHash:  string(hashedPassword),
		RoleID:        suite.testRole.ID,
		FirstName:     "Find",
		LastName:      "ByID",
		FirstNameKana: "ファインド",
		LastNameKana:  "バイアイディ",
	}

	// Insert user directly for test setup
	suite.DB.Create(originalUser)

	// Test finding by ID
	foundUser, err := suite.userRepo.FindByID(ctx, originalUser.ID)

	// Assertions
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), foundUser)
	assert.Equal(suite.T(), originalUser.ID, foundUser.ID)
	assert.Equal(suite.T(), originalUser.Email, foundUser.Email)
	assert.Equal(suite.T(), originalUser.FirstName, foundUser.FirstName)
	assert.Equal(suite.T(), originalUser.LastName, foundUser.LastName)
	assert.Equal(suite.T(), originalUser.RoleID, foundUser.RoleID)
	assert.NotNil(suite.T(), foundUser.Role, "Role should be loaded")
}

func (suite *UserRepositoryTestSuite) TestFindByEmail() {
	ctx := context.Background()

	// Create a test user
	email := "find-by-email@example.com"
	password := "testpassword"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	originalUser := &models.User{
		Email:         email,
		PasswordHash:  string(hashedPassword),
		RoleID:        suite.testRole.ID,
		FirstName:     "Find",
		LastName:      "ByEmail",
		FirstNameKana: "ファインド",
		LastNameKana:  "バイイーメール",
	}

	// Insert user directly for test setup
	suite.DB.Create(originalUser)

	// Test finding by email
	foundUser, err := suite.userRepo.FindByEmail(ctx, email)

	// Assertions
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), foundUser)
	assert.Equal(suite.T(), originalUser.ID, foundUser.ID)
	assert.Equal(suite.T(), email, foundUser.Email)
	assert.Equal(suite.T(), originalUser.FirstName, foundUser.FirstName)
	assert.Equal(suite.T(), originalUser.LastName, foundUser.LastName)
	assert.NotNil(suite.T(), foundUser.Role, "Role should be loaded")
}

func (suite *UserRepositoryTestSuite) TestUpdate() {
	ctx := context.Background()

	// Create a test user
	email := "update-test@example.com"
	password := "testpassword"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := &models.User{
		Email:         email,
		PasswordHash:  string(hashedPassword),
		RoleID:        suite.testRole.ID,
		FirstName:     "Original",
		LastName:      "Name",
		FirstNameKana: "オリジナル",
		LastNameKana:  "ネーム",
	}

	// Insert user directly for test setup
	suite.DB.Create(user)

	// Update user details
	user.FirstName = "Updated"
	user.LastName = "User"
	user.FirstNameKana = "アップデート"
	user.LastNameKana = "ユーザー"

	// Test updating the user
	err := suite.userRepo.Update(ctx, user)

	// Assertions
	assert.NoError(suite.T(), err)

	// Verify user was updated in the database
	var updatedUser models.User
	suite.DB.First(&updatedUser, user.ID)
	assert.Equal(suite.T(), "Updated", updatedUser.FirstName)
	assert.Equal(suite.T(), "User", updatedUser.LastName)
	assert.Equal(suite.T(), "アップデート", updatedUser.FirstNameKana)
	assert.Equal(suite.T(), "ユーザー", updatedUser.LastNameKana)
}

func (suite *UserRepositoryTestSuite) TestDelete() {
	ctx := context.Background()

	// Create a test user
	email := "delete-test@example.com"
	password := "testpassword"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := &models.User{
		Email:         email,
		PasswordHash:  string(hashedPassword),
		RoleID:        suite.testRole.ID,
		FirstName:     "Delete",
		LastName:      "User",
		FirstNameKana: "デリート",
		LastNameKana:  "ユーザー",
	}

	// Insert user directly for test setup
	suite.DB.Create(user)

	// Test deleting the user
	err := suite.userRepo.Delete(ctx, user.ID)

	// Assertions
	assert.NoError(suite.T(), err)

	// Verify user is not found after deletion
	var count int64
	suite.DB.Model(&models.User{}).Where("id = ?", user.ID).Count(&count)
	assert.Equal(suite.T(), int64(0), count, "User should be completely deleted")

	// Verify user is not found with repository method
	foundUser, err := suite.userRepo.FindByID(ctx, user.ID)
	assert.NoError(suite.T(), err)
	assert.Nil(suite.T(), foundUser, "Deleted user should not be found")
}

func (suite *UserRepositoryTestSuite) TestList() {
	ctx := context.Background()

	// Create multiple test users for pagination testing
	password := "testpassword"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	for i := 1; i <= 15; i++ {
		user := &models.User{
			Email:         fmt.Sprintf("user%d@example.com", i),
			PasswordHash:  string(hashedPassword),
			RoleID:        suite.testRole.ID,
			FirstName:     fmt.Sprintf("User%d", i),
			LastName:      "Test",
			FirstNameKana: "ユーザー",
			LastNameKana:  "テスト",
		}
		suite.DB.Create(user)
	}

	// Test first page (5 users per page)
	users, totalPages, err := suite.userRepo.List(ctx, 1, 5)

	// Assertions for first page
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), users, 5)
	assert.Equal(suite.T(), 3, totalPages) // 15 users / 5 per page = 3 pages

	// Test second page
	usersPage2, _, err := suite.userRepo.List(ctx, 2, 5)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), usersPage2, 5)

	// Check that users on page 2 are different from users on page 1
	assert.NotEqual(suite.T(), users[0].ID, usersPage2[0].ID)
}

func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}
