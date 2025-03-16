package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTCPTransport(t *testing.T) {
	opts := TCPTransportOps{
		ListenAddr:    ":3000",
		HandshakeFunc: NOPHandshakeFunc,
		Decoder:       &DefaultDecoder{},
	}
	listenAddr := ":3000"
	tcpTransport := NewTCPTransport(opts)

	assert.Equal(t, tcpTransport.ListenAddr, listenAddr)

	// Server
	// tcpTransport.Start()
	assert.Nil(t, tcpTransport.ListenAndAccept())

}
