package cryptolib

import (
	"fmt"
	"github.com/iotaledger/iota.go/v3"
	stronghold "github.com/lmoe/stronghold.rs/bindings/native/go"
)

type StrongholdAddressSigner struct {
	stronghold   *stronghold.StrongholdNative
	addressIndex uint32
}

func NewStrongholdAddressSigner(strongholdInstance *stronghold.StrongholdNative, addressIndex uint32) *StrongholdAddressSigner {
	return &StrongholdAddressSigner{
		stronghold:   strongholdInstance,
		addressIndex: addressIndex,
	}
}

func (s *StrongholdAddressSigner) Sign(address iotago.Address, msg []byte) (signature iotago.Signature, err error) {
	strongholdAddress, err := s.stronghold.GetAddress(s.addressIndex)
	ed25519Address := (*iotago.Ed25519Address)(&strongholdAddress)

	if !address.Equal(ed25519Address) {
		return nil, fmt.Errorf("stronghold Address: [%v] mismatches supplied address: [%v]", ed25519Address, address)
	}

	signed, err := s.stronghold.SignForDerived(s.addressIndex, msg)

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
	stronghold   *stronghold.StrongholdNative
	addressIndex uint32
}

func NewStrongholdKeyPair(strongholdInstance *stronghold.StrongholdNative, addressIndex uint32) *StrongholdKeyPair {
	return &StrongholdKeyPair{
		stronghold:   strongholdInstance,
		addressIndex: addressIndex,
	}
}

func (kp *StrongholdKeyPair) GetPublicKey() *PublicKey {
	publicKeyBytes, _ := kp.stronghold.GetPublicKeyFromDerived(kp.addressIndex)
	publicKey, _ := NewPublicKeyFromBytes(publicKeyBytes[:])

	return publicKey
}

func (kp *StrongholdKeyPair) Address() *iotago.Ed25519Address {
	address, _ := kp.stronghold.GetAddress(kp.addressIndex)
	return (*iotago.Ed25519Address)(&address)
}

func (kp *StrongholdKeyPair) Sign(data []byte) []byte {
	recordPath := fmt.Sprintf("seed.%d", kp.addressIndex)
	signature, _ := kp.stronghold.Sign(recordPath, data)

	return signature[:]
}

func (kp *StrongholdKeyPair) AsAddressSigner() iotago.AddressSigner {
	return NewStrongholdAddressSigner(kp.stronghold, kp.addressIndex)
}

func (kp *StrongholdKeyPair) AddressKeysForEd25519Address(addr *iotago.Ed25519Address) iotago.AddressKeys {
	return iotago.AddressKeys{Address: addr}
}
