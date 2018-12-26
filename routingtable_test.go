package kademlia

import "testing"

func TestXor(t *testing.T) {
	rt := routingTable{}

	k1 := fill(&KadID{}, 0x00)
	k2 := fill(&KadID{}, 0xFF)
	k2[0] = 0xFE
	xor := rt.xor(k1, k2)

	ans := fill(&KadID{}, 0xFF)
	ans[0] = 0xFE

	if *xor != *ans {
		t.Fatal("xor calculation error")
	}
}

func TestIndex(t *testing.T) {
	ownId := fill(&KadID{}, 0x00)
	rt := routingTable{}
	rt.ownId = *ownId

	// in case that first bit is 1.
	k1 := fill(&KadID{}, 0x00)
	k1[0] = 0xFF
	index1 := rt.index(k1)
	if index1 != 0 {
		t.Fatalf("index calculation error: index is %d", index1)
	}

	// in case that last bit is 1.
	k2 := fill(&KadID{}, 0x00)
	k2[len(k2)-1] = 0x01
	index2 := rt.index(k2)
	if index2 != KadIdLen-1 {
		t.Fatalf("index calculation error: index is %d", index2)
	}

	// in case that middle bit is 1.
	k3 := fill(&KadID{}, 0x00)
	k3[10] = 0x0F
	index3 := rt.index(k3)
	if index3 != 8*10+4 {
		t.Fatalf("index calculation error: index is %d", index3)
	}
}

func fill(id *KadID, b byte) *KadID {
	for i := range id {
		id[i] = b
	}
	return id
}
