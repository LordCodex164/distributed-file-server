package main

import (
	"fmt"
	"log"

	"github.com/lordcodex164/fstore/p2p"
)

func On_Peer(peer p2p.Peer) error {
	peer.Close()
	return nil
}

func main(){
	fmt.Println("starting the transport protocol")
	tcp_options := p2p.TCPTransportOpts{
		Listen_address: ":4000",
		ShakeHands: p2p.NOPHandshakeFunc,
		Decoder: p2p.DefaultDecoder{},
		OnPeer: On_Peer,
	}
	tcp_t := p2p.NewTcpTransport(tcp_options, )
	if err := tcp_t.ListenAndAccept(); err != nil  {
		log.Fatal(err)
	}
	go func(){
		for {
			rpc_ch := tcp_t.Consume()
			msg := <-rpc_ch
			fmt.Println("here is a message from a peer:", string(msg.Payload))
		}
	}()
	select{}
}