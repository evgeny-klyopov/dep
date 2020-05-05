package events

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type task struct {
	name string
	method func(app *event)error
}
type defaultTasks map[string]task
type Tasks []task

func (e *event) GetOrderTasks(action string)[]string {
	var orderTasks []string
	r := reflect.ValueOf(e.deployConfig.OrderTasks)
	k := r.FieldByName(action)

	orderTasks, _ = k.Interface().([]string)

	return orderTasks
}
func (e *event) CreateTasks(orderTasks []string) Tasks {
	var tasks Tasks



	var taskObject task

	defaultTasks := taskObject.defaultTasks()
	fmt.Println(defaultTasks)

	//fmt.Println(action)


	return tasks
}

func(t *task) defaultTasks() defaultTasks {
	var tasks defaultTasks
	tasks = make(map[string]task)
	tasks["deploy:check"] = task{"deploy:check", t.deployCheck}
	tasks["deploy:release"] = task{"deploy:release",  t.deployRelease}
	tasks["deploy:update_code"] = task{"deploy:update_code", t.deployUpdateCode}
	tasks["deploy:symlink"] = task{"deploy:symlink", t.deploySymlink}
	//tasks["cleanup"] = task{"cleanup", t.cleanup}
	//
	//tasks["local:create-archive"] = task{"local:create-archive", false, t.localCreateArchive}
	//tasks["deploy:extract-archive"] = task{"deploy:extract-archive", true, t.deployExtractArchive}
	//tasks["local:send-archive"] = task{"local:send-archive", false, t.localSendArchive}
	//
	//
	//tasks["rollback:check"] = task{"rollback:check", true, t.rollbackCheck}
	//tasks["rollback:symlink"] = task{"rollback:symlink", true, t.rollbackSymlink}
	//tasks["rollback:cleanup"] = task{"rollback:cleanup", true, t.rollbackCleanup}
	//
	//
	//tasks["deploy:shared"] = task{"deploy:shared", true, t.deployShared}
	//tasks["deploy:writable"] = task{"deploy:writable", true, t.deployWritable}

	return tasks
}

func(t *task) deployCheck(e *event) error {
	host := e.getHost()

	return (*e.bash).MultiRun([]string{
		"if [ ! -d " + host.DeployPath + " ]; then mkdir -p " + host.DeployPath + "; fi",
		"cd " + host.DeployPath + " && if [ ! -d .dep ]; then mkdir .dep; fi",
		"cd " + host.DeployPath + " && if [ ! -d releases ]; then mkdir releases; fi",
		"cd " + host.DeployPath + " && if [ ! -d shared ]; then mkdir shared; fi",

		"cd " + host.DeployPath + " && if [ -e release ]; then rm release; fi",
		"cd " + host.DeployPath + " && if [ -h release ]; then rm release; fi",
	})
}

type release struct {
	path string
	name string
	number int64
	previousReleasePath *string
}

func(t *task) deployRelease(e *event) error {
	host := e.getHost()
	meta := (*e.config).GetMetaSetting()

	rewDate, err := (*e.bash).Run(`date +"%Y-%m-%d %H:%M:%S"`)
	if err != nil {
		return err
	}
	date := rewDate.ToString()

	rewCount, err := (*e.bash).Run("cd " + host.DeployPath + "/releases && ls -l ./| grep ^d | wc -l")
	if err != nil {
		return err
	}
	count := rewCount.ToString()

	if *count == "0" {
		e.release.path = host.DeployPath + "/releases/1"
		e.release.name = "1"
		e.release.number = 1
	} else {
		rawLast, err := (*e.bash).Run("cd " + host.DeployPath + "/.dep && tail -n 1 " + meta.Name)
		if err != nil {
			return err
		}
		last := rawLast.ToString()

		previousRelease := strings.Split(*last, meta.Separator)[1]
		lastRelease, _ := strconv.ParseInt(previousRelease, 10, 64)

		e.release.name = strconv.FormatInt(lastRelease + 1, 10)
		e.release.number = lastRelease + 1

		e.release.path = host.DeployPath + "/releases/" + e.release.name

		previousReleasePath := host.DeployPath + "/releases/" + previousRelease
		e.release.previousReleasePath = &previousReleasePath
	}

	return (*e.bash).MultiRun([]string{
		"cd " + host.DeployPath + `/.dep && echo "` + *date + `, ` + e.release.name  + `" >> ` + meta.Name,
		"mkdir -p " + e.release.path,
		"ln -nfs " + e.release.path  + " " + host.DeployPath + "/release",
	})
}

func(t *task) deployUpdateCode(e *event) error {
	host := e.getHost()
	reference := ""
	if e.release.previousReleasePath != nil {
		reference += " --reference " + *e.release.previousReleasePath
	}

	return (*e.bash).MultiRun([]string{
		"git clone -b " + *host.Branch + " --recursive  " + reference +" --dissociate " + *e.deployConfig.Repository + " " + e.release.path + " 2>&1",
	})
}

func(t *task) getCommandSymlink(from string, to string) string {
	return "mv -T " + from + " " + to
}
func(t *task) deploySymlink(e *event) error {
	host := e.getHost()
	return (*e.bash).MultiRun([]string{
		t.getCommandSymlink(host.DeployPath + "/release", host.DeployPath + "/current"),
	})
}

//func(t *task) getReleases(e *event, count int) (*[]string, error) {
//	meta, err := app.Bash.runOutput("cd " + app.Release.DeployPath + "/.dep && tail -n " + strconv.Itoa(count) + " " + app.Meta.Name)
//	if err != nil {
//		return nil, err
//	}
//
//	lastReleases := strings.Split(*meta, "\n")
//
//	var releases []string
//	for _, str := range lastReleases {
//		releases = append(releases, strings.Split(str, app.Meta.Separator)[1])
//	}
//
//	return &releases, err
//}
//
//func(t *task) cleanup(e *event) error {
//	kp := e.deployConfig.KeepReleases
//	releases, err := t.getReleases(e, *kp)
//	if err != nil {
//		return err
//	}
//
//	scan, err := app.Bash.runOutput("cd " + app.Release.DeployPath + " && ls -xm releases/")
//	if err != nil {
//		return err
//	}
//
//	dirs := strings.Split(strings.Replace(*scan, "\r\n", " ", -1), app.Meta.Separator)
//
//	var commands []string
//	for _, dir := range dirs {
//		exists, _ := InArray(dir, *releases)
//
//		if exists == false {
//			commands = append(commands, "cd " + app.Release.DeployPath + "/releases && rm -rf " + dir)
//		}
//	}
//
//	commands = append(commands, "cd " + app.Release.DeployPath + " && if [ -e release ]; then rm release; fi")
//	commands = append(commands, "cd " + app.Release.DeployPath + " && if [ -h release ]; then rm release; fi")
//
//	return app.Bash.multiRun(commands)
//}