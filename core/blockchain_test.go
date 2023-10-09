package core

import (
	"fmt"
	"os"
	"testing"

	"github.com/go-kit/log"
	"github.com/raja-dettex/modular-blockchain/types"
	"github.com/stretchr/testify/assert"
)

func TestAddBlockWitVM(t *testing.T) {
	bc, err := NewBlockChain(log.NewLogfmtLogger(os.Stderr), RandomBlock(0, t, types.Hash{}))
	assert.Nil(t, err)
	txx := []*Transaction{}
	for i := 0; i < 10; i++ {
		tx := generateRandomTransactionWithSignatureforVM(t)
		txx = append(txx, tx)
	}
	hash := getPrevBlockHash(t, bc, int32(1))
	block := RandomBlockWithtxx(int32(1), t, hash, txx)
	err = bc.AddBlock(block)
	assert.Nil(t, err)

}

func TestBlockchainValidator(t *testing.T) {
	bc, err := NewBlockChain(log.NewLogfmtLogger(os.Stderr), RandomBlock(0, t, types.Hash{}))
	assert.Nil(t, err)
	assert.NotNil(t, bc.Validator)
}

func TestHasBlock(t *testing.T) {
	bc := newBlockchainwithgenesisBlock(t)
	assert.True(t, bc.HasBlock(0))
}

func TestAddSingleBlock(t *testing.T) {
	bc := newBlockchainwithgenesisBlock(t)
	hash := getPrevBlockHash(t, bc, int32(1))
	block := RandomBlock(int32(1), t, hash)
	err := bc.AddBlock(block)
	assert.Nil(t, err)
}

func TestAddBlock(t *testing.T) {
	bc := newBlockchainwithgenesisBlock(t)

	for i := 0; i < 1000; i++ {
		hash := getPrevBlockHash(t, bc, int32(i+1))

		err := bc.AddBlock(RandomBlock(int32(i+1), t, hash))
		assert.Nil(t, err)
	}
	// assert.NotNil(t, bc.AddBlock(RandomBlockWithSignature(90, t)))
	assert.Equal(t, bc.Height(), uint32(1000))
	assert.Equal(t, len(bc.Headers), 1001)
}

func TestGetHeader(t *testing.T) {
	bc := newBlockchainwithgenesisBlock(t)
	for i := 0; i < 1000; i++ {
		block := RandomBlock(int32(i+1), t, getPrevBlockHash(t, bc, int32(i+1)))
		err := bc.AddBlock(block)
		assert.Nil(t, err)
		header, err := bc.GetHeader(int32(i + 1))
		assert.Nil(t, err)
		assert.Equal(t, header, block.Header)

	}
}

func TestValidatorWithGarbageHeight(t *testing.T) {
	bc := newBlockchainwithgenesisBlock(t)
	err := bc.AddBlock(RandomBlock(1, t, getPrevBlockHash(t, bc, 1)))
	assert.Nil(t, err)
	fmt.Println(bc.Height())
	err = bc.AddBlock(RandomBlock(3, t, types.Hash{}))
	assert.NotNil(t, err)
}

func newBlockchainwithgenesisBlock(t *testing.T) *Blockchain {
	bc, err := NewBlockChain(log.NewLogfmtLogger(os.Stderr), RandomBlock(0, t, types.Hash{}))
	assert.Nil(t, err)
	return bc
}

func getPrevBlockHash(t *testing.T, bc *Blockchain, height int32) types.Hash {
	prevheader, err := bc.GetHeader(height - 1)
	assert.Nil(t, err)
	return BlockHasher{}.Hash(prevheader)

}
