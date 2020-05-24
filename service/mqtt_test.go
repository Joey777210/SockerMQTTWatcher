package service

import (
	"fmt"
	"testing"
)

func TestReplace(t *testing.T) {
	Replace("bird")

}

func TestSockerImp_RunNewContainer(t *testing.T) {
	socker := sockerImp{}
	order := Order{
		Order:   "run",
		Content: Marshal1(),
	}
	socker.RunNewContainer(order)
}

func TestMarshal1(t *testing.T) {
	fmt.Print(Marshal1())
}

func TestFillContainerInfo(t *testing.T) {
	var arr = []string{"8080:80", "498"}
	container := ContainerImp{
		Name:        "bird",
		Command:     "top -b",
		Image:       "ubuntu",
		Memory:      "512",
		CpuSet:      "1",
		CpuShare:    "512",
		PortMapping: "80:80",
	}
	FillContainerInfo(&container)
}

func TestSockerImp_ConatainerLs(t *testing.T) {
	s := sockerImp{}
	order := Order{
		Order:   "run",
		Content: Marshal1(),
	}
	//order := Order{
	//	Order:   "run",
	//	Content: Marshal1(),
	//}
	s.RunNewContainer(order)
	s.ContainerLs(Client)
}

func TestSockerImp_ImageLs(t *testing.T) {
	s := sockerImp{}
	s.ImageLs(Client)
}

func TestSockerImp_ContainerStop(t *testing.T) {
	order := Order{
		Order:   "run",
		Content: Marshal1(),
	}
	s := sockerImp{}
	s.ContainerStop(order)
}

func TestSockerImp_ContainerLogs(t *testing.T) {
	s := sockerImp{}
	order := Order{
		Order:   "run",
		Content: Marshal1(),
	}
	s.ContainerLogs(Client, order)
}