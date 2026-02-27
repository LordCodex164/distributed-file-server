package main

import (
	//"fmt"
	"log"
	"github.com/lordcodex164/fstore/p2p"
)

func On_Peer(peer p2p.Peer) error {
	//peer.Close()
	return nil
}

func makeServer(listenAddr string, nodes ...string) *FileServer {
	tcp_options := p2p.TCPTransportOpts{
		Listen_address: listenAddr,
		ShakeHands: p2p.NOPHandshakeFunc,
		Decoder: p2p.DefaultDecoder{},
		OnPeer: On_Peer,
	}
	tcp_t := p2p.NewTcpTransport(tcp_options)
	fileServerOpts := FileServerOpts{
		StorageRoot: listenAddr + "_network",
		Transport: tcp_t,
		BootStrapNodes: nodes,
	}
	s := NewServer(fileServerOpts)
	// go func ()  {
	// 	time.Sleep(3 * time.Second)
	// 	s.Quit()	
	// }()
	return s
}

func main(){
	s1 := makeServer(":3000", "") //first server
	s2 := makeServer(":4000", ":3000")
	go func ()  {
		log.Fatal(s1.Start())
	}()
	go func ()  {
		log.Println("connecting to server 2")
		s2.Start()
	}()
	select{}
}