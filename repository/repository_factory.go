package repository

import "github.com/jmoiron/sqlx"

type RepositoryFactory struct {
	MasterUserRepository MasterUserRepository
}

func NewRepositoryFactory(db *sqlx.DB) (*RepositoryFactory, error) {
	masterUserRepository, err := NewMasterUserRepository(db)
	if err != nil {
		return nil, err
	}

	return &RepositoryFactory{
		MasterUserRepository: masterUserRepository,
	}, nil
}
