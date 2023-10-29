package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"time"

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
	Transactions []*Transaction
	Validator    crypto.PublicKey
	Signature    *crypto.Signature

	// cache the hash of the header
	hash types.Hash
}

func (h *Header) Bytes() []byte {
	buff := &bytes.Buffer{}
	enc := gob.NewEncoder(buff)
	enc.Encode(h)
	return buff.Bytes()
}

func NewBlock(h *Header, txx []*Transaction) (*Block, error) {
	return &Block{
		Header:       h,
		Transactions: txx,
	}, nil
}

func NewBlockFromPrevHeader(prevHeader *Header, txx []*Transaction) (*Block, error) {
	dataHash, err := CalculateDataHash(txx)
	if err != nil {
		return nil, err
	}
	header := &Header{
		Version:   1,
		DataHash:  dataHash,
		PrevBlock: BlockHasher{}.Hash(prevHeader),
		Timestamp: time.Now().UnixNano(),
		Height:    prevHeader.Height + 1,
	}
	return NewBlock(header, txx)
}
func (b *Block) AddTransacation(tx *Transaction) {
	b.Transactions = append(b.Transactions, tx)
}

func (b *Block) Sign(privKey crypto.PrivateKey) error {
	pubKey := privKey.GeneratePublicKey()
	sig, err := privKey.Sign(b.Header.Bytes())
	if err != nil {
		return err
	}
	b.Validator = pubKey
	b.Signature = sig
	return nil
}
func (b *Block) Verify() error {
	if b.Signature == nil {
		return fmt.Errorf("signature is empty")
	}
	if !b.Signature.Verify(b.Validator, b.Header.Bytes()) {
		return fmt.Errorf("the block has invalid signature")
	}
	for _, tx := range b.Transactions {
		if err := tx.Verify(); err != nil {
			return err
		}
	}
	txxHash, err := CalculateDataHash(b.Transactions)
	if err != nil {
		return err
	}
	if txxHash != b.Header.DataHash {
		return fmt.Errorf("data hash %v is invalid and tx Hash %v", b.Header.DataHash, txxHash)
	}
	return nil
}
func (b *Block) Hash(hasher Hasher[*Header]) types.Hash {
	if b.hash.IsZero() {
		b.hash = hasher.Hash(b.Header)
	}
	return b.hash

}

func (b *Block) Encoder(enc Encoder[*Block]) error {
	return enc.Encode(b)
}
func (b *Block) Decoder(dec Decoder[*Block]) error {
	return dec.Decode(b)
}

func CalculateDataHash(txx []*Transaction) (types.Hash, error) {
	buff := &bytes.Buffer{}
	for _, tx := range txx {
		if err := tx.Encode(NewGobTxEncoder(buff)); err != nil {
			return types.Hash{}, err
		}
	}
	return sha256.Sum256(buff.Bytes()), nil
}
