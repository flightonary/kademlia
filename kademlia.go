package kademlia

import "net"

type ID [160]byte

func (id *ID) toString() {

}

type Node struct {
	ID []byte
	IP net.IP
	Port int
}

type Kbucket struct {
	Nodes []Node
}

type RoutingTable struct {
	buckets []Kbucket
}

type Kademlia struct {
	rt RoutingTable
	store map[string][]byte
}

func NewKademlia() Kademlia {
	return Kademlia {}
}

func (kad *Kademlia) Bootstrap(node Node) {

}

func (kad *Kademlia) Ping(node Node) {

}

func (kad *Kademlia) Store(key []byte, value []byte) {

}

func (kad *Kademlia) FindNode(id []byte) {

}

func (kad *Kademlia) FindValue(key []byte) {

}

func (kad *Kademlia) Status() string {
	return "OK"
}