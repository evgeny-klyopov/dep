package actions

import (
	"fmt"
	"github.com/evgeny-klyopov/bashColor"
	"github.com/evgeny-klyopov/dep/app/config"
	"github.com/evgeny-klyopov/dep/app/events"
	"github.com/urfave/cli/v2"
)

const DeployAction = "deploy"
const RollbackAction = "rollback"

type action struct {
	debug             bool
	event *events.Eventer
	color             *bashColor.Colorer
	command           *cli.Command
	appConfig config.AppConfigurator
}

type Commander interface {
	getEvent() *events.Eventer
	initCommand() *cli.Command
	getCommand() *cli.Command
	setCommand(Commander) Commander
	setDebugFlag(c *cli.Context)
	run(c *cli.Context) error
}

func GetCommands(color *bashColor.Colorer, version string) []*cli.Command {
	appConfig := config.NewAppConfig(version, color)
	return []*cli.Command{

		newCommand(&deploy{
			action{
				color:             color,
				event: events.NewEvent(&appConfig),
				appConfig: appConfig,
			},
		}),
	}
}

func newCommand(c Commander) *cli.Command {
	return c.setCommand(c).getCommand()
}

func (a *action) setDebugFlag(c *cli.Context) {
	a.debug = c.Bool("debug")
}

func (a *action) getEvent() *events.Eventer{
	return a.event
}

func (a *action) setCommand(c Commander) Commander {
	a.command = c.initCommand()
	a.command.Action = func(ctx *cli.Context) error {
		c.setDebugFlag(ctx)
		return c.run(ctx)
	}
	return c
}

func (a *action) getFlagDebug() *cli.BoolFlag {
	return &cli.BoolFlag{
		Name:    "debug",
		Usage:   "Debug mode",
		Aliases: []string{"d"},
	}
}
func (a *action) getCommand() *cli.Command {
	return a.command
}

func (a *action) printTaskName(task string) {
	color := *(a.color)
	fmt.Print(color.GetColor(bashColor.Green) + "➤" + color.GetColor(bashColor.Default) + " Executing task " + color.GetColor(bashColor.Green) + task + color.GetColor(bashColor.Default))
}
