package repository

import (
	"context"
	"github.com/aivyss/password-manager/entity"
	"github.com/aivyss/password-manager/pwmErr"
	"github.com/jmoiron/sqlx"
)

const (
	InsertPasswordKey = iota
	GetPasswordByUserPkAndKeyKey
	UpdatePasswordByUserPkAndKeyKey
	GetAllPasswordsKey
)

type PasswordListRepository interface {
	Insert(ctx context.Context, userPk int, key string, password string) error
	GetPasswordByUserPkAndKey(ctx context.Context, userPk int, key string) (entity.PasswordList, error)
	GetAllPasswords(ctx context.Context, userPk int) ([]entity.PasswordListWithDescription, error)
	UpdatePasswordByUserPkAndKey(ctx context.Context, userPk int, key string, password string) error
}

type passwordListRepository struct {
	db       *sqlx.DB
	queryMap map[int]*sqlx.NamedStmt
}

func (p *passwordListRepository) Insert(ctx context.Context, userPk int, key string, password string) error {
	if _, err := GetStatement(ctx, p.queryMap[InsertPasswordKey]).ExecContext(ctx, map[string]any{
		"userPk":   userPk,
		"key":      key,
		"password": password,
	}); err != nil {
		return pwmErr.InsertDB
	}

	return nil
}

func (p *passwordListRepository) GetPasswordByUserPkAndKey(ctx context.Context, userPk int, key string) (entity.PasswordList, error) {
	var result entity.PasswordList
	password := new(entity.PasswordList)

	if err := p.queryMap[GetPasswordByUserPkAndKeyKey].GetContext(ctx, password, map[string]any{
		"userPk": userPk,
		"key":    key,
	}); err != nil {
		return result, pwmErr.NotRegisteredKey
	}
	result = *password

	return result, nil
}

func (p *passwordListRepository) GetAllPasswords(ctx context.Context, userPk int) ([]entity.PasswordListWithDescription, error) {
	var results []entity.PasswordListWithDescription
	if err := p.queryMap[GetAllPasswordsKey].SelectContext(ctx, &results, map[string]any{
		"userPk": userPk,
	}); err != nil {
		return nil, pwmErr.NoRecord
	}

	return results, nil
}

func (p *passwordListRepository) UpdatePasswordByUserPkAndKey(ctx context.Context, userPk int, key string, password string) error {
	if _, err := p.queryMap[UpdatePasswordByUserPkAndKeyKey].ExecContext(ctx, map[string]any{
		"userPk":   userPk,
		"key":      key,
		"password": password,
	}); err != nil {
		return pwmErr.FailUpdatePw
	}

	return nil
}

func NewPasswordListRepository(db *sqlx.DB) (PasswordListRepository, error) {
	insertPasswordQuery, err := db.PrepareNamed(InsertPassword)
	if err != nil {
		return nil, err
	}
	getPasswordByUserPkAndKeyQuery, err := db.PrepareNamed(GetPasswordByUserPkAndKey)
	if err != nil {
		return nil, err
	}
	updatePasswordByUserPkAndKeyQuery, err := db.PrepareNamed(UpdatePasswordByUserPkAndKey)
	if err != nil {
		return nil, err
	}
	getAllPasswordsQuery, err := db.PrepareNamed(GetAllPasswords)
	if err != nil {
		return nil, err
	}

	return &passwordListRepository{
		db: db,
		queryMap: map[int]*sqlx.NamedStmt{
			InsertPasswordKey:               insertPasswordQuery,
			GetPasswordByUserPkAndKeyKey:    getPasswordByUserPkAndKeyQuery,
			UpdatePasswordByUserPkAndKeyKey: updatePasswordByUserPkAndKeyQuery,
			GetAllPasswordsKey:              getAllPasswordsQuery,
		},
	}, nil
}
