package main

import (
	"strconv"
	"strings"
)

type task struct {
	name string
	method func(app *App)error
}

type defaultTasks struct {
	Remote map[string]func(app *App) error
}

type Tasks struct {
	Remote []task
}



var defaultTasksOrder = []string{"deploy:check","deploy:release","deploy:update_code","deploy:symlink","restart:service","cleanup"}

func NewTasks(sequence []string, configTasks ConfigTasks) Tasks {
	var tasks Tasks
	var taskObject task

	defaultTasks := taskObject.defaultTasks()

	for _, name := range sequence {
		if _, ok := defaultTasks.Remote[name]; ok {
			tasks.Remote = append(tasks.Remote, task{
				name,
				defaultTasks.Remote[name],
			})
		} else {
			for _, row := range configTasks.Remote {
				if row.Name == name {
					tasks.Remote = append(tasks.Remote, task{
						name,
						func(app *App) error {
							_, err := app.Bash.run(row.Command)

							return err
						},
					})
					break
				}
			}
		}
	}

	return tasks
}

func(t *task) defaultTasks() defaultTasks {
	var tasks defaultTasks
	tasks.Remote = make(map[string]func(app *App) error)
	tasks.Remote["deploy:check"] = t.deployCheck
	tasks.Remote["deploy:release"] = t.deployRelease
	tasks.Remote["deploy:update_code"] = t.deployUpdateCode
	tasks.Remote["deploy:symlink"] = t.deploySymlink
	tasks.Remote["cleanup"] = t.cleanup

	return tasks
}
func(t *task) deployCheck(app *App) error {
	return app.Bash.multiRun([]string{
		"if [ ! -d " + app.Release.DeployPath + " ]; then mkdir -p " + app.Release.DeployPath + "; fi",
		"cd " + app.Release.DeployPath + " && if [ ! -d .dep ]; then mkdir .dep; fi",
		"cd " + app.Release.DeployPath + " && if [ ! -d releases ]; then mkdir releases; fi",
		"cd " + app.Release.DeployPath + " && if [ ! -d shared ]; then mkdir shared; fi",

		"cd " + app.Release.DeployPath + " && if [ -e release ]; then rm release; fi",
		"cd " + app.Release.DeployPath + " && if [ -h release ]; then rm release; fi",
	})
}
func(t *task) deployRelease(app *App) error {
	date, err := app.Bash.runOutput(`date +"%Y-%m-%d %H:%M:%S"`)
	if err != nil {
		return err
	}

	count, err := app.Bash.runOutput("cd " + app.Release.DeployPath + "/releases && ls -l ./| grep ^d | wc -l")
	if err != nil {
		return err
	}

	if *count == "0" {
		app.Release.Path = app.Release.DeployPath + "/releases/1"
		app.Release.Name = "1"
		app.Release.Number = 1
	} else {
		last, err := app.Bash.runOutput("cd " + app.Release.DeployPath + "/.dep && tail -n 1 " + app.Meta.Name)
		if err != nil {
			return err
		}

		previousRelease := strings.Split(*last, app.Meta.Separator)[1]
		lastRelease, _ := strconv.ParseInt(previousRelease, 10, 64)

		app.Release.Name = strconv.FormatInt(lastRelease + 1, 10)
		app.Release.Number = lastRelease + 1

		app.Release.Path = app.Release.DeployPath + "/releases/" + app.Release.Name

		previousReleasePath := app.Release.DeployPath + "/releases/" + previousRelease
		app.Release.PreviousReleasePath = &previousReleasePath
	}

	return app.Bash.multiRun([]string{
		"cd " + app.Release.DeployPath + `/.dep && echo "` + *date + `, ` + app.Release.Name  + `" >> ` + app.Meta.Name,
		"mkdir -p " + app.Release.Path,
		"ln -nfs " + app.Release.Path + " " + app.Release.DeployPath + "/release",
	})
}
func(t *task) deployUpdateCode(app *App) error {
	reference := ""
	if app.Release.PreviousReleasePath != nil {
		reference += " --reference " + *app.Release.PreviousReleasePath
	}

	return app.Bash.multiRun([]string{
		"git clone -b " + app.Release.Branch + " --recursive  " + reference +" --dissociate " + app.Release.Repository + " " + app.Release.Path + " 2>&1",
	})
}
func(t *task) deploySymlink(app *App) error {
	return app.Bash.multiRun([]string{
		"mv -T " + app.Release.DeployPath + "/release " + app.Release.DeployPath + "/current",
	})
}
func(t *task) cleanup(app *App) error {
	meta, err := app.Bash.runOutput("cd " + app.Release.DeployPath + "/.dep && tail -n " + strconv.Itoa(app.Release.KeepReleases) + " " + app.Meta.Name)
	if err != nil {
		return err
	}

	lastReleases := strings.Split(*meta, "\n")

	var releases []string
	for _, str := range lastReleases {
		releases = append(releases, strings.Split(str, app.Meta.Separator)[1])
	}

	scan, err := app.Bash.runOutput("cd " + app.Release.DeployPath + " && ls -xm releases/")
	if err != nil {
		return err
	}

	dirs := strings.Split(strings.Replace(*scan, "\r\n", " ", -1), app.Meta.Separator)

	var commands []string
	for _, dir := range dirs {
		exists, _ := InArray(dir, releases)

		if exists == false {
			commands = append(commands, "cd " + app.Release.DeployPath + "/releases && rm -rf " + dir)
		}
	}

	commands = append(commands, "cd " + app.Release.DeployPath + " && if [ -e release ]; then rm release; fi")
	commands = append(commands, "cd " + app.Release.DeployPath + " && if [ -h release ]; then rm release; fi")

	return app.Bash.multiRun(commands)
}


