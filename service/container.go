package service

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"os"
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

func (c *ContainerImp) Stop(containerName string) error {
	//TODO
	return nil
}

func (c *ContainerImp) Remove(containerName string) error {
	//TODO
	return nil
}

func (c *ContainerImp) Commit(containerName string) error {
	//TODO
	return nil
}

func (c *ContainerImp) Log(containerName string) error {
	//TODO
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