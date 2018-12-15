package main

import (
	"log"

	"github.com/flightonary/kademlia"
)

func main() {
	node := kademlia.Node{}
	node.Id = kademlia.GenerateRondomId()
	node.IP, _ = kademlia.GetHostIp()
	node.Port = 7001

	kad := kademlia.NewKademlia(node)

	err := kad.Bootstrap("localhost", 7009)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("ok")
}
