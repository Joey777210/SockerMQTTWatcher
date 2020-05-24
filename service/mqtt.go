package service

import (

	"crypto/tls"
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"time"
)

var Client mqtt.Client

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

	Client = mqtt.NewClient(opts)

	var flag = 0
	for {
		if flag == 0 {
			if token := Client.Connect(); token.Wait() && token.Error() != nil {
				flag = 1
			} else {
				break
			}
		} else if flag == 1 {
			if token := Client.Connect(); token.Wait() && token.Error() != nil {

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
	log.Info("onconnect")
	if token := client.Subscribe(GetTopic(SysOrderSub), 0, onMessageReceived); token.Wait() && token.Error() != nil {
		log.Errorf("mqtt Client subscribe topic %s Error %v", GetTopic(SysOrderSub), token.Error())
	}
}

func onMessageReceived(client mqtt.Client, message mqtt.Message) {
	msg := message.Payload()
	order := Order{}
	err := json.Unmarshal(msg, &order)
	if err != nil {
		log.Errorf("Unmarshal order error %v", err)
	}
	log.Info(string(msg))
	log.Info(order.Name + "\n" + order.Target + "\n" + order.Content + "\n")
	log.Infof("content!!!!!!!!!! %s", order.Content)

	socker := sockerImp{}
	switch order.Target {
	case "container":
		switch order.Order {
		case "run":
			//TODO
			//listen new container and send data
			//return new container info
			socker.RunNewContainer(order)
		case "stop":
			socker.ContainerStop(order)
		case "commit":
			socker.ContainerCommit(order)
		case "remove":
			socker.ContainerRemove(order)
		case "ls":
			socker.ContainerLs(client)
		case "log":
			socker.ContainerLogs(client, order)
		}
	case "image":
		switch order.Order {
		case "ls":
			socker.ImageLs(client)
		case "remove":
			socker.ImageRm(order)
		}
	case "socker":
		switch order.Order {
		case "newNetwork":
		case "delete":
			socker.Delete()
		}
	}
}

