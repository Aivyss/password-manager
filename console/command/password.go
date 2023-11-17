package command

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/aivyss/password-manager/pwmErr"
	"github.com/aivyss/password-manager/pwmOs"
	"github.com/aivyss/password-manager/repository"
	"github.com/aivyss/typex/util"
	"github.com/urfave/cli/v2"
	"io"
)

type PasswordCommandHandler struct {
	passwordListRepository repository.PasswordListRepository
	userPk                 int
	cryptoKey              []byte
}

func (h *PasswordCommandHandler) SetPassword(c *cli.Context) error {
	defer pwmOs.ClearTerminalBuffer()

	key := c.String("k")
	password := c.String("pw")

	if util.IsBlank(key) || util.IsBlank(password) {
		return pwmErr.InvalidOpt
	}

	// encrypt password
	encrypt, err := h.encrypt(password)
	if err != nil {
		return err
	}

	// persist password
	ctx := context.Background()
	if err = h.passwordListRepository.Insert(ctx, h.userPk, key, *encrypt); err != nil {
		return err
	}

	return nil
}

func (h *PasswordCommandHandler) GetPassword(c *cli.Context) error {
	key := c.String("k")
	if util.IsBlank(key) {
		return pwmErr.InvalidOpt
	}

	ctx := context.Background()
	passwordEntity, err := h.passwordListRepository.GetPasswordByUserPkAndKey(ctx, h.userPk, key)
	if err != nil {
		return err
	}

	decrypt, err := h.decrypt(passwordEntity.Password)
	if err != nil {
		return err
	}

	fmt.Printf("[pwm][main console] check your password: %s\n", *decrypt)
	fmt.Printf("[pwm][main console] please 'clear' command after checking your password\n")
	return nil
}

func (h *PasswordCommandHandler) decrypt(encryptedPw string) (*string, error) {
	block, err := aes.NewCipher(h.cryptoKey)
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

func (h *PasswordCommandHandler) encrypt(password string) (*string, error) {
	block, err := aes.NewCipher(h.cryptoKey)
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

func NewPasswordCommandHandler(userPk int, password string) (*PasswordCommandHandler, error) {
	repositoryFactory, err := repository.GetRepositoryFactory()
	if err != nil {
		return nil, err
	}

	return &PasswordCommandHandler{
		passwordListRepository: repositoryFactory.PasswordListRepository,
		userPk:                 userPk,
		cryptoKey:              []byte(password)[:16],
	}, nil
}
