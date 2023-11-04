package core

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"testing"

	"github.com/raja-dettex/modular-blockchain/crypto"
	"github.com/stretchr/testify/assert"
)

func TestTransaction(t *testing.T) {
	fromPrivKey := crypto.GeneratePrivateKey()
	toPrivKey := crypto.GeneratePrivateKey()
	tx := &Transaction{
		Data:  []byte("foo"),
		To:    toPrivKey.GeneratePublicKey(),
		Value: 666,
	}
	err := tx.Sign(fromPrivKey)
	assert.Nil(t, err)
	hash := tx.Hash(TransactionHashesr{})
	fmt.Println(tx)
	fmt.Println(hash)

}

func TestNFTTransaction(t *testing.T) {
	collectiontx := CollectionTx{
		Fee:      200,
		MetaData: []byte("collection nft"),
	}

	tx := &Transaction{
		TxInnner: collectiontx,
	}
	privKey := crypto.GeneratePrivateKey()
	tx.Sign(privKey)
	buff := new(bytes.Buffer)

	err := gob.NewEncoder(buff).Encode(tx)
	assert.Nil(t, err)
	txDecoded := &Transaction{}
	err = gob.NewDecoder(buff).Decode(txDecoded)
	assert.Nil(t, err)
	assert.Equal(t, tx, txDecoded)
}

func TestSignTransaction(t *testing.T) {
	privKey := crypto.GeneratePrivateKey()
	tx := Transaction{
		Data: []byte("hello"),
	}
	assert.Nil(t, tx.Sign(privKey))
	assert.NotNil(t, tx.Signaure)
}

func TestVerifyTransaction(t *testing.T) {
	privKey := crypto.GeneratePrivateKey()
	tx := Transaction{
		Data: []byte("hello"),
	}
	tx.Sign(privKey)
	assert.Nil(t, tx.Verify())
}

func TestTransactionEncoding(t *testing.T) {
	tx := generateRandomTransactionWithSignature(t)
	buff := &bytes.Buffer{}
	encoder := NewGobTxEncoder(buff)
	decoder := NewGobTxDecoder(buff)
	err := tx.Encode(encoder)
	assert.Nil(t, err)
	txDecoded := new(Transaction)
	txDecoded.Decode(decoder)
	assert.Equal(t, tx, txDecoded)
}

func generateRandomTransactionWithSignature(t *testing.T) *Transaction {
	privKey := crypto.GeneratePrivateKey()
	tx := NewTransaction([]byte("Foo"))
	tx.Sign(privKey)
	return tx
}
func generateRandomTransactionWithSignatureforVM(t *testing.T) *Transaction {
	privKey := crypto.GeneratePrivateKey()
	tx := NewTransaction(Contract())
	tx.Sign(privKey)
	return tx
}
