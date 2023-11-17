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
	if err := repository.InitDB(db); err != nil {
		fmt.Println(err.Error())
	}
	factory, err := repository.NewRepositoryFactory(db)
	if err != nil {
		fmt.Println(pwmErr.FailToCreateRepository.Error())
	}
	masterUserCommandHandler := command.NewMasterUserCommandHandler(factory.MasterUserRepository)

	// command mapping
	app := cli.App{
		Name: "pwm",
		Commands: []*cli.Command{
			masterUserCommandHandler.Command(),
		},
		Description: "store your passwords safely",
	}

	args := os.Args
	fmt.Println(args)
	if err := app.Run(args); err != nil {
		fmt.Println(err.Error())
	}
}
