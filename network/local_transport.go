package network

import (
	"bytes"
	"fmt"
	"net"
	"sync"
)

type LocalTransport struct {
	addr      net.Addr
	consumeCh chan RPC
	peers     map[net.Addr]*LocalTransport
	lock      sync.RWMutex
}

func NewLocalTransport(addr net.Addr) *LocalTransport {
	return &LocalTransport{
		addr:      addr,
		consumeCh: make(chan RPC, 1024),
		peers:     make(map[net.Addr]*LocalTransport),
	}
}

func (lt *LocalTransport) Addr() net.Addr {
	return lt.addr
}

func (lt *LocalTransport) Consume() <-chan RPC {
	return lt.consumeCh
}

func (lt *LocalTransport) Connect(local Transport) error {
	lt.lock.Lock()
	defer lt.lock.Unlock()
	if lt.Addr() == local.Addr() {
		return fmt.Errorf("can not connect to its own, local peer %v, remote peer %v\n", lt.Addr(), local.Addr())
	}
	lt.peers[local.Addr()] = local.(*LocalTransport)
	return nil
}

func (lt *LocalTransport) SendMessage(addr net.Addr, payload []byte) error {
	lt.lock.Lock()
	defer lt.lock.Unlock()
	peer, ok := lt.peers[addr]
	if !ok {
		return fmt.Errorf("unable to send from %s to %s\n", lt.addr, addr)
	}
	peer.consumeCh <- RPC{
		From:    lt.Addr(),
		Payload: bytes.NewReader(payload),
	}
	return nil
}

func (lt *LocalTransport) Broadcast(payload []byte) error {
	for _, peer := range lt.peers {
		if err := lt.SendMessage(peer.Addr(), payload); err != nil {
			return err
		}
	}
	return nil
}
