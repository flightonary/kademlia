package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/flightonary/kademlia"
)

func main() {
	if len(os.Args) < 5 {
		fmt.Printf("Usage: %s LISTEN_IP LISTEN_PORT BOOTSTRAP_IP BOOTSTRAP_PORT\n", os.Args[0])
		os.Exit(1)
	}

	myIp, resolveIpError := net.ResolveIPAddr("ip", os.Args[1])
	myPort, atoiError1 := strconv.Atoi(os.Args[2])
	bsIp := os.Args[3]
	bsPort, atoiError2 := strconv.Atoi(os.Args[4])

	if resolveIpError != nil {
		fmt.Println(resolveIpError)
		os.Exit(1)
	}
	if atoiError1 != nil || atoiError2 != nil {
		fmt.Println("LISTEN_PORT and BOOTSTRAP_PORT must be integer")
		os.Exit(1)
	}

	kademlia.SetLogLevelDebug()

	id := kademlia.GenerateRandomKadId()
	node := kademlia.NewNode(id, myIp.IP, myPort)
	kad := kademlia.NewKademlia(node)

	err := kad.Bootstrap(bsIp, bsPort)
	if err != nil {
		fmt.Println(err)
	}

	stdin := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	for stdin.Scan(){
		input := stdin.Text()
		switch input {
		case "help":
			fmt.Println(`command list:
  help:              show help
  show rt:           show routing table
  store <key:value>: store <key:value>
  find <key>:        find value of <key>
  quit:              terminate program
`)
		case "quit":
			os.Exit(0)
		case "":
		default:
			fmt.Println("unknown command")
		}
		fmt.Print("> ")
	}
}
