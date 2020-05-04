package main

import (
	"fmt"
	"github.com/evgeny-klyopov/bashColor"
	"github.com/evgeny-klyopov/dep/app"
	"log"
	"os"
)

var version string

func init() {
	version = "v2.0.0-alpha"
}

func main() {
	dep := app.NewApp(version)

	color := *dep.GetColor()

	err := dep.GetCliApp().Run(os.Args)

	if err != nil {
		color.Print(color.Fatal, "Errors:")
		fmt.Print(color.GetColor(bashColor.Red))
		log.Fatal(err)
	}
}
