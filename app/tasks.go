package app

import (
	"strconv"
	"strings"
)

type task struct {
	name string
	remote bool
	method func(app *App)error
}

type defaultTasks map[string]task
type Tasks []task

var defaultTasksOrder = []string{"deploy:check","deploy:release","deploy:update_code","deploy:symlink","restart:service","cleanup"}

func NewTasks(sequence []string, configTasks ConfigTasks) Tasks {
	var tasks Tasks
	var taskObject task

	defaultTasks := taskObject.defaultTasks()

	for _, name := range sequence {
		if _, ok := defaultTasks[name]; ok {
			tasks = append(tasks, defaultTasks[name])
		} else {
			isFind := false
			for _, row := range configTasks.Remote {
				if row.Name == name {
					isFind = true
					tasks = append(tasks, task{
						name,
						true,
						func(app *App) error {
							_, err := app.Bash.run(row.Command)

							return err
						},
					})
					break
				}
			}

			if isFind == false {
				for _, row := range configTasks.Local {
					if row.Name == name {
						isFind = true
						tasks = append(tasks, task{
							name,
							false,
							func(app *App) error {
								args := RegexSplit(row.Command, `\s+`)
								command := args[0]
								args = args[1:]

								return app.cmd(command, args...)
							},
						})
						break
					}
				}
			}
		}
	}

	return tasks
}

