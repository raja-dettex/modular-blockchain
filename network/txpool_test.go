package network

import (
	"fmt"
	"testing"

	"github.com/raja-dettex/modular-blockchain/core"
	"github.com/stretchr/testify/assert"
)

func TestTxpool(t *testing.T) {
	p := NewTxPool(100)
	assert.Equal(t, 0, p.AllCount())
}

func TestPoolAdd(t *testing.T) {
	p := NewTxPool(100)
	tx := core.NewTransaction([]byte("fooo"))
	p.Add(tx)
	assert.Equal(t, p.AllCount(), 1)
	// err = p.Add(tx)
	// fmt.Println(err)
	// assert.NotNil(t, err)
}

func TestPoolFlush(t *testing.T) {
	pool := NewTxPool(100)
	for i := 0; i < 100; i++ {
		tx := core.NewTransaction([]byte(fmt.Sprintf("test tx [%v]", i)))
		txHash := tx.Hash(core.TransactionHashesr{})
		fmt.Println(txHash)
		fmt.Println(pool.Contains(txHash))
		assert.False(t, pool.Contains(txHash))
		pool.Add(tx)
		assert.True(t, pool.Contains(txHash))
	}

	assert.Equal(t, 100, pool.AllCount())
	assert.Equal(t, 100, pool.PendingCount())
	pool.ClearPending()
	assert.Equal(t, 0, pool.PendingCount())
}

func TestTXpoolSort(t *testing.T) {
	p := NewTxPool(100)
	len := 100
	for i := 0; i < len; i++ {
		tx := core.NewTransaction([]byte(fmt.Sprintf("tran %v", i)))
		tx.SetFirstSeen(int64(i * 100))
		p.Add(tx)
	}
	// assert.Equal(t, p.Len(), len)
	// for i := 0; i < len-1; i++ {
	// 	val := p.Transactions()[i].FirstSeen() < p.Transactions()[i+1].FirstSeen()
	// 	assert.True(t, val)
	// }
}
