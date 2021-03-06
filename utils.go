package kademlia

import (
	"errors"
	"math/rand"
	"net"
	"time"
)

func GetHostIp() (net.IP, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				kadlog.debug("host ip is ", ipnet.IP)
				return ipnet.IP, nil
			}
		}
	}
	return nil, errors.New("can not find ip address other than loopback")
}

func GenerateRandomKadId() KadID {
	rand.Seed(time.Now().UnixNano())
	kid := KadID{}
	rand.Read(kid[:])
	return kid
}
