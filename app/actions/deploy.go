package actions

import (

	"github.com/urfave/cli/v2"
)

type deploy struct {
	action
}

func (d deploy) initCommand() *cli.Command {
	return &cli.Command{
		Name:        "deploy",
		Aliases:     []string{"dep"},
		ArgsUsage:   "[<stage>]",
		Description: "Deploy your project",
		Flags: []cli.Flag{
			d.getFlagDebug(),
		},
		Usage: "Deploy your project",
	}
}

func (d deploy) run(c *cli.Context) error {
	err := (*d.event).Prepare(c, 0)
	if err != nil {
		return err
	}

	orderTasks := (*d.event).GetOrderTasks(DeployAction)
	_ = (*d.event).CreateTasks(orderTasks)



	//event.InitTasks(DeployAction)
	//tasks := event.InitTasks(DeployAction)

	//event.InitTasks()





	//event.SetArguments(c.Args())
	//event.SetStage(0)
	//
	//err := event.CheckStage()
	//
	//if err != nil {
	//	return err
	//}
	//
	//
	////err := event.SetInputParams(c)
	////if  {
	////
	////}
	////event.setStage(c)
	//
	//
	//setting, err := event.GetSetting(d.appConfig.GetFilePathSetting())
	//
	////fmt.Println()
	//fmt.Println(setting)
	//fmt.Println(err)


	//fmt.Println(event)
	//d.printTaskName("deploy:prepare\n")

	//fmt.Println(d.appConfig)

	//
	//err := app.prepare(c)

	return nil
}
