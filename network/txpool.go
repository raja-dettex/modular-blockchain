package network

import (
	"sync"

	"github.com/raja-dettex/modular-blockchain/core"
	"github.com/raja-dettex/modular-blockchain/types"
)

type TxPool struct {
	all     *TxSortedMap
	pending *TxSortedMap

	maxLength int
}

func NewTxPool(maxLength int) *TxPool {
	return &TxPool{
		all:       NewTxSortedMap(),
		pending:   NewTxSortedMap(),
		maxLength: maxLength,
	}
}

func (pool *TxPool) Add(tx *core.Transaction) {
	if pool.all.Count() == pool.maxLength {
		first := pool.all.First()
		pool.all.Remove(first.Hash(core.TransactionHashesr{}))
	}
	if !pool.all.Contains(tx.Hash(core.TransactionHashesr{})) {
		pool.all.Add(tx)
		pool.pending.Add(tx)
	}
}

func (pool *TxPool) Contains(hash types.Hash) bool {
	return pool.all.Contains(hash)
}

func (pool *TxPool) Pending() []*core.Transaction {
	return pool.pending.txx.Data
}

func (pool *TxPool) ClearPending() {
	pool.pending.Clear()
}

func (pool *TxPool) PendingCount() int {
	return pool.pending.Count()
}

func (pool *TxPool) AllCount() int {
	return pool.all.Count()
}

type TxSortedMap struct {
	lock   sync.RWMutex
	lookup map[types.Hash]*core.Transaction
	txx    *types.List[*core.Transaction]
}

func NewTxSortedMap() *TxSortedMap {
	return &TxSortedMap{
		lookup: make(map[types.Hash]*core.Transaction),
		txx:    types.NewList[*core.Transaction](),
	}
}

func (txMap *TxSortedMap) First() *core.Transaction {
	txMap.lock.Lock()
	defer txMap.lock.Unlock()
	first := txMap.txx.Get(0)
	return txMap.lookup[first.Hash(core.TransactionHashesr{})]

}

func (txMap *TxSortedMap) Add(tx *core.Transaction) {
	hash := tx.Hash(core.TransactionHashesr{})
	txMap.lock.Lock()
	defer txMap.lock.Unlock()
	if _, ok := txMap.lookup[hash]; !ok {
		txMap.lookup[hash] = tx
		txMap.txx.Insert(tx)
	}
}

func (txMap *TxSortedMap) Remove(hash types.Hash) {
	txMap.lock.Lock()
	defer txMap.lock.Unlock()
	txMap.txx.Remove(txMap.lookup[hash])
	delete(txMap.lookup, hash)

}

func (txMap *TxSortedMap) Count() int {
	return len(txMap.lookup)
}

func (txMap *TxSortedMap) Contains(h types.Hash) bool {
	txMap.lock.Lock()
	defer txMap.lock.Unlock()
	_, ok := txMap.lookup[h]
	return ok
}
func (txMap *TxSortedMap) Clear() {
	txMap.lock.Lock()
	defer txMap.lock.Unlock()

	txMap.lookup = make(map[types.Hash]*core.Transaction)
	txMap.txx.Clear()
}
