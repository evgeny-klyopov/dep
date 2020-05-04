package app

import (
	"github.com/evgeny-klyopov/bashColor"
	"github.com/evgeny-klyopov/dep/app/actions"
	"github.com/evgeny-klyopov/dep/app/config"
	"github.com/urfave/cli/v2"
	"sort"
)

type deploy struct {
	color *bashColor.Colorer
	config *config.CliConfigurator
}


func (d *deploy) GetConfig() *config.CliConfigurator {
	return d.config
}

func (d *deploy) GetColor() *bashColor.Colorer {
	return d.color
}

type Deployer interface {
	GetCommands() []*cli.Command
	GetConfig() *config.CliConfigurator
	GetColor() *bashColor.Colorer
	GetCliApp() *cli.App
}

func NewApp(version string) Deployer {
	var dep Deployer

	color := bashColor.NewColor()
	cfg := config.NewCliConfig(version, &color)
	dep = &deploy{
		color:  &color,
		config: &cfg,
	}

	return dep
}

func (d * deploy) GetCommands() []*cli.Command {
	return actions.GetCommands(d.GetColor(), (*d.GetConfig()).GetVersion())
}

func (d *deploy) GetCliApp() *cli.App {
	cfg := *d.GetConfig()
	cli.AppHelpTemplate = cfg.GetAppHelpTemplate()
	cli.CommandHelpTemplate = cfg.GetCommandHelpTemplate()

	cliApp := &cli.App{
		Name:     cfg.GetName(),
		Usage:     cfg.GetUsage(),
		Version:  cfg.GetVersion(),
		Authors:  cfg.GetAuthors(),
		Commands: d.GetCommands(),
	}

	sort.Sort(cli.FlagsByName(cliApp.Flags))
	sort.Sort(cli.CommandsByName(cliApp.Commands))

	return cliApp
}