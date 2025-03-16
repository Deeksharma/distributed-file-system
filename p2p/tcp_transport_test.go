package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTCPTransport(t *testing.T) {
	listenAddr := ":4000"
	tcpTransport := NewTCPTransport(listenAddr)

	assert.Equal(t, tcpTransport.listenAddress, listenAddr)

	// Server
	// tcpTransport.Start()
	assert.Nil(t, tcpTransport.ListenAndAccept())

}
