package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
)



func main() {
	//PrintMessage(Magenta, "Test")

	//var color Color

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


	// EXAMPLE: Replace the `HelpPrinter` func
	//cli.HelpPrinter = func(w io.Writer, templ string, data interface{}) {
	//	fmt.Println("Ha HA.  I pwnd the help!!1")
	//}
	var authors[]*cli.Author

	authors = append(authors, &cli.Author{
		Name:  "ス",
		Email: "mail@klepov.info",
	})
	a :=  NewApp(color)

	fmt.Println(a)

	//
	//
	//app := &cli.App{
	//	Name: "Deployer",
	//	Usage: "A deployment tool",
	//	Version: "v1.0.0",
	//	Authors: authors,
	//	Commands: []*cli.Command{
	//		{
	//			Name:    "deploy",
	//			Aliases: []string{"dep"},
	//			Flags: []cli.Flag{
	//				&cli.BoolFlag{Name: "debug", Aliases: []string{"d"}},
	//			},
	//			Usage:   "Deploy",
	//			Action:  func(c *cli.Context) error {
	//				return a.deploy(c)
	//			},
	//		},
	//	},
	//}
	//
	//sort.Sort(cli.FlagsByName(app.Flags))
	//sort.Sort(cli.CommandsByName(app.Commands))
	//
	//err := app.Run(os.Args)
	//if err != nil {
	//	color.Print(color.Fatal, "Errors:")
	//	fmt.Print(color.Code.Red)
	//	log.Fatal(err)
	//}
}