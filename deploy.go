package main

import (
	"github.com/urfave/cli/v2"
)



func (app *App) deploy(c *cli.Context) error {
	app.printTaskName("deploy:prepare")

	err := app.prepare(c)

	if err != nil {
		return err
	}


	tasks := NewTasks(app.TasksOrder, app.ConfigTasks)

	err = app.run(tasks)

	if err != nil {
		return err
	}

	app.Color.Print(app.Color.Green, "Successfully deployed!")

	return nil
}





