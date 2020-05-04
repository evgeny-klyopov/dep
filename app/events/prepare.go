package events

import (
	"github.com/urfave/cli/v2"
)

func(e *event) Prepare(c *cli.Context, numberStageArguments int) error {
	e.printEventName("deploy:prepare\n")

	e.SetArguments(c.Args())
	e.setStage(numberStageArguments)

	err := e.checkStage()
	if err != nil {
		return err
	}

	err = e.readDeployConfig()
	if err != nil {
		return err
	}

	err = e.setHost()
	if err != nil {
		return err
	}

	e.setKeepRelease()

	err = e.checkRepositoryConfig()
	if err != nil {
		return err
	}

	e.setVariables()

	return nil
}

