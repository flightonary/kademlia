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
			kadlog.debugln("host ip is", ipnet.IP)
			return ipnet.IP, nil
		}
	}
	return nil, errors.New("can not find ip address other than loopback")
}

func GenerateRandomId() []byte {
	rand.Seed(time.Now().UnixNano())
	buff := make([]byte, KadIdLen / 8)
	rand.Read(buff)
	return buff
}