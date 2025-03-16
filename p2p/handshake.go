package p2p

// HandshakeFunc is the handshaker function between two remote machines while connection establishment
type HandshakeFunc func(any) error

func NOPHandshakeFunc(any) error {
	return nil
}
