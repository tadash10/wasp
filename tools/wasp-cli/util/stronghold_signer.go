package util

import (
	iotago "github.com/iotaledger/iota.go/v3"
	stronghold "github.com/lmoe/stronghold.rs/bindings/native/go"
)

type StrongholdSigner struct {
	stronghold   *stronghold.StrongholdNative
	addressIndex int
}

func NewStrongholdSigner(stronghold *stronghold.StrongholdNative, addressIndex int) *StrongholdSigner {
	return &StrongholdSigner{
		stronghold:   stronghold,
		addressIndex: addressIndex,
	}
}

func (s *StrongholdSigner) Sign(addr iotago.Address, msg []byte) (iotago.Signature, error) {
	signature := &iotago.Ed25519Signature{}

	// TODO: Add stronghold sign here

	return signature, nil
}
