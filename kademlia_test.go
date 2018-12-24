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

	err := kad1.Bootstrap("127.0.0.1", 9999)
	if err != nil {
		t.Fatalf("failed kad1.Bootstrap %#v", err)
	}
	time.Sleep(1 * time.Second)

	err = kad2.Bootstrap("127.0.0.1", 7005)
	if err != nil {
		t.Fatalf("failed kad1.Bootstrap %#v", err)
	}
	time.Sleep(1 * time.Second)

	kad1.Leave()
	kad2.Leave()
}