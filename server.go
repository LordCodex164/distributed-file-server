package main

import (
	"fmt"

	"github.com/lordcodex164/fstore/p2p"
)

type FileServerOpts struct {
	ListenAddr     string
	StorageRoot    string
	Transport      p2p.Transport
	BootStrapNodes []string
}

type FileServer struct {
	FileServerOpts
	store  *Store
	quitCh chan struct{}
}

func NewServer(opts FileServerOpts) *FileServer {
	storeOpts := StoreOpts{
		PathTransformFunc: DefaultPathTransformFunc,
		Root:              "src",
	}
	return &FileServer{
		FileServerOpts: opts,
		store:          NewStore(storeOpts),
		quitCh:         make(chan struct{}),
	}
}

func (s *FileServer) Quit() {
	close(s.quitCh)
}

func (s *FileServer) bootStrapNetwork() {
	for _, addr := range s.BootStrapNodes {
		if len(addr) == 0 {
			continue
		}
		go func(addr string) {
			fmt.Println("attempting to connect to:", addr)
			if err := s.Transport.Dial(addr); err != nil {
				fmt.Println("panic error", err)
				return
			}
			fmt.Println("listening on the addr", addr)
		}(addr)

	}
}

func (s *FileServer) Start() error {
	if err := s.Transport.ListenAndAccept(); err != nil {
		return err
	}
	s.bootStrapNetwork()
	s.loop()
	return nil
}

func (s *FileServer) loop() {
	defer func() {
		fmt.Println("file server stopped due to error or user quit action")
		//close the transport
		s.Transport.Close()
	}()
	for {
		select {
		case msg := <-s.Transport.Consume():
			fmt.Println("msg", msg)
		case <-s.quitCh:
			return
		}
	}
}
