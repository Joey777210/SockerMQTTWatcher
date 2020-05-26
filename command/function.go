package command

import (
	"SockerMQTTWatcher/service"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"os"
	"os/exec"
	"time"
)

func start(gatewayName string) {
	//log.SetMylog()
	loadContainers()
	InitSubnet()
	service.Connect(gatewayName)
	go service.LogAutoPub()
	go containerLiveCheck()

	go service.Listen(":8888", service.Client)
}

func stop() {

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

func loadContainers() {
	filePath := service.DefaultMQTTPath + "/containers.log"
	file, err := os.OpenFile(filePath, os.O_RDWR, 0666)
	if err != nil {
		log.Errorf("Open file %s error %v", filePath, err)
	}
	buf := make([]byte, 1024*1024)
	n, _ := file.Read(buf)

	err = json.Unmarshal(buf[:n], &service.Containers)
	if err != nil {
		log.Errorf("Json Unmarshal file %s error %v", filePath, err)
	}
}