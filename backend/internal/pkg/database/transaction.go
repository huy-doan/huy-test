package database

import (
	"context"
	"errors"

	"github.com/huydq/test/internal/pkg/logger"
	"gorm.io/gorm"
)

// Transaction is an interface for database transactions
type Transaction[T any] interface {
	Transact(ctx context.Context, f func(ctx context.Context) (T, error)) (result T, err error)
}

// txKey is used to store the transaction in the context
type txKey struct{}

// tx implements Transaction interface
type tx[T any] struct {
	db *gorm.DB
}

// NewTransaction creates a new transaction
func NewTransaction[T any](db *gorm.DB) Transaction[T] {
	return &tx[T]{db: db}
}

// Transact executes a transaction
func (t *tx[T]) Transact(ctx context.Context, f func(ctx context.Context) (T, error)) (result T, err error) {
	tx := t.db.Begin()
	err = tx.Error
	if err != nil {
		var zero T

		return zero, err
	}

	ctx = context.WithValue(ctx, txKey{}, tx)

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			logger.GetLogger().Error("panic Rollback: ", map[string]any{
				"error": p,
			})
			panic(p)
		} else if err != nil {
			tx.Rollback()
			logger.GetLogger().Error("Rollback: ", map[string]any{
				"error": err,
			})
		} else if err = ctx.Err(); err != nil {
			logger.GetLogger().Error("context error occurred: ", map[string]any{
				"error": err.Error(),
			})
		} else if err = tx.Commit().Error; err != nil {
			tx.Rollback()
			logger.GetLogger().Error("Commit error: ", map[string]any{
				"error": err.Error(),
			})
		}
	}()

	v, err := f(ctx)
	if err != nil {
		var zero T
		return zero, err
	}

	return v, nil
}

// NewTx creates a new transaction
func NewTx[T any](ctx context.Context) (Transaction[T], error) {
	conn, err := GetDB(ctx)
	if err != nil {
		return nil, err
	}

	return NewTransaction[T](conn), nil
}

// GetTxOrDB returns the transaction or the database connection
func GetTxOrDB(ctx context.Context) (*gorm.DB, error) {
	tx, ok := getTx(ctx)
	if ok {
		return tx, nil
	}

	return GetDB(ctx)
}

// GetTx returns the transaction
func GetTx(ctx context.Context) (*gorm.DB, error) {
	tx, ok := getTx(ctx)
	if !ok {
		return nil, errors.New("tx not found")
	}
	return tx, nil
}

// getTx returns the transaction from the context
func getTx(ctx context.Context) (*gorm.DB, bool) {
	if tx, ok := ctx.Value(txKey{}).(*gorm.DB); ok {
		return tx, ok
	}
	return nil, false
}
