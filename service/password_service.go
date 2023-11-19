package service

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"github.com/aivyss/password-manager/entity"
	"github.com/aivyss/password-manager/pwmErr"
	"github.com/aivyss/password-manager/repository"
	"github.com/aivyss/typex/util"
	"golang.org/x/crypto/bcrypt"
	"io"
)

type PasswordService interface {
	SetPassword(key, plainPw, description string) error
	GetPassword(key string) (*string, error)
	UpdatePassword(key, plainPwForPersist, plainUserPw string) error
	GetAllPasswords() ([]entity.PasswordList, error)
}

type passwordListService struct {
	masterUserRepository         repository.MasterUserRepository
	passwordListRepository       repository.PasswordListRepository
	passwordListDetailRepository repository.PasswordListDetailRepository
	txManager                    repository.TxManager
	userPk                       int
	cryptoKey                    []byte
}

func (s *passwordListService) GetAllPasswords() ([]entity.PasswordList, error) {
	ctx := context.Background()
	return s.passwordListRepository.GetAllPasswords(ctx, s.userPk)
}

func (s *passwordListService) UpdatePassword(key, plainPwForPersist, plainUserPw string) error {
	ctx := context.Background()
	user, err := s.masterUserRepository.GetUserById(ctx, s.userPk)
	if err != nil {
		return err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(plainUserPw+user.UserName)); err != nil {
		return pwmErr.NoUser
	}

	encryptedPw, err := s.encrypt(plainPwForPersist)
	if err != nil {
		return err
	}

	return s.passwordListRepository.UpdatePasswordByUserPkAndKey(ctx, s.userPk, key, *encryptedPw)
}

func (s *passwordListService) GetPassword(key string) (*string, error) {
	ctx := context.Background()
	passwordEntity, err := s.passwordListRepository.GetPasswordByUserPkAndKey(ctx, s.userPk, key)
	if err != nil {
		return nil, err
	}

	return s.decrypt(passwordEntity.Password)
}

func (s *passwordListService) SetPassword(key, plainPw, description string) error {
	// encrypt password
	encrypt, err := s.encrypt(plainPw)
	if err != nil {
		return err
	}

	// persist password
	ctx := context.Background()
	return s.txManager.Txx(ctx, func(ctx context.Context) error {
		if err := s.passwordListRepository.Insert(ctx, s.userPk, key, *encrypt); err != nil {
			return err
		}

		e, err := s.passwordListRepository.GetPasswordByUserPkAndKey(ctx, s.userPk, key)
		if err != nil {
			return err
		}

		if err := s.passwordListDetailRepository.Insert(ctx, e.Id, description); err != nil {
			return err
		}

		return nil
	})
}

func (s *passwordListService) decrypt(encryptedPw string) (*string, error) {
	block, err := aes.NewCipher(s.cryptoKey)
	if err != nil {
		return nil, err
	}

	decodedCipher, err := base64.URLEncoding.DecodeString(encryptedPw)
	if err != nil {
		return nil, err
	}

	iv := decodedCipher[:aes.BlockSize]
	cipherBytes := decodedCipher[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherBytes, cipherBytes)

	password := string(cipherBytes)
	return &password, nil
}

func (s *passwordListService) encrypt(password string) (*string, error) {
	block, err := aes.NewCipher(s.cryptoKey)
	if err != nil {
		return nil, err
	}

	cipherText := make([]byte, aes.BlockSize+len(password))
	iv := cipherText[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], []byte(password))

	return util.MustPointer(base64.URLEncoding.EncodeToString(cipherText)), nil
}

func NewPasswordService(userPk int, cryptoKey []byte) (PasswordService, error) {
	repositoryFactory, err := repository.GetRepositoryFactory()
	if err != nil {
		return nil, err
	}

	return &passwordListService{
		masterUserRepository:         repositoryFactory.MasterUserRepository,
		passwordListRepository:       repositoryFactory.PasswordListRepository,
		passwordListDetailRepository: repositoryFactory.PasswordListDetailRepository,
		txManager:                    repositoryFactory.TxManager,
		userPk:                       userPk,
		cryptoKey:                    cryptoKey,
	}, nil
}
