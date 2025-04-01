package p2p

import "net"

// Peer is an interface that represents the remote node.
type Peer interface {
	//Conn() net.Conn
	Send([]byte) error
	//RemoteAddr() net.Addr
	//Close() error
	net.Conn // we can just embbed net.Connn and all the functions will work
}

// Transport is anything that handles the communication between the nodes in the network.
// This can be of the form (TCP.md, UDP, websockets, ...)
type Transport interface {
	ListenAndAccept() error
	Consume() <-chan RPC
	Close() error
	Dial(addr string) error
	ListenAddress() string
}
