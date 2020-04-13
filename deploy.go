package main

import (
	"fmt"
	"github.com/evgeny-klyopov/dep/app"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"sort"
)

func main() {
	dep :=  app.NewApp()

	appHelp, commandHelp := dep.HelpTemplate()

	cli.AppHelpTemplate = appHelp
	cli.CommandHelpTemplate = commandHelp

	var authors[]*cli.Author
	authors = append(authors, &cli.Author{
		Name:  "ス",
		Email: "mail@klepov.info",
	})

	cliApp := &cli.App{
		Name: "Deployer",
		Usage: "A deployment tool",
		Version: dep.GetVersion(),
		Authors: authors,
		Commands: []*cli.Command{
			{
				Name:    "deploy",
				Aliases: []string{"dep"},
				ArgsUsage: "[<stage>]",
				Description: "Deploy your project",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name: "debug",
						Usage:   "Debug mode",
						Aliases: []string{"d"},
					},
				},
				Usage:   "Deploy your project",
				Action:  func(c *cli.Context) error {
					return dep.Deploy(c)
				},
			},
			{
				Name:    "rollback",
				Aliases: []string{"rol"},
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "debug", Aliases: []string{"d"}},
				},
				Usage:   "Rollback to previous release",
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