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