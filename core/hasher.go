package core

import (
	"crypto/sha256"
	"encoding/binary"

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
	buff := make([]byte, 8)
	binary.LittleEndian.PutUint32(buff, uint32(tx.Nonce))
	data := append(buff, tx.Data...)
	return sha256.Sum256(data)
}
