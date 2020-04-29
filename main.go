package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
	"SockerMQTTWatcher/command"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "mqtt"
	app.Commands = []*cli.Command{
		&command.StartCommand,
		&command.StopCommand,
	}


	//init logrus
	app.Before = func(context *cli.Context) error {
		//set log as JSON formatter
		log.SetFormatter(&log.JSONFormatter{})
		log.SetOutput(os.Stdout)
		return nil
	}
	//output args you had just type-in
	sockerCommand := os.Args
	log.Printf("args: %s", sockerCommand)


	if err := app.Run(os.Args); err != nil{
		log.Fatal(err)
	}
}
