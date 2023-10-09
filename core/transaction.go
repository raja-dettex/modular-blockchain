package core

import (
	"fmt"

	"github.com/raja-dettex/modular-blockchain/crypto"
	"github.com/raja-dettex/modular-blockchain/types"
)

type Transaction struct {
	Data     []byte
	From     crypto.PublicKey
	Signaure *crypto.Signature

	// cache the hash of transaction
	hash types.Hash

	// first seen timestamp
	firstSeen int64
}

func NewTransaction(data []byte) *Transaction {
	return &Transaction{
		Data: data,
	}
}

func (tx *Transaction) Sign(privKey crypto.PrivateKey) error {
	sig, err := privKey.Sign(tx.Data)
	pubKey := privKey.GeneratePublicKey()
	if err != nil {
		return err
	}
	tx.From = pubKey
	tx.Signaure = sig
	return nil
}

func (tx *Transaction) Hash(hasher Hasher[*Transaction]) types.Hash {
	if tx.hash.IsZero() {
		tx.hash = hasher.Hash(tx)
	}
	return tx.hash
}

func (tx *Transaction) Verify() error {
	if tx.Signaure == nil {
		return fmt.Errorf("transaction has no signature")
	}
	if !tx.Signaure.Verify(tx.From, tx.Data) {
		return fmt.Errorf("invalid transaction")
	}
	return nil
}

func (tx *Transaction) SetFirstSeen(t int64) {
	tx.firstSeen = t
}

func (tx *Transaction) FirstSeen() int64 {
	return tx.firstSeen
}

func (tx *Transaction) Encode(enc Encoder[*Transaction]) error {
	return enc.Encode(tx)
}
func (tx *Transaction) Decode(dec Decoder[*Transaction]) error {
	return dec.Decode(tx)
}
