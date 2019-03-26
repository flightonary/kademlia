package kademlia

import (
	"net"
)

const KadIdLen = 160 // bits
const KadIdLenByte = KadIdLen / 8 // bytes

type KadID [KadIdLenByte]byte

type Node struct {
	Id   KadID
	IP   net.IP
	Port int
}

func NewNode(id KadID, ip net.IP, port int) *Node {
	return &Node{id, ip, port}
}