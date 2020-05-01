package app

import (
	"github.com/urfave/cli/v2"
)


func (app *App) Deploy(c *cli.Context) error {
	app.printTaskName("deploy:prepare\n")

	err := app.prepare(c)

	if err != nil {
		return err
	}


	tasks := NewTasks(app.TasksOrder, app.ConfigTasks, app.Release)

	logTasks, errRun := app.run(tasks)

	if app.ConfigNotifications != nil {
		app.Color.Print(app.Color.Info, "Send notifications ...")
		sender := NewSender(*app.ConfigNotifications, messageProperties{
			"Deploy",
			*logTasks,
			app.Release.Path,
			app.Release.Stage,
			app.Release.Name,
			app.Bash.Host,
			app.Bash.Port,
		})
		errSender := sender.Send()
		if errSender != nil {
			app.Color.Print(app.Color.Fatal, "Error send notifications")
		} else {
			app.Color.Print(app.Color.Green, "Notifications send")
		}
	}

	if errRun != nil {
		return errRun
	}

	app.Color.Print(app.Color.Green, "Successfully deployed!")

	return nil
}





