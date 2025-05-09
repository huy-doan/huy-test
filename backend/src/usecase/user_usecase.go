package usecase

import (
	"context"
	"errors"

	validator "github.com/huydq/test/src/api/http/validator/user"
	models "github.com/huydq/test/src/domain/models"
	"github.com/huydq/test/src/domain/repositories"
	"github.com/huydq/test/src/infrastructure/auth"
	"golang.org/x/crypto/bcrypt"
)

// UserUsecase handles user-related business logic
type UserUsecase struct {
	userRepo   repositories.UserRepository
	roleRepo   repositories.RoleRepository
	jwtService *auth.JWTService
}

// NewUserUseCase creates a new UserUsecase
func NewUserUseCase(
	userRepo repositories.UserRepository,
	roleRepo repositories.RoleRepository,
	jwtService *auth.JWTService,
) *UserUsecase {
	return &UserUsecase{
		userRepo:   userRepo,
		roleRepo:   roleRepo,
		jwtService: jwtService,
	}
}

// LoginRequest represents a login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest represents a user registration request
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
}

// UpdateProfileRequest represents a profile update request
type UpdateProfileRequest struct {
	FullName string `json:"full_name" binding:"required"`
}

// LoginResponse represents a login response with token
type LoginResponse struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}

// Login authenticates a user and returns a JWT token
func (uc *UserUsecase) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	user, err := uc.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if user == nil || !user.VerifyPassword(req.Password) {
		return nil, errors.New("login.failed")
	}

	token, err := uc.jwtService.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token: token,
		User:  user,
	}, nil
}

// Register handles user registration
func (uc *UserUsecase) Register(ctx context.Context, req RegisterRequest) (*models.User, error) {
	// Check if email already exists
	existingUser, err := uc.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("email.already_exists")
	}

	// Get the normal user role
	role, err := uc.roleRepo.FindByCode(ctx, string(models.RoleCodeNormalUser))
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, errors.New("role.not_found")
	}

	// Create new user
	user, err := models.NewUser(req.Email, req.Password, req.FullName, role.ID)
	if err != nil {
		return nil, err
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Set the Role for the response
	user.Role = role

	return user, nil
}

// GetUserByID retrieves a user by ID
func (uc *UserUsecase) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	return uc.userRepo.FindByID(ctx, id)
}

// UpdateUserProfile updates a user's profile
func (uc *UserUsecase) UpdateUserProfile(ctx context.Context, userID int, req UpdateProfileRequest) (*models.User, error) {
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	if err := user.UpdateProfile(req.FullName); err != nil {
		return nil, err
	}

	if err := uc.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// ChangePassword changes a user's password
func (uc *UserUsecase) ChangePassword(ctx context.Context, userID int, currentPassword, newPassword string) error {
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	if !user.VerifyPassword(currentPassword) {
		return errors.New("current password is incorrect")
	}

	if err := user.ChangePassword(newPassword); err != nil {
		return err
	}

	return uc.userRepo.Update(ctx, user)
}

// ListUsers lists users with pagination
func (uc *UserUsecase) ListUsers(ctx context.Context, filter validator.UserListFilter) ([]*models.User, int, int, error) {
	const (
		defaultPage     = 1
		defaultPageSize = 10
	)

	if filter.Page <= 0 {
		filter.Page = defaultPage
	}

	if filter.PageSize <= 0 {
		filter.PageSize = defaultPageSize
	}

	repoParams := validator.UserListFilter{
		Page:      filter.Page,
		PageSize:  filter.PageSize,
		Search:    filter.Search,
		RoleID:    filter.RoleID,
		SortField: filter.SortField,
		SortOrder: filter.SortOrder,
	}

	return uc.userRepo.List(ctx, repoParams)
}

// GetJWTService returns the JWT service
func (uc *UserUsecase) GetJWTService() *auth.JWTService {
	return uc.jwtService
}

// UpdateUser updates a user's profile by admin
func (u *UserUsecase) UpdateUser(ctx context.Context, userID int, req validator.UpdateUserRequest) (*models.User, error) {
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("account.not_found")
	}

	if req.Email != nil && *req.Email != user.Email {
		existingUser, err := u.userRepo.FindByEmail(ctx, *req.Email)
		if existingUser != nil && err == nil {
			return nil, errors.New("email.already_exists")
		}
	}

	if req.RoleID != nil {
		role, err := u.roleRepo.FindByID(ctx, *req.RoleID)
		if role == nil && err == nil {
			return nil, errors.New("role.not_found")
		}
		user.Role = role

		if err != nil {
			return nil, err
		}
	}

	// Update fields only if admin provided
	u.updateUserFields(user, req)
	if err := u.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// updateUserFields updates the user fields based on the request
func (u *UserUsecase) updateUserFields(user *models.User, req validator.UpdateUserRequest) {
	if req.FullName != nil {
		user.FullName = *req.FullName
	}
	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.EnabledMFA != nil {
		user.EnabledMFA = *req.EnabledMFA
	}

	if req.Password != nil {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)

		if err != nil {
			return
		}

		user.PasswordHash = string(hashedPassword)
	}
}

// CreateUser creates a new user with the given data
func (u *UserUsecase) CreateUser(ctx context.Context, req *validator.CreateUserRequest) (*models.User, error) {
	existingUser, err := u.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("email.already_exists")
	}

	role, err := u.roleRepo.FindByID(ctx, req.RoleID)
	if err != nil {
		return nil, err
	}

	if role == nil {
		return nil, errors.New("role.not_found")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	newUser := &models.User{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		RoleID:       req.RoleID,
		EnabledMFA:   req.EnabledMFA,
		FullName:     req.FullName,
	}

	if err := u.userRepo.Create(ctx, newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}

func (u *UserUsecase) DeleteUser(ctx context.Context, userID int) error {
	return u.userRepo.Delete(ctx, userID)
}

// ResetPassword resets a user's password by user ID (admin only)
func (u *UserUsecase) ResetPassword(ctx context.Context, userID int, newPassword string) error {
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	if user == nil {
		return errors.New("account.not_found")
	}

	if err := user.ChangePassword(newPassword); err != nil {
		return err
	}

	return u.userRepo.Update(ctx, user)
}

// GetUserByEmail retrieves a user by email
func (uc *UserUsecase) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return uc.userRepo.FindByEmail(ctx, email)
}
