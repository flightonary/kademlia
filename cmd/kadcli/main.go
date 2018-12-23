package main

import (
	"log"

	"github.com/flightonary/kademlia"
)

func main() {
	kademlia.SetLogLevelDebug()

	id := kademlia.GenerateRandomId()
	ip, _ := kademlia.GetHostIp()
	port := 7001
	node := kademlia.NewNode(id, ip, port)

	kad := kademlia.NewKademlia(node)

	err := kad.Bootstrap("localhost", 7009)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("done")
}
