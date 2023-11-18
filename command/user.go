package command

import (
	"fmt"
	"github.com/aivyss/password-manager/console"
	"github.com/aivyss/password-manager/csv"
	"github.com/aivyss/password-manager/pwmErr"
	"github.com/aivyss/password-manager/service"
	"github.com/aivyss/password-manager/view"
	"github.com/aivyss/typex/util"
	"github.com/urfave/cli/v2"
	"time"
)

type MasterUserCommandHandler struct {
	masterUserService service.MasterUserService
}

func NewMasterUserCommandHandler(masterUserService service.MasterUserService) *MasterUserCommandHandler {
	return &MasterUserCommandHandler{
		masterUserService: masterUserService,
	}
}

func (h *MasterUserCommandHandler) CreateUser(c *cli.Context) error {
	name := c.String("name")
	password := c.String("pw")

	if util.IsBlank(name) || len(password) < 16 {
		return pwmErr.InvalidOpt
	}

	return h.masterUserService.CreateUser(name, password)
}

func (h *MasterUserCommandHandler) Login(c *cli.Context) error {
	name := c.String("name")
	password := c.String("pw")
	if util.IsBlank(name) || len(password) < 16 {
		return pwmErr.InvalidOpt
	}

	user, err := h.masterUserService.Login(name, password)
	if err != nil {
		return err
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

func (h *MasterUserCommandHandler) GetUsers(c *cli.Context) error {
	head := c.Int("head")
	tail := c.Int("tail")

	users, err := h.masterUserService.GetUsers(head, tail)
	if err != nil {
		return err
	}

	type userCsvBindObject struct {
		UserName  string    `csv:"Username"`
		CreatedAt time.Time `csv:"Created Date"`
		UpdatedAt time.Time `csv:"Last Updated Date"`
	}
	objects := make([]userCsvBindObject, 0, len(users))
	for _, user := range users {
		objects = append(objects, userCsvBindObject{
			UserName:  user.UserName,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}

	if len(objects) > 0 {
		csvLines := csv.CreateCsvLines(objects)
		view.StdoutTableView(csvLines[0], csvLines[1:])
		return nil
	}

	fmt.Println("[pwm] no user")
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
						Name:     "name",
						Value:    "${user_name}",
						Usage:    "set user name",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "pw",
						Value:    "${user_password}",
						Usage:    "set user password",
						Required: true,
					},
				},
				Action: h.CreateUser,
			},
			{
				Name:  "login",
				Usage: "login user",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "name",
						Value:    "${user_name}",
						Usage:    "set user name",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "pw",
						Value:    "${user_password}",
						Usage:    "set user password",
						Required: true,
					},
				},
				Action: h.Login,
			},
			{
				Name:        "list",
				Description: "list up users (default: 10)",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name: "head",
						Aliases: []string{
							"H",
						},
						Usage:    "head n",
						Required: false,
					},
					&cli.IntFlag{
						Name: "tail",
						Aliases: []string{
							"t",
						},
						Usage:    "tail n",
						Required: false,
					},
				},
				Action: h.GetUsers,
			},
		},
	}
}
