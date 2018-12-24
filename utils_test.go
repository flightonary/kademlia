package kademlia

import (
	"testing"
)

func TestGetHostIp(t *testing.T) {
	SetLogLevelDebug()

	ip, err := GetHostIp()
	if err != nil {
		t.Fatalf("GetHostIp invokes error, %#v", err)
	}

	if ip.To4() == nil {
		t.Fatalf("IP doesn't ipv4, %#v", ip)
	}

	if ip.IsLoopback() {
		t.Fatalf("IP is loopback, %#v", ip)
	}
}

func TestGenerateRandomKadId(t *testing.T) {
	kid := GenerateRandomKadId()
	for i := 0; i < 100; i++ {
		tmp := GenerateRandomKadId()
		if kid == tmp {
			t.Fatal("id is not random")
		}
		kid = tmp
	}
}
