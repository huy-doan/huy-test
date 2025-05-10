package auth

import (
	"context"
	"errors"

	"github.com/huydq/test/internal/datastructure/inputdata"
	"github.com/huydq/test/internal/datastructure/outputdata"
	"github.com/huydq/test/internal/domain/model/user"
	userRepo "github.com/huydq/test/internal/domain/repository/user"
	"github.com/huydq/test/internal/infrastructure/adapter/auth"
)

type AuthUsecase struct {
	userRepo   userRepo.UserRepository
	jwtService *auth.JWTService
}

func NewAuthUsecase(
	userRepo userRepo.UserRepository,
	jwtService *auth.JWTService,
) *AuthUsecase {
	return &AuthUsecase{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

func (uc *AuthUsecase) Login(ctx context.Context, input *inputdata.LoginInputData) (*outputdata.LoginOutputData, error) {
	user, err := uc.userRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}

	if user == nil || !user.VerifyPassword(input.Password) {
		return nil, errors.New("login.failed")
	}

	token, err := uc.jwtService.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	return &outputdata.LoginOutputData{
		Token: token,
		User:  user,
	}, nil
}

// GetMe retrieves a user by their ID
func (uc *AuthUsecase) GetMe(ctx context.Context, userID int) (*outputdata.UserProfileOutputData, error) {
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user.not_found")
	}

	return &outputdata.UserProfileOutputData{
		User: user,
	}, nil
}

// FindUserByEmail finds a user by their email address
func (uc *AuthUsecase) FindUserByEmail(ctx context.Context, email string) (*user.User, error) {
	return uc.userRepo.FindByEmail(ctx, email)
}
