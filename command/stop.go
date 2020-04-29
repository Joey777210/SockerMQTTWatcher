package command

import (
	"crypto/tls"
	"fmt"
	log "github.com/Sirupsen/logrus"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"mqttWatcher/mqttStruct"
	"os"
	"strconv"
	"syscall"
)

func stopMqttWatcher(containerName string, pid int) {
	log.Infof("kill process %d", pid)
	processKill(pid)
	MqttOffline(containerName)
}

func processKill(pid int) {
	if err := syscall.Kill(pid, 9); err != nil {
		log.Errorf("Kill %d process error %v", pid, err)
	}
}

func MqttOffline(containerName string) {
	mqttStruct.SetMqttClient(&mqttStruct.MqttClient)
	fmt.Println(mqttStruct.MqttClient.Server)
	opts := mqtt.NewClientOptions().AddBroker(mqttStruct.MqttClient.Server)
	opts.SetCleanSession(true)
	opts.SetClientID(containerName)

	tlsConfig := &tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert}
	opts.SetTLSConfig(tlsConfig)

	mqttStruct.Replace(containerName)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Errorf("Mqtt offline connect error %v", token.Error())
	}

	if token := client.Publish(mqttStruct.GetTopic(mqttStruct.SysOnLinePub), 0, false, mqttStruct.OffLine); token.Wait() && token.Error() != nil {
		log.Errorf("client publish error %v\n", token.Error())
	}
	log.Infof("send message: %s, %s", mqttStruct.GetTopic(mqttStruct.SysOnLinePub), mqttStruct.OffLine)
	client.Disconnect(250)
}

func ReadPid(containerName string) int {
	dir := fmt.Sprintf(DefaultInfoLocation, containerName)
	filePath := dir + "/mqttPid"
	log.Info(filePath)
	file, err := os.Open(filePath)
	if err != nil {
		log.Errorf("Open mqtt pid file error %v", err)
	}
	defer file.Close()
	pidBytes := make([]byte, 1024)
	i, err := file.Read(pidBytes)
	if err != nil {
		log.Errorf("Read pid file error %v", err)
	}
	pidStr := string(pidBytes[:i])
	log.Info(pidStr)
	pid, _ := strconv.Atoi(pidStr)
	log.Info(pid)
	return pid
}