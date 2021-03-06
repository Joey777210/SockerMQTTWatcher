package service

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"os"
)

const DefaultMQTTLogDir string = "/var/run/mqtt"

func ackPublic(client mqtt.Client, ackMsg string) error {
	if token := client.Publish(GetTopic(SysACKPub), 0, false, ackMsg); token.Wait() && token.Error() != nil {
		log.Errorf("Client publish on Topic %s error %v\n", GetTopic(SysACKPub), token.Error())
		return token.Error()
	}
	return nil
}


func MessagePublic(client mqtt.Client, topic string, message string) error {
	if token := client.Publish(topic, 0, false, message); token.Wait() && token.Error() != nil {
		log.Errorf("Client publish on Topic %s error %v\n", topic, token.Error())
		return token.Error()
	}
	return nil
}

//example "container:bird:stop:1"
func AckMsgFormat(target string, name string, command string, ack int) string {
	ackMsg := "%s:%s:%s:%d"
	return	fmt.Sprintf(ackMsg, target, name, command, ack)
}

func SaveContainers() error {
	_, err := os.Stat(DefaultMQTTLogDir)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(DefaultMQTTLogDir, 0777); err != nil {
			log.Errorf("NewParentProcess mkdir %s error %v", DefaultMQTTLogDir, err)
			return err
		}
	}

	stdErrFilePath := DefaultMQTTLogDir + "/containers.log"
	_, err = os.Stat(stdErrFilePath)
	if os.IsNotExist(err) {
		_, err := os.Create(stdErrFilePath)
		if err != nil {
			log.Errorf("%v", err)
		}
	}
	stdLogFile, err := os.OpenFile(stdErrFilePath, os.O_RDWR | os.O_TRUNC, 0777)
	if err != nil {
		log.Errorf("NewParentProcess create file %s error %v", stdErrFilePath, err)
		return err
	}


	containerJson, _ := json.Marshal(Containers)
	_, err = stdLogFile.Write(containerJson)
	if err != nil {
		log.Errorf("Write err to file %s error %v", stdErrFilePath, err)
		return err
	}
	return nil
}