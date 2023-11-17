package main

import (
	"fmt"
	"github.com/aivyss/password-manager/command"
	"github.com/aivyss/password-manager/pwmErr"
	"github.com/aivyss/password-manager/repository"
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
	repository.InitDB(db)
	factory, err := repository.NewRepositoryFactory(db)
	if err != nil {
		fmt.Println(pwmErr.FailToCreateRepository.Error())
	}
	masterUserCommandHandler := command.NewMasterUserCommandHandler(factory.MasterUserRepository)

	app := cli.App{
		Name: "pwm",
		Commands: []*cli.Command{
			masterUserCommandHandler.Command(),
		},
		Description: "store your passwords safely",
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err.Error())
	}
}
