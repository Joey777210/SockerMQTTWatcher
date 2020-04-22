package command

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"mqttWatcher/mqttStruct"
	"os"
	"strconv"
)

const (
	DefaultInfoLocation = "/var/run/socker/%s"
)

func startMqttWatcher(containerName string) {
	savePid(containerName)
	mqttClient(containerName)
}

func savePid(containerName string) {
	//save pid and stop it by pid
	dir := fmt.Sprintf(DefaultInfoLocation, containerName)
	filePath := dir + "/mqttPid"
	file, err := os.Create(filePath)
	defer file.Close()
	if err != nil {
		log.Errorf("Create pid file error %v", err)
	}
	pid := os.Getpid()
	_, err = file.WriteString(strconv.Itoa(pid))
	if err != nil {
		log.Errorf("Write pid file error %v", err)
	}
}

func mqttClient(containerName string) {
	log.Info("mqtt watcher start!")
	mq := mqttStruct.MqttImpl{}
	if err := mq.Connect(containerName); err != nil {
		log.Errorf("mqtt open error: %v", err)
	}
}
