package service

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"os"
	"os/exec"
)

var Containers = make(map[string]ContainerImp)

type Container interface {
	Stop(containerName string) error
	Remove(containerName string) error
	Commit(containerName string) error
	Log(containerName string) error
}


type ContainerImp struct {
	ID		string			`json:"id"`
	Name	string			`json:"name"`
	Pid		string			`json:"pid"`
	Status	string			`json:"status"`
	Command	string			`json:"command"`
	CreatTime	string		`json:"createTime"`
	Image	string			`json:"image"`
	Memory	string			`json:"memory"`
	CpuSet	string			`json:"cpuset"`
	CpuShare	string		`json:"cpushare"`
	PortMapping	[]string	`json:"portmapping"`
}


func (c *ContainerImp) Run(order Order) error {
	content := order.Content
	err := json.Unmarshal([]byte(content), &c)
	if err != nil {
		log.Errorf("Json Unmarshal run order error %v", err)
	}

	_, ok := Containers[c.Name]
	if ok {
		if Containers[c.Name].Status == "runnning" {
			err := errors.New("container already running")
			ErrorPublic(err)
			return err
		}
	}

	base := "sudo socker run -d -mqtt %s -net testbridge -p %s --name %s %s %s"
	resource := ""
	if c.Memory != ""{
		resource += "-m " + c.Memory
	}
	if c.CpuSet != ""{
		resource += " -cpuset " + c.CpuSet
	}
	if c.CpuShare != ""{
		resource += " -cpushare " + c.CpuShare
	}
	runCmd := fmt.Sprintf(base, resource, c.PortMapping[0], c.Name, c.Image, c.Command)

	cmd := exec.Command("/bin/sh", "-c", runCmd)
	err = cmd.Run()
	if err != nil {
		log.Errorf("Start command %s error %v", runCmd, err)
		ErrorPublic(err)
		return err
	}
	FillContainerInfo(c)
	Containers[c.Name] = *c


	//TODO
	//Listen on contianer Topic
	//ThreadPool
	return nil
}

func (c *ContainerImp) Stop(containerName string) error {
	if containerName == "" {
		return errors.New("empty name error")
	}
	c.Name = containerName
	_, ok := Containers[c.Name]
	if !ok {
		err := errors.New("no such container")
		ErrorPublic(err)
		return err
	}

	command1 := "sudo /bin/sh -c \"socker stop " + c.Name + " > " + DefaultLogPath +"\""
	cmd := exec.Command("/bin/sh", "-c", command1)

	err := cmd.Run()
	if err != nil {
		log.Errorf("Exec command %s error %v", cmd.String(), err)
		ErrorPublic(err)
		return err
	}
	c.Status = "stopped"

	return nil
}

func (c *ContainerImp) Remove(containerName string) error {
	if containerName == "" {
		return errors.New("empty name error")
	}
	c.Name = containerName
	err := isExist(c.Name)
	if err != nil {
		return err
	}

	if c.Status == "running" {
		err = errors.New("you can't remove a running container, please stop the container first!")
		ErrorPublic(err)
		return err
	} else if c.Status == "stopped" {
		command := "sudo sh -c \"sudo socker remove" + c.Name + "\""
		cmd := exec.Command("bin/sh", "-c", command)
		err = cmd.Run()
		if err != nil {
			ErrorPublic(err)
			return err
		}
		delete(Containers, c.Name)
	}
	return nil
}

func (c *ContainerImp) Commit(containerName string) error {
	if containerName == "" {
		return errors.New("empty name error")
	}
	c.Name = containerName
	err := isExist(c.Name)
	if err != nil {
		return err
	}

	command := "sudo /bin/sh -c \"socker commit " + c.Name + " > " + DefaultLogPath +"\""
	cmd := exec.Command("/bin/sh", "-c", command)
	err = cmd.Run()
	if err != nil {
		ErrorPublic(err)
		return err
	}
	s := sockerImp{}
	s.ImageLs(Client)
	return nil
}

func (c *ContainerImp) Logs(client mqtt.Client, containerName string) error {
	if containerName == "" {
		return errors.New("empty name error")
	}

	c.Name = containerName
	err := isExist(c.Name)
	if err != nil {
		return err
	}

	logPath := fmt.Sprintf(DefaultInfoPath, c.Name) + "container.log"
	logs, err := readFile(logPath)
	err = MessagePublic(client, GetTopicCN(SysLogPub, c.Name), logs)
	if err != nil {
		err := errors.New(fmt.Sprintf("Logs public error %v", err))
		ErrorPublic(err)
		return err
	}

	return nil
}

func FillContainerInfo(container *ContainerImp) {
	confPath := fmt.Sprintf(DefaultInfoPath, container.Name) + "/config.json"
	buf := make([]byte, 1024)
	file, err := os.Open(confPath)
	n, err := file.Read(buf)
	conf := buf[:n]
	if err != nil {
		log.Errorf("Read file %s error %v", confPath, err)
	}
	fmt.Println(string(conf))
	err = json.Unmarshal(conf, &container)
	if err != nil {
		log.Errorf("Json unmarshal error %v", err)
	}
	fmt.Printf("%s %s %s %s %s %s", container.CreatTime, container.Status, container.Pid, container.ID, container.Name, container.Command)
}

func Marshal1() string {
	container := ContainerImp{
		Name:        "bird",
		Command:     "top -b",
		Image:       "ubuntu",
		Memory:      "512",
		CpuSet:      "1",
		CpuShare:    "512",
		PortMapping: []string{"8080:80", "9090:90"},
	}
	bytes, err := json.Marshal(&container)
	if err != nil {
		log.Errorf("%v", err)
	}
	return string(bytes)
}

func isExist(containerName string) error {
	_, ok := Containers[containerName]
	if !ok {
		err := errors.New("no such container")
		ErrorPublic(err)
		return err
	}
	return nil
}