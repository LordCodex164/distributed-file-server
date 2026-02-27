package p2p

import (
	"errors"
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

type TCPTransportOpts struct {
	ShakeHands     HandshakeFunc
	Decoder        Decoder
	Listen_address string
	OnPeer         func(Peer) error
}

type TcpTransport struct {
	listener net.Listener
	rpcCh    chan RPC
	mu       sync.Mutex
	TCPTransportOpts
}

func NewTcpTransport(opts TCPTransportOpts) *TcpTransport {
	return &TcpTransport{
		TCPTransportOpts: opts,
		rpcCh:            make(chan RPC),
	}
}

func NewTcpPeer(conn net.Conn, outbound bool) *TcpPeer {
	return &TcpPeer{
		Conn:     conn,
		outbound: outbound,
		wg:       &sync.WaitGroup{},
	}
}

func (p *TcpPeer) CloseStream() {
	p.wg.Done()
}

func (p *TcpPeer) Send(b []byte) error {
	_, err := p.Conn.Write(b)
	return err
}

func (p *TCPTransportOpts) Addr() string {
	return p.Listen_address
}

func (t *TcpTransport) Consume() <-chan RPC {
	return t.rpcCh
}

func (t *TcpTransport) Dial(addr string) error {
	conn, err := net.Dial("tcp", addr)
	fmt.Println("conn", conn, addr)
	if err != nil {
		return err
	}
	go t.handleConn(conn, true)
	return nil
}

func (t *TcpTransport) ListenAndAccept() error {
	var err error

	t.listener, err = net.Listen("tcp", t.Listen_address)
	if err != nil {
		return err
	}
	go t.acceptHandShakeLoop()
	log.Printf("TCP transport listening on port: %s\n", t.Listen_address)
	return nil
}

func (t *TcpTransport) acceptHandShakeLoop() {

	for {
		conn, err := t.listener.Accept()
		if errors.Is(err, net.ErrClosed) {
			return
		}
		if err != nil {
			fmt.Println("tcp accept error:", err, conn)
		}
		//handle connections
		go t.handleConn(conn, false)
	}

}

type Temp struct{}

func (t *TcpTransport) handleConn(conn net.Conn, outbound bool) {

	var err error

	peer := NewTcpPeer(conn, outbound)

	if err = t.ShakeHands(peer); err != nil {
		conn.Close()
		fmt.Println("tcp handshake error:", err)
		return
	}
	fmt.Println("handling transport connection", conn)

	if t.OnPeer != nil {
		if err := t.OnPeer(peer); err != nil {
			fmt.Println("connection close", conn)
			peer.Close()
			return
		}
	}

	msg := RPC{}

	//the encoding of the rpc between two peers
	for {
		err = t.Decoder.Decode(conn, &msg)
		if err != nil {
			fmt.Println("err", err)
			return
		}
		msg.From = conn.RemoteAddr().String()
		//write to the channel
		t.rpcCh <- msg
	}
}

func (t *TcpTransport) Close() error {
	return t.listener.Close()
}
