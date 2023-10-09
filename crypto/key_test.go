package crypto

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"errors"
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

type P256Curve struct{}

func (P256Curve) GobEncode() ([]byte, error) {
	return []byte("P-256"), nil
}

func (P256Curve) GobDecode(data []byte) error {
	if string(data) != "P-256" {
		return errors.New("invalid data for P256Curve")
	}
	return nil
}

func TestEncodingPublicKey(t *testing.T) {
	privKey := GeneratePrivateKey()
	pbKey := privKey.GeneratePublicKey()
	buff := &bytes.Buffer{}

	gob.Register(elliptic.P256())
	err := gob.NewEncoder(buff).Encode(pbKey)
	fmt.Println(err)
	assert.Nil(t, err)
	fmt.Println(buff)
}
