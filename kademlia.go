package kademlia



type Kademlia struct {
	myself   *node
	nodes    []*node
	mainChan *chan interface{}
}

func (kad *Kademlia) Bootstrap(entryNode string) {

}
