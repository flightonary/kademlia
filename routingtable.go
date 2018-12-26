package kademlia

import "container/list"

const BucketSize = 20

type routingInterface interface {
	add(node *Node) (bool, *Node)
	del(kid *KadID)
	closest(kid *KadID) []*Node
}

type routingTable struct {
	ownId KadID
	table [KadIdLen]list.List
}

func (*routingTable) add(node *Node) (bool, *Node) {
	panic("implement me")
}

func (*routingTable) del(kid *KadID) {
	panic("implement me")
}

func (*routingTable) closest(kid *KadID) []*Node {
	panic("implement me")
}

func (rt *routingTable) index(kid *KadID) int {
	distance := rt.xor(&rt.ownId, kid)
	firstBitIndex := 0
	for _, v := range distance {
		if v == 0 {
			firstBitIndex += 8
			continue
		}
		for i := 0; i < 8; i++ {
			if v & (0x80 >> uint(i)) != 0 {
				break
			}
			firstBitIndex++
		}
		break
	}
	return firstBitIndex
}

func (*routingTable) xor(kid1 *KadID, kid2 *KadID) *KadID {
	xor := &KadID{}
	for i := range kid1 {
		xor[i] = kid1[i] ^ kid2[i]
	}
	return xor
}