package service

import (
	"fmt"
	"strings"
)

var (
	Server	= "tcp://121.40.101.210:1883"
)

var (
	SysDataPub		= "sysDataPub"
	SysOrderSub		= "sysOrderSub"
	SysStatusPub	= "sysStatusPub"
	SysLogPub		= "sysLogPub"
	SysGWSub		= "sysGWSub"
	SysCtnlsPub		= "sysCtnlsPub"
	SysImglsPub		= "sysImglsPub"
	SysGWErrPub		= "sysGWErrPub"
	//SysCNErrPub		= "sysCNErrPub"
)

var topics = map[string]string {
	"sysDataPub"	: "sys/{GW}/{CN}/msg",			//data up
	"sysOrderSub"	: "sys/{GW}/{CN}/order",		//order
	"sysStatusPub"	: "sys/{GW}/{CN}/online",		//online
	"sysLogPub"		: "sys/{GW}/{CN}/log",			//log up
	"sysGWSub"		: "sys/{GW}/order",				//gateway order
	"sysCtnlsPub"	: "sys/{GW}/Ctnls",				//pub container ls
	"sysImglsPub"	: "sys/{GW}/Imgls",				//pub container ls
	"sysGWErrPub"	: "sys/{GW}/err",				//gateway err
	//"sysCNErrPub"	: "sys/{GW}/{CN}/err",			//container err
}

type Order struct {
	Order   string		`json:"order"`
	Content string		`json:"content"`
}

func GetTopic(key string) string {
	return topics[key]
}

func Replace(gatewayName string) {
	for n := range topics {
		topics[n] = strings.Replace(topics[n], "{GW}", gatewayName, -1)
	}
	for n := range topics{
		fmt.Print(topics[n] + "\n")
	}
}
