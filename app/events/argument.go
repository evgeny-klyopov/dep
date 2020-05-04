package events

import "github.com/urfave/cli/v2"

type argument struct {
	arguments cli.Args
}

func (a *argument) SetArguments(args cli.Args) {
	a.arguments = args
}
type Argumenter interface {
	SetArguments(args cli.Args)
}