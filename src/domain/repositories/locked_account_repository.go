package repositories

import (
	"context"

	"github.com/vnlab/makeshop-payment/src/domain/models"
)

// LockedAccountRepository defines the interface for locked account operations in domain layer
type LockedAccountRepository interface {
	Create(ctx context.Context, account *models.LockedAccount) error

	Update(ctx context.Context, account *models.LockedAccount) error

	GetByEmail(ctx context.Context, email string) (*models.LockedAccount, error)

	GetByUserID(ctx context.Context, userID int) (*models.LockedAccount, error)

	GetByID(ctx context.Context, id int) (*models.LockedAccount, error)

	List(ctx context.Context, page, pageSize int) ([]*models.LockedAccount, int, error)
}
