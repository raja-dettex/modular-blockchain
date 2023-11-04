package core

import (
	"encoding/gob"
	"fmt"
	"math/rand"

	"github.com/raja-dettex/modular-blockchain/crypto"
	"github.com/raja-dettex/modular-blockchain/types"
)

type TxType byte

const (
	TxTypeCollection TxType = iota
	TxTypeMint
)

type CollectionTx struct {
	Fee      int64
	MetaData []byte
}

type MintTx struct {
	Fee             int64
	NFT             types.Hash
	Collection      types.Hash
	MetaData        []byte
	CollectionOwner crypto.PublicKey
	Signature       crypto.Signature
}

type Transaction struct {
	// only used for native nft logic on our chain
	TxInnner any

	// data or smart contract that will be executed by vm
	Data []byte

	To       crypto.PublicKey
	Value    uint64
	From     crypto.PublicKey
	Signaure *crypto.Signature
	Nonce    int64
	// cache the hash of transaction
	hash types.Hash

	// first seen timestamp
	firstSeen int64
}

func NewTransaction(data []byte) *Transaction {
	return &Transaction{
		Data:  data,
		Nonce: rand.Int63n(100000000000),
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

func init() {
	gob.Register(CollectionTx{})
	gob.Register(MintTx{})
}
