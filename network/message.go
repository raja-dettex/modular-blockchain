package network

import "github.com/raja-dettex/modular-blockchain/core"

type StatusMessage struct {
	ID            string
	CurrentHeight uint32
}

type GetStatusMessage struct {
}

type GetBlocksMessage struct {
	From uint32
	To   uint32
}

type BlocksMessage struct {
	Blocks []*core.Block
}
