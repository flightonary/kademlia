package kademlia

import (
	"net"
)

type ctrlCmd interface{}

type findNodeCmd struct {
	ip   net.IP
	port int
}

type Kademlia struct {
	own         *Node
	nodes       []*Node
	mainChan    chan ctrlCmd
	transporter transporter
}

func NewKademlia(own *Node) *Kademlia {
	kad := &Kademlia{}
	kad.own = own
	kad.nodes = []*Node{}
	kad.mainChan = make(chan ctrlCmd, 10)
	return kad
}

func (kad *Kademlia) Bootstrap(entryNodeAddr string, entryNodePort int) error {
	addr, err := net.ResolveIPAddr("ip", entryNodeAddr)
	if err != nil {
		return err
	}

	go kad.mainRoutine()

	fnCmd := findNodeCmd{addr.IP, entryNodePort}
	kad.mainChan <- fnCmd

	return nil
}

func (kad *Kademlia) mainRoutine() {

}
