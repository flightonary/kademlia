package kademlia

import "net"

const (
	KadIdLen = 160 // bits
)

type KadId [KadIdLen / 8]byte

type Node struct {
	Id   KadId
	IP   net.IP
	Port int
}
