package command

import (
	"context"
	"github.com/aivyss/password-manager/console"
	"github.com/aivyss/password-manager/pwmErr"
	"github.com/aivyss/password-manager/repository"
	"github.com/aivyss/typex/util"
	"github.com/urfave/cli/v2"
	"golang.org/x/crypto/bcrypt"
)

type MasterUserCommandHandler struct {
	masterUserRepository repository.MasterUserRepository
}

func NewMasterUserCommandHandler(masterUserRepository repository.MasterUserRepository) *MasterUserCommandHandler {
	return &MasterUserCommandHandler{
		masterUserRepository: masterUserRepository,
	}
}

func (h *MasterUserCommandHandler) CreateUser(c *cli.Context) error {
	name := c.String("name")
	password := c.String("pw")

	if util.IsBlank(name) || len(password) < 16 {
		return pwmErr.InvalidOpt
	}

	hashedPw, err := bcrypt.GenerateFromPassword([]byte(password+name), 15)
	if err != nil {
		return pwmErr.GeneratePw
	}

	ctx := context.Background()
	if err := h.masterUserRepository.Insert(ctx, name, string(hashedPw)); err != nil {
		return err
	}

	return nil
}

func (h *MasterUserCommandHandler) Login(c *cli.Context) error {
	name := c.String("name")
	password := c.String("pw")
	if util.IsBlank(name) || len(password) < 16 {
		return pwmErr.InvalidOpt
	}

	ctx := context.Background()
	user, err := h.masterUserRepository.GetUserByUserName(ctx, name)
	if err != nil {
		return err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password+name)); err != nil {
		return pwmErr.NoUser
	}

	mainConsole, err := console.NewMainConsole(user.Id, password)
	if err != nil {
		return err
	}

	if err = mainConsole.Run(); err != nil {
		return err
	}

	return nil
}

func (h *MasterUserCommandHandler) Command() *cli.Command {
	return &cli.Command{
		Name:        "user",
		Description: "user domain",
		Subcommands: []*cli.Command{
			{
				Name:  "create",
				Usage: "create user",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "name",
						Value: "${user_name}",
						Usage: "set user name",
					},
					&cli.StringFlag{
						Name:  "pw",
						Value: "${user_password}",
						Usage: "set user password",
					},
				},
				Action: h.CreateUser,
			},
			{
				Name:  "login",
				Usage: "login user",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "name",
						Value: "${user_name}",
						Usage: "set user name",
					},
					&cli.StringFlag{
						Name:  "pw",
						Value: "${user_password}",
						Usage: "set user password",
					},
				},
				Action: h.Login,
			},
		},
	}
}
