package network

type NetAdddr string

type RPC struct {
	From    NetAdddr
	Payload []byte
}

type Transport interface {
	Addr() NetAdddr
	Connect(Transport) error
	Consume() <-chan RPC
	SendMessage(NetAdddr, []byte) error
}
