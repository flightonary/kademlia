package kademlia

import (
	"net"
	"testing"
	"time"
)

func TestBootstrap(t *testing.T) {
	SetLogLevelDebug()

	addr, _ := net.ResolveIPAddr("ip", "127.0.0.1")
	kadNode1 := NewNode(GenerateRandomKadId(), addr.IP, 7005)
	kadNode2 := NewNode(GenerateRandomKadId(), addr.IP, 7006)
	kad1 := NewKademlia(kadNode1)
	kad2 := NewKademlia(kadNode2)

	err1 := kad1.Bootstrap("127.0.0.1", 9999)
	if err1 != nil {
		t.Fatalf("failed kad1.Bootstrap %#v", err1)
	}
	time.Sleep(1 * time.Second)

	err2 := kad2.Bootstrap("127.0.0.1", 7005)
	if err2 != nil {
		t.Fatalf("failed kad1.Bootstrap %#v", err2)
	}
	time.Sleep(1 * time.Second)

	kad1.Leave()
	kad2.Leave()
}