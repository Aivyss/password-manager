package main

import (
	"fmt"
	"github.com/aivyss/bean"
	"github.com/aivyss/jsonx"
	"github.com/aivyss/password-manager/command"
	"github.com/aivyss/password-manager/options"
	"github.com/aivyss/password-manager/pwmContext"
	"github.com/aivyss/password-manager/repository"
	"github.com/aivyss/password-manager/service"
	"github.com/aivyss/password-manager/validator"
	"github.com/aivyss/typex/pointer"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/urfave/cli/v2"
	"os"
)

const dataSource = "password_manager.db"

var buildSecretKey string

func main() {
	// Set BuildSecretKey
	pwmContext.NewGlobalContext(buildSecretKey)

	// DB
	db, err := sqlx.Connect("sqlite3", dataSource)
	if err != nil {
		panic(err)
	}
	if err := repository.InitDB(db); err != nil {
		fmt.Println(err.Error())
		return
	}
	beanBuff := bean.GetBeanBuffer()
	beanBuff.RegisterBean(func() *sqlx.DB { return db })
	beanBuff.RegisterBean(repository.NewAppVersionRepository)
	beanBuff.RegisterBean(repository.NewMasterUserRepository)
	beanBuff.RegisterBean(repository.NewPasswordListRepository)
	beanBuff.RegisterBean(repository.NewPasswordListDetailRepository)
	beanBuff.RegisterBean(repository.NewTxManager)
	beanBuff.RegisterBean(service.NewMasterUserService)
	beanBuff.RegisterBean(command.NewMasterUserCommandHandler)
	beanBuff.RegisterBean(command.NewAppVersionCommandHandler)
	err = beanBuff.Buffer()
	if err != nil {
		fmt.Println(err.Error())
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
	masterUserCommandHandler, err := bean.GetBean[*command.MasterUserCommandHandler]()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	appVersionCommandHandler, err := bean.GetBean[*command.AppVersionCommandHandler]()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

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
	jsonx.RegisterValidator[options.UserCreateOptNameOptPw](pointer.MustPointer(validator.UserCreateOptNameOptPwValidator(1)))
	jsonx.RegisterValidator[options.UserLoginOptNameOptPw](pointer.MustPointer(validator.UserLoginOptNameOptPwValidator(1)))
}
