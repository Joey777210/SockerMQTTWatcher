package command

import (
	"fmt"
	"github.com/urfave/cli"
	"os"
)

var StartCommand = cli.Command{
	Name:                   "start",
	Usage:                  "start mqtt watcher",
	Action: func(context *cli.Context) error {
		if context.Args().Len() < 1 {
			return fmt.Errorf("Missing mqtt command")
		}
		var cmdArray []string
		for i := 0; i < context.Args().Len(); i++ {
			arg := context.Args().Get(i)
			cmdArray = append(cmdArray, arg)
		}
		gatewayName := cmdArray[0]
		start(gatewayName)
		return nil
	},
}


var StopCommand = cli.Command{
	Name:                   "stop",
	Usage:                  "stop mqtt watcher",
	Action: func(context *cli.Context) error {
		if context.Args().Len() < 1 {
			return fmt.Errorf("Missing mqtt command")
		}
		var cmdArray []string
		for i := 0; i < context.Args().Len(); i++ {
			arg := context.Args().Get(i)
			cmdArray = append(cmdArray, arg)
		}
		containerName := cmdArray[0]
		pid := ReadPid(containerName)
		stopMqttWatcher(containerName, pid)
		os.Exit(0)
		return nil
	},
}
