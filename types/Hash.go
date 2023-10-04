package types

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

type Hash [32]uint8

func (h Hash) IsZero() bool {
	for i := 0; i < 32; i++ {
		if h[i] != 0 {
			return false
		}
	}
	return true
}

func (h Hash) ToByteSlice() []byte {
	buff := make([]byte, 32)
	for i := 0; i < 32; i++ {
		buff[i] = h[i]
	}
	return buff
}

func (h Hash) String() string {
	return hex.EncodeToString(h.ToByteSlice())
}

func HashFromBytes(b []byte) Hash {
	if len(b) != 32 {
		msg := fmt.Sprintf("not enough byte size %v", len(b))
		panic(msg)
	}
	var value [32]uint8
	for i := 0; i < 32; i++ {
		value[i] = b[i]
	}
	return Hash(value)
}

func RandomHash() Hash {
	b := make([]byte, 32)
	rand.Read(b)
	return HashFromBytes(b)
}
