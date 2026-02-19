package main

import (
	"fmt"
	"log"

	"github.com/lordcodex164/fstore/p2p"
)


func main(){
	fmt.Println("starting the transport protocol")
	tcp_t := p2p.NewTcpTransport(":3000")
	if err := tcp_t.ListenAndAccept(); err != nil  {
		log.Fatal(err)
	}
	select{}
}