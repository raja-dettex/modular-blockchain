package network

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"

	"github.com/raja-dettex/modular-blockchain/core"
)

type MessageType byte

const (
	MessageTypeTX    MessageType = 0x1
	MessageTypeBlock MessageType = 0x2
)

type RPC struct {
	From    NetAdddr
	Payload io.Reader
}

type Message struct {
	Header MessageType
	Data   []byte
}

func NewMessage(t MessageType, data []byte) Message {
	return Message{
		Header: t,
		Data:   data,
	}
}

func (m Message) Bytes() ([]byte, error) {
	buff := &bytes.Buffer{}
	if err := gob.NewEncoder(buff).Encode(m); err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}

type DecodedMessage struct {
	From NetAdddr
	Data any
}

type RPCDecodeFunc func(RPC) (*DecodedMessage, error)

func DefaultRPCDecodeFunc(rpc RPC) (*DecodedMessage, error) {
	msg := Message{}
	if err := gob.NewDecoder(rpc.Payload).Decode(&msg); err != nil {
		return nil, err
	}
	switch msg.Header {
	case MessageTypeTX:
		tx := new(core.Transaction)
		if err := tx.Decode(core.NewGobTxDecoder(bytes.NewReader(msg.Data))); err != nil {
			return nil, err
		}
		return &DecodedMessage{
			From: rpc.From,
			Data: tx,
		}, nil
	case MessageTypeBlock:
		block := new(core.Block)
		if err := block.Decoder(core.NewGobBlockDecoder(bytes.NewReader(msg.Data))); err != nil {
			return nil, err
		}
		return &DecodedMessage{
			From: rpc.From,
			Data: block,
		}, nil

	default:
		return nil, fmt.Errorf("handler other than transaction not registered")
	}
}

type RPCProcessor interface {
	ProcessMessage(*DecodedMessage) error
}
