package p2p

import (
	"errors"
	"fmt"
	"log"
	"net"
	"reflect"
)

// TCPPeer represents the remote node over a TCP established connection.
type TCPPeer struct {
	// conn is the underlying connection of the peer which in this case a tcp connection
	net.Conn
	// if we dial a conn and retrieve a conn => outbound == true
	// if we accept a conn => outbound == false
	outbound bool // tcp transport dial direction
}

// Send implements the Peer interface
func (p *TCPPeer) Send(b []byte) error {
	_, err := p.Write(b)
	return err
}

//// RemoteAddr implements the Peer interface.
//func (p *TCPPeer) RemoteAddr() net.Addr {
//	return p.conn.RemoteAddr()
//}
//
//// Close implements the Peer interface.
//func (p *TCPPeer) Close() error {
//	return p.conn.Close()
//}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		Conn:     conn,
		outbound: outbound,
	}
}

// every kind of communication should have a channel, we are going to ass peers for node, then transmit data to all the peers

type TCPTransport struct {
	TCPTransportOps
	//listenAddress string
	listener net.Listener
	rpcchan  chan RPC
	//shakeHands    HandshakeFunc
	//decoder       Decoder

	// a server should be responsible for the peers not the transport, but the transport should be aware if there is a new peer
	// a notification should be sent to the sever
	//mu    sync.RWMutex // put mutex above the variable that you want to protect
	//peers map[net.Addr]*Peer
}

type TCPTransportOps struct {
	ListenAddr    string
	HandshakeFunc HandshakeFunc
	Decoder       Decoder
	// if this function returns an error then we are going drop the peer
	OnPeer func(Peer) error
}

func NewTCPTransport(opts TCPTransportOps) *TCPTransport {
	return &TCPTransport{
		TCPTransportOps: opts,
		rpcchan:         make(chan RPC),
	}
}

// Consume implements the Transport interface, which will return a readonly channel for reading the incoming messages received from another peer in the network.
func (t *TCPTransport) Consume() <-chan RPC { // from this channel only ead is possible
	return t.rpcchan
}

// ListenAddress implements the Transport interface.
func (t *TCPTransport) ListenAddress() string {
	return t.ListenAddr
}

// ListenAndAccept listens and accept
func (t *TCPTransport) ListenAndAccept() error {
	var err error

	// connection establishment
	t.listener, err = net.Listen("tcp", t.ListenAddr)
	if err != nil {
		return err
	}

	// data processing/transmission
	go t.startAcceptLoop()

	log.Printf("ListenAndAccept %s ok", t.ListenAddr)
	return nil
}

func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if errors.Is(err, net.ErrClosed) {
			return
		}
		if err != nil {
			fmt.Printf("TCP accept error: %s\n", err.Error())
		}
		fmt.Printf("new incoming connection %+v \n", conn)
		go t.handleConn(conn, false)
	}
}

func (t *TCPTransport) handleConn(conn net.Conn, outbound bool) {
	//defer conn.Close()
	var err error
	defer func() {
		fmt.Printf("dropping the peer connection %s \n", err)
		conn.Close()
	}()

	peer := NewTCPPeer(conn, true)

	if err = t.HandshakeFunc(peer); err != nil {
		// we need to drop the connection if there is an error in connection
		//conn.Close()
		//fmt.Printf("TCP handshake error: %s\n", err.Error())
		return
	}

	if t.OnPeer != nil {
		if err = t.OnPeer(peer); err != nil {
			return
		}
	}

	// Read loop
	//buf := make([]byte, 2000)
	rpc := RPC{}
	for {
		//n, err := conn.Read(buf)
		//if err != nil {
		//	fmt.Printf("TCP read error: %s\n", err.Error())
		//}

		err = t.Decoder.Decode(conn, &rpc)
		if err != nil {
			fmt.Println(reflect.TypeOf(err))

			if errors.Is(err, net.ErrClosed) || errors.Is(err, net.ErrWriteToConnected) {
				return
			}
			fmt.Printf("TCP decode error: %s\n", err.Error()) // we are keep looping if there is an error
			continue
		}

		rpc.From = conn.RemoteAddr()
		fmt.Printf("TCP receive from %+v \n", rpc.From.String())
		fmt.Printf("RPC: %+v\n", string(rpc.Payload))
		//fmt.Printf("message %+v\n", buf[:n])

		t.rpcchan <- rpc // we are inserting the message in the channel which the server is reading
	}

}

// Dial implements the Transport interface.
// it checks whether the node is reachable or not
func (t *TCPTransport) Dial(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	go t.handleConn(conn, true)

	return nil
}

// Close implements the Transport interface.
func (t *TCPTransport) Close() error {
	return t.listener.Close()
}
