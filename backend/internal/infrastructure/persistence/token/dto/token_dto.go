package dto

import (
	"time"

	"github.com/huydq/test/internal/domain/model/token"
	persistence "github.com/huydq/test/internal/infrastructure/persistence/util"
)

type TokenDTO struct {
	ID        int       `gorm:"column:id;primaryKey"`
	Token     string    `gorm:"column:token"`
	IsActive  bool      `gorm:"column:is_active"`
	ExpiredAt time.Time `gorm:"column:expired_at"`
	persistence.BaseColumnTimestamp
}

// TableName returns the table name for GORM
func (TokenDTO) TableName() string {
	return "token"
}

func (d *TokenDTO) ToTokenModel() *token.Token {
	if d == nil {
		return nil
	}

	model := &token.Token{
		ID:        d.ID,
		Token:     d.Token,
		IsActive:  d.IsActive,
		ExpiredAt: d.ExpiredAt,
	}

	return model
}

func ToTokenDTO(model *token.Token) *TokenDTO {
	if model == nil {
		return nil
	}

	dto := &TokenDTO{
		ID:        model.ID,
		Token:     model.Token,
		IsActive:  model.IsActive,
		ExpiredAt: model.ExpiredAt,
	}

	return dto
}
