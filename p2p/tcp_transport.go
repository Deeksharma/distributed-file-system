package p2p

import (
	"fmt"
	"net"
	"sync"
)

// TCPPeer represents the remote node over a TCP established connection.
type TCPPeer struct {
	// conn is the underlying connection of the peer
	conn net.Conn
	// if we dial a conn and retrieve a conn => outbound == true
	// if we accept a conn => outbound == false
	outbound bool // tcp transport dial direction
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

// every kind of communication should have a channel, we are going to ass peers for node, then transmit data to all the peers

type TCPTransport struct {
	listenAddress string
	listener      net.Listener
	shakeHands    HandshakeFunc
	decoder       Decoder

	mu    sync.RWMutex // put mutex above the variable that you want to protect
	peers map[net.Addr]*Peer
}

func NewTCPTransport(listenAddr string) *TCPTransport {
	return &TCPTransport{
		listenAddress: listenAddr,
		shakeHands:    NOPHandshakeFunc,
	}
}

// ListenAndAccept listens and accept
func (t *TCPTransport) ListenAndAccept() error {
	var err error

	// connection establishment
	t.listener, err = net.Listen("tcp", t.listenAddress)
	if err != nil {
		return err
	}

	// data processing/transmission
	go t.startAcceptLoop()

	return nil
}

func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Printf("TCP accept error: %s\n", err.Error())
		}
		fmt.Printf("new incoming connection %+v \n", conn)
		go t.handleConn(conn)
	}
}

type Temp struct{}

func (t *TCPTransport) handleConn(conn net.Conn) {
	//defer conn.Close()
	peer := NewTCPPeer(conn, true)
	if err := t.shakeHands(peer); err != nil {
		
	}

	// Read loop
	//buf := new(bytes.Buffer)
	msg := &Temp{}
	for {
		//n, _ := conn.Read(buf.Bytes())
		if err := t.decoder.Decode(conn, msg); err != nil {
			fmt.Printf("TCP decode error: %s\n", err.Error())
			continue
		}
	}

}
