package repository

import (
	"context"
	"github.com/aivyss/password-manager/pwmErr"
	"github.com/jmoiron/sqlx"
)

type PasswordListDetailRepository interface {
	Insert(ctx context.Context, passwordListKey int, detail string) error
	UpdateDescriptionByPasswordListKey(ctx context.Context, passwordListKey int, description string) error
}

func (r *passwordListDetailRepository) Insert(ctx context.Context, passwordListKey int, detail string) error {
	if _, err := GetStatement(ctx, r.insertDetailByPasswordListKey).ExecContext(ctx, map[string]any{
		"passwordListKey": passwordListKey,
		"description":     detail,
	}); err != nil {
		return pwmErr.InsertDB
	}

	return nil
}

type passwordListDetailRepository struct {
	db                            *sqlx.DB
	insertDetailByPasswordListKey *sqlx.NamedStmt
	updateDetailByPasswordListKey *sqlx.NamedStmt
}

func (r *passwordListDetailRepository) UpdateDescriptionByPasswordListKey(ctx context.Context, passwordListKey int, description string) error {
	if _, err := r.updateDetailByPasswordListKey.ExecContext(ctx, map[string]any{
		"passwordListKey": passwordListKey,
		"description":     description,
	}); err != nil {
		return pwmErr.FailUpdatePwDescription
	}

	return nil
}

func NewPasswordListDetailRepository(db *sqlx.DB) (PasswordListDetailRepository, error) {
	insertDetailByPasswordListKey, err := db.PrepareNamed(InsertDetailByPasswordListKey)
	if err != nil {
		return nil, err
	}
	updateDetailByPasswordListKey, err := db.PrepareNamed(UpdateDetailByPasswordListKey)
	if err != nil {
		return nil, err
	}

	return &passwordListDetailRepository{
		db:                            db,
		insertDetailByPasswordListKey: insertDetailByPasswordListKey,
		updateDetailByPasswordListKey: updateDetailByPasswordListKey,
	}, nil
}
