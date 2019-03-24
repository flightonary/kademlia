package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

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

	kad.SetFindValueCallback(func(key string, value []byte){
		fmt.Printf("key: %s, value: %s\n", key, string(value))
	})

	stdin := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	for stdin.Scan(){
		input := stdin.Text()
		inputs := strings.Split(input, " ")
		switch inputs[0] {
		case "help":
			fmt.Println(`command list:
  help               - show help
  show               - show routing table
  store <key:value>  - store <key:value>
  find <key>         - find value of <key>
  quit               - terminate program
`)
		case "show":
			rt := kad.GetRoutingTable()
			for i := 0; i < kademlia.KadIdLen; i++ {
				if rt[i].Len() > 0 {
					fmt.Printf("distance 2^%02d:", kademlia.KadIdLen - i)
					for e := rt[i].Front(); e != nil; e = e.Next() {
						fmt.Printf(" %x", e.Value.(*kademlia.Node).Id)
					}
					fmt.Println("")
				}
			}
		case "store":
			kv := strings.Split(inputs[1], ":")
			kad.Store(kv[0], []byte(kv[1]))
		case "find":
			kad.FindValue(inputs[1])
		case "quit":
			os.Exit(0)
		case "":
		default:
			fmt.Println("unknown command")
		}
		fmt.Print("> ")
	}
}
