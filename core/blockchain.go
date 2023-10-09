package core

import (
	"fmt"
	"sync"

	"github.com/go-kit/log"
)

type Blockchain struct {
	Logger    log.Logger
	Store     Storage
	lock      sync.RWMutex
	Headers   []*Header
	Validator Validator
}

func NewBlockChain(logger log.Logger, genesis *Block) (*Blockchain, error) {
	bc := &Blockchain{
		Logger:  logger,
		Store:   MemoryStorage{},
		Headers: []*Header{},
	}
	bc.Validator = NewBlockValidator(bc)
	return bc, bc.addBlockWithoutValidation(genesis)
}

func (bc *Blockchain) GetHeader(height int32) (*Header, error) {
	if height > int32(bc.Height()) {
		return nil, fmt.Errorf("height has no block %v", height)
	}
	bc.lock.Lock()
	defer bc.lock.Unlock()
	return bc.Headers[height], nil
}

func (bc *Blockchain) Height() uint32 {
	bc.lock.Lock()
	defer bc.lock.Unlock()
	return uint32(len(bc.Headers) - 1)
}

func (bc *Blockchain) HasBlock(height uint32) bool {
	return height <= bc.Height()
}

func (bc *Blockchain) SetValidator(v Validator) {
	bc.Validator = v
}

func (bc *Blockchain) AddBlock(b *Block) error {
	// validate
	err := bc.Validator.Validate(b)
	if err != nil {
		// log the error and return
		bc.Logger.Log("error", err)
		return err
	}
	for _, tx := range b.Transactions {
		bc.Logger.Log("msg", "exeuting code on vm", "tx_hash", tx.Hash(TransactionHashesr{}))
		vm := NewVM(tx.Data)
		if err := vm.Run(); err != nil {
			return err
		}
		bc.Logger.Log("vm result", vm.firstNode())
	}
	return bc.addBlockWithoutValidation(b)
}

func (bc *Blockchain) addBlockWithoutValidation(b *Block) error {
	bc.lock.Lock()
	bc.Headers = append(bc.Headers, b.Header)
	bc.lock.Unlock()
	bc.Logger.Log("msg", "new block", "hash", b.Hash(BlockHasher{}), "height ", b.Header.Height)
	return bc.Store.Put(b)
}
