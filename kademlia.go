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
	mainChan    chan ctrlCmd
	endChan     chan interface{}
	transporter transporter
	rt          *routingTable
	querySN     int64 // TODO: change to transaction Id
}

func NewKademlia(own *Node) *Kademlia {
	kad := &Kademlia{}
	kad.own = own
	kad.mainChan = make(chan ctrlCmd, 10)
	kad.endChan = make(chan interface{})
	kad.transporter = newUdpTransporter()
	kad.rt = newRoutingTable(&own.Id)
	kad.querySN = 0
	return kad
}

func (kad *Kademlia) Bootstrap(entryNodeAddr string, entryNodePort int) error {
	// TODO: take net.IP from args directly
	addr, err := net.ResolveIPAddr("ip", entryNodeAddr)
	if err != nil {
		return err
	}
	err = kad.transporter.run(kad.own.IP, kad.own.Port)
	if err != nil {
		return err
	}

	var bsCmd ctrlCmd = bootstrapCmd{addr.IP, entryNodePort}
	kad.mainChan <- bsCmd

	go kad.mainRoutine()

	return nil
}

func (kad *Kademlia) Leave() {
	close(kad.endChan)
	kad.transporter.stop()
}

func (kad *Kademlia) mainRoutine() {
	for {
		select {
		// handle internal control command
		case ctrl := <-kad.mainChan:
			switch cmd := ctrl.(type) {
			case bootstrapCmd:
				kadlog.debug("receive bootstrapCmd")
				err := kad.sendFindNodeQuery(cmd.ip, cmd.port, kad.own.Id)
				if err != nil {
					kadlog.debug(err)
				}
			default:
				kadlog.debug("receive unknown ctrlCmd")
			}
		// handle received message from outside
		case rcvMsg := <-kad.transporter.receiveChannel():
			kadMsg, err := deserializeKademliaMessage(rcvMsg.data)
			if err != nil {
				kadlog.debug(err)
			} else {
				switch query := kadMsg.Body.(type) {
				case *findNodeQuery:
					kadlog.debugf("receive FindNodeQuery from Node(0x%x)", kadMsg.Origin.Id[0])
					// reply with closest nodes
					closest := kad.rt.closest(&query.Target)
					err := kad.sendFindNodeReply(rcvMsg.srcIp, rcvMsg.srcPort, kadMsg.QuerySN, closest)
					if err != nil {
						kadlog.debug(err)
					}
					// add source node to routing table
					kad.rt.add(kadMsg.Origin)
				case *findNodeReply:
					kadlog.debugf("receive FindNodeReply from Node(0x%x)", kadMsg.Origin.Id[0])
					// add source node to routing table
					kad.rt.add(kadMsg.Origin)
					// add new node to routing table and send FindNodeQuery if it is unknown
					for _, node := range query.Closest {
						if kad.rt.find(&node.Id) == nil {
							err := kad.sendFindNodeQuery(node.IP, node.Port, kad.own.Id)
							if err != nil {
								kadlog.debug(err)
							}
						} else {
							// TODO: move node to last of list
						}
						kad.rt.add(node)
					}
				default:
					kadlog.debug("receive unknown kademlia message")
				}
			}
		case <-kad.endChan:
			kadlog.debug("leave from kademlia cluster")
			return
		}
	}
}

func (kad *Kademlia) baseKademliaMessage() *kademliaMessage {
	return &kademliaMessage{
		Origin: kad.own,
	}
}

func (kad *Kademlia) sendFindNodeQuery(ip net.IP, port int, target KadID) error {
	kadMsg := kad.baseKademliaMessage()
	kadMsg.QuerySN = kad.newSN()
	kadMsg.TypeId = typeFindNodeQuery
	kadMsg.Body = &findNodeQuery{target}
	data, err := serializeKademliaMessage(kadMsg)
	if err != nil {
		return err
	}
	msg := &sendMsg{ip, port, data}
	kad.transporter.send(msg)
	kadlog.debugf("send FindNodeQuery to Node(0x%x)", target[0])
	return nil
}

func (kad *Kademlia) sendFindNodeReply(ip net.IP, port int, sn int64, closest []*Node) error {
	kadMsg := kad.baseKademliaMessage()
	kadMsg.TypeId = typeFindNodeReply
	kadMsg.QuerySN = sn
	kadMsg.Body = &findNodeReply{closest}
	data, err := serializeKademliaMessage(kadMsg)
	if err != nil {
		return err
	}
	msg := &sendMsg{ip, port, data}
	kad.transporter.send(msg)
	kadlog.debug("send FindNodeReply")
	return nil
}

func (kad *Kademlia) newSN() int64 {
	kad.querySN++
	return kad.querySN
}
