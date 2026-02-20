package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPTransport(t *testing.T){
	
	tc := TCPTransportOpts{
		Listen_address: ":3000",
		Decoder: GOBDecoder{},
	}
	
	listenAdr := ":3000"
	
	tc_transport := NewTcpTransport(tc)

	assert.Equal(t, tc_transport.Listen_address, listenAdr)

	assert.Nil(t, tc_transport.ListenAndAccept())
}