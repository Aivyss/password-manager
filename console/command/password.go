package command

import (
	"bufio"
	"fmt"
	"github.com/aivyss/password-manager/csv"
	"github.com/aivyss/password-manager/options"
	"github.com/aivyss/password-manager/pwmOs"
	"github.com/aivyss/password-manager/service"
	"github.com/aivyss/password-manager/view"
	"github.com/urfave/cli/v2"
	"os"
	"time"
)

type PasswordCommandHandler struct {
	passwordService service.PasswordService
}

func (h *PasswordCommandHandler) SetPassword(c *cli.Context) error {
	defer pwmOs.ClearTerminalBuffer()

	opts, err := options.NewMainSaveOptKOptPw(c)
	if err != nil {
		return err
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

	fmt.Printf("[pwm][main console] check your password: %s\n", *password)
	fmt.Printf("[pwm][main console] please 'clear' command after checking your password\n")
	return nil
}

func (h *PasswordCommandHandler) UpdatePassword(c *cli.Context) error {
	opts, err := options.NewMainUpdateOptKOptPw(c)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("[pwm][main console] please enter your master password again: ")
	scanner.Scan()
	plainMasterUserPassword := scanner.Text()

	return h.passwordService.UpdatePassword(opts.Key, opts.Password, plainMasterUserPassword)
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
