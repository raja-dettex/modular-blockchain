package crypto

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"math/big"

	"github.com/raja-dettex/modular-blockchain/types"
)

type PublicKey []byte

type PrivateKey struct {
	Key *ecdsa.PrivateKey
}

func GeneratePrivateKey() PrivateKey {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	return PrivateKey{
		Key: key,
	}
}

func (k PrivateKey) GeneratePublicKey() PublicKey {
	return elliptic.MarshalCompressed(k.Key.PublicKey, k.Key.X, k.Key.Y)
}

func (k PrivateKey) Sign(data []byte) (*Signature, error) {
	r, s, err := ecdsa.Sign(rand.Reader, k.Key, data)
	if err != nil {
		return nil, err
	}
	return &Signature{
		R: r,
		S: s,
	}, nil
}

func (k PublicKey) Address() types.Address {
	h := sha256.Sum256(k)
	return types.AddressFromByte(h[len(h)-20:])
}

type Signature struct {
	R, S *big.Int
}

func (s *Signature) String() string {
	buff := new(bytes.Buffer)
	if err := gob.NewEncoder(buff).Encode(s); err != nil {
		return ""
	}
	return hex.EncodeToString(buff.Bytes())
}

func (s *Signature) Verify(pubKey PublicKey, data []byte) bool {
	x, y := elliptic.UnmarshalCompressed(elliptic.P256(), pubKey)
	key := &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}
	return ecdsa.Verify(key, data, s.R, s.S)
}
