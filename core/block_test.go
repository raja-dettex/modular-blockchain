package core

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/raja-dettex/modular-blockchain/crypto"
	"github.com/raja-dettex/modular-blockchain/types"
	"github.com/stretchr/testify/assert"
)

func RandomBlock(height int32, t *testing.T, hash types.Hash) *Block {
	privKey := crypto.GeneratePrivateKey()
	header := &Header{
		Version:   1,
		PrevBlock: hash,
		Timestamp: time.Now().UnixNano(),
		Height:    height,
	}
	tx := generateRandomTransactionWithSignature(t)
	block, err := NewBlock(header, []*Transaction{tx})
	assert.Nil(t, err)
	blockHash, err := CalculateDataHash(block.Transactions)
	assert.Nil(t, err)
	block.Header.DataHash = blockHash
	block.Sign(privKey)
	errp := block.Verify()
	assert.Nil(t, errp)
	assert.NotNil(t, block)
	return block
}
func RandomBlockWithtxx(height int32, t *testing.T, hash types.Hash, txx []*Transaction) *Block {
	privKey := crypto.GeneratePrivateKey()
	header := &Header{
		Version:   1,
		PrevBlock: hash,
		Timestamp: time.Now().UnixNano(),
		Height:    height,
	}
	block, err := NewBlock(header, txx)
	assert.Nil(t, err)
	blockHash, err := CalculateDataHash(block.Transactions)
	assert.Nil(t, err)
	block.Header.DataHash = blockHash
	block.Sign(privKey)
	errp := block.Verify()
	assert.Nil(t, errp)
	assert.NotNil(t, block)
	return block
}

func TestBlockEncodeDecode(t *testing.T) {
	block := RandomBlock(0, t, types.Hash{})
	buff := &bytes.Buffer{}
	assert.Nil(t, block.Encoder(NewGobBlockEncoder(buff)))
	bDecoded := new(Block)
	assert.Nil(t, bDecoded.Decoder(NewGobBlockDecoder(buff)))
	assert.Equal(t, bDecoded, block)
}

//	func RandomBlockWithSignature(height int32, t *testing.T, hash types.Hash) *Block {
//		b := RandomBlock(height, t, hash)
//		privKey := crypto.GeneratePrivateKey()
//		assert.Nil(t, b.Sign(privKey))
//		return b
//	}
func randomTransactionWithSignature(t *testing.T, tx Transaction) *Transaction {
	privKey := crypto.GeneratePrivateKey()
	err := tx.Sign(privKey)
	assert.Nil(t, err)
	return &tx

}

func TestHashBlock(t *testing.T) {
	block := RandomBlock(0, t, types.Hash{})
	fmt.Println(block.Hash(BlockHasher{}))
}

func TestVerifyBlock(t *testing.T) {
	block := RandomBlock(0, t, types.Hash{})
	tx := &Transaction{Data: []byte("hello")}
	transaction := randomTransactionWithSignature(t, *tx)
	block.AddTransacation(transaction)
	//assert.Nil(t, block.Verify())
	assert.NotNil(t, block.Signature)
}
