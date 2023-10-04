package network

import (
	"fmt"
	"time"
)

type ServerOpts struct {
	Transports []Transport
}
type Server struct {
	Opts   ServerOpts
	rpcCh  chan RPC
	quitCh chan struct{}
}

func NewServer(opts ServerOpts) *Server {
	return &Server{
		Opts:   opts,
		rpcCh:  make(chan RPC),
		quitCh: make(chan struct{}),
	}
}

func (s *Server) Start() {
	s.InitTransport()
	ticker := time.NewTicker(5 * time.Second)
free:
	for {
		select {
		case rpc := <-s.rpcCh:
			fmt.Printf("message %v\n", rpc)
		case <-s.quitCh:
			break free
		case <-ticker.C:
			fmt.Println("do stuff")
		}
	}
	fmt.Println("server is shutting down")
}

func (s *Server) InitTransport() {
	for _, tra := range s.Opts.Transports {
		go func(tr Transport) {
			for rpc := range tr.Consume() {
				s.rpcCh <- rpc
			}
		}(tra)
	}
}
