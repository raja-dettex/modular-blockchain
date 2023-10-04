package core

import (
	"fmt"
	"testing"
	"time"

	"github.com/raja-dettex/modular-blockchain/types"
)

func RandomBlock(height int32) *Block {
	header := &Header{
		Version:   1,
		PrevBlock: types.RandomHash(),
		Timestamp: time.Now().UnixNano(),
		Height:    height,
	}
	tx := Transaction{
		Data: []byte("hello"),
	}
	block := NewBlock(header, []Transaction{tx})
	return block
}

func TestHashBlock(t *testing.T) {
	block := RandomBlock(0)
	fmt.Println(block.Hash(BlockHasher{}))
}
