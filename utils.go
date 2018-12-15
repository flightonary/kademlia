package kademlia

import (
	"errors"
	"net"
)

func GetHostIp() (net.IP, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			return ipnet.IP, nil
		}
	}
	return nil, errors.New("can not find ip address other than loopback")
}

func GenerateRondomId() KadId {
	return KadId{}
}