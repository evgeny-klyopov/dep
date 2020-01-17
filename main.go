package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"sort"
)



func main() {
	color := NewColor()
	info := color.White(`{{.Name}} - {{.Usage}}`)
	info += color.Green(`{{if .Version}} {{.Version}}{{end}}`)

	// EXAMPLE: Override a template
	cli.AppHelpTemplate = info + `

` + color.Yellow("USAGE:") + `
  {{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}{{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}
  {{if .Commands}}
` + color.Yellow("COMMANDS:") + `
{{range .Commands}}{{if not .HideHelp}}` + color.Code.Green + `{{join .Names ", "}}` + color.Code.Default + `{{ "\t"}}{{.Usage}}{{ "\n" }}{{end}}{{end}}{{end}}{{if .VisibleFlags}}
` + color.Yellow("GLOBAL OPTIONS:") + `
  {{range .VisibleFlags}}{{.}}
  {{end}}{{end}}{{if len .Authors}}
` + color.Yellow("AUTHOR:") + `
  {{range .Authors}}{{ . }}{{end}}
  {{end}}
`

	var authors[]*cli.Author

	authors = append(authors, &cli.Author{
		Name:  "ス",
		Email: "mail@klepov.info",
	})

	app :=  NewApp(color)

	cliApp := &cli.App{
		Name: "Deployer",
		Usage: "A deployment tool",
		Version: "v1.0.0",
		Authors: authors,
		Commands: []*cli.Command{
			{
				Name:    "deploy",
				Aliases: []string{"dep"},
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "debug", Aliases: []string{"d"}},
				},
				Usage:   "Deploy",
				Action:  func(c *cli.Context) error {
					return app.deploy(c)
				},
			},
		},
	}

	sort.Sort(cli.FlagsByName(cliApp.Flags))
	sort.Sort(cli.CommandsByName(cliApp.Commands))

	err := cliApp.Run(os.Args)
	if err != nil {
		color.Print(color.Fatal, "Errors:")
		fmt.Print(color.Code.Red)
		log.Fatal(err)
	}
}