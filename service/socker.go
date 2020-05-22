package service

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var (
	DefaultInfoPath = "/var/run/socker/%s"
	DefaultRootPath = "/root"
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
	content := order.Content
	container := ContainerImp{}
	err := json.Unmarshal([]byte(content), &container)
	if err != nil {
		log.Errorf("Json Unmarshal run order error %v", err)
	}

	_, ok := Containers[container.Name]
	if ok {
		ErrorPublic(err)
		return
	}

	base := "sudo socker run -d -mqtt %s -net testbridge -p %s --name %s %s %s"
	resource := ""
	if container.Memory != ""{
		resource += "-m " + container.Memory
	}
	if container.CpuSet != ""{
		resource += " -cpuset " + container.CpuSet
	}
	if container.CpuShare != ""{
		resource += " -cpushare " + container.CpuShare
	}
	runCmd := fmt.Sprintf(base, resource, container.PortMapping[0], container.Name, container.Image, container.Command)

	cmd := exec.Command("/bin/sh", "-c", runCmd)
	err = cmd.Run()

	FillContainerInfo(&container)
	if err != nil {
		log.Errorf("Start command %s error %v", runCmd, err)
	}
	Containers[container.Name] = container

	//TODO
	//Listen on contianer Topic
	//ThreadPool
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
	content := order.Content
	image := image{}
	err := json.Unmarshal([]byte(content), &image)
	if err != nil {
		log.Errorf("Json unmarshal error %v", err)
		ErrorPublic(err)
	}

	imageName := image.Name + ".tar"
	imagePath := DefaultRootPath + "/" +imageName
	err = os.Remove(imagePath)
	if err != nil {
		log.Errorf("Remove Image %s error %v", imagePath, err)
		ErrorPublic(err)
	}
}

func (s *sockerImp)ContainerStop(order Order) {
	content := order.Content
	container := ContainerImp{}
	err := json.Unmarshal([]byte(content), &container)
	if err != nil {
		log.Errorf("Json unmarshal error %v", err)
		ErrorPublic(err)
	}
	cmd := exec.Command("/bin/sh", "-c", "sudo socker stop > /var/run/socker/sockerlog", container.Name)
	err = cmd.Run()
	if err != nil {
		log.Errorf("Exec command %s error %v", cmd.String(), err)
		ErrorPublic(err)
	}


}