package dto

import (
	"github.com/huydq/test/internal/domain/model/user"
	roleConvert "github.com/huydq/test/internal/infrastructure/persistence/role/convert"
	roleDto "github.com/huydq/test/internal/infrastructure/persistence/role/dto"
	persistence "github.com/huydq/test/internal/infrastructure/persistence/util"
)

// User represents the data transfer object for user entities
type User struct {
	ID           int           `gorm:"column:id;primaryKey" json:"id"`
	Email        string        `gorm:"column:email" json:"email"`
	PasswordHash string        `gorm:"column:password_hash" json:"-"`
	RoleID       int           `gorm:"column:role_id" json:"role_id"`
	Role         *roleDto.Role `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	EnabledMFA   bool          `gorm:"column:enabled_mfa" json:"enabled_mfa"`
	MFAType      int           `gorm:"column:mfa_type" json:"mfa_type"`
	FullName     string        `gorm:"column:full_name" json:"full_name"`
	persistence.BaseColumnTimestamp
}

// TableName specifies the database table name
func (User) TableName() string {
	return "`user`"
}

// ToUserModel converts a User to a domain model User.
func (dto *User) ToUserModel() *user.User {
	userModel := &user.User{
		ID:           dto.ID,
		Email:        dto.Email,
		PasswordHash: dto.PasswordHash,
		FullName:     dto.FullName,
		EnabledMFA:   dto.EnabledMFA,
		MFAType:      dto.MFAType,
		RoleID:       dto.RoleID,
	}

	userModel.CreatedAt = dto.CreatedAt
	userModel.UpdatedAt = dto.UpdatedAt

	if dto.Role != nil {
		userModel.Role = roleConvert.ToRoleModel(dto.Role)
	}

	return userModel
}

// ToUserDTO converts a domain model User to a User.
func ToUserDTO(u *user.User) *User {
	userDTO := &User{
		ID:           u.ID,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		FullName:     u.FullName,
		EnabledMFA:   u.EnabledMFA,
		MFAType:      u.MFAType,
		RoleID:       u.RoleID,
	}

	userDTO.CreatedAt = u.CreatedAt
	userDTO.UpdatedAt = u.UpdatedAt

	if u.Role != nil {
		userDTO.Role = roleConvert.ToRoleDTO(u.Role)
	}

	return userDTO
}
