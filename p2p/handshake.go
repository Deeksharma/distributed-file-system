package p2p

import "errors"

// ErrInvalidHandshake is returned if the handshake between the local and remote node could not be established.
var ErrInvalidHandshake = errors.New("invalid handshake")

// HandshakeFunc is the handshaker function between two remote machines while connection establishment
type HandshakeFunc func(Peer) error

func NOPHandshakeFunc(Peer) error {
	return nil
}
