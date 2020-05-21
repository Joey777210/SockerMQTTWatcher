package mqttStruct

import (
	log "github.com/Sirupsen/logrus"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"io"
	"net"
)

func Listen(port string, client mqtt.Client) {
	service := port
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	lisenter, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Errorf("listen port 8888 error %v", err)
	}

	for {
		conn, err := lisenter.Accept()
		if err != nil {
			log.Errorf("listen port 8888 error %v", err)
		}

		go connectionHandler(conn, client)
	}
}

func connectionHandler(connection net.Conn,client mqtt.Client) {
	buf := make([]byte, 1024)
	var data string
	for {
		n, err := connection.Read(buf)
		if err == io.EOF {
			break
		}
		data += string(buf[:n])
	}
	msgPub(client, data)
}
