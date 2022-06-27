package wallet

import (
	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/wasp/packages/cryptolib"
	"github.com/iotaledger/wasp/tools/wasp-cli/config"
	"github.com/iotaledger/wasp/tools/wasp-cli/log"
	"github.com/mr-tron/base58"
	"github.com/spf13/cobra"
)

type WalletConfig struct {
	Seed []byte
}

var addressIndex int

type Wallet struct {
	KeyPair cryptolib.VariantKeyPair
}

func (w *Wallet) PrivateKey() *cryptolib.PrivateKey {
	kp, ok := w.KeyPair.(*cryptolib.KeyPair)

	if ok {
		return kp.GetPrivateKey()
	}

	return nil
}

func (w *Wallet) Address() iotago.Address {
	return w.KeyPair.GetPublicKey().AsEd25519Address()
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new wallet",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		if config.IsPlainScheme() {
			err = config.Store.GenerateAndStorePlainSeed()
		} else {
			err = config.Store.InitializeNewStronghold()
		}

		log.Check(err)

		log.Printf("Initialized wallet seed, saved in key chain.\n")
		log.Printf("\nIMPORTANT: wasp-cli is alpha phase. The seed is currently being stored " +
			"in a plain text file which is NOT secure. Do not use this seed to store funds " +
			"in the mainnet!\n")
	},
}

func Load() *Wallet {
	if config.IsPlainScheme() {
		return initializePlainWallet()
	}

	return initializeStrongholdWallet()
}

func initializeStrongholdWallet() *Wallet {
	strongholdPtr, err := config.Store.OpenStronghold(uint32(addressIndex))
	if err != nil {
		log.Fatalf("[%s] call `init` first", err)
	}

	keyPair := cryptolib.NewStrongholdKeyPair(strongholdPtr, uint32(addressIndex))

	return &Wallet{KeyPair: keyPair}
}

func initializePlainWallet() *Wallet {
	seedb58, err := config.Store.Seed()

	log.Check(err)

	seedEnclave, err := seedb58.Open()
	defer seedEnclave.Destroy()

	if err != nil {
		log.Fatalf("call `init` first")
	}

	seedBytes, err := base58.Decode(seedEnclave.String())
	log.Check(err)
	seed := cryptolib.NewSeedFromBytes(seedBytes)
	kp := cryptolib.NewKeyPairFromSeed(seed)

	return &Wallet{KeyPair: kp}
}
