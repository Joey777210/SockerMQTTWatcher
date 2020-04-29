package command

import (
	"SockerMQTTWatcher/mqttStruct"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"os"
	"strconv"
	"syscall"
)

func stopMqttWatcher(containerName string, pid int) {
	log.Infof("kill process %d", pid)
	processKill(pid)
	mqttStruct.MqttOffline(containerName)
}

func processKill(pid int) {
	if err := syscall.Kill(pid, 9); err != nil {
		log.Errorf("Kill %d process error %v", pid, err)
	}
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