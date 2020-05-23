package command

import (
	"SockerMQTTWatcher/service"
	log "github.com/Sirupsen/logrus"
	"os/exec"
)

func start(gatewayName string) {
	InitSubnet()
	go service.LogAutoPub()
	service.Connect(gatewayName)
	//service.Containers = make(map[string]service.ContainerImp)
}

func InitSubnet() {
	createCmd := "sudo socker network create --driver bridge --subnet 192.168.10.1/24 testbridge"
	cmd := exec.Command("/bin/sh", "-c", createCmd)
	err := cmd.Start()
	if err != nil {
		log.Errorf("Start command %s error %v", createCmd, err)
	}
}