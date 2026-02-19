package p2p

import (
	"fmt"
	"log"
	"net"
	"sync"
)

type TcpPeer struct {
	// The underlying connection of the peer. Which in this case
	// is a TCP connection.
	net.Conn
	// if we dial and retrieve a conn => outbound == true
	// if we accept and retrieve a conn => outbound == false
	outbound bool
	wg       *sync.WaitGroup
}

type TcpTransport struct {
	listen_address string
	listener net.Listener
	rpcCh    chan RPC
	mu       sync.Mutex
	shakeHands HandshakeFunc
}

func NewTcpTransport(listenAddr string) *TcpTransport {
	return &TcpTransport{
		listen_address: listenAddr,
	}
}

func NewTcpPeer(conn net.Conn, outbound bool) *TcpPeer {
	return &TcpPeer{
		Conn: conn,
		outbound: outbound,
		wg: &sync.WaitGroup{},
	}
}

func (p *TcpPeer) CloseStream() {
	p.wg.Done()
}

func (p *TcpPeer) Send(b []byte) error {
	_, err := p.Conn.Write(b)
	return err
}

func (t *TcpTransport) ListenAndAccept() error {
	var err error

	t.listener, err = net.Listen("tcp", t.listen_address)

	if err != nil {
		return err
	}
	fmt.Println("listener", t.listener)
	go t.acceptHandShakeLoop()
	log.Printf("TCP transport listening on port: %s\n", t.listen_address)
	return nil
}

func (t *TcpTransport) acceptHandShakeLoop() {
	conn, err := t.listener.Accept()
	if(err != nil){
		return
	}
	//handle 
	go t.handleConn(conn, false)
}

func (t *TcpTransport) handleConn(conn net.Conn, outbound bool) {
	peer := NewTcpPeer(conn, outbound)

	if err := t.shakeHands(peer); err != nil {
		return 
	}
	fmt.Println("handling transport connection", conn)
}
