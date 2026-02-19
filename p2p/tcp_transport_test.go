package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPTransport(t *testing.T){
	listenAdr := ":3000"
	
	tc_transport := NewTcpTransport(listenAdr)

	assert.Equal(t, tc_transport.listen_address, listenAdr)

	assert.Nil(t, tc_transport.ListenAndAccept())
}