package main

import (
	"strconv"
	"strings"
)

type Commands struct {
	Local map[string]func(app *App) error
	Remote map[string]func(app *App, in chan<- string, out <-chan string) error
}



func NewCommands() Commands {
	return Commands{}
}


func (c Commands) init(tasksOrder [] string, tasks Tasks) Commands {
	defaultTasks := c.getTasks()
	c.Remote = make(map[string]func(app *App, in chan<- string, out <-chan string) error)
	for _, name := range tasksOrder {
		if _, ok := defaultTasks.Remote[name]; ok {
			c.Remote[name] = defaultTasks.Remote[name]
		} else {
			for _, row := range tasks.Remote {
				if row.Name == name {
					c.Remote[name] = func(app *App, in chan<- string, out <-chan string) error {
						app.sshExecute(row.Command, in, out)
						return nil
					}
					break
				}
			}
		}
	}

	return c
}

func (c Commands) getTasks() Commands {
	var tasks Commands

	tasks.Remote = make(map[string]func(app *App, in chan<- string, out <-chan string) error)

	tasks.Remote["deploy:check"] = func(app *App, in chan<- string, out <-chan string) error {
		app.sshExecute("if [ ! -d " + app.Host.DeployPath + " ]; then mkdir -p " + app.Host.DeployPath + "; fi", in, out)
		app.sshExecute("cd " + app.Host.DeployPath + " && if [ ! -d .dep ]; then mkdir .dep; fi", in, out)
		app.sshExecute("cd " + app.Host.DeployPath + " && if [ ! -d releases ]; then mkdir releases; fi", in, out)
		//app.sshExecute("cd " + app.Host.DeployPath + " && if [ ! -d shared ]; then mkdir shared; fi", in, out)

		app.sshExecute("cd " + app.Host.DeployPath + " && if [ -e release ]; then rm release; fi", in, out)
		app.sshExecute("cd " + app.Host.DeployPath + " && if [ -h release ]; then rm release; fi", in, out)

		return nil
	}

	// Create release directory
	tasks.Remote["deploy:release"] = func(app *App, in chan<- string, out <-chan string) error {
		date := app.sshExecute(`date +"%Y-%m-%d %H:%M:%S"`, in, out)
		count := app.sshExecute("cd " + app.Host.DeployPath + "/releases && ls -l ./| grep ^d | wc -l", in, out)

		if count == "0" {
			app.Release = Release{
				Path:   app.Host.DeployPath + "/releases/1",
				Name:   "1",
				Number: 1,
			}
		} else {
			last := app.sshExecute("cd " + app.Host.DeployPath + "/.dep && tail -n 1 releases", in, out)
			previousRelease := strings.Split(last, app.MetaSeparator)[1]
			lastRelease, _ := strconv.ParseInt(previousRelease, 10, 64)

			app.Release = Release{
				Name:  strconv.FormatInt(lastRelease + 1, 10),
				Number: lastRelease + 1,
			}
			app.Release.Path = app.Host.DeployPath + "/releases/" + app.Release.Name

			previousReleasePath := app.Host.DeployPath + "/releases/" + previousRelease
			app.Release.PreviousReleasePath = &previousReleasePath
		}

		app.sshExecute("cd " + app.Host.DeployPath + `/.dep && echo "` + date + `, ` + app.Release.Name  + `" >> releases`, in, out)
		app.sshExecute("mkdir -p " + app.Release.Path, in, out)
		app.sshExecute("ln -nfs " + app.Release.Path + " " + app.Host.DeployPath + "/release", in, out)

		return nil
	}

	tasks.Remote["deploy:update_code"] = func(app *App, in chan<- string, out <-chan string) error {
		reference := ""
		if app.Release.PreviousReleasePath != nil {
			reference += " --reference " + *app.Release.PreviousReleasePath
		}

		app.sshExecute("git clone -b " + app.Host.Branch + " --recursive  " + reference +" --dissociate " + app.Repository + " " + app.Release.Path + " 2>&1", in, out)

		return nil
	}

	tasks.Remote["deploy:symlink"] = func(app *App, in chan<- string, out <-chan string) error {
		app.sshExecute("mv -T " + app.Host.DeployPath + "/release " + app.Host.DeployPath + "/current", in, out)

		return nil
	}

	tasks.Remote["cleanup"] = func(app *App, in chan<- string, out <-chan string) error {
		meta := app.sshExecute("cd " + app.Host.DeployPath + "/.dep && tail -n " + strconv.Itoa(app.KeepReleases) + " releases", in, out)
		lastReleases := strings.Split(meta, "\r\n")

		var releases []string
		for _, str := range lastReleases {
			releases = append(releases, strings.Split(str, app.MetaSeparator)[1])
		}


		scan := app.sshExecute("cd " + app.Host.DeployPath + " && ls -xm releases/", in, out)

		dirs := strings.Split(strings.Replace(scan, "\r\n", " ", -1), app.MetaSeparator)
		for _, dir := range dirs {
			exists, _ := InArray(dir, releases)

			if exists == false {
				app.sshExecute("cd " + app.Host.DeployPath + "/releases && rm -rf " + dir, in, out)
			}
		}

		app.sshExecute("cd " + app.Host.DeployPath + " && if [ -e release ]; then rm release; fi", in, out)
		app.sshExecute("cd " + app.Host.DeployPath + " && if [ -h release ]; then rm release; fi", in, out)

		return nil
	}

	return tasks
}

