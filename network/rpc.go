package network

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io"
	"net"

	"github.com/raja-dettex/modular-blockchain/core"
)

type MessageType byte

const (
	MessageTypeTX               MessageType = 0x1
	MessageTypeBlock            MessageType = 0x2
	MessageTypeGetBlocks        MessageType = 0x3
	MessageTypeStatusMessage    MessageType = 0x4
	MessageTypeGetStatusMessage MessageType = 0x5
	MessageTypeSyncBlocks       MessageType = 0x06
)

type RPC struct {
	From    net.Addr
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
	From net.Addr
	Data any
}

type RPCDecodeFunc func(RPC) (*DecodedMessage, error)

func DefaultRPCDecodeFunc(rpc RPC) (*DecodedMessage, error) {
	msg := Message{}
	if err := gob.NewDecoder(rpc.Payload).Decode(&msg); err != nil {
		return nil, fmt.Errorf("failed to decode message %v", err)
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
	case MessageTypeStatusMessage:
		statusMessage := new(StatusMessage)
		if err := gob.NewDecoder(bytes.NewReader(msg.Data)).Decode(statusMessage); err != nil {
			return nil, err
		}
		return &DecodedMessage{
			From: rpc.From,
			Data: statusMessage,
		}, nil
	case MessageTypeGetStatusMessage:
		return &DecodedMessage{
			From: rpc.From,
			Data: &GetStatusMessage{},
		}, nil
	case MessageTypeGetBlocks:
		getBlocksMessage := new(GetBlocksMessage)
		if err := gob.NewDecoder(bytes.NewReader(msg.Data)).Decode(getBlocksMessage); err != nil {
			return nil, err
		}
		return &DecodedMessage{
			From: rpc.From,
			Data: getBlocksMessage,
		}, nil
	case MessageTypeSyncBlocks:
		blocksMessge := new(BlocksMessage)
		if err := gob.NewDecoder(bytes.NewReader(msg.Data)).Decode(blocksMessge); err != nil {
			return nil, err
		}
		return &DecodedMessage{
			From: rpc.From,
			Data: blocksMessge,
		}, nil
	default:
		return nil, fmt.Errorf("handler other than transaction not registered")
	}
}

type RPCProcessor interface {
	ProcessMessage(*DecodedMessage) error
}

func init() {
	gob.Register(elliptic.P256())
}
