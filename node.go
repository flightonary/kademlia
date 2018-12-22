package kademlia

import (
	"net"
)

const KadIdLen = 160 // bits

type KadID [KadIdLen / 8]byte

type Node struct {
	Id   KadID
	IP   net.IP
	Port int
}

func NewNode(id KadID, ip net.IP, port int) (*Node, error) {
	return &Node{id, ip, port}, nil
}
