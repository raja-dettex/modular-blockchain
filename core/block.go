package core

import (
	"io"

	"github.com/raja-dettex/modular-blockchain/crypto"
	"github.com/raja-dettex/modular-blockchain/types"
)

type Header struct {
	Version   uint32
	DataHash  types.Hash
	PrevBlock types.Hash
	Timestamp int64
	Height    int32
}

type Block struct {
	Header       *Header
	Transactions []Transaction
	Validator    crypto.PublicKey
	Signature    crypto.Signature

	// cache the hash of the header
	hash types.Hash
}

func NewBlock(h *Header, txx []Transaction) *Block {
	return &Block{
		Header:       h,
		Transactions: txx,
	}
}

func (b *Block) Hash(hasher Hasher[*Block]) types.Hash {
	if b.hash.IsZero() {
		b.hash = hasher.Hash(b)
	}
	return b.hash

}

func (b *Block) Encoder(w io.Writer, enc Encoder[*Block]) error {
	return enc.Encode(w, b)
}
func (b *Block) Decoder(r io.Reader, dec Decoder[*Block]) error {
	return dec.Decode(r, b)
}
