package kademlia

import (
	"net"
	"testing"
)

func TestPingQuery(t *testing.T) {
	SetLogLevelDebug()

	msg := commonKademliaMessage()
	msg.TypeId = typePingQuery
	msg.Body = &pingQuery{KadID{}}

	deserialized := commonCheck(t, msg)
	switch b := deserialized.Body.(type) {
	case *pingQuery:
		p := msg.Body.(*pingQuery)
		if p.Target != b.Target {
			t.Fatalf("failed to compere target")
		}
	default:
		t.Fatalf("failed to cast pingQuery")
	}
}

func commonCheck(t *testing.T, msg *kademliaMessage) *kademliaMessage {
	serialized, err := serializeKademliaMessage(msg)
	if err != nil {
		t.Fatalf("failed to serialize KademliaMessage, %#v", err)
	}
	deserialized, err := deserializeKademliaMessage(serialized)
	if err != nil {
		t.Fatalf("failed to deserialize KademliaMessage, %#v", err)
	}

	if msg.Origin.Id != deserialized.Origin.Id {
		t.Fatalf("failed to compere kadId")
	}
	if msg.Origin.IP.String() != deserialized.Origin.IP.String() {
		t.Fatalf("failed to compere kadId")
	}
	if msg.Origin.Port != deserialized.Origin.Port {
		t.Fatalf("failed to compere kadId")
	}
	if msg.TypeId != deserialized.TypeId {
		t.Fatalf("failed to compere typeId")
	}
	if msg.QuerySN != deserialized.QuerySN {
		t.Fatalf("failed to compere querySN")
	}

	return deserialized
}

func commonKademliaMessage() *kademliaMessage {
	addr, _ := net.ResolveIPAddr("ip", "0.0.0.0")
	node := NewNode(KadID{}, addr.IP, 9999)
	return &kademliaMessage{
		Origin:  node,
		QuerySN: 0,
	}
}