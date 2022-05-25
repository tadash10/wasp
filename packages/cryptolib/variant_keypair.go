package cryptolib

import (
	"fmt"
	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/lmoe/stronghold.rs/bindings/native/go"
)

// VariantKeyPair originates from cryptolib.KeyPair
type VariantKeyPair interface {
	GetPublicKey() *PublicKey
	Address() *iotago.Ed25519Address

	Sign(data []byte) []byte
}

type StrongholdKeyPair struct {
	stronghold *stronghold.StrongholdNative
	index      uint32
}

func NewStrongholdKeyPair(strongholdInstance *stronghold.StrongholdNative, index uint32) *StrongholdKeyPair {
	return &StrongholdKeyPair{
		stronghold: strongholdInstance,
		index:      index,
	}
}

func (kp *StrongholdKeyPair) GetPublicKey() []byte {
	publicKey, _ := kp.stronghold.GetPublicKeyFromDerived(kp.index)

	return publicKey[:]
}

func (kp *StrongholdKeyPair) Address() *iotago.Ed25519Address {
	address, _ := kp.stronghold.GetAddress(kp.index)
	return (*iotago.Ed25519Address)(&address)
}

func (kp *StrongholdKeyPair) Sign(data []byte) []byte {
	recordPath := fmt.Sprintf("seed.%d", kp.index)
	signature, _ := kp.stronghold.Sign(recordPath, data)

	return signature[:]
}
