package core

import (
	"bytes"
	"testing"

	"github.com/raja-dettex/modular-blockchain/crypto"
	"github.com/stretchr/testify/assert"
)

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
