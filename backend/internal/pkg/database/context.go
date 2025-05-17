package database

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type contextKey string

const (
	ctxDBKey contextKey = "ctxDBKey"
)

// Get DB from context
func GetDB(ctx context.Context) (*gorm.DB, error) {
	val := ctx.Value(ctxDBKey)
	db, ok := val.(gorm.DB)
	if !ok {
		return nil, errors.New("Couldn't get gorm.DB from context")
	}

	return &db, nil
}

// Set db to context
func SetDB(ctx context.Context, db *gorm.DB) (context.Context, error) {
	if db == nil {
		return nil, errors.New("gorm.DBがnilです")
	}
	ctx = context.WithValue(ctx, ctxDBKey, *db)

	return ctx, nil
}
