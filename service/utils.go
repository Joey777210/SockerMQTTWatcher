package service

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

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