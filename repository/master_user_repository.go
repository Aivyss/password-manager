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
	GetMasterUserByIdKey
	GetAllUsersKey
)

type MasterUserRepository interface {
	Insert(ctx context.Context, userName string, userPassword string) error
	GetUserByUserName(ctx context.Context, userName string) (*entity.MasterUser, error)
	GetUserById(ctx context.Context, id int) (*entity.MasterUser, error)
	GetAllUsers(ctx context.Context) ([]entity.MasterUser, error)
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
	getUserById, err := db.PrepareNamed(GetMasterUserById)
	if err != nil {
		return nil, err
	}
	getAllUsers, err := db.PrepareNamed(GetAllUsers)
	if err != nil {
		return nil, err
	}

	return &masterUserRepository{
		db: db,
		queryMap: map[int]*sqlx.NamedStmt{
			InsertMasterUserKey:        insertQuery,
			GetMasterUserByUserNameKey: getByUserName,
			GetMasterUserByIdKey:       getUserById,
			GetAllUsersKey:             getAllUsers,
		},
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

func (r *masterUserRepository) GetUserById(ctx context.Context, id int) (*entity.MasterUser, error) {
	user := new(entity.MasterUser)
	ctx = context.Background()

	if err := r.queryMap[GetMasterUserByIdKey].GetContext(ctx, user, map[string]any{
		"id": id,
	}); err != nil {
		return nil, pwmErr.NoUser
	}

	return user, nil
}

func (r *masterUserRepository) GetAllUsers(ctx context.Context) ([]entity.MasterUser, error) {
	var users []entity.MasterUser
	if err := r.queryMap[GetAllUsersKey].SelectContext(ctx, &users, map[string]any{}); err != nil {
		return nil, pwmErr.NoUser
	}

	return users, nil
}
