package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"

	"github.com/raja-dettex/modular-blockchain/types"
)

type Hasher[T any] interface {
	Hash(T) types.Hash
}

type BlockHasher struct {
}

func (bh BlockHasher) Hash(b *Block) types.Hash {
	buff := &bytes.Buffer{}
	encoder := gob.NewEncoder(buff)
	if err := encoder.Encode(b.Header); err != nil {
		panic(err)
	}
	h := sha256.Sum256(buff.Bytes())
	return types.Hash(h)
}
