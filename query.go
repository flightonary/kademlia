package kademlia

import (
	"bytes"
	"encoding/gob"
)

const (
	typePingQuery = iota
	typeStoreQuery
	typeFindNodeQuery
	typeFindValueQuery
	typePingReply
	typeStoreReply
	typeFindNodeReply
	typeFindValueReply
)

type kademliaMessage struct {
	origin  *Node
	queryId int64
	typeId  int
	body    interface{}
}

type pingQuery struct {
	target KadID
}

type findNodeQuery struct {
	target KadID
}

type findValueQuery struct {
	target KadID
}

type dataStoreQuery struct {
	data KadID
}

type pingReply struct {
	target KadID
}

type findNodeReply struct {
	closest []*Node
}

type findValueReply struct {
	closest  []*Node
	hasValue bool
	value    []byte
}

type dataStoreReply struct {
	success bool
}

func init() {
	gob.Register(&pingQuery{})
	gob.Register(&findNodeQuery{})
	gob.Register(&findValueQuery{})
	gob.Register(&dataStoreQuery{})
	gob.Register(&pingReply{})
	gob.Register(&findNodeReply{})
	gob.Register(&findValueReply{})
	gob.Register(&dataStoreReply{})
}

func serializeKademliaMessage(msg *kademliaMessage) ([]byte, error) {
	var msgBuffer bytes.Buffer
	enc := gob.NewEncoder(&msgBuffer)
	err := enc.Encode(msg)
	if err != nil {
		return nil, err
	}
	return msgBuffer.Bytes(), nil
}

func deserializeKademliaMessage(rawMsg []byte) (*kademliaMessage, error) {
	reader := bytes.NewBuffer(rawMsg)
	msg := &kademliaMessage{}
	dec := gob.NewDecoder(reader)
	err := dec.Decode(msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}
