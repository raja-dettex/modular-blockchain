package crypto

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeyPair(t *testing.T) {
	privateKey := GeneratePrivateKey()
	publicKey := privateKey.GeneratePublicKey()
	address := publicKey.Address()
	fmt.Println(address)
}

func TestKeyPairSign_Verify(t *testing.T) {
	type testCase struct {
		name     string
		privKey  PrivateKey
		pubKey   PublicKey
		data     []byte
		expected bool
	}

	privKey := GeneratePrivateKey()
	pubKey := privKey.GeneratePublicKey()
	data := []byte("hello")
	oPrivKey := GeneratePrivateKey()
	oPubKey := oPrivKey.GeneratePublicKey()

	testCases := []testCase{
		{
			name:     "valid case",
			privKey:  privKey,
			pubKey:   pubKey,
			data:     data,
			expected: true,
		},
		{
			name:     "invalid case",
			privKey:  privKey,
			pubKey:   oPubKey,
			data:     data,
			expected: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sig, err := tc.privKey.Sign(tc.data)
			assert.Nil(t, err)
			assert.Equal(t, tc.expected, sig.Verify(tc.pubKey, tc.data))
		})
	}
}
