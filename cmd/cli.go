package main

import (
	"fmt"
	"github.com/flightonary/kademlia"
)

func main() {
	kad := kademlia.NewKademlia()
	fmt.Print(kad.Status())
}
