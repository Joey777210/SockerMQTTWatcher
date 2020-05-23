package mqttStruct

import (
	"crypto/tls"
	"fmt"
	log "github.com/Sirupsen/logrus"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"os"
	"time"
)

const (
	OnLine = "online"
	OffLine = "offline"
)

type Imqtt interface {
	//connect mqtt
	Connect() error
}

type MqttImpl struct {
}

var CN string

//mqtt connect
func (m *MqttImpl) Connect(containerName string) error {
	CN = containerName
	SetMqttClient(&MqttClient)
	fmt.Println(MqttClient.Server)
	opts := mqtt.NewClientOptions().AddBroker(MqttClient.Server)
	opts.SetCleanSession(true)
	opts.SetClientID(containerName)
	opts.OnConnect = OnConnect
	opts.OnConnectionLost = OnConnectLost
	opts.SetWill(GetTopic(SysOnLinePub), OffLine, 1, true)

	//replace {CN} with containerName
	Replace(containerName)

	tlsConfig := &tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert}
	opts.SetTLSConfig(tlsConfig)

	client := mqtt.NewClient(opts)

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
	log.Infoln("onconnect  + " + GetTopic(SysOrderSub))

	if token := client.Publish(GetTopic(SysOnLinePub), 0, false, OnLine); token.Wait() && token.Error() != nil {
		log.Errorf("client publish error %v\n", token.Error())
	}
	if token := client.Subscribe(GetTopic(SysOrderSub), 0, onMessageReceived); token.Wait() && token.Error() != nil {
		log.Errorf("client subscribe message Error %v", token.Error())
	}

	go Listen(":8888", client)
	//watch file change and send message
	//sendMessage(client)
}

func OnConnectLost(client mqtt.Client, err error) {
	log.Error("mqtt client lost!")
}

func onMessageReceived(client mqtt.Client, message mqtt.Message) {
	log.Infof("Received message on topic: %s \t Message: %s\n", message.Topic(), message.Payload())

	if string(message.Payload()) == "stop" {
		stop(CN)
	} else if string(message.Payload()) == "log" {
		showLog(client, CN)
	}

	//dirURL := fmt.Sprintf("/root/mergeDir/%s", CN)
	//
	//fileName := dirURL + "/mqttSub"
	//file, err := os.Create(fileName)
	//defer file.Close()
	//if err != nil {
	//	fmt.Printf("Create file %s error %v \n", fileName, err)
	//}
	//jsonStr := string(message.Payload())
	//_, err = file.WriteString(jsonStr)
	//if err != nil {
	//	log.Errorf("Write json str error %v", err)
	//}
}

func sendMessage(client mqtt.Client) {
	//dirURL := fmt.Sprintf("/root/mergeDir/%s", CN)
	//fileName := dirURL + "/mqttPub"
	//log.Info(fileName)
	//isExist, _ := PathExists(fileName)
	//if !isExist {
	//	if _, err := os.Create(fileName); err != nil {
	//		log.Errorf("Create file %s error %v", fileName, err)
	//	}
	//}
	//
	//watcher, err := fsnotify.NewWatcher()
	//if err != nil {
	//	log.Errorf("New watcher error %v", err)
	//}
	//defer watcher.Close()
	//
	////done := make(chan bool)
	//go func() {
	//	for {
	//		select {
	//		case event, ok := <-watcher.Events:
	//			if !ok {
	//				return
	//			}
	//			log.Infoln("event: ", event)
	//			if event.Op&fsnotify.Write == fsnotify.Write {
	//				message := readFile(fileName)
	//				msgPub(client, message)
	//			}
	//		case err, ok := <-watcher.Errors:
	//			if !ok {
	//				return
	//			}
	//			log.Infof("Watch file error1 %v", err)
	//		}
	//	}
	//}()
	//
	//err = watcher.Add(fileName)
	//if err != nil {
	//	log.Errorf("Watch file error2 %v", err)
	//}
	////<-done
	////循环
	//select {}
}

