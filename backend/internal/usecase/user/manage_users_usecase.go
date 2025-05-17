package user

import (
	"context"
	"errors"

	"github.com/huydq/test/internal/datastructure/inputdata"
	"github.com/huydq/test/internal/domain/model/user"
	passwordPkg "github.com/huydq/test/internal/domain/object/password"
	roleRepo "github.com/huydq/test/internal/domain/repository/role"
	userRepo "github.com/huydq/test/internal/domain/repository/user"
)

type UserManagementUsecase interface {
	GetUserByID(ctx context.Context, id int) (*user.User, error)
	GetUserByEmail(ctx context.Context, email string) (*user.User, error)
	ChangePassword(ctx context.Context, userID int, input *inputdata.ChangePasswordInputData) error
	ListUsers(ctx context.Context, input *inputdata.UserListInputData) ([]*user.User, int, int, error)
	CreateUser(ctx context.Context, input *inputdata.CreateUserInputData) (*user.User, error)
	UpdateUser(ctx context.Context, userID int, input *inputdata.UpdateUserInputData) (*user.User, error)
	DeleteUser(ctx context.Context, userID int) error
	ResetPassword(ctx context.Context, userID int, input *inputdata.ResetPasswordInputData) error
}

type ManageUsersUsecase struct {
	userRepo userRepo.UserRepository
	roleRepo roleRepo.RoleRepository
}

func NewManageUsersUsecase(
	userRepo userRepo.UserRepository,
	roleRepo roleRepo.RoleRepository,
) *ManageUsersUsecase {
	return &ManageUsersUsecase{
		userRepo: userRepo,
		roleRepo: roleRepo,
	}
}

// ChangePassword changes a user's password
func (uc *ManageUsersUsecase) ChangePassword(ctx context.Context, userID int, input *inputdata.ChangePasswordInputData) error {
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	if !user.VerifyPassword(input.CurrentPassword) {
		return errors.New("current password is incorrect")
	}

	if err := user.ChangePassword(input.NewPassword); err != nil {
		return err
	}

	return uc.userRepo.Update(ctx, user)
}

// ListUsers lists users
func (uc *ManageUsersUsecase) ListUsers(ctx context.Context, input *inputdata.UserListInputData) ([]*user.User, int, int, error) {
	const (
		defaultPage     = 1
		defaultPageSize = 10
	)

	if input.Page <= 0 {
		input.Page = defaultPage
	}

	if input.PageSize <= 0 {
		input.PageSize = defaultPageSize
	}

	return uc.userRepo.List(ctx, input)
}

// CreateUser creates a new user
func (uc *ManageUsersUsecase) CreateUser(ctx context.Context, input *inputdata.CreateUserInputData) (*user.User, error) {
	existingUser, err := uc.userRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("email.already_exists")
	}

	role, err := uc.roleRepo.FindByID(ctx, input.RoleID)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, errors.New("role.not_found")
	}

	hashedPassword, err := passwordPkg.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	newUser := &user.User{
		Email:        input.Email,
		PasswordHash: hashedPassword,
		RoleID:       input.RoleID,
		EnabledMFA:   input.EnabledMFA,
		FullName:     input.FullName,
	}

	if err := uc.userRepo.Create(ctx, newUser); err != nil {
		return nil, err
	}

	newUser.Role = role

	return newUser, nil
}

// UpdateUser updates a user
func (uc *ManageUsersUsecase) UpdateUser(ctx context.Context, userID int, input *inputdata.UpdateUserInputData) (*user.User, error) {
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("account.not_found")
	}

	if input.Email != nil && *input.Email != user.Email {
		existingUser, err := uc.userRepo.FindByEmail(ctx, *input.Email)
		if existingUser != nil && err == nil {
			return nil, errors.New("email.already_exists")
		}
	}

	if input.RoleID != nil {
		role, err := uc.roleRepo.FindByID(ctx, *input.RoleID)
		if role == nil && err == nil {
			return nil, errors.New("role.not_found")
		}
		user.Role = role
		user.RoleID = *input.RoleID

		if err != nil {
			return nil, err
		}
	}

	uc.updateUserFields(user, input)

	if err := uc.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// ResetPassword resets a user's password
func (uc *ManageUsersUsecase) ResetPassword(ctx context.Context, userID int, input *inputdata.ResetPasswordInputData) error {
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	if user == nil {
		return errors.New("account.not_found")
	}

	if err := user.ChangePassword(input.NewPassword); err != nil {
		return err
	}

	return uc.userRepo.Update(ctx, user)
}

// updateUserFields updates the fields of a user
func (uc *ManageUsersUsecase) updateUserFields(user *user.User, input *inputdata.UpdateUserInputData) {
	if input.FullName != nil {
		user.FullName = *input.FullName
	}
	if input.Email != nil {
		user.Email = *input.Email
	}
	if input.EnabledMFA != nil {
		user.EnabledMFA = *input.EnabledMFA
	}

	if input.Password != nil {
		hashedPassword, err := passwordPkg.HashPassword(*input.Password)
		if err != nil {
			return
		}
		user.PasswordHash = hashedPassword
	}
}

// GetUserByEmail returns a user by email
func (uc *ManageUsersUsecase) GetUserByEmail(ctx context.Context, email string) (*user.User, error) {
	return uc.userRepo.FindByEmail(ctx, email)
}

// DeleteUser deletes a user by ID
func (uc *ManageUsersUsecase) DeleteUser(ctx context.Context, userID int) error {
	return uc.userRepo.Delete(ctx, userID)
}

// GetUserByID returns a user by ID
func (uc *ManageUsersUsecase) GetUserByID(ctx context.Context, id int) (*user.User, error) {
	return uc.userRepo.FindByID(ctx, id)
}
