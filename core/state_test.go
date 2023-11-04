package core

import (
	"testing"

	"github.com/raja-dettex/modular-blockchain/crypto"
	"github.com/stretchr/testify/assert"
)

func TestAccountTransferToFail(t *testing.T) {
	accountState := NewAccountState()
	from := crypto.GeneratePrivateKey().GeneratePublicKey().Address()

	to := crypto.GeneratePrivateKey().GeneratePublicKey().Address()
	assert.NotNil(t, accountState.Transfer(from, to, 90))

}
func TestAccountTransferToSucceed(t *testing.T) {
	accountState := NewAccountState()
	from := crypto.GeneratePrivateKey().GeneratePublicKey().Address()

	to := crypto.GeneratePrivateKey().GeneratePublicKey().Address()
	err := accountState.AddBalance(from, 100)
	assert.Nil(t, err)
	assert.Nil(t, accountState.Transfer(from, to, 90))

}
