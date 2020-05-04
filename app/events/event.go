package events

import (
	"fmt"
	"github.com/evgeny-klyopov/bashColor"
	"github.com/evgeny-klyopov/dep/app/config"
	"github.com/urfave/cli/v2"
)

type event struct {
	setting
	argument
	stage string
	config *config.AppConfigurator
}
//
//func (e event) CheckInputParams() error {
//	return nil
//}





func NewEvent(cfg *config.AppConfigurator) *Eventer {
	var evt Eventer
	evt = &event{
		config: cfg,
	}
	return &evt
}



func (e *event) printEventName(name string) {
	color := *(*(e.config)).GetColor()

	fmt.Print(color.GetColor(bashColor.Green) +
		"➤" + color.GetColor(bashColor.Default) +
		" Executing task " + color.GetColor(bashColor.Green) +
		name + color.GetColor(bashColor.Default))
}
func (e *event) checkStage() error {
	if len(e.stage) < 1 {
		return e.getError(ErrorEmptyStage, false, nil, nil)
	}

	return nil
}
func (e *event) setStage(numberStageArguments int) {
	e.stage = e.arguments.Get(numberStageArguments)
}




type Eventer interface {
	settinger
	Error
	Argumenter
	setStage(n int)
	checkStage() error
	printEventName(eventName string)
	Prepare(c *cli.Context, n int) error
	InitTasks(typeOperation string)
}
