package console

import (
	"bufio"
	"fmt"
	"github.com/aivyss/password-manager/console/command"
	"github.com/aivyss/password-manager/options/parser"
	"github.com/aivyss/password-manager/pwmErr"
	"github.com/aivyss/password-manager/pwmOs"
	"github.com/urfave/cli/v2"
	"os"
	"sync"
)

var welcomeOnce sync.Once

const prefixCommandName = "main"

type MainConsole struct {
	app cli.App
}

func (m *MainConsole) Run() error {
	scanner := bufio.NewScanner(os.Stdin)
	pwmOs.ClearTerminalBuffer()

	for {
		welcomeOnce.Do(func() {
			fmt.Println("[pwm] you entered password-manager console")
		})

		fmt.Print("[pwm][main console] > ")

		// parse command line
		scanner.Scan()
		text := scanner.Text()
		if text == "" {
			continue
		}

		nonPrefixCommand := parser.ParseCommand(text)
		cmd := make([]string, 0, len(nonPrefixCommand)+1)
		cmd = append(cmd, prefixCommandName)
		cmd = append(cmd, nonPrefixCommand...)

		if err := m.app.Run(cmd); err != nil {
			if err == pwmErr.ExitErr {
				return nil
			}

			fmt.Printf("[pwm][main console][error]%s\n", err.Error())
		}
	}
}

func NewMainConsole(userPk int, password string) (MainConsole, error) {
	passwordCommandHandler, err := command.NewPasswordCommandHandler(userPk, password)
	if err != nil {
		return MainConsole{}, err
	}

	mainConsoleApp := cli.App{
		Name: prefixCommandName,
		Commands: []*cli.Command{
			{
				Name: "save",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "k",
						Usage:    "${example: google}",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "description",
						Aliases:  []string{"d"},
						Usage:    "site https://www.google.com",
						Required: false,
					},
				},
				Action: passwordCommandHandler.SetPassword,
			},
			{
				Name: "get",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "k",
						Usage: "${example: google}",
					},
				},
				Action: passwordCommandHandler.GetPassword,
			},
			{
				Name:    "list",
				Aliases: []string{"l"},
				Action:  passwordCommandHandler.GetAllKeys,
			},
			{
				Name:        "update",
				Description: "update your password",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "k",
						Usage: "${example: google}",
					},
				},
				Action: passwordCommandHandler.UpdatePassword,
			},
			{
				Name: "clear",
				Action: func(context *cli.Context) error {
					pwmOs.ClearTerminalBuffer()
					return nil
				},
			},
			{
				Name:    "description",
				Aliases: []string{"d"},
				Subcommands: []*cli.Command{
					{
						Name:    "update",
						Aliases: []string{"u"},
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "key",
								Aliases: []string{"k"},
							},
							&cli.StringFlag{
								Name:    "value",
								Aliases: []string{"v"},
							},
						},
						Action: passwordCommandHandler.UpdateDescription,
					},
				},
			},
			{
				Name: "exit",
				Action: func(context *cli.Context) error {
					return pwmErr.ExitErr
				},
			},
		},
		Description: "store your passwords safely",
	}

	return MainConsole{
		app: mainConsoleApp,
	}, nil
}
