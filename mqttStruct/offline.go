package mqttStruct

import (
	"crypto/tls"
	"fmt"
	log "github.com/Sirupsen/logrus"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func MqttOffline(containerName string) {
	SetMqttClient(&MqttClient)
	fmt.Println(MqttClient.Server)
	opts := mqtt.NewClientOptions().AddBroker(MqttClient.Server)
	opts.SetCleanSession(true)
	opts.SetClientID(containerName)

	tlsConfig := &tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert}
	opts.SetTLSConfig(tlsConfig)

	Replace(containerName)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Errorf("Mqtt offline connect error %v", token.Error())
	}

	if token := client.Publish(GetTopic(SysOnLinePub), 0, false, OffLine); token.Wait() && token.Error() != nil {
		log.Errorf("client publish error %v\n", token.Error())
	}
	log.Infof("send message: %s, %s", GetTopic(SysOnLinePub), OffLine)
	client.Disconnect(250)
}
