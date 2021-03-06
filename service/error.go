package service

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"os"
	"time"
)

func ErrorPublic(err error) {
	errSave(err)
	if token := Client.Publish(GetTopic(SysGWErrPub), 0, false, err.Error()); token.Wait() && token.Error() != nil {
		log.Errorf("Client publish on Topic %s error %v\n", GetTopic(SysGWErrPub), token.Error())
		errSave(token.Error())
	}
}

func ErrorMsgPublic(message string) {
	errMsgSave(message)
	if token := Client.Publish(GetTopic(SysGWErrPub), 0, false, message); token.Wait() && token.Error() != nil {
		log.Errorf("Client publish on Topic %s error %v\n", GetTopic(SysGWErrPub), token.Error())
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
			file.Chmod(0777)
		}
	}
	file, _ = os.OpenFile(errfilePath, os.O_RDWR, 0666)
	errMsg := fmt.Sprintf("%s Error: %v", time.Now().String(), error)
	_, err = file.Write([]byte(errMsg))
	if err != nil {
		log.Errorf("\n Write err to file error %v", err)
	}
}


func errMsgSave(message string) {
	errfilePath := fmt.Sprintf(DefaultInfoPath, "mqtterr")
	var file *os.File
	file, err := os.OpenFile(errfilePath, os.O_RDWR, 0777)
	if err != nil {
		if os.IsNotExist(err) {
			file, _ = os.Create(errfilePath)
			file.Chmod(0777)
		}
	}

	errMsg := fmt.Sprintf("%s Error: %s", time.Now().String(), message)
	_, err = file.WriteString(errMsg)
	if err != nil {
		log.Errorf("Write err to file error %v", err)
	}
}