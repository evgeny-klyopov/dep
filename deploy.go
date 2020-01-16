package main

import (
	"github.com/urfave/cli/v2"
)



func (app *App) deploy(c *cli.Context) error {
	app.printMessageTask("deploy:prepare")
	err := app.prepare(c)

	if err != nil {
		return err
	}

	commands := NewCommands()
	commands = commands.init(app.TaskOrder, app.Tasks)

	//fmt.Println(commands)
	//
	//
	//
	//return errors.New("test")

	//commands := GetCommands()
	//
	err = app.runCommand(commands)

	if err != nil {
		return err
	}

	app.Color.Print(app.Color.Green, "Successfully deployed!")

	return nil
}