func(t *task) defaultTasks() defaultTasks {
	var tasks defaultTasks
	tasks = make(map[string]task)
	tasks["deploy:check"] = task{"deploy:check", true, t.deployCheck}
	tasks["deploy:release"] = task{"deploy:release", true, t.deployRelease}
	tasks["deploy:update_code"] = task{"deploy:update_code", true, t.deployUpdateCode}
	tasks["deploy:symlink"] = task{"deploy:symlink", true, t.deploySymlink}
	tasks["cleanup"] = task{"cleanup", true, t.cleanup}

	tasks["local:create-archive"] = task{"local:create-archive", false, t.localCreateArchive}
	tasks["deploy:extract-archive"] = task{"deploy:extract-archive", true, t.deployExtractArchive}
	tasks["local:send-archive"] = task{"local:send-archive", false, t.localSendArchive}


	tasks["rollback:check"] = task{"rollback:check", true, t.rollbackCheck}
	tasks["rollback:symlink"] = task{"rollback:symlink", true, t.rollbackSymlink}
	tasks["rollback:cleanup"] = task{"rollback:cleanup", true, t.rollbackCleanup}


	tasks["deploy:shared"] = task{"deploy:shared", true, t.deployShared}
	tasks["deploy:writable"] = task{"deploy:writable", true, t.deployWritable}

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



func(t *task) getCommandSymlink(from string, to string) string {
	return "mv -T " + from + " " + to
}
func(t *task) deploySymlink(app *App) error {
	return app.Bash.multiRun([]string{
		t.getCommandSymlink(app.Release.DeployPath + "/release", app.Release.DeployPath + "/current"),
	})
}
func(t *task) localCreateArchive(app *App) error {
	args := []string{"-cvzf", app.ArchiveName, app.Release.LocalDeployPath}
	return app.cmd("tar", args...)
}
func(t *task) localSendArchive(app *App) error {
	args := []string{"-P" + strconv.Itoa(app.Bash.Port), "-o StrictHostKeyChecking=no", app.ArchiveName, app.Bash.User + "@" + app.Bash.Host + ":" + app.Release.Path}
	err := app.cmd("scp", args...)
	if err != nil {
		return err
	}

	return app.cmd("rm",  app.ArchiveName)
}
func(t *task) deployExtractArchive(app *App) error {
	return app.Bash.multiRun([]string{
		"cd " + app.Release.Path + " && tar xvzf " + app.ArchiveName,
		"cd " + app.Release.Path + " && rm " + app.ArchiveName,
	})
}


func(t *task) rollbackCheck(app *App) error {
	releases, err := t.getReleases(app, 2)

	if err != nil {
		return err
	}

	if 2 > len(*releases) {
		return app.error(NotFoundPreviousRelease, false, nil, nil)
	}

	app.Release.Name = (*releases)[0]
	app.Release.Rollback = (*releases)[1]
	app.Release.Number, _ = strconv.ParseInt((*releases)[0], 10, 64)
	app.Release.Path = app.Release.DeployPath + "/releases/" + app.Release.Name

	return nil
}
func(t *task) rollbackSymlink(app *App) error {
	return app.Bash.multiRun([]string{
		"ln -sfn " + app.Release.Path + " " + app.Release.DeployPath + "/current",
	})
}
func(t *task) rollbackCleanup(app *App) error {
	return app.Bash.multiRun([]string{
		"rm -rf " + app.Release.DeployPath + "/releases/" + app.Release.Rollback,
	})
}
func(t *task) rollbackCleanMeta(app *App) error {
	return app.Bash.multiRun([]string{
		"cd " + app.Release.DeployPath + "/.dep && sed -i '$d' " + app.Meta.Name,
	})
}


func(t *task) getReleases(app *App, count int) (*[]string, error) {
	meta, err := app.Bash.runOutput("cd " + app.Release.DeployPath + "/.dep && tail -n " + strconv.Itoa(count) + " " + app.Meta.Name)
	if err != nil {
		return nil, err
	}

	lastReleases := strings.Split(*meta, "\n")

	var releases []string
	for _, str := range lastReleases {
		releases = append(releases, strings.Split(str, app.Meta.Separator)[1])
	}

	return &releases, err
}
func(t *task) cleanup(app *App) error {
	releases, err := t.getReleases(app, app.Release.KeepReleases)
	if err != nil {
		return err
	}

	scan, err := app.Bash.runOutput("cd " + app.Release.DeployPath + " && ls -xm releases/")
	if err != nil {
		return err
	}

	dirs := strings.Split(strings.Replace(*scan, "\r\n", " ", -1), app.Meta.Separator)

	var commands []string
	for _, dir := range dirs {
		exists, _ := InArray(dir, *releases)

		if exists == false {
			commands = append(commands, "cd " + app.Release.DeployPath + "/releases && rm -rf " + dir)
		}
	}

	commands = append(commands, "cd " + app.Release.DeployPath + " && if [ -e release ]; then rm release; fi")
	commands = append(commands, "cd " + app.Release.DeployPath + " && if [ -h release ]; then rm release; fi")

	return app.Bash.multiRun(commands)
}
func(t *task) deployShared(app *App) error {
	var err error
	for _, dir := range app.Release.Shared {
		if dir.IsDir == true {
			err = app.Bash.multiRun([]string{
				"mkdir -p " + app.Release.DeployPath + "/shared/" + dir.Path,
				"ln -s " + app.Release.DeployPath + "/shared/" + dir.Path + " " + app.Release.Path + "/" + dir.Path,
			})
			if err != nil {
				break
			}
		} else {
			var out *string
			out, err = app.Bash.runOutput("cd " + app.Release.DeployPath + "/shared/" + " && [[ -f " + dir.Path + " ]] && echo 'exist' || echo 'not_exist' ")
			if err != nil {
				break
			}

			if *out == "exist" {
				_, err = app.Bash.run("ln -s " + app.Release.DeployPath + "/shared/" + dir.Path + " " + app.Release.Path + "/" + dir.Path)
				if err != nil {
					break
				}
			}
		}
	}

	return err
}
func(t *task) deployWritable(app *App) error {
	var err error
	for _, dirPath := range app.Release.Writable {
		err = app.Bash.multiRun([]string{
			"chmod 0755 " + app.Release.DeployPath + "/shared/" + dirPath,
		})

		if err != nil {
			break
		}
	}

	return err
}

