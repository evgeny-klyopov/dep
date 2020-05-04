package events

import (
	"errors"
	"fmt"
)

const (
	ErrorEmptyStage           = "empty.stage"
	NotFoundConfigurationFile = "not.found.configuration.file"
	NotValidConfigurationFile = "not.valid.configuration.file"
	NoSetBranch               = "not.set.branch"
	NotFoundBranch            = "not.found.branch"
	NotFoundStage            = "not.found.stage"
	NotFoundPreviousRelease   = "not.found.previous.release"
)

var message = map[string]string{
	ErrorEmptyStage: "you need to specify at least one host or stage",
	NotFoundConfigurationFile: "Not found configuration file [%s]",
	NotValidConfigurationFile: "Not valid configuration file [%s]",
	NoSetBranch: "Not set branch from stage[%s] in configuration file [%s]",
	NotFoundBranch: "Not found stage[%s] configuration file [%s]",
	NotFoundStage: "Not found stage[%s] configuration file [%s]",
	NotFoundPreviousRelease: "Not found previous release",
}

type Error interface {
	errorMessage(code string, args ...interface{}) string
	getError(code string, print bool, err error, args ...interface{}) error
}

func (e *event) errorMessage(code string, args ...interface{}) string {
	if args[0] == nil {
	return message[code]
	}

	return fmt.Sprintf(message[code], args...)
}

func (e *event) getError(code string, print bool, err error, args ...interface{}) error {
	if code != "" {
		message := e.errorMessage(code, args...)

		if print == true {
			color := *(*(e.config)).GetColor()
			color.Print(color.Info, message)
		}

		if err == nil {
			err = errors.New(message)
		}
	}

	return err
}