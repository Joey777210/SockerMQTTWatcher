package service

import (
	"fmt"
	"strings"
)

var (
	Server	= "tcp://121.40.101.210:1883"
)

var (
	SysOrderSub		= "sysOrderSub"
	SysStatusPub	= "sysStatusPub"
	SysLogPub		= "sysLogPub"
	SysGWLogPub		= "sysGWLogPub"
	SysDataPub		= "sysDataPub"
	SysCtnlsPub		= "sysCtnlsPub"
	SysImglsPub		= "sysImglsPub"
	SysGWErrPub		= "sysGWErrPub"
	SysACKPub		= "sysACKPub"
	//SysCNErrPub		= "sysCNErrPub"
)

var topics = map[string]string {
	"sysOrderSub"	: "sys/{GW}/order",				//order
	"sysDataPub"	: "sys/{GW}/{CN}/msg",			//data up
	"sysStatusPub"	: "sys/{GW}/{CN}/online",		//online
	"sysLogPub"		: "sys/{GW}/{CN}/log",			//container log
	"sysGWLogPub"	: "sys/{GW}/log",				//gateway (socker) log
	"sysCtnlsPub"	: "sys/{GW}/ctnls",				//pub container ls
	"sysImglsPub"	: "sys/{GW}/imgls",				//pub container ls
	"sysGWErrPub"	: "sys/{GW}/err",				//gateway err
	"sysACKPub"		: "sys/{GW}/ack",
	//"sysCNErrPub"	: "sys/{GW}/{CN}/err",			//container err
}

type Order struct {
	Target	string		`json:"target"`		//container/image/network
	Order   string		`json:"order"`		//run/stop/ls...
	Name	string		`json:"name"`		//bird...
	Content string		`json:"content"`	//memory....
}

func GetTopic(key string) string {
	return topics[key]
}

func GetTopicCN(key string, containerName string) string {
	topic := topics[key]
	topic = strings.Replace(topic, "{CN}", containerName, -1)
	return topic
}


func Replace(gatewayName string) {
	for n := range topics {
		topics[n] = strings.Replace(topics[n], "{GW}", gatewayName, -1)
	}
	for n := range topics{
		fmt.Print(topics[n] + "\n")
	}
}

