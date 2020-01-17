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
	tasks := NewTasks(app.TasksOrder, ConfigTasks{})

	err = app.run(tasks)

	if err != nil {
		return err
	}

	app.Color.Print(app.Color.Green, "Rollback to " + app.Release.Name + " release was successful.")

	return nil
}