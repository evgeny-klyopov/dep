package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/evgeny-klyopov/bashColor"
	"github.com/go-playground/validator/v10"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Release struct {
	DeployPath          string `validate:"required,min=5"`
	PreviousReleasePath *string
	Path                string
	Name                string
	Rollback            string
	Number              int64
	Stage               string `validate:"required,min=2"`
	Branch              string
	Repository          string
	KeepReleases        int `validate:"required,min=5"`
	LocalObjectPath     []string
	Shared              []Shared
	Writable            []string
	Variables           map[string]string
}

type App struct {
	Color       bashColor.Colorer
	Bash        Bash
	Release     Release
	Debug       bool
	Config      string
	ArchiveName string
	TasksOrder  []string
	Meta        Meta
	ConfigTasks ConfigTasks
}

type Host struct {
	Host       string `json:"host" validate:"required,min=5"`
	Branch     string `json:"branch"`
	Stage      string `json:"stage" validate:"required,min=1"`
	User       string `json:"user" validate:"required,min=1"`
	Port       int    `json:"port" validate:"required,min=1"`
	DeployPath string `json:"deploy_path" validate:"required,min=1"`
}

type Shared struct {
	Path  string `json:"path"`
	IsDir bool   `json:"is_dir"`
}

type Config struct {
	Repository      string            `json:"repository"`
	LocalObjectPath []string          `json:"local_object_path"`
	Hosts           []Host            `json:"hosts" validate:"required,min=1"`
	KeepReleases    int               `json:"keep_releases"`
	TasksOrder      []string          `json:"tasks_order" validate:"required,min=1"`
	Shared          []Shared          `json:"shared"`
	Writable        []string          `json:"writable"`
	ConfigTasks     ConfigTasks       `json:"tasks"`
	Variables       map[string]string `json:"variables"`
}
type ConfigTask struct {
	Name    string `json:"name" validate:"required,min=1"`
	Command string `json:"command" validate:"required,min=1"`
}
type ConfigTasks struct {
	Remote []ConfigTask `json:"remote"`
	Local  []ConfigTask `json:"local"`
}

type Meta struct {
	Name      string
	Separator string
}

func NewApp() App {
	return App{
		Color:       bashColor.NewColor(),
		Debug:       false,
		TasksOrder:  defaultTasksOrder,
		Config:      "deploy.json",
		ArchiveName: "deploy.tar.gz",
		Meta: Meta{
			Name:      "meta",
			Separator: ", ",
		},
		Release: Release{
			KeepReleases: 10,
		},
	}
}

func (app *App) GetVersion() string {
	return "v1.0.7"
}
func (app *App) HelpTemplate() (appHelp string, commandHelp string) {
	info := app.Color.White(`{{.Name}} - {{.Usage}}`)
	info += app.Color.Green(`{{if .Version}} {{.Version}}{{end}}`)

	appHelp = info + `

` + app.Color.Yellow("Usage:") + `
	{{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}{{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}
 {{if .Commands}}
` + app.Color.Yellow("Commands:") + `
{{range .Commands}}{{if not .HideHelp}}` + "	" + app.Color.GetColor(bashColor.Green) + `{{join .Names ", "}}` + app.Color.GetColor(bashColor.Default) + `{{ "\t"}}{{.Usage}}{{ "\n" }}{{end}}{{end}}{{end}}{{if .VisibleFlags}}
` + app.Color.Yellow("Global options:") + `
{{range .VisibleFlags}}  {{.}}
{{end}}{{end}}`

	commandHelp = app.Color.Yellow("Description:") + ` 
   {{.Usage}}

` + app.Color.Yellow("Usage:") + `
   {{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}}{{if .VisibleFlags}} [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}
{{if .VisibleFlags}}
` + app.Color.Yellow("Arguments:") + `
	` + app.Color.GetColor(bashColor.Green) + `stage` + app.Color.GetColor(bashColor.Default) + `{{ "\t"}}{{ "\t"}}{{ "\t"}}{{ "\t"}} Stage or hostname

` + app.Color.Yellow("Options:") + `
   {{range .VisibleFlags}}{{.}}
   {{end}}{{end}}
`
	return appHelp, commandHelp

}
func (app *App) printTaskName(task string) {
	fmt.Print(app.Color.GetColor(bashColor.Green) + "➤" +  app.Color.GetColor(bashColor.Default) + " Executing task " + app.Color.GetColor(bashColor.Green) + task +  app.Color.GetColor(bashColor.Default))
}

func (app *App) cmd(command string, args ...string) error {
	prefixDebug := "[local]" + " " + app.Color.GetColor(bashColor.Yellow) + "> " + app.Color.GetColor(bashColor.Default)

	if app.Debug == true {
		fmt.Println(prefixDebug + command + " " + strings.Join(args, " "))
	}

	cmd := exec.Command(command, args...)

	return cmd.Run()
}
func (app *App) printTimer(prefix string, timer time.Time) {
	fmt.Println(prefix + app.Color.GetColor(bashColor.Purple) + fmt.Sprintf("%v", time.Since(timer)) + app.Color.GetColor(bashColor.Default))
}
func (app *App) run(tasks Tasks) error {
	prefixTaskTime := " - "
	totalTimer := time.Now()
	taskTimer := time.Now()

	if app.Debug == true {
		prefixTaskTime = "Time: "
		fmt.Print("\n")
	}

	for _, task := range tasks {
		app.printTaskName(task.name)
		if app.Debug == true {
			fmt.Print("\n")
		}

		if err := task.method(app); err != nil {
			_ = app.Bash.close()

			return err
		}
		app.printTimer(prefixTaskTime, taskTimer)

		taskTimer = time.Now()
	}
	app.printTimer("Total - ", totalTimer)

	return app.Bash.close()
}

func (app *App) error(code string, print bool, err error, args ...interface{}) error {
	if code != "" {
		message := ErrorMessage(code, args...)

		if print == true {
			app.Color.Print(app.Color.Info, message)
		}

		if err == nil {
			err = errors.New(message)
		}
	}

	return err
}

func (app *App) prepare(c *cli.Context) error {
	stage := c.Args().First()
	app.Debug = c.Bool("debug")

	if len(stage) < 1 {
		return app.error(ErrorEmptyStage, false, nil, nil)
	}

	app.Release.Stage = stage

	jsonFile, err := os.Open(app.Config)
	if err != nil {
		return app.error(NotFoundConfigurationFile, true, err, app.Config)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var config Config
	err = json.Unmarshal(byteValue, &config)
	if err != nil {
		return app.error(NotValidConfigurationFile, true, err, app.Config)
	}

	var validate = validator.New()
	var errs = validate.Struct(config)
	if errs != nil {
		return app.error(NotValidConfigurationFile, true, errs, app.Config)
	}

	var errHosts []string
	for i, host := range config.Hosts {
		var err = validate.Struct(host)

		if host.Stage == app.Release.Stage {
			prefixDebug := "[" + host.Host + "]" + " " + app.Color.GetColor(bashColor.Teal) + "> " + app.Color.GetColor(bashColor.Default)

			app.Bash = NewBash(host.User, host.Host, host.Port, app.Debug, prefixDebug)
			app.Release.DeployPath = host.DeployPath
			app.Release.Branch = host.Branch

		}

		if err != nil {
			errHosts = append(errHosts, err.Error()+" [host number = "+strconv.Itoa(i)+"]")
		}
	}
	app.Release.Repository = config.Repository
	if config.KeepReleases > 0 {
		app.Release.KeepReleases = config.KeepReleases
	}

	if len(errHosts) > 0 {
		return app.error(NotValidConfigurationFile, true, errors.New(strings.Join(errHosts, "\n")), app.Config)
	}

	if validate.Struct(app.Release) != nil {
		return app.error(NotFoundBranch, false, nil, app.Release.Stage, app.Config)
	}

	if app.Release.Branch == "" && app.Release.Repository != "" {
		return app.error(NoSetBranch, false, nil, app.Release.Stage, app.Config)
	}

	app.TasksOrder = config.TasksOrder
	app.ConfigTasks = config.ConfigTasks
	app.Release.LocalObjectPath = config.LocalObjectPath
	app.Release.Shared = config.Shared
	app.Release.Writable = config.Writable

	app.Release.Variables = make(map[string]string)
	app.Release.Variables["{{release_path}}"] = app.Release.DeployPath + "/release"
	app.Release.Variables["{{stage}}"] = app.Release.Stage
	for k, v := range config.Variables {
		app.Release.Variables[k] = v
	}

	return nil
}
