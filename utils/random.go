package utils

import (
	"crypto/rand"

	"github.com/raja-dettex/modular-blockchain/types"
)

func RandomBytes(size int) []byte {
	bytes := make([]byte, size)
	rand.Read(bytes)
	return bytes
}

func RandomHash(buff []byte) types.Hash {
	return types.HashFromBytes(buff)
}
