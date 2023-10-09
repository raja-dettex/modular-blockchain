package network

type NetAdddr string

type Transport interface {
	Addr() NetAdddr
	Connect(Transport) error
	Consume() <-chan RPC
	Broadcast([]byte) error
	SendMessage(NetAdddr, []byte) error
}
