package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"

	"github.com/raja-dettex/modular-blockchain/types"
)

type Hasher[T any] interface {
	Hash(T) types.Hash
}

type BlockHasher struct {
}

func (bh BlockHasher) Hash(h *Header) types.Hash {
	hHash := sha256.Sum256(h.Bytes())
	return types.Hash(hHash)
}

type TransactionHashesr struct{}

func (tHahser TransactionHashesr) Hash(tx *Transaction) types.Hash {
	buff := &bytes.Buffer{}
	if err := gob.NewEncoder(buff).Encode(tx); err != nil {
		fmt.Println(err)
	}
	return sha256.Sum256(buff.Bytes())
}
