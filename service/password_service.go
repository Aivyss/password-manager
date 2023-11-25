package service

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"github.com/aivyss/bean"
	"github.com/aivyss/password-manager/entity"
	"github.com/aivyss/password-manager/pwmContext"
	"github.com/aivyss/password-manager/pwmErr"
	"github.com/aivyss/password-manager/repository"
	"golang.org/x/crypto/bcrypt"
	"io"
)

type PasswordService interface {
	SetPassword(key, plainPw, description string) error
	GetPassword(key string) (string, error)
	UpdatePassword(key, plainPwForPersist, plainUserPw string) error
	GetAllPasswords() ([]entity.PasswordListWithDescription, error)
	UpdateDescription(key, userPlainPassword, description string) error
}

type passwordListService struct {
	masterUserRepository         repository.MasterUserRepository
	passwordListRepository       repository.PasswordListRepository
	passwordListDetailRepository repository.PasswordListDetailRepository
	txManager                    repository.TxManager
	userPk                       int
	cryptoKey                    []byte
}

func (s *passwordListService) UpdateDescription(key, userPlainPassword, description string) error {
	ctx := context.Background()

	user, err := s.masterUserRepository.GetUserById(ctx, s.userPk)
	if err != nil {
		return pwmErr.NoUser
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userPlainPassword+user.UserName)); err != nil {
		return pwmErr.WrongPw
	}

	passwordEntity, err := s.passwordListRepository.GetPasswordByUserPkAndKey(ctx, s.userPk, key)
	if err != nil {
		return err
	}

	return s.passwordListDetailRepository.UpdateDescriptionByPasswordListKey(ctx, passwordEntity.Id, description)
}

func (s *passwordListService) GetAllPasswords() ([]entity.PasswordListWithDescription, error) {
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

	return s.passwordListRepository.UpdatePasswordByUserPkAndKey(ctx, s.userPk, key, encryptedPw)
}

func (s *passwordListService) GetPassword(key string) (string, error) {
	ctx := context.Background()
	passwordEntity, err := s.passwordListRepository.GetPasswordByUserPkAndKey(ctx, s.userPk, key)
	if err != nil {
		return "", err
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
		if err := s.passwordListRepository.Insert(ctx, s.userPk, key, encrypt); err != nil {
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

func (s *passwordListService) decrypt(encryptedPw string) (string, error) {
	cryptoKey := make([]byte, 0, 32)
	cryptoKey = append(cryptoKey, s.cryptoKey...)
	cryptoKey = append(cryptoKey, []byte(pwmContext.GetGlobalContext().BuildSecretKey)[:16]...)
	block, err := aes.NewCipher(cryptoKey)
	if err != nil {
		return "", err
	}

	decodedCipher, err := base64.URLEncoding.DecodeString(encryptedPw)
	if err != nil {
		return "", err
	}

	iv := decodedCipher[:aes.BlockSize]
	cipherBytes := decodedCipher[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherBytes, cipherBytes)

	password := string(cipherBytes)
	return password, nil
}

func (s *passwordListService) encrypt(password string) (string, error) {
	cryptoKey := make([]byte, 0, 32)
	cryptoKey = append(cryptoKey, s.cryptoKey...)
	cryptoKey = append(cryptoKey, []byte(pwmContext.GetGlobalContext().BuildSecretKey)[:16]...)
	block, err := aes.NewCipher(cryptoKey)
	if err != nil {
		return "", err
	}

	cipherText := make([]byte, aes.BlockSize+len(password))
	iv := cipherText[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], []byte(password))

	return base64.URLEncoding.EncodeToString(cipherText), nil
}

func NewPasswordService(userPk int, cryptoKey []byte) (PasswordService, error) {
	masterUserRepository, err := bean.GetBean[repository.MasterUserRepository]()
	if err != nil {
		return nil, err
	}
	passwordListRepository, err := bean.GetBean[repository.PasswordListRepository]()
	if err != nil {
		return nil, err
	}
	passwordListDetailRepository, err := bean.GetBean[repository.PasswordListDetailRepository]()
	if err != nil {
		return nil, err
	}
	txManager, err := bean.GetBean[repository.TxManager]()
	if err != nil {
		return nil, err
	}

	return &passwordListService{
		masterUserRepository:         masterUserRepository,
		passwordListRepository:       passwordListRepository,
		passwordListDetailRepository: passwordListDetailRepository,
		txManager:                    txManager,
		userPk:                       userPk,
		cryptoKey:                    cryptoKey,
	}, nil
}
