package service

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"go/types"
	"time"
)

var Containers types.Map

type container interface {
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
	Command	string			`json:"cmd"`
	CreatTime	string		`json:"time"`
	Image	string			`json:"image"`
	Memory	string			`json:"memory"`
	CpuSet	string			`json:"cpuset"`
	CpuShare	string		`json:"cpushare"`
	PortMapping	string		`json:"portmapping"`
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

func fillContainerInfo(container ContainerImp) {
	container.CreatTime = time.Now().String()
	//container.ID
	//container.Pid
	//container.Status =
}

func Marshal1() string {
	container := ContainerImp{
		Name:        "bird",
		Command:     "top -b",
		Image:       "ubuntu",
		Memory:      "512",
		CpuSet:      "1",
		CpuShare:    "512",
		PortMapping: "8080:80",
	}
	bytes, err := json.Marshal(&container)
	if err != nil {
		log.Errorf("%v", err)
	}
	return string(bytes)
}