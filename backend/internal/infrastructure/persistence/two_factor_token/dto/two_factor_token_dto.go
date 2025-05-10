package dto

import (
	"time"

	"github.com/huydq/test/internal/domain/model/two_factor_token"
	"github.com/huydq/test/internal/infrastructure/persistence/user/dto"
	persistence "github.com/huydq/test/internal/infrastructure/persistence/util"
)

// TwoFactorTokenDTO is the data transfer object for two factor tokens
type TwoFactorTokenDTO struct {
	ID        int          `gorm:"column:id;primaryKey"`
	UserID    int          `gorm:"column:user_id"`
	Token     string       `gorm:"column:token"`
	MFAType   int          `gorm:"column:mfa_type"`
	User      *dto.UserDTO `gorm:"foreignKey:UserID"`
	IsUsed    bool         `gorm:"column:is_used"`
	ExpiredAt time.Time    `gorm:"column:expired_at"`
	persistence.BaseColumnTimestamp
}

// TableName returns the table name for GORM
func (TwoFactorTokenDTO) TableName() string {
	return "two_factor_token"
}

// ToTwoFactorTokenModel converts DTO to domain model
func (d *TwoFactorTokenDTO) ToTwoFactorTokenModel() *two_factor_token.TwoFactorToken {
	if d == nil {
		return nil
	}

	model := &two_factor_token.TwoFactorToken{
		ID:        d.ID,
		UserID:    d.UserID,
		Token:     d.Token,
		MFAType:   d.MFAType,
		IsUsed:    d.IsUsed,
		ExpiredAt: d.ExpiredAt,
	}

	if d.User != nil {
		model.User = d.User.ToUserModel()
	}

	return model
}

// ToTwoFactorTokenDTO converts domain model to DTO
func ToTwoFactorTokenDTO(model *two_factor_token.TwoFactorToken) *TwoFactorTokenDTO {
	if model == nil {
		return nil
	}

	dto := &TwoFactorTokenDTO{
		ID:        model.ID,
		UserID:    model.UserID,
		Token:     model.Token,
		MFAType:   model.MFAType,
		IsUsed:    model.IsUsed,
		ExpiredAt: model.ExpiredAt,
	}

	return dto
}
