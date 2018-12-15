package kademlia

import "net"

const (
	KadIdLen = 160 // bits
)

type kadId [KadIdLen/8]byte

type node struct {
	id kadId
	ip net.IP
	port int
}