package actions

import (
	"fmt"
	"github.com/evgeny-klyopov/bashColor"
	"github.com/evgeny-klyopov/dep/app/config"
	"github.com/evgeny-klyopov/dep/app/events"
	"github.com/urfave/cli/v2"
)

const DeployAction = "Deploy"
const RollbackAction = "Rollback"

type action struct {
	event *events.Eventer
	color             *bashColor.Colorer
	command           *cli.Command
	appConfig config.AppConfigurator
}

type Commander interface {
	//getEvent() *events.Eventer
	initCommand() *cli.Command
	getCommand() *cli.Command
	setCommand(Commander) Commander
	setFlags(c *cli.Context)
	run(c *cli.Context) error
}

func GetCommands(color *bashColor.Colorer, version string) []*cli.Command {
	appConfig := config.NewAppConfig(version, color)

	action := action{
		color:             color,
		appConfig: appConfig,
		event: events.NewEvent(&appConfig),
	}

	return []*cli.Command{
		newCommand(&deploy{action}),
	}
}

func newCommand(c Commander) *cli.Command {
	return c.setCommand(c).getCommand()
}

func (a *action) setFlags(c *cli.Context) {
	a.appConfig.SetDebugFlag(c.Bool("debug"))
}

//func (a *action) getEvent() *events.Eventer{
//	return a.event
//}

func (a *action) setCommand(c Commander) Commander {
	a.command = c.initCommand()
	a.command.Action = func(ctx *cli.Context) error {
		c.setFlags(ctx)





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
