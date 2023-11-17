package repository

import (
	"context"
	"github.com/aivyss/password-manager/entity"
	"github.com/aivyss/password-manager/pwmErr"
	"github.com/jmoiron/sqlx"
)

const (
	InsertMasterUserKey = iota
	GetMasterUserByUserNameKey
)

type MasterUserRepository interface {
	Insert(ctx context.Context, userName string, userPassword string) error
	GetUserByUserName(ctx context.Context, userName string) (*entity.MasterUser, error)
}

type masterUserRepository struct {
	db       *sqlx.DB
	queryMap map[int]*sqlx.NamedStmt
}

func NewMasterUserRepository(db *sqlx.DB) (MasterUserRepository, error) {
	insertQuery, err := db.PrepareNamed(InsertMasterUser)
	if err != nil {
		return nil, err
	}
	getByUserName, err := db.PrepareNamed(GetMasterUserByUserName)
	if err != nil {
		return nil, err
	}

	queryMap := map[int]*sqlx.NamedStmt{}
	queryMap[InsertMasterUserKey] = insertQuery
	queryMap[GetMasterUserByUserNameKey] = getByUserName

	return &masterUserRepository{
		db:       db,
		queryMap: queryMap,
	}, nil
}

func (r *masterUserRepository) Insert(ctx context.Context, userName string, userPassword string) error {
	_, err := r.queryMap[InsertMasterUserKey].ExecContext(ctx, map[string]any{
		"userName":     userName,
		"userPassword": userPassword,
	})
	if err != nil {
		return pwmErr.InsertDB
	}

	return nil
}

func (r *masterUserRepository) GetUserByUserName(ctx context.Context, userName string) (*entity.MasterUser, error) {
	result := new(entity.MasterUser)

	if err := r.queryMap[GetMasterUserByUserNameKey].GetContext(ctx, result, map[string]any{
		"userName": userName,
	}); err != nil {
		return nil, pwmErr.NoUser
	}

	return result, nil
}
