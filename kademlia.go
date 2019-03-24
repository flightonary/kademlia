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
	findValueCallback func(string, []byte)
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

func (kad *Kademlia) Store(key string, value []byte) {
	kad.store.Put(key, []byte(value))

	// store data in each closest node
	keyKid := kad.toKadID(key)
	closest := kad.rt.closest(keyKid)
	for _, node := range closest {
		err := kad.sendStoreQuery(node, key)
		if err != nil {
			kadlog.debug(err)
		}
	}
}

func (kad *Kademlia) FindValue(key string) {
	if kad.store.Exist(key) {
		if kad.findValueCallback != nil {
			kad.findValueCallback(key, kad.store.Get(key))
		}
	} else {
		closest := kad.rt.closest(kad.toKadID(key))
		if len(closest) > 0 {
			// inquiry to the closest node
			err := kad.sendFindValueQuery(closest[0], key)
			if err != nil {
				kadlog.debug(err)
			}
		}
	}
}

func (kad *Kademlia) SetFindValueCallback(fn func(string, []byte)) {
	kad.findValueCallback = fn
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
				err := kad.SendFindNodeQuery(cmd.ip, cmd.port, kad.own.Id)
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
					// TODO: move node to last of list when Origin is known
					kad.rt.add(kadMsg.Origin)
					// add new node to routing table and send FindNodeQuery if it is unknown
					for _, node := range query.Closest {
						if kad.isNotSameHost(node) && kad.rt.find(&node.Id) == nil {
							// TODO: own Id should be used in case of bootstrap but are there other FindNodeQuery use-case?
							err := kad.SendFindNodeQuery(node.IP, node.Port, kad.own.Id)
							if err != nil {
								kadlog.debug(err)
							}
						}
					}
				case *storeQuery:
					kadlog.debugf("receive StoreQuery from Node(%x)", kadMsg.Origin.Id)
					// TODO: add node to routing table if it's unknown
					// TODO: move node to last of list if it's known
					kad.store.Put(query.Key, query.Data)
					// TODO: send StoreReply
				case *findValueQuery:
					kadlog.debugf("receive FindValueQuery from Node(%x)", kadMsg.Origin.Id)
					hasValue := kad.store.Exist(query.Key)
					data := kad.store.Get(query.Key)
					var closest []*Node
					if hasValue {closest = kad.rt.closest(kad.toKadID(query.Key))}
					err := kad.sendFindValueReply(rcvMsg.srcIp, rcvMsg.srcPort, kadMsg.QuerySN, query.Key, hasValue, data, closest)
					if err != nil {
						kadlog.debug(err)
					}
				case *findValueReply:
					kadlog.debugf("receive FindValueReply from Node(%x)", kadMsg.Origin.Id)
					if query.HasValue {
						if kad.findValueCallback != nil {
							kad.findValueCallback(query.Key, query.Value)
						}
					} else {
						if len(query.Closest) > 0 {
							nextInquiryNode := query.Closest[0]
							err := kad.sendFindValueQuery(nextInquiryNode, query.Key)
							if err != nil {
								kadlog.debug(err)
							}
						}
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

func (kad *Kademlia) isNotSameHost(node *Node) bool {
	return bytes.Compare(kad.own.IP, node.IP) != 0 || kad.own.Port != node.Port
}

func (kad *Kademlia) baseKademliaMessage() *kademliaMessage {
	return &kademliaMessage{Origin: kad.own}
}

func (kad *Kademlia) sendKadMsg(ip net.IP, port int, kadMsg *kademliaMessage) error {
	data, err := serializeKademliaMessage(kadMsg)
	if err != nil {
		return err
	}
	msg := &sendMsg{ip, port, data}
	kad.transporter.send(msg)
	return nil
}

func (kad *Kademlia) SendFindNodeQuery(ip net.IP, port int, target KadID) error {
	kadMsg := kad.baseKademliaMessage()
	kadMsg.QuerySN = kad.newSN()
	kadMsg.TypeId = typeFindNodeQuery
	kadMsg.Body = &findNodeQuery{target}
	kadlog.debugf("send FindNodeQuery to Node(%x)", target)
	return kad.sendKadMsg(ip, port, kadMsg)
}

func (kad *Kademlia) sendFindNodeReply(ip net.IP, port int, sn int64, closest []*Node) error {
	kadMsg := kad.baseKademliaMessage()
	kadMsg.TypeId = typeFindNodeReply
	kadMsg.QuerySN = sn
	kadMsg.Body = &findNodeReply{closest}
	kadlog.debug("send FindNodeReply")
	return kad.sendKadMsg(ip, port, kadMsg)
}

func (kad *Kademlia) sendStoreQuery(node *Node, key string) error {
	kadMsg := kad.baseKademliaMessage()
	kadMsg.TypeId = typeStoreQuery
	kadMsg.QuerySN = kad.newSN()
	kadMsg.Body = &storeQuery{key, kad.store.Get(key)}
	kadlog.debugf("send StoreQuery to Node(%x)", node.Id)
	return kad.sendKadMsg(node.IP, node.Port, kadMsg)
}

func (kad *Kademlia) sendFindValueQuery(node *Node, key string) error {
	kadMsg := kad.baseKademliaMessage()
	kadMsg.TypeId = typeFindValueQuery
	kadMsg.QuerySN = kad.newSN()
	kadMsg.Body = &findValueQuery{key}
	kadlog.debugf("send FindValueQuery to Node(%x)", node.Id)
	return kad.sendKadMsg(node.IP, node.Port, kadMsg)
}

func (kad *Kademlia) sendFindValueReply(ip net.IP, port int, sn int64, key string, hasValue bool, data []byte, closest []*Node) error {
	kadMsg := kad.baseKademliaMessage()
	kadMsg.TypeId = typeFindValueReply
	kadMsg.QuerySN = sn
	kadMsg.Body = &findValueReply{key, hasValue, data, closest}
	kadlog.debug("send FindValueReply")
	return kad.sendKadMsg(ip, port, kadMsg)
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
