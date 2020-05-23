package command

import (
	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

var StartCommand = cli.Command{
	Name:  "start",
	Usage: "start mqtt watcher",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:	"name",
			Usage:	"Set gateway name",
		},
	},
	Action: func(context *cli.Context) error {
		gatewayName := context.String("name")
		log.Infof(gatewayName)
		start(gatewayName)
		return nil
	},
}

var StopCommand = cli.Command{
	Name:  "stop",
	Usage: "stop mqtt watcher",
	Action: func(context *cli.Context) error {
		stop()
		return nil
	},
}