//
//func Init(tasks []string) Commands {
//	var Commands Commands
//
//	fmt.Println(tasks)
//
//	return Commands
//}



/*

func GetCommands() Commands {
	var Commands Commands

	// Create main directory if not exists
	Commands.Remote = append(Commands.Remote, func(app *App, in chan<- string, out <-chan string) error {
		app.printMessageTask("deploy:check")
		app.sshExecute("if [ ! -d " + app.Host.DeployPath + " ]; then mkdir -p " + app.Host.DeployPath + "; fi", in, out)
		app.sshExecute("cd " + app.Host.DeployPath + " && if [ ! -d .dep ]; then mkdir .dep; fi", in, out)
		app.sshExecute("cd " + app.Host.DeployPath + " && if [ ! -d releases ]; then mkdir releases; fi", in, out)
		//app.sshExecute("cd " + app.Host.DeployPath + " && if [ ! -d shared ]; then mkdir shared; fi", in, out)

		app.sshExecute("cd " + app.Host.DeployPath + " && if [ -e release ]; then rm release; fi", in, out)
		app.sshExecute("cd " + app.Host.DeployPath + " && if [ -h release ]; then rm release; fi", in, out)

		return nil
	})


	// Create release directory
	Commands.Remote = append(Commands.Remote, func(app *App, in chan<- string, out <-chan string) error {
		app.printMessageTask("deploy:release")

		date := app.sshExecute(`date +"%Y-%m-%d %H:%M:%S"`, in, out)
		count := app.sshExecute("cd " + app.Host.DeployPath + "/releases && ls -l ./| grep ^d | wc -l", in, out)

		if count == "0" {
			app.Release = Release{
				Path:   app.Host.DeployPath + "/releases/1",
				Name:   "1",
				Number: 1,
			}
		} else {
			last := app.sshExecute("cd " + app.Host.DeployPath + "/.dep && tail -n 1 releases", in, out)
			previousRelease := strings.Split(last, app.MetaSeparator)[1]
			lastRelease, _ := strconv.ParseInt(previousRelease, 10, 64)

			app.Release = Release{
				Name:  strconv.FormatInt(lastRelease + 1, 10),
				Number: lastRelease + 1,
			}
			app.Release.Path = app.Host.DeployPath + "/releases/" + app.Release.Name

			previousReleasePath := app.Host.DeployPath + "/releases/" + previousRelease
			app.Release.PreviousReleasePath = &previousReleasePath
		}

		app.sshExecute("cd " + app.Host.DeployPath + `/.dep && echo "` + date + `, ` + app.Release.Name  + `" >> releases`, in, out)
		app.sshExecute("mkdir -p " + app.Release.Path, in, out)
		app.sshExecute("ln -nfs " + app.Release.Path + " " + app.Host.DeployPath + "/release", in, out)

		return nil
	})

	Commands.Remote = append(Commands.Remote, func(app *App, in chan<- string, out <-chan string) error {
		app.printMessageTask("deploy:update_code")
		reference := ""
		if app.Release.PreviousReleasePath != nil {
			reference += " --reference " + *app.Release.PreviousReleasePath
		}

		app.sshExecute("git clone -b " + app.Host.Branch + " --recursive  " + reference +" --dissociate " + app.Repository + " " + app.Release.Path + " 2>&1", in, out)

		return nil
	})

	Commands.Remote = append(Commands.Remote, func(app *App, in chan<- string, out <-chan string) error {
		app.printMessageTask("deploy:symlink")

		app.sshExecute("mv -T " + app.Host.DeployPath + "/release " + app.Host.DeployPath + "/current", in, out)

		return nil
	})

	Commands.Remote = append(Commands.Remote, func(app *App, in chan<- string, out <-chan string) error {
		app.printMessageTask("cleanup")

		meta := app.sshExecute("cd " + app.Host.DeployPath + "/.dep && tail -n " + strconv.Itoa(app.KeepReleases) + " releases", in, out)
		lastReleases := strings.Split(meta, "\r\n")

		var releases []string
		for _, str := range lastReleases {
			releases = append(releases, strings.Split(str, app.MetaSeparator)[1])
		}


		scan := app.sshExecute("cd " + app.Host.DeployPath + " && ls -xm releases/", in, out)

		dirs := strings.Split(strings.Replace(scan, "\r\n", " ", -1), app.MetaSeparator)
		for _, dir := range dirs {
			exists, _ := InArray(dir, releases)

			if exists == false {
				app.sshExecute("cd " + app.Host.DeployPath + "/releases && rm -rf " + dir, in, out)
			}
		}

		app.sshExecute("cd " + app.Host.DeployPath + " && if [ -e release ]; then rm release; fi", in, out)
		app.sshExecute("cd " + app.Host.DeployPath + " && if [ -h release ]; then rm release; fi", in, out)

		return nil
	})


	return Commands
}

*/