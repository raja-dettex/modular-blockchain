package core

import (
	"fmt"
	"sync"

	"github.com/go-kit/log"
	"github.com/raja-dettex/modular-blockchain/types"
)

type Blockchain struct {
	Logger        log.Logger
	Store         Storage
	lock          sync.RWMutex
	Headers       []*Header
	Blocks        []*Block
	BlockStore    map[types.Hash]*Block
	TxStore       map[types.Hash]*Transaction
	Validator     Validator
	ContractState *State
}

func NewBlockChain(logger log.Logger, genesis *Block) (*Blockchain, error) {
	bc := &Blockchain{
		Logger:        logger,
		Store:         MemoryStorage{},
		Headers:       []*Header{},
		BlockStore:    make(map[types.Hash]*Block),
		TxStore:       make(map[types.Hash]*Transaction),
		ContractState: NewState(),
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

func (bc *Blockchain) GetBlock(height uint32) (*Block, error) {
	// bc.lock.RLock()
	// defer bc.lock.RUnlock()
	if height > bc.Height() {
		return nil, fmt.Errorf("block height %v too high", height)
	}
	return bc.Blocks[height], nil

}

func (bc *Blockchain) GetTxByHash(hash types.Hash) (*Transaction, error) {
	bc.lock.Lock()
	defer bc.lock.Unlock()
	tx, ok := bc.TxStore[hash]
	if !ok {
		return nil, fmt.Errorf("Transaction with hash %v not found", hash)
	}
	return tx, nil
}

func (bc *Blockchain) GetBlockByHash(blockHash types.Hash) (*Block, error) {
	bc.lock.Lock()
	defer bc.lock.Unlock()
	block, ok := bc.BlockStore[blockHash]
	if !ok {
		return nil, fmt.Errorf("block with hash %v not found", blockHash)
	}
	return block, nil
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
		//bc.Logger.Log("error", err)
		return err
	}
	for _, tx := range b.Transactions {
		bc.Logger.Log("msg", "exeuting code on vm", "tx_hash", tx.Hash(TransactionHashesr{}))
		vm := NewVM(tx.Data, bc.ContractState)
		if err := vm.Run(); err != nil {
			return err
		}
		val, _ := bc.ContractState.Get([]byte("da"))

		fmt.Printf("contract state -> %v\n", val)
	}
	return bc.addBlockWithoutValidation(b)
}

func (bc *Blockchain) addBlockWithoutValidation(b *Block) error {
	bc.lock.Lock()
	bc.Headers = append(bc.Headers, b.Header)
	bc.Blocks = append(bc.Blocks, b)
	blockHash := b.Hash(BlockHasher{})
	bc.BlockStore[blockHash] = b
	for _, tx := range b.Transactions {
		bc.TxStore[tx.Hash(TransactionHashesr{})] = tx
	}
	bc.lock.Unlock()
	bc.Logger.Log("msg", "new block", "hash", b.Hash(BlockHasher{}), "height ", b.Header.Height, "transactions", len(b.Transactions))
	return bc.Store.Put(b)
}
