package app

import "github.com/urfave/cli/v2"

func (app *App) Rollback(c *cli.Context) error {
	app.printTaskName("rollback:prepare")
	err := app.prepare(c)

	if err != nil {
		return err
	}

	tasksOrder := []string{
		"rollback:check",
		"rollback:symlink",
		"rollback:cleanup",
		"rollback:clean-meta",
	}
	app.TasksOrder = tasksOrder
	tasks := NewTasks(app.TasksOrder, ConfigTasks{}, app.Release)

	logTasks, errRun := app.run(tasks)

	if app.ConfigNotifications != nil {
		app.Color.Print(app.Color.Info, "Send notifications ...")
		sender := NewSender(*app.ConfigNotifications, messageProperties{
			"Rollback",
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

	app.Color.Print(app.Color.Green, "Rollback to " + app.Release.Name + " release was successful.")

	return nil
}