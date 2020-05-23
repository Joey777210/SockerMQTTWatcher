package service

import (
	"SockerMQTTWatcher/log"
	"fmt"
	"os"
	"time"
)

func ErrorPublic(err error) {
	errSave(err)
	if token := Client.Publish(GetTopic(SysGWErrPub), 0, false, err.Error()); token.Wait() && token.Error() != nil {
		log.Mylog.Errorf("Client publish on Topic %s error %v\n", GetTopic(SysGWErrPub), token.Error())
		errSave(token.Error())
	}
}

func ErrorMsgPublic(message string) {
	errMsgSave(message)
	if token := Client.Publish(GetTopic(SysGWErrPub), 0, false, message); token.Wait() && token.Error() != nil {
		log.Mylog.Errorf("Client publish on Topic %s error %v\n", GetTopic(SysGWErrPub), token.Error())
		errSave(token.Error())
	}
}

func errSave(error error) {
	errfilePath := fmt.Sprintf(DefaultInfoPath, "mqtterr")
	var file *os.File
	file, err := os.Open(errfilePath)
	if err != nil {
		if os.IsNotExist(err) {
			file, _ = os.Create(errfilePath)
		}
	}
	errMsg := fmt.Sprintf("%s Error: %v", time.Now().String(), error)
	_, err = file.WriteString(errMsg)
	if err != nil {
		log.Mylog.Errorf("Write err to file error %v", err)
	}
}


func errMsgSave(message string) {
	errfilePath := fmt.Sprintf(DefaultInfoPath, "/mqtterr")
	var file *os.File
	file, err := os.Open(errfilePath)
	if err != nil {
		if os.IsNotExist(err) {
			file, _ = os.Create(errfilePath)
		}
	}
	errMsg := fmt.Sprintf("%s Error: %s", time.Now().String(), message)
	_, err = file.WriteString(errMsg)
	if err != nil {
		log.Mylog.Errorf("Write err to file error %v", err)
	}
}