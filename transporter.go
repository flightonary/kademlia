package kademlia

import (
	"net"
)

type sendMsg struct {
	destIp   net.IP
	destPort int
	data     []byte
}

type rcvMsg struct {
	srcIp   net.IP
	srcPort int
	data    []byte
}

type transporter interface {
	run(net.IP, int) error
	stop()
	send(*sendMsg)
	receiveChannel() chan *rcvMsg
}

func newUdpTransporter() *udpTransporter {
	rcvChan := make(chan *rcvMsg, 10)
	return &udpTransporter{nil, rcvChan, true, nil}
}

type udpTransporter struct {
	sendChan chan *sendMsg
	rcvChan  chan *rcvMsg
	stopFlag bool
	conn     *net.UDPConn
}

func (ut *udpTransporter) run(listenIp net.IP, listenPort int) error {
	src := &net.UDPAddr{listenIp, listenPort, ""}
	conn, err := net.ListenUDP("udp", src)
	if err != nil {
		return err
	}

	ut.stopFlag = false
	ut.conn = conn
	ut.sendChan = make(chan *sendMsg, 10)

	// receiver routine
	go func() {
		var buf [1500]byte

		for {
			// TODO: use Read or ReadMsgIP - https://golang.org/src/net/iprawsock.go?s=207:773#L2
			// TODO: use SetTimer
			n, addr, err := conn.ReadFromUDP(buf[0:])
			if err != nil {
				if ut.stopFlag {
					kadlog.debug("stopped running udpTransporter @ receiver goroutine")
					break
				}
				// TODO: return err via rcvChan or errChan
				continue
			}

			data := make([]byte, n)
			copy(data, buf[:n])
			ut.rcvChan <- &rcvMsg{addr.IP, addr.Port, data}
		}

	}()

	// sender routine
	go func() {
		for sendMsg := range ut.sendChan {
			dst := &net.UDPAddr{sendMsg.destIp, sendMsg.destPort, ""}
			_, err = conn.WriteToUDP(sendMsg.data, dst)
			if err != nil {
				kadlog.debug(err)
				// TODO: return err via rcvChan or errChan
			}
		}
		kadlog.debug("stopped running udpTransporter @ sender goroutine")
	}()

	return nil
}

func (ut *udpTransporter) stop() {
	ut.stopFlag = true
	close(ut.sendChan)
	_ = ut.conn.Close()
	ut.conn = nil
}

func (ut *udpTransporter) send(msg *sendMsg) {
	ut.sendChan <- msg
}

func (ut *udpTransporter) receiveChannel() chan *rcvMsg {
	return ut.rcvChan
}
