package repository

import (
	"github.com/aivyss/password-manager/pwmErr"
	"github.com/jmoiron/sqlx"
	"sync"
)

var factoryOnce sync.Once
var repositoryFactory *RepositoryFactory

type RepositoryFactory struct {
	MasterUserRepository   MasterUserRepository
	PasswordListRepository PasswordListRepository
}

func NewRepositoryFactory(db *sqlx.DB) (*RepositoryFactory, error) {
	var factory *RepositoryFactory
	var err error

	factoryOnce.Do(func() {
		masterUserRepo, e := NewMasterUserRepository(db)
		if e != nil {
			err = e
			return
		}

		passwordListRepo, e := NewPasswordListRepository(db)
		if err != nil {
			err = e
			return
		}

		factory = &RepositoryFactory{
			MasterUserRepository:   masterUserRepo,
			PasswordListRepository: passwordListRepo,
		}
		repositoryFactory = factory
	})

	return factory, err
}

func GetRepositoryFactory() (*RepositoryFactory, error) {
	if repositoryFactory == nil {
		return nil, pwmErr.Unknown
	}

	return repositoryFactory, nil
}
