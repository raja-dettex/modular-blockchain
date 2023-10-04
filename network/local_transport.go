package network

import (
	"fmt"
	"sync"
)

type LocalTransport struct {
	addr      NetAdddr
	consumeCh chan RPC
	peers     map[NetAdddr]*LocalTransport
	lock      sync.RWMutex
}

func NewLocalTransport(addr NetAdddr) *LocalTransport {
	return &LocalTransport{
		addr:      addr,
		consumeCh: make(chan RPC, 1024),
		peers:     make(map[NetAdddr]*LocalTransport),
	}
}

func (lt *LocalTransport) Addr() NetAdddr {
	return lt.addr
}

func (lt *LocalTransport) Consume() <-chan RPC {
	return lt.consumeCh
}

func (lt *LocalTransport) Connect(local Transport) error {
	lt.lock.Lock()
	defer lt.lock.Unlock()
	lt.peers[local.Addr()] = local.(*LocalTransport)
	return nil
}

func (lt *LocalTransport) SendMessage(addr NetAdddr, payload []byte) error {
	peer, ok := lt.peers[addr]
	if !ok {
		return fmt.Errorf("unable to send from %s to %s", lt.addr, addr)
	}
	peer.consumeCh <- RPC{
		From:    lt.Addr(),
		Payload: payload,
	}
	return nil
}
