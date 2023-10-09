package core

import (
	"crypto/sha256"

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
	return sha256.Sum256(tx.Data)
}
