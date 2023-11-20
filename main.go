package main

import (
	"fmt"
	"github.com/aivyss/jsonx"
	"github.com/aivyss/password-manager/command"
	"github.com/aivyss/password-manager/options"
	"github.com/aivyss/password-manager/pwmErr"
	"github.com/aivyss/password-manager/repository"
	"github.com/aivyss/password-manager/service"
	"github.com/aivyss/password-manager/validator"
	"github.com/aivyss/typex/util"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/urfave/cli/v2"
	"os"
)

const dataSource = "password_manager.db"

func main() {
	// DB
	db, err := sqlx.Connect("sqlite3", dataSource)
	if err != nil {
		panic(err)
	}
	if err := repository.InitDB(db); err != nil {
		fmt.Println(err.Error())
		return
	}
	factory, err := repository.NewRepositoryFactory(db)
	if err != nil {
		fmt.Println(pwmErr.FailToCreateRepository.Error())
		return
	}

	// App Version
	if err = repository.CheckAppVersion(); err != nil {
		fmt.Println(err.Error())
		return
	}

	// Register Validations
	registerValidations()

	// Create Handlers
	masterUserCommandHandler := command.NewMasterUserCommandHandler(service.NewMasterUserService(factory.MasterUserRepository))
	appVersionCommandHandler := command.NewAppVersionCommandHandler(factory.AppVersionRepository)

	// Command Mapping
	app := cli.App{
		Name: "pwm",
		Commands: []*cli.Command{
			masterUserCommandHandler.Command(),
			appVersionCommandHandler.Command(),
		},
		Description: "store your passwords safely",
	}

	args := os.Args
	if err := app.Run(args); err != nil {
		fmt.Println(err.Error())
	}
}

func registerValidations() {
	jsonx.RegisterValidator[options.UserCreateOptNameOptPw](util.MustPointer(validator.UserCreateOptNameOptPwValidator(1)))
	jsonx.RegisterValidator[options.UserLoginOptNameOptPw](util.MustPointer(validator.UserLoginOptNameOptPwValidator(1)))
}
