package console

import (
	"bufio"
	"fmt"
	"github.com/aivyss/password-manager/console/command"
	"github.com/aivyss/password-manager/pwmErr"
	"github.com/aivyss/password-manager/pwmOs"
	"github.com/urfave/cli/v2"
	"os"
	"regexp"
	"strings"
	"sync"
)

var welcomeOnce sync.Once

const prefixCommandName = "main"

type MainConsole struct {
	app *cli.App
}

func (m *MainConsole) Run() error {
	patternMultipleSpace, err := regexp.Compile("\\b\\s+\\b")

	scanner := bufio.NewScanner(os.Stdin)
	pwmOs.ClearTerminalBuffer()

	for {
		welcomeOnce.Do(func() {
			fmt.Println("[pwm] you entered password-manager console")
		})

		fmt.Print("[pwm][main console] > ")

		// parse command line
		scanner.Scan()
		if err != nil {
			return err
		}
		commandLine := scanner.Text()
		commandLine = string(patternMultipleSpace.ReplaceAll([]byte(commandLine), []byte(" ")))
		args := strings.Split(commandLine, " ")
		inputArgs := make([]string, 0, len(args)+1)
		inputArgs = append(inputArgs, prefixCommandName)
		for _, arg := range args {
			inputArgs = append(inputArgs, arg)
		}

		if err := m.app.Run(inputArgs); err != nil {
			if err == pwmErr.ExitErr {
				return nil
			}

			return err
		}
	}
}

func NewMainConsole(userPk int, password string) (*MainConsole, error) {
	passwordCommandHandler, err := command.NewPasswordCommandHandler(userPk, password)
	if err != nil {
		return nil, err
	}

	mainConsoleApp := cli.App{
		Name: prefixCommandName,
		Commands: []*cli.Command{
			{
				Name: "save",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "k",
						Usage: "${example: google}",
					},
					&cli.StringFlag{
						Name:  "pw",
						Usage: "${your_password}",
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
					&cli.StringFlag{
						Name:  "pw",
						Usage: "${your_password}",
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
				Name: "exit",
				Action: func(context *cli.Context) error {
					return pwmErr.ExitErr
				},
			},
		},
		Description: "store your passwords safely",
	}

	return &MainConsole{
		app: &mainConsoleApp,
	}, nil
}
