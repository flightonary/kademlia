package kademlia

import "container/list"

const BucketSize = 20

type routingTable struct {
	ownId *KadID
	table [KadIdLen]list.List
}

func newRoutingTable(ownId *KadID) *routingTable {
	return &routingTable{ownId: ownId}
}

func (rt *routingTable) add(node *Node) (bool, *Node) {
	if &node.Id != rt.ownId && rt.find(&node.Id) == nil {
		index := rt.index(&node.Id)
		rt.table[index].PushBack(node)
		// TODO: check if list len is longer than 20
	}
	return true, nil
}

func (rt *routingTable) del(kid *KadID) {
	index := rt.index(kid)
	for e := rt.table[index].Front(); e != nil; e = e.Next() {
		if e.Value.(*Node).Id == *kid {
			rt.table[index].Remove(e)
			return
		}
	}
}

func (rt *routingTable) find(kid *KadID) *Node {
	index := rt.index(kid)
	for e := rt.table[index].Front(); e != nil; e = e.Next() {
		if e.Value.(*Node).Id == *kid {
			return e.Value.(*Node)
		}
	}
	return nil
}

func (rt *routingTable) closest(kid *KadID) []*Node {
	// TODO: make sorted array of 20 closest nodes
	nodes := make([]*Node, 0)
	for i := 0; i < KadIdLen; i++ {
		for e := rt.table[i].Front(); e != nil; e = e.Next() {
			nodes = append(nodes, e.Value.(*Node))
		}
	}
	return nodes
}

func (rt *routingTable) index(kid *KadID) int {
	distance := rt.xor(kid)
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

func (rt *routingTable) xor(kid *KadID) *KadID {
	xor := &KadID{}
	for i := range kid {
		xor[i] = rt.ownId[i] ^ kid[i]
	}
	return xor
}