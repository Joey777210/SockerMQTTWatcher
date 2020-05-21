package service

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"os/exec"
)

type socker interface {
	Run(order Order) error
	containerls() error
}

type sockerImp struct {

}

func (s *sockerImp) RunNewContainer(order Order) {
	content := order.Content
	container := ContainerImp{}
	err := json.Unmarshal([]byte(content), &container)
	if err != nil {
		log.Errorf("Json Unmarshal run order error %v", err)
	}
	fmt.Print(container.Memory, container.Name)
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
	runCmd := fmt.Sprintf(base, resource, container.PortMapping, container.Name, container.Image, container.Command)
	//TODO
	//fillContainerInfo(container)
	fmt.Print(runCmd)
	cmd := exec.Command("/bin/sh", "-c", runCmd)
	err = cmd.Start()
	if err != nil {
		log.Errorf("Start command %s error %v", runCmd, err)
	}
}

