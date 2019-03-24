package kademlia

import (
	"bytes"
	"container/list"
	"crypto/sha1"
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
	store       *dataStore
	querySN     int64 // TODO: change to transaction Id
}

func NewKademlia(own *Node) *Kademlia {
	kad := &Kademlia{}
	kad.own = own
	kad.mainChan = make(chan ctrlCmd, 10)
	kad.endChan = make(chan interface{})
	kad.transporter = newUdpTransporter()
	kad.rt = newRoutingTable(&own.Id)
	kad.store = newDataStore()
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

	kadlog.debugf("run bootstrap, own kid: %x", kad.own.Id)

	return nil
}

func (kad *Kademlia) Leave() {
	close(kad.endChan)
	kad.transporter.stop()
}

func (kad *Kademlia) Store(key string, value string) {
	keyKid := kad.toKadID(key)
	kad.store.Put(keyKid, []byte(value))

	// store data in each closest node
	closest := kad.rt.closest(keyKid)
	for _, node := range closest {
		err := kad.sendStoreQuery(node, keyKid)
		if err != nil {
			kadlog.debug(err)
		}
	}
}

func (kad *Kademlia) FindValue(key string) string {
	// TODO: query FindValue to closest nodes when the value is unknown.
	return 	kad.store.GetAsString(kad.toKadID(key))
}

// For debug purpose
func (kad *Kademlia) GetRoutingTable() [KadIdLen]list.List {
	return kad.rt.table
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
					// TODO: check if the message is for me
					kadlog.debugf("receive FindNodeQuery from Node(%x)", kadMsg.Origin.Id)
					// reply with closest nodes
					closest := kad.rt.closest(&query.Target)
					err := kad.sendFindNodeReply(rcvMsg.srcIp, rcvMsg.srcPort, kadMsg.QuerySN, closest)
					if err != nil {
						kadlog.debug(err)
					}
					// add source node to routing table
					kad.rt.add(kadMsg.Origin)
				case *findNodeReply:
					kadlog.debugf("receive FindNodeReply from Node(%x)", kadMsg.Origin.Id)
					// add source node to routing table
					kad.rt.add(kadMsg.Origin)
					// add new node to routing table and send FindNodeQuery if it is unknown
					for _, node := range query.Closest {
						if kad.isNotSameHost(node) && kad.rt.find(&node.Id) == nil {
							err := kad.sendFindNodeQuery(node.IP, node.Port, node.Id)
							if err != nil {
								kadlog.debug(err)
							}
						} else {
							// TODO: move node to last of list
						}
						kad.rt.add(node)
					}
				case *storeQuery:
					kadlog.debugf("receive StoreQuery from Node(%x)", kadMsg.Origin.Id)
					// TODO: add node to routing table if it's unknown
					// TODO: move node to last of list if it's known
					kad.store.Put(&query.Key, query.Data)
					// TODO: send StoreReply
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

func (kad *Kademlia) isNotSameHost(node *Node) bool {
	return bytes.Compare(kad.own.IP, node.IP) != 0 && kad.own.Port != node.Port
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
	kadlog.debugf("send FindNodeQuery to Node(%x)", target)
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

func (kad *Kademlia) sendStoreQuery(node *Node, key *KadID) error {
	kadMsg := kad.baseKademliaMessage()
	kadMsg.TypeId = typeStoreQuery
	kadMsg.QuerySN = kad.newSN()
	kadMsg.Body = &storeQuery{*key, kad.store.Get(key)}
	data, err := serializeKademliaMessage(kadMsg)
	if err != nil {
		return err
	}
	msg := &sendMsg{node.IP, node.Port, data}
	kad.transporter.send(msg)
	kadlog.debug("send StoreQuery")
	return nil
}

func (kad *Kademlia) toKadID(key string) *KadID {
	bs := sha1.Sum([]byte(key))
	kid := KadID{}
	copy(kid[:], bs[:])
	return &kid
}

func (kad *Kademlia) newSN() int64 {
	kad.querySN++
	return kad.querySN
}
