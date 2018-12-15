package kademlia

import (
	"fmt"
	"net"
)

const (
	KadIdLen = 160 // bits
)

type Node struct {
	id   []byte
	ip   net.IP
	port int
}

func NewNode(id []byte, ip net.IP, port int) (*Node, error) {
	if len(id) != KadIdLen / 8 {
		return nil, fmt.Errorf("length of id must be %d bits", KadIdLen)
	}
	return &Node{id, ip, port}, nil
}