package command

import (
	"context"
	"fmt"
	"github.com/aivyss/password-manager/console"
	"github.com/aivyss/password-manager/csv"
	"github.com/aivyss/password-manager/entity"
	"github.com/aivyss/password-manager/pwmErr"
	"github.com/aivyss/password-manager/repository"
	"github.com/aivyss/password-manager/view"
	"github.com/aivyss/typex/util"
	"github.com/urfave/cli/v2"
	"golang.org/x/crypto/bcrypt"
	"sort"
	"time"
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

func (h *MasterUserCommandHandler) GetUsers(c *cli.Context) error {
	head := c.Int("head")
	tail := c.Int("tail")

	ctx := context.Background()
	users, err := h.masterUserRepository.GetAllUsers(ctx)
	if err != nil {
		return err
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
		result = users[len(users)-tail:]
	} else {
		result = users[:]
		if len(users) > 10 {
			result = users[:10]
		}
	}

	type userCsvBindObject struct {
		UserName  string    `csv:"Username"`
		CreatedAt time.Time `csv:"Created Date"`
		UpdatedAt time.Time `csv:"Last Updated Date"`
	}
	objects := make([]userCsvBindObject, 0, len(result))
	for _, user := range result {
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
