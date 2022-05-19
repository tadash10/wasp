package chainclient

import (
	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/wasp/packages/cryptolib"
	"github.com/lmoe/stronghold.rs/bindings/native/go"
)

type Owner interface {
	PublicKey() string
	Address() iotago.Address
	Sign() bool
}

type KeyPairOwner struct {
	keyPair *cryptolib.KeyPair
}

func NewKeyPairOwner(keyPair *cryptolib.KeyPair) *KeyPairOwner {
	return &KeyPairOwner{keyPair: keyPair}
}

type StrongholdOwner struct {
	stronghold *stronghold.StrongholdNative
}

func NewStrongholdOwner(strongholdInstance *stronghold.StrongholdNative) *StrongholdOwner {
	return &StrongholdOwner{
		stronghold: strongholdInstance,
	}
}
