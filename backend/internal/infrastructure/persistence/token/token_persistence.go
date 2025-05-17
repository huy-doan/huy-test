package two_factor_token

import (
	"context"
	"errors"

	tokenModel "github.com/huydq/test/internal/domain/model/token"
	tokenRepo "github.com/huydq/test/internal/domain/repository/token"
	"github.com/huydq/test/internal/infrastructure/persistence/token/dto"
	"github.com/huydq/test/internal/pkg/database"
	"gorm.io/gorm"
)

type TokenRepositoryImpl struct {
	db *gorm.DB
}

func NewTokenRepository(db *gorm.DB) tokenRepo.TokenRepository {
	return &TokenRepositoryImpl{db: db}
}

func (r *TokenRepositoryImpl) Create(ctx context.Context, token *tokenModel.Token) error {
	db, err := database.GetTxOrDB(ctx)
	if err != nil {
		return err
	}
	tokenDTO := dto.ToTokenDTO(token)
	return db.Create(tokenDTO).Error
}

func (r *TokenRepositoryImpl) Update(ctx context.Context, token *tokenModel.Token) error {
	db, err := database.GetTxOrDB(ctx)
	if err != nil {
		return err
	}
	tokenDTO := dto.ToTokenDTO(token)
	return db.Model(tokenDTO).Updates(tokenDTO).Error
}

func (r *TokenRepositoryImpl) FindByToken(ctx context.Context, token string) (*tokenModel.Token, error) {
	db, err := database.GetTxOrDB(ctx)
	if err != nil {
		return nil, err
	}
	var tokenDTO dto.TokenDTO

	err = db.Where("token = ?", token).First(&tokenDTO).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return tokenDTO.ToTokenModel(), nil
}
