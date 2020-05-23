package command

import (
	"SockerMQTTWatcher/service"
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"os"
	"os/exec"
	"time"
)

func start(gatewayName string) {
	InitSubnet()
	service.Connect(gatewayName)
	go service.LogAutoPub()
	go containerLiveCheck()

	//TODO listen port:8888, resend data
}

func InitSubnet() {
	createCmd := "sudo socker network create --driver bridge --subnet 192.168.10.1/24 testbridge"
	cmd := exec.Command("/bin/sh", "-c", createCmd)
	err := cmd.Start()
	if err != nil {
		log.Errorf("Start command %s error %v", createCmd, err)
	}
}

func containerLiveCheck() {
	for {
		for i := range service.Containers {
			container := service.Containers[i]
			_, err := os.Stat("/proc" + container.Pid)
			if err == nil {
				service.MessagePublic(service.Client, service.GetTopicCN(service.SysStatusPub, container.Name), "online")
			}
			if os.IsNotExist(err) {
				service.MessagePublic(service.Client, service.GetTopicCN(service.SysStatusPub, container.Name), "offline")
				err := errors.New(fmt.Sprintf("container %s disconnected for unknown reasons", container.Name))
				service.ErrorPublic(err)
				container.Status = "stopped"
			}
		}

		time.Sleep(100*time.Second)
	}



}