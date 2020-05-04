package events

import "fmt"

type task struct {
	name string
	//remote bool
	//method func(app *App)error
}
type defaultTasks map[string]task
type Tasks []task

func (e *event) InitTasks(action string) {
	//var tasks Tasks
	var taskObject task

	defaultTasks := taskObject.defaultTasks()
	fmt.Println(defaultTasks)

	fmt.Println(action)
}

func(t *task) defaultTasks() defaultTasks {
	var tasks defaultTasks
	tasks = make(map[string]task)
	//tasks["deploy:check"] = task{"deploy:check", true, t.deployCheck}
	//tasks["deploy:release"] = task{"deploy:release", true, t.deployRelease}
	//tasks["deploy:update_code"] = task{"deploy:update_code", true, t.deployUpdateCode}
	//tasks["deploy:symlink"] = task{"deploy:symlink", true, t.deploySymlink}
	//tasks["cleanup"] = task{"cleanup", true, t.cleanup}
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
	//return app.Bash.multiRun([]string{
	//	"if [ ! -d " + app.Release.DeployPath + " ]; then mkdir -p " + app.Release.DeployPath + "; fi",
	//	"cd " + app.Release.DeployPath + " && if [ ! -d .dep ]; then mkdir .dep; fi",
	//	"cd " + app.Release.DeployPath + " && if [ ! -d releases ]; then mkdir releases; fi",
	//	"cd " + app.Release.DeployPath + " && if [ ! -d shared ]; then mkdir shared; fi",
	//
	//	"cd " + app.Release.DeployPath + " && if [ -e release ]; then rm release; fi",
	//	"cd " + app.Release.DeployPath + " && if [ -h release ]; then rm release; fi",
	//})

	return nil
}