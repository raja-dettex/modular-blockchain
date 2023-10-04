package types

import "encoding/hex"

type Address [20]uint8

func (a Address) ToSlice() []byte {
	buff := make([]byte, 20)
	for i := 0; i < 20; i++ {
		buff[i] = a[i]
	}
	return buff
}

func (a Address) String() string {
	return hex.EncodeToString(a.ToSlice())
}

func AddressFromByte(b []byte) Address {
	if len(b) != 20 {
		panic("byte size can not be other than 20")
	}
	var value [20]uint8
	for i := 0; i < 20; i++ {
		value[i] = b[i]
	}
	return Address(value)
}
