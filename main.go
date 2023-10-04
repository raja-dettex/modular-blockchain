package main

import (
	"time"

	"github.com/raja-dettex/modular-blockchain/network"
)

func main() {
	traLocal := network.NewLocalTransport("LOCAL")
	traRemote := network.NewLocalTransport("REMOTE")
	traLocal.Connect(traRemote)
	traRemote.Connect(traLocal)
	go func() {
		for {
			traRemote.SendMessage(traLocal.Addr(), []byte("hello world"))
			time.Sleep(time.Second * 1)
		}
	}()
	opts := network.ServerOpts{
		Transports: []network.Transport{traLocal},
	}
	server := network.NewServer(opts)
	server.Start()
}
