package config

import (
	//"github.com/evgeny-klyopov/bashColor"
	//"github.com/evgeny-klyopov/dep/app"

	"github.com/evgeny-klyopov/bashColor"
	"github.com/urfave/cli/v2"
)

type cliConfig struct {
	abstract
	name                string
	usage               string
	appHelpTemplate     string
	commandHelpTemplate string
	authors             []*cli.Author
}

type CliConfigurator interface {
	Configurator
	GetName() string
	GetUsage() string
	setAppHelpTemplate() CliConfigurator
	GetAppHelpTemplate() string
	setCommandHelpTemplate() CliConfigurator
	GetCommandHelpTemplate() string
	GetAuthors() []*cli.Author
}

func NewCliConfig(version string, color *bashColor.Colorer) CliConfigurator {
	return (&cliConfig{
		abstract: abstract{
			version: version,
			color:   color,
		},
		name:    "Deployer",
		usage:   "A deployment tool",
		authors: []*cli.Author{
			&cli.Author{
				Name:  "ス",
				Email: "mail@klepov.info",
			},
		},
	}).setAppHelpTemplate().setCommandHelpTemplate()
}

func (c *cliConfig) setAppHelpTemplate() CliConfigurator {
	color := *(c.color)
	c.appHelpTemplate = color.White(`{{.Name}} - {{.Usage}}`)
	c.appHelpTemplate += color.Green(`{{if .Version}} {{.Version}}{{end}}`)
	c.appHelpTemplate += color.Yellow("\n\nUsage:") + `
	{{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}{{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}
{{if .Commands}}
` + color.Yellow("Commands:") + `
{{range .Commands}}{{if not .HideHelp}}` + "	" + color.GetColor(bashColor.Green) + `{{join .Names ", "}}` + color.GetColor(bashColor.Default) + `{{ "\t"}}{{.Usage}}{{ "\n" }}{{end}}{{end}}{{end}}{{if .VisibleFlags}}
` + color.Yellow("Global options:") + `
{{range .VisibleFlags}}  {{.}}
{{end}}{{end}}`

	return c
}

func (c *cliConfig) setCommandHelpTemplate() CliConfigurator {
	color := *(c.color)

	c.commandHelpTemplate = color.Yellow("Description:") + `
 {{.Usage}}

` + color.Yellow("Usage:") + `
 {{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}}{{if .VisibleFlags}} [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}
{{if .VisibleFlags}}
` + color.Yellow("Arguments:") + `
	` + color.GetColor(bashColor.Green) + `stage` + color.GetColor(bashColor.Default) + `{{ "\t"}}{{ "\t"}}{{ "\t"}}{{ "\t"}} Stage or hostname

` + color.Yellow("Options:") + `
 {{range .VisibleFlags}}{{.}}
 {{end}}{{end}}
`
	return c
}

func (c *cliConfig) GetAppHelpTemplate() string {
	return c.appHelpTemplate
}

func (c *cliConfig) GetCommandHelpTemplate() string {
	return c.commandHelpTemplate
}
func (c *cliConfig) GetName() string {
	return c.name
}

func (c *cliConfig) GetUsage() string {
	return c.usage
}

func (c *cliConfig) GetAuthors() []*cli.Author {
	return c.authors
}