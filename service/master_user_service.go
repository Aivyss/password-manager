package service

import (
	"context"
	"github.com/aivyss/password-manager/entity"
	"github.com/aivyss/password-manager/pwmErr"
	"github.com/aivyss/password-manager/repository"
	"golang.org/x/crypto/bcrypt"
	"sort"
)

type MasterUserService interface {
	CreateUser(name, password string) error
	Login(name, password string) (*entity.MasterUser, error)
	GetUsers(head, tail int) ([]entity.MasterUser, error)
}

type masterUserService struct {
	masterUserRepository repository.MasterUserRepository
}

func (s *masterUserService) GetUsers(head, tail int) ([]entity.MasterUser, error) {
	ctx := context.Background()
	users, err := s.masterUserRepository.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	sort.Slice(users, func(i, j int) bool {
		return users[i].Id < users[j].Id
	})

	var result []entity.MasterUser
	if head > 0 {
		result = users[:]
		if len(users) > 10 {
			result = users[:head]
		}
	} else if tail > 0 {
		result = users[:]
		if len(users)-tail >= 0 {
			result = users[len(users)-tail:]
		}
	} else {
		result = users[:]
		if len(users) > 10 {
			result = users[:10]
		}
	}

	return result, nil
}

func (s *masterUserService) Login(name, password string) (*entity.MasterUser, error) {
	ctx := context.Background()
	user, err := s.masterUserRepository.GetUserByUserName(ctx, name)
	if err != nil {
		return nil, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password+name)); err != nil {
		return nil, pwmErr.NoUser
	}

	return user, nil
}

func (s *masterUserService) CreateUser(name, password string) error {
	hashedPw, err := bcrypt.GenerateFromPassword([]byte(password+name), 15)
	if err != nil {
		return pwmErr.GeneratePw
	}

	ctx := context.Background()
	if err := s.masterUserRepository.Insert(ctx, name, string(hashedPw)); err != nil {
		return err
	}

	return nil
}

func NewMasterUserService(masterUserRepository repository.MasterUserRepository) MasterUserService {
	return &masterUserService{
		masterUserRepository: masterUserRepository,
	}
}
