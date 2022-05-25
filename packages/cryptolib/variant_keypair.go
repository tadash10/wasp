package cryptolib

import (
	"fmt"
	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/lmoe/stronghold.rs/bindings/native/go"
	"github.com/pkg/errors"
)

// VariantKeyPair originates from cryptolib.KeyPair
type VariantKeyPair interface {
	GetPublicKey() *PublicKey
	Address() *iotago.Ed25519Address
	AsAddressSigner() iotago.AddressSigner
	AddressKeysForEd25519Address(addr *iotago.Ed25519Address) iotago.AddressKeys
	Sign(data []byte) []byte
}

type StrongholdAddressSigner struct {
	stronghold *stronghold.StrongholdNative
	index      uint32
}

func NewStrongholdAddressSigner(strongholdInstance *stronghold.StrongholdNative, index uint32) *StrongholdAddressSigner {
	return &StrongholdAddressSigner{
		stronghold: strongholdInstance,
		index:      index,
	}
}

func (s *StrongholdAddressSigner) Sign(address iotago.Address, msg []byte) (signature iotago.Signature, err error) {
	strongholdAddress, err := s.stronghold.GetAddress(s.index)
	ed25519Address := (*iotago.Ed25519Address)(&strongholdAddress)

	if !address.Equal(ed25519Address) {
		return nil, errors.Errorf("Stronghold Address: [%v] mismatches supplied address: [%v]", ed25519Address, address)
	}

	signed, err := s.stronghold.SignForDerived(s.index, msg)

	if err != nil {
		return nil, err
	}

	publicKey, err := s.stronghold.GetPublicKeyFromDerived(0)

	if err != nil {
		return nil, err
	}

	ed25519Sig := &iotago.Ed25519Signature{}
	copy(ed25519Sig.Signature[:], signed[:])
	copy(ed25519Sig.PublicKey[:], publicKey[:])

	return ed25519Sig, nil
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

func (kp *StrongholdKeyPair) AsAddressSigner() *StrongholdAddressSigner {
	return NewStrongholdAddressSigner(kp.stronghold, kp.index)
}

func (kp *StrongholdKeyPair) AddressKeysForEd25519Address(addr *iotago.Ed25519Address) iotago.AddressKeys {
	return iotago.AddressKeys{Address: addr}
}
