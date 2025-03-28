package usecase

import (
	"context"
	"errors"

	models "github.com/huydq/ddd-project/src/domain/models"
	"github.com/huydq/ddd-project/src/domain/repositories"
	"github.com/huydq/ddd-project/src/infrastructure/auth"
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
	Email         string `json:"email" binding:"required,email"`
	Password      string `json:"password" binding:"required,min=6"`
	FirstName     string `json:"first_name" binding:"required"`
	LastName      string `json:"last_name" binding:"required"`
	FirstNameKana string `json:"first_name_kana" binding:"required"`
	LastNameKana  string `json:"last_name_kana" binding:"required"`
}

// UpdateProfileRequest represents a profile update request
type UpdateProfileRequest struct {
	FirstName     string `json:"first_name" binding:"required"`
	LastName      string `json:"last_name" binding:"required"`
	FirstNameKana string `json:"first_name_kana" binding:"required"`
	LastNameKana  string `json:"last_name_kana" binding:"required"`
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

	if user == nil {
		return nil, errors.New("invalid email or password 1")
	}

	if !user.VerifyPassword(req.Password) {
		return nil, errors.New("invalid email or password 2")
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

// Register creates a new user
func (uc *UserUsecase) Register(ctx context.Context, req RegisterRequest) (*models.User, error) {
	// Check if email already exists
	existingUser, err := uc.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("email already exists")
	}

	// Get customer role
	customerRole, err := uc.roleRepo.FindByCode(ctx, string(models.RoleCodeNormalUser))
	if err != nil {
		return nil, err
	}
	if customerRole == nil {
		return nil, errors.New("customer role not found")
	}

	// Create new user with customer role
	user, err := models.NewUser(
		req.Email,
		req.Password,
		req.FirstName,
		req.LastName,
		req.FirstNameKana,
		req.LastNameKana,
		customerRole.ID,
	)
	if err != nil {
		return nil, err
	}

	// Save user to database
	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Reload user to get the role relationship
	user, err = uc.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

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

	if err := user.UpdateProfile(req.FirstName, req.LastName, req.FirstNameKana, req.LastNameKana); err != nil {
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
func (uc *UserUsecase) ListUsers(ctx context.Context, page, pageSize int) ([]*models.User, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	return uc.userRepo.List(ctx, page, pageSize)
}

// GetJWTService returns the JWT service
func (uc *UserUsecase) GetJWTService() *auth.JWTService {
	return uc.jwtService
}
