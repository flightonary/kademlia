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
	Origin  *Node
	QuerySN int64
	TypeId  int
	Body    interface{}
}

type pingQuery struct {
	Target KadID
}

type findNodeQuery struct {
	Target KadID
}

type findValueQuery struct {
	Target KadID
}

type dataStoreQuery struct {
	Data KadID
}

type pingReply struct {
	Target KadID
}

type findNodeReply struct {
	Closest []*Node
}

type findValueReply struct {
	Closest  []*Node
	HasValue bool
	Value    []byte
}

type dataStoreReply struct {
	Success bool
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
