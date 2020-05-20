package mqttStruct

import (
	log "github.com/Sirupsen/logrus"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"io/ioutil"
	"os/exec"
)

func stop(containerName string) {
	cmd := exec.Command("/bin/sh", "-c", "sudo socker stop", containerName)
	err := cmd.Start()
	if err != nil {
		log.Errorf("Stop container %s error %v", containerName, err)
	}
}


func showLog(client mqtt.Client, containerName string) {
	logPath := "/var/run/socker/" + containerName + "/container.log"
	logs, err := ioutil.ReadFile(logPath)
	if err != nil {
		log.Errorf("Read file %s error %v", logPath, err)
		return
	}
	msgPub(client, string(logs))
}