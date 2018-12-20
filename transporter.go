package kademlia

import (
	"fmt"
	"net"
)

type sendMsg struct {
	requestId int64
	destIp    net.IP
	destPort  int
	data      []byte
}

type rcvMsg struct {
	srcIp   net.IP
	srcPort int
	data    []byte
}

type transporter interface {
	run(net.IP, int) error
	stop()
}

func newUdpTransporter() *udpTransporter {
	sendChan := make(chan *sendMsg, 10)
	rcvChan := make(chan *rcvMsg, 10)
	return &udpTransporter{sendChan, rcvChan, true, nil}
}

type udpTransporter struct {
	sendChan chan *sendMsg
	rcvChan  chan *rcvMsg
	stopFlag bool
	conn     *net.UDPConn
}

func (ut *udpTransporter) run(listenIp net.IP, listenPort int) error {
	src, err := net.ResolveUDPAddr("udp", joinHostPort(listenIp, listenPort))
	if err != nil {
		return err
	}

	conn, err := net.ListenUDP("udp", src)
	if err != nil {
		return err
	}

	ut.stopFlag = false
	ut.conn = conn

	go func() {
		var buf [1500]byte

		for {
			// TODO: use Read or ReadMsgIP - https://golang.org/src/net/iprawsock.go?s=207:773#L2
			// TODO: user SetTimer
			c, addr, err := conn.ReadFromUDP(buf[0:])
			if err != nil {
				if ut.stopFlag {
					kadlog.debugln("stopped running udpTransporter @ receiver goroutine")
					break
				}
				// TODO: return err via rcvChan or errChan
				continue
			}

			data := make([]byte, c)
			copy(data, buf[:c])
			ut.rcvChan <- &rcvMsg{addr.IP, addr.Port, data}
		}

	}()

	go func() {
		for {
			sendMsg, ok := <-ut.sendChan
			if ok {
				dst, err := net.ResolveUDPAddr("udp", joinHostPort(sendMsg.destIp, sendMsg.destPort))
				if err != nil {
					// TODO: return err via rcvChan or errChan
					continue
				}

				_, err = conn.WriteToUDP(sendMsg.data, dst)
				if err != nil {
					// TODO: return err via rcvChan or errChan
					continue
				}
			} else {
				kadlog.debugln("stopped running udpTransporter @ sender goroutine")
				break
			}
		}
	}()

	return nil
}

func (ut *udpTransporter) stop() {
	ut.stopFlag = true
	close(ut.sendChan)
	close(ut.rcvChan)
	_ = ut.conn.Close()
	ut.conn = nil
}

func joinHostPort(ip net.IP, port int) string {
	return fmt.Sprintf("%s:%d", ip.String(), port)
}
