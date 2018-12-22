package kademlia

import (
	"net"
)

type ctrlCmd interface{}

type bootstrapCmd struct {
	ip   net.IP
	port int
}

type Kademlia struct {
	own         *Node
	nodes       []*Node
	mainChan    chan *ctrlCmd
	endChan     chan *interface{}
	transporter transporter
	querySN	    int64
}

func NewKademlia(own *Node) *Kademlia {
	kad := &Kademlia{}
	kad.own = own
	kad.nodes = []*Node{}
	kad.mainChan = make(chan *ctrlCmd, 10)
	kad.endChan = make(chan *interface{})
	kad.transporter = newUdpTransporter()
	kad.querySN = 0
	return kad
}

func (kad *Kademlia) Bootstrap(entryNodeAddr string, entryNodePort int) error {
	err := kad.transporter.run(kad.own.ip, kad.own.port)
	if err != nil {
		return err
	}

	go kad.mainRoutine()

	addr, err := net.ResolveIPAddr("ip", entryNodeAddr)
	if err != nil {
		return err
	}
	var bsCmd ctrlCmd = bootstrapCmd{addr.IP, entryNodePort}
	kad.mainChan <- &bsCmd

	return nil
}

func (kad *Kademlia) Leave() {
	close(kad.endChan)
}

func (kad *Kademlia) mainRoutine() {
	select {
		case cmd := <-kad.mainChan:
			switch c := (*cmd).(type) {
			case bootstrapCmd:
				kadlog.debug("receive cmd", c)
				err := kad.sendFindNodeQuery(c.ip, c.port, kad.own.id)
				if err != nil {
					kadlog.debug(err)
				}
			default:
				kadlog.debug("receive cmd", c)
			}
		case <- kad.endChan:
			kadlog.debug("leave from kademlia cluster")
			return
	}
}

func (kad *Kademlia) baseKademliaMessage() *kademliaMessage {
	kad.querySN++
	return &kademliaMessage{
		origin: kad.own,
		querySN: kad.querySN,
	}
}

func (kad *Kademlia) sendFindNodeQuery(ip net.IP, port int, target KadID) error {
	kadMsg := kad.baseKademliaMessage()
	kadMsg.typeId = typeFindNodeQuery
	kadMsg.body = &findNodeQuery{target}
	data, err := serializeKademliaMessage(kadMsg)
	if err != nil {
		return err
	}
	msg := &sendMsg{ip, port, data}
	kad.transporter.send(msg)
	return nil
}