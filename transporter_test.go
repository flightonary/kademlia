package kademlia

import (
	"net"
	"testing"
)

func TestSendAndReceive(t *testing.T) {
	SetLogLevelDebug()

	listenIp := net.ParseIP("0.0.0.0")
	local := net.ParseIP("127.0.0.1")
	serverPort := 7001
	clientPort := 7002

	server := newUdpTransporter()
	client := newUdpTransporter()

	err := server.run(listenIp, serverPort)
	if err != nil {
		t.Fatalf("failed server.run %#v", err)
	}

	err = client.run(listenIp, clientPort)
	if err != nil {
		t.Fatalf("failed client.run %#v", err)
	}

	sendData := []byte("Hello")
	client.send(&sendMsg{local, serverPort, sendData})

	rcvMsg, ok := <- server.receiveChannel()
	if !ok {
		t.Fatalf("failed server.rcvChan")
	}

	if string(rcvMsg.data) != "Hello" {
		t.Fatalf("failed string(rcvMsg.data) != Hello")
	}

	server.stop()
	client.stop()
}