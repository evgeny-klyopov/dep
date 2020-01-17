package main

import (
	"deploy/app"
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"sort"

	//"fmt"
	//"github.com/urfave/cli/v2"
	//"log"
	//"os"
	//"sort"
)



func main() {
	dep :=  app.NewApp()



	info := dep.Color.White(`{{.Name}} - {{.Usage}}`)
	info += dep.Color.Green(`{{if .Version}} {{.Version}}{{end}}`)

	// EXAMPLE: Override a template
	cli.AppHelpTemplate = info + `

` + dep.Color.Yellow("USAGE:") + `
 {{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}{{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}
 {{if .Commands}}
` + dep.Color.Yellow("COMMANDS:") + `
{{range .Commands}}{{if not .HideHelp}}` + dep.Color.Code.Green + `{{join .Names ", "}}` + dep.Color.Code.Default + `{{ "\t"}}{{.Usage}}{{ "\n" }}{{end}}{{end}}{{end}}{{if .VisibleFlags}}
` + dep.Color.Yellow("GLOBAL OPTIONS:") + `
 {{range .VisibleFlags}}{{.}}
 {{end}}{{end}}{{if len .Authors}}
` + dep.Color.Yellow("AUTHOR:") + `
 {{range .Authors}}{{ . }}{{end}}
 {{end}}
`

	var authors[]*cli.Author

	authors = append(authors, &cli.Author{
		Name:  "ス",
		Email: "mail@klepov.info",
	})


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
					return dep.Deploy(c)
				},
			},
			//'Rollback to previous release'
			{
				Name:    "rollback",
				Aliases: []string{"rol"},
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "debug", Aliases: []string{"d"}},
				},
				Usage:   "Rollback",
				Action:  func(c *cli.Context) error {
					return dep.Rollback(c)
				},

			},
		},
	}

	sort.Sort(cli.FlagsByName(cliApp.Flags))
	sort.Sort(cli.CommandsByName(cliApp.Commands))

	err := cliApp.Run(os.Args)
	if err != nil {
		dep.Color.Print(dep.Color.Fatal, "Errors:")
		fmt.Print(dep.Color.Code.Red)
		log.Fatal(err)
	}
}