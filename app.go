package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/urfave/cli/v2"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
)

type Release struct {
	PreviousReleasePath *string
	Path                string
	Name                string
	Number              int64
	Stage               string
	Branch              string
	Repository          string
	KeepReleases        int
}



type App struct {
	Color Color
	Bash  Bash
	Release		  Release

	Debug bool `default:"false"`
	Config string `default:"deploy.json"`
	MetaSeparator string `default:", "`
}

type Host struct {
	Host       string `json:"host" validate:"required,min=5"`
	Branch     string `json:"branch" validate:"required,min=1"`
	Stage      string `json:"stage" validate:"required,min=1"`
	User       string `json:"user" validate:"required,min=1"`
	Port       int    `json:"port" validate:"required,min=1"`
	DeployPath string `json:"deploy_path" validate:"required,min=1"`
}



type Task struct {
	Name string `json:"name"`
	Command string `json:"command"`
}
type Tasks struct {
	Remote []Task  `json:"remote"`
	Local []Task `json:"local"`
}
type Config struct {
	Repository string `json:"repository" validate:"required,min=10"`
	Hosts      []Host `json:"hosts" validate:"required,min=1"`
	KeepReleases int `json:"keep_releases"`
	TaskOrder []string `json:"task_order" validate:"required,min=1"`
	Tasks Tasks `json:"tasks"`
	//Commands Commands
}


func NewApp(color Color) App  {
	return App{
		Color:color,
		//Config: "deploy.json",
		//MetaSeparator: ", ",
		//KeepReleases: 10,
	}
}

func(app *App) defaultKeyPath() string {
	home := os.Getenv("HOME")
	if len(home) > 0 {
		return path.Join(home, ".ssh/id_rsa")
	}
	return ""
}

func(app *App) printMessageTask(task string) {
	fmt.Println(app.Color.Code.Green + "➤" + app.Color.Code.Default + " Executing task " + app.Color.Code.Green + task + app.Color.Code.Default)
}
func (app *App) sshExecute(command string, in chan<- string, out <-chan string) string {
	if app.Debug == true {
		fmt.Println("[" + app.Host.Host + "]" + " " + app.Color.Code.Teal + "> " + app.Color.Code.Default + command)
	}

	in <- command

	output := <-out

	outputLen := len(output)
	if outputLen > 4 {
		output = output[:outputLen-4]
	}

	return output
}
func(app *App) runCommand(commands Commands) error {
	var pk = app.defaultKeyPath()

	key, err := ioutil.ReadFile(pk)
	if err != nil {
		return errors.New("Not found ssh key file[" + pk + "]")
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return errors.New("Not valid ssh key file[" + pk + "]")
	}

	config := &ssh.ClientConfig{
		User:            app.Host.User,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
	}

	addr := fmt.Sprintf("%s:%d", app.Host.Host, app.Host.Port)

	client, err := ssh.Dial("tcp", addr, config)

	if err != nil {
		return err
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		return err
	}

	w, err := session.StdinPipe()
	if err != nil {
		return err
	}
	r, err := session.StdoutPipe()
	if err != nil {
		return err
	}

	in, out := app.ssh(w, r)
	if err := session.Start("/bin/sh"); err != nil {
		return err
	}

	<-out //ignore the shell output
	for name, callback := range commands.Remote {
		app.printMessageTask(name)
		if err := callback(app, in, out); err != nil {
			in <- "exit"
			return err
		}
	}
	in <- "exit"

	err = session.Wait()

	if err != nil {
		return err
	}

	return nil
}


func(app *App) ssh(w io.Writer, r io.Reader) (chan<- string, <-chan string) {
	var wg sync.WaitGroup



	//sync.WaitGroup.
	in := make(chan string, 1)
	out := make(chan string, 1)

	wg.Add(1)



	go func() {
		for cmd := range in {
			wg.Add(1)
			w.Write([]byte(cmd + "\n"))
			wg.Wait()
		}
	}()

	go func() {
		var (
			buf [65 * 1024]byte
			t   int
		)
		for {
			n, err := r.Read(buf[t:])
			if err != nil {
				close(in)
				close(out)
				return
			}
			t += n
			if buf[t-2] == '$' {
				out <- string(buf[:t])

				//fmt.Println("---------------")
				//fmt.Println(out)
				//fmt.Println("---------------")

				t = 0
				wg.Done()
			}
		}
	}()

	return in, out
}

func(app *App) prepare(c *cli.Context) error{
	stage := c.Args().First()
	app.Debug = c.Bool("debug")

	if len(stage) < 1 {
		return errors.New("you need to specify at least one host or stage")
	}

	app.Stage = stage

	jsonFile, err := os.Open(app.Config)
	if err != nil {
		app.Color.Print(app.Color.Info, "Not found configuration file [" + app.Config + "]")
		return err
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var deploy Config

	_ = json.Unmarshal(byteValue, &deploy)


	//fmt.Println(deploy)
	//
	//return errors.New("Not valid configuration file")

	var validate = validator.New()
	var errs = validate.Struct(deploy)

	if errs != nil {
		app.Color.Print(app.Color.Info, "Not valid configuration file [" + app.Config + "]")
		return errs
	}

	var errHosts[]string
	for i, row := range deploy.Hosts {
		var err = validate.Struct(row)

		if row.Stage == app.Stage {
			app.Host = row
		}

		if err != nil {
			errHosts = append(errHosts, err.Error() + " [host number = " + strconv.Itoa(i) + "]")
		}
	}

	if len(errHosts) > 0 {
		app.Color.Print(app.Color.Info, "Not valid configuration file [" + app.Config + "]")
		return errors.New(strings.Join(errHosts, "\n"))
	}

	if validate.Struct(app.Host) != nil {
		return errors.New("Not found stage[" + app.Stage + "] configuration file [" + app.Config + "]")
	}


	//command.Commands = deploy.Commands
	if deploy.KeepReleases > 0  {
		app.KeepReleases = deploy.KeepReleases
	}

	app.Repository = deploy.Repository
	app.TaskOrder = deploy.TaskOrder
	app.Tasks = deploy.Tasks



	return nil
}


