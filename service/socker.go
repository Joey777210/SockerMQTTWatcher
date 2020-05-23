package service

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/fsnotify/fsnotify"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

var (
	DefaultInfoPath = "/var/run/socker/%s"
	DefaultRootPath = "/root"
	DefaultLogPath	= "/var/run/socker/sockerlog"
)

type socker interface {
	RunNewContainer(order Order) error
	ContainerLs(client mqtt.Client)
	ImageLs(client mqtt.Client)
	Delete()
}

type sockerImp struct {

}

//
func (s *sockerImp) RunNewContainer(order Order) {
	container := ContainerImp{}
	err := container.Run(order)
	if err != nil {
		log.Errorf("Run new container %s error %v", order.Target, err)
	}
}

func (s *sockerImp) ContainerLs(client mqtt.Client) {

	bytes, err := json.Marshal(Containers)
	if err != nil {
		log.Errorf("json marshal error %v", err)
	}
	message := string(bytes)
	err = MessagePublic(client, GetTopic(SysCtnlsPub), message)
	if err != nil {
		ErrorPublic(err)
	}
}

func (s *sockerImp) ImageLs(client mqtt.Client) {
	err := findImages()
	if err != nil {
		log.Errorf("Find images error %v", err)
		ErrorPublic(err)
	}

	bytes, err := json.Marshal(Images)
	if err != nil {
		log.Errorf("Json marshal images error %v", err)
		ErrorPublic(err)
	}
	message := string(bytes)
	err = MessagePublic(client, GetTopic(SysImglsPub), message)
	if err != nil {
		log.Errorf("Message Public image ls error %v", err)
		ErrorPublic(err)
	}
}

func findImages() error {
	files, err := ioutil.ReadDir(DefaultRootPath)
	if err != nil {
		log.Errorf("Open dir %s error %v", DefaultRootPath, err)
		ErrorPublic(err)
		return err
	}
	//get all image
	for _, f := range files {
		strs := strings.Split(f.Name(), ".")
		if len(strs) == 2 && strs[1] == "tar"{
			image := getImage(f)
			Images[image.Name] = image
		}
	}
	return nil
}

func getImage(f os.FileInfo) image {
	var image image
	name := strings.Split(f.Name(), ".")
	modTime := f.ModTime().String()
	times := strings.Split(modTime, " ")
	hms := strings.Split(times[1], ".")

	image.Name = name[0]
	image.ModTime = times[0] + " " + hms[0]
	image.Size = strconv.Itoa(int(f.Size()/1024/1024)) + "MB"
	return image
}

func (s *sockerImp)Delete() {
	//Do nothing
}

func (s *sockerImp)ImageRm(order Order) {
	image := image{}
	err := image.Remove(order)
	if err != nil {
		log.Errorf("Remove image %s error %v", order.Target, err)
	}
}

func (s *sockerImp)ContainerStop(order Order) {
	container := ContainerImp{}
	err := container.Stop(order)
	if err != nil {
		log.Errorf("Stop container %s error %v", order.Target, err)
	}
}

func LogAutoPub() {
	isExist, _ := PathExists(DefaultLogPath)
	if !isExist {
		if _, err := os.Create(DefaultLogPath); err != nil {
			log.Errorf("Create file %s error %v", DefaultLogPath, err)
			ErrorMsgPublic(fmt.Sprintf("Create file %s error %v", DefaultLogPath, err))
		}
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Errorf("New watcher error %v", err)
	}
	defer watcher.Close()

	//done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Infoln("event: ", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					message := readFile(DefaultLogPath)
					fmt.Println(message)
					//err := MessagePublic(client, GetTopic(SysGWLogPub), message)
					if err != nil {
						log.Errorf("Send message error %v", err)
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Infof("Watch file error1 %v", err)
			}
		}
	}()

	err = watcher.Add(DefaultLogPath)
	if err != nil {
		log.Errorf("Watch file error2 %v", err)
	}
	//<-done
	//循环
	select {}
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func readFile(fileName string) string {

	var message []byte
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Open file %s error %v \n", fileName, err)
	}
	defer file.Close()
	buf := make([]byte, 1024)
	for {
		n, _ := file.Read(buf)
		if 0 == n {
			break
		}
		message = append(message, buf[:n]...)
	}
	////clear file
	//_ = os.Truncate(fileName, 0)
	//_, _ = file.Seek(0, 0)
	//after test, don't need clear!!
	return string(message)
}

func (s *sockerImp) ContainerCommit(order Order) {
	container := ContainerImp{}
	err := container.Commit(order)
	if err != nil {
		log.Errorf("Commit container error %v", err)
	}
}

func (s *sockerImp) ContainerRemove(order Order) {
	container := ContainerImp{}
	err := container.Remove(order)
	if err != nil {
		log.Errorf("Remove container error %v", err)
	}
}