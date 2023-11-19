package command

import (
	"fmt"
	"github.com/aivyss/password-manager/csv"
	"github.com/aivyss/password-manager/options"
	"github.com/aivyss/password-manager/pwmErr"
	"github.com/aivyss/password-manager/service"
	"github.com/aivyss/password-manager/view"
	"github.com/atotto/clipboard"
	"github.com/urfave/cli/v2"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"time"
)

type PasswordCommandHandler struct {
	passwordService service.PasswordService
}

func (h *PasswordCommandHandler) SetPassword(c *cli.Context) error {
	fmt.Print("[pwm] enter password: ")
	password, _ := terminal.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()

	object, err := options.NewmainSaveOptKOptPw(c)
	if err != nil {
		return err
	}

	opts, err := object.ToEntity(string(password))
	if err != nil {
		return err
	}

	if entity, _ := h.passwordService.GetPassword(opts.Key); entity != nil {
		return pwmErr.AlreadyExistKey
	}

	if err := h.passwordService.SetPassword(opts.Key, opts.Password); err != nil {
		return err
	}

	return nil
}

func (h *PasswordCommandHandler) GetPassword(c *cli.Context) error {
	opts, err := options.NewMainGetOptK(c)
	if err != nil {
		return err
	}

	password, err := h.passwordService.GetPassword(opts.Key)
	if err != nil {
		return err
	}

	clipboard.WriteAll(*password)

	fmt.Println("[pwm][main console] your password is written in clipboard.")
	return nil
}

func (h *PasswordCommandHandler) UpdatePassword(c *cli.Context) error {
	fmt.Print("[pwm][main console] please enter master user password again: ")
	userPassword, _ := terminal.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()

	fmt.Print("[pwm][main console] please enter password for save: ")
	password, _ := terminal.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()

	object, err := options.NewmainUpdateOptKOptPw(c)
	if err != nil {
		return err
	}

	opts, err := object.ToEntity(string(password))
	if err != nil {
		return err
	}

	return h.passwordService.UpdatePassword(opts.Key, opts.Password, string(userPassword))
}

func (h *PasswordCommandHandler) GetAllKeys(_ *cli.Context) error {
	passwords, err := h.passwordService.GetAllPasswords()
	if err != nil {
		return err
	}

	type passwordListCsvBindObject struct {
		Key       string    `csv:"Key"`
		CreatedAt time.Time `csv:"Created Date"`
		UpdatedAt time.Time `csv:"Last Updated Date"`
	}

	objects := make([]passwordListCsvBindObject, 0, len(passwords))
	for _, password := range passwords {
		objects = append(objects, passwordListCsvBindObject{
			Key:       password.Key,
			CreatedAt: password.CreatedAt,
			UpdatedAt: password.UpdatedAt,
		})
	}

	if len(objects) > 0 {
		csvLines := csv.CreateCsvLines(objects)
		view.StdoutTableView(csvLines[0], csvLines[1:])
		return nil
	}

	fmt.Println("[pwm][main console] there is no record")

	return nil
}

func NewPasswordCommandHandler(userPk int, password string) (*PasswordCommandHandler, error) {
	passwordService, err := service.NewPasswordService(userPk, []byte(password)[:16])
	if err != nil {
		return nil, err
	}

	return &PasswordCommandHandler{
		passwordService: passwordService,
	}, nil
}
