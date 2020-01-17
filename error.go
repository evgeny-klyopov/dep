package main

import "fmt"

const (
	ErrorEmptyStage = "empty.stage"
	NotFoundConfigurationFile = "not.found.configuration.file"
	NotValidConfigurationFile = "not.valid.configuration.file"
	NoSetBranch = "not.set.branch"
	NotFoundBranch = "not.found.branch"
)

var message = map[string]string{
	ErrorEmptyStage: "you need to specify at least one host or stage",
	NotFoundConfigurationFile: "Not found configuration file [%s]",
	NotValidConfigurationFile: "Not valid configuration file [%s]",
	NoSetBranch: "Not set branch from stage[%s] in configuration file [%s]",
	NotFoundBranch: "Not found stage[%s] configuration file [%s]",
}

func ErrorMessage(code string, args ...interface{}) string {
	if args[0] == nil {
		return message[code]
	}

	return fmt.Sprintf(message[code], args...)
}


