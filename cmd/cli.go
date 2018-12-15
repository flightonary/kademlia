package main

import (
	"fmt"
	"github.com/flightonary/kademlia"
)

func main() {
	kad := &kademlia.Kademlia{}
	kad.Bootstrap("localhost")
	fmt.Print("ok")
}
