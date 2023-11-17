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
)

type PasswordListRepository interface {
	Insert(ctx context.Context, userPk int, key string, password string) error
	GetPasswordByUserPkAndKey(ctx context.Context, userPk int, key string) (*entity.PasswordList, error)
}

type passwordListRepository struct {
	db       *sqlx.DB
	queryMap map[int]*sqlx.NamedStmt
}

func (p *passwordListRepository) Insert(ctx context.Context, userPk int, key string, password string) error {
	if _, err := p.queryMap[InsertPasswordKey].ExecContext(ctx, map[string]any{
		"userPk":   userPk,
		"key":      key,
		"password": password,
	}); err != nil {
		return pwmErr.InsertDB
	}

	return nil
}

func (p *passwordListRepository) GetPasswordByUserPkAndKey(ctx context.Context, userPk int, key string) (*entity.PasswordList, error) {
	password := new(entity.PasswordList)

	if err := p.queryMap[GetPasswordByUserPkAndKeyKey].GetContext(ctx, password, map[string]any{
		"userPk": userPk,
		"key":    key,
	}); err != nil {
		return nil, pwmErr.NotRegisteredKey
	}

	return password, nil
}

func NewPasswordListRepository(db *sqlx.DB) (PasswordListRepository, error) {
	queryMap := map[int]*sqlx.NamedStmt{}
	insertPasswordQuery, err := db.PrepareNamed(InsertPassword)
	if err != nil {
		return nil, err
	}
	getPasswordByUserPkAndKeyQuery, err := db.PrepareNamed(GetPasswordByUserPkAndKey)
	if err != nil {
		return nil, err
	}

	queryMap[InsertPasswordKey] = insertPasswordQuery
	queryMap[GetPasswordByUserPkAndKeyKey] = getPasswordByUserPkAndKeyQuery

	return &passwordListRepository{
		db:       db,
		queryMap: queryMap,
	}, nil
}
