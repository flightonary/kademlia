package kademlia

import (
	"container/list"
	"fmt"
	"testing"
)

func TestXor(t *testing.T) {
	rt := routingTable{}

	k1 := fill(&KadID{}, 0x00)
	k2 := fill(&KadID{}, 0xFF)
	k2[0] = 0xFE
	rt.ownId = k1
	xor := rt.xor(k2)

	ans := fill(&KadID{}, 0xFF)
	ans[0] = 0xFE

	if *xor != *ans {
		t.Fatal("xor calculation error")
	}
}

func TestIndex(t *testing.T) {
	ownId := fill(&KadID{}, 0x00)
	rt := routingTable{}
	rt.ownId = ownId

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

func Test_routingTable_add(t *testing.T) {
	ownId := fill(&KadID{}, 0x00)
	ownId[0] = 0x01
	node := Node{}
	node.Id = *fill(&KadID{}, 0x00)
	node.Id[19] = 0x01

	type fields struct {
		ownId *KadID
		table [KadIdLen]list.List
	}
	type args struct {
		node *Node
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
		finally func(*routingTable) (bool, string)
	}{
		{
			name: "simple add node",
			fields: fields{
				ownId: ownId,
				table: [KadIdLen]list.List{},
			},
			args: args{node: &node},
			want: true,
			finally: func(rt *routingTable) (bool, string) {
				index := rt.index(&node.Id)
				e := rt.table[index].Front()
				if e == nil {
					return false, "Node is not added in the list."
				}
				if e.Value.(*Node) != &node {
					return false, fmt.Sprintf("Added node kid = %v, want %v", e.Value.(*Node).Id, node.Id)
				}
				return true, ""
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt := &routingTable{
				ownId: tt.fields.ownId,
				table: tt.fields.table,
			}
			got := rt.add(tt.args.node)
			if got != tt.want {
				t.Errorf("routingTable.add() got = %v, want %v", got, tt.want)
			}
			if tt.finally != nil {
				success, errMsg := tt.finally(rt)
				if !success {
					t.Error(errMsg)
				}
			}
		})
	}
}
