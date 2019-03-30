package kademlia

import (
	"container/list"
	"sort"
)

const BucketSize = 20

type routingTable struct {
	ownId *KadID
	table [KadIdLen]list.List
}

func newRoutingTable(ownId *KadID) *routingTable {
	return &routingTable{ownId: ownId}
}

func (rt *routingTable) add(node *Node) bool {
	index := rt.index(&node.Id)
	if &node.Id != rt.ownId && rt.find(&node.Id) == nil && rt.table[index].Len() <= BucketSize {
		rt.table[index].PushBack(node)
		return true
	}
	return false
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

func (rt *routingTable) closer(kid *KadID) []*Node {
	list2slice := func(li list.List) []*Node {
		nodes := make([]*Node, 0)
		for e := li.Front(); e != nil; e = e.Next() {
			nodes = append(nodes, e.Value.(*Node))
		}
		return nodes
	}

	closestIndex := rt.index(kid)
	nodes := list2slice(rt.table[closestIndex])
	for i := 1; i < KadIdLen; i++ {
		upper := closestIndex + i
		lower := closestIndex - i
		tmp := make([]*Node, 0)
		if upper < KadIdLen {
			tmp = append(tmp, list2slice(rt.table[upper])...)
		}
		if lower >= 0 {
			tmp = append(tmp, list2slice(rt.table[lower])...)
		}
		sort.Slice(tmp, func(i, j int) bool {
			iXor := xor(&tmp[i].Id, kid)
			jXor := xor(&tmp[j].Id, kid)
			for ii := 0; ii < KadIdLenByte; ii++ {
				if iXor[ii] == jXor[ii] {
					continue
				}
				return iXor[ii] < jXor[ii]
			}
			return true
		})
		nodes = append(nodes, tmp...)
		if len(nodes) >= BucketSize {
			return nodes[:BucketSize]
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
	return xor(rt.ownId, kid)
}

func xor(kid1 *KadID, kid2 *KadID) *KadID {
	xor := &KadID{}
	for i := range kid1 {
		xor[i] = kid1[i] ^ kid2[i]
	}
	return xor
}