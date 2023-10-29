package network

import "net"

type NetAdddr string

type Transport interface {
	Addr() net.Addr
	Connect(Transport) error
	Consume() <-chan RPC
	Broadcast([]byte) error
	SendMessage(net.Addr, []byte) error
}
