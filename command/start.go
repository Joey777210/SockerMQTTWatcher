package command

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"SockerMQTTWatcher/mqttStruct"
	"os"
	"strconv"
)

const (
	DefaultInfoLocation = "/var/run/socker/%s"
)


func startMqttWatcher(containerName string) {
	//mqttStruct.Listen(":8888", msg)
	savePid(containerName)
	mqttClient(containerName)
}

func savePid(containerName string) {
	//save pid and stop it by pid
	dir := fmt.Sprintf(DefaultInfoLocation, containerName)
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0777)
		if err != nil {
			log.Errorf("Create dir %s error %v", dir, err)
		}
	}
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
