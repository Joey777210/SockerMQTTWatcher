package service

import (
	"crypto/tls"
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"time"
)

var client mqtt.Client

//Only for gateway topic.    Container topic is in another connection
func Connect(gatewayName string) {
	opts := mqtt.NewClientOptions().AddBroker(Server)
	opts.SetCleanSession(true)
	opts.SetClientID(gatewayName)
	opts.OnConnect = OnConnect
	//opts.OnConnectionLost = OnConnectLost

	Replace(gatewayName)

	tlsConfig := &tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert}
	opts.SetTLSConfig(tlsConfig)

	client = mqtt.NewClient(opts)

	var flag = 0
	for {
		if flag == 0 {
			if token := client.Connect(); token.Wait() && token.Error() != nil {
				flag = 1
			} else {
				break
			}
		} else if flag == 1 {
			if token := client.Connect(); token.Wait() && token.Error() != nil {

			} else {
				flag = 0
				break
			}
		}
		time.Sleep(1 * time.Second)
	}
	for {
		time.Sleep(1 * time.Second)
	}
}

func OnConnect(client mqtt.Client) {
	if token := client.Subscribe(GetTopic(SysGWSub), 0, onMessageReceived); token.Wait() && token.Error() != nil {
		log.Errorf("mqtt client subscribe topic %s Error %v", GetTopic(SysGWSub), token.Error())
	}
}

func onMessageReceived(client mqtt.Client, message mqtt.Message) {
	msg := message.Payload()
	order := Order{}
	err := json.Unmarshal(msg, &order)
	if err != nil {
		log.Errorf("Unmarshal order error %v", err)
	}
	//add container
	if order.Order == "run" {
		socker := sockerImp{}
		socker.RunNewContainer(order)
	} else if order.Order == "containerls" {
		socker := sockerImp{}
		socker.ContainerLs(client)
	} else if order.Order == "imagels" {
		socker := sockerImp{}
		socker.ImageLs(client)
	} else if order.Order == "delete" {
		//Delete this gateway from web
		//Do nothing
		socker := sockerImp{}
		socker.Delete()
	}
}

func MessagePublic(client mqtt.Client, topic string, message string) error {
	if token := client.Publish(topic, 0, false, message); token.Wait() && token.Error() != nil {
		log.Errorf("Client publish on Topic %s error %v\n", topic, token.Error())
		return token.Error()
	}
	return nil
}