package wallet

import (
	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/wasp/packages/cryptolib"
	"github.com/iotaledger/wasp/tools/wasp-cli/config"
	"github.com/iotaledger/wasp/tools/wasp-cli/log"
	stronghold "github.com/lmoe/stronghold.rs/bindings/native/go"
	"github.com/mr-tron/base58"
	"github.com/spf13/cobra"
	"os"
	"path"
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
		err := config.Store.StoreNewSeed()

		log.Check(err)

		log.Printf("Initialized wallet seed, saved in key chain.\n")
		log.Printf("\nIMPORTANT: wasp-cli is alpha phase. The seed is currently being stored " +
			"in a plain text file which is NOT secure. Do not use this seed to store funds " +
			"in the mainnet!\n")
	},
}

func Load() *Wallet {
	if config.WalletScheme() == config.WalletSchemePlain {
		return initializePlainWallet()
	}

	return initializeStrongholdWallet()
}

func initializeStrongholdWallet() *Wallet {
	key, err := config.Store.StrongholdKey()

	log.Check(err)

	stronghold := stronghold.NewStronghold(key)

	// TODO: Make configurable
	cwd, _ := os.Getwd()
	vaultPath := path.Join(cwd, "wasp-cli.vault")
	//

	success, err := stronghold.OpenOrCreate(vaultPath)

	log.Check(err)

	if !success {
		log.Fatalf("failed to open vault with an unknown error")
	}

	_, err = stronghold.DeriveSeed(uint32(addressIndex))

	log.Check(err)

	keyPair := cryptolib.NewStrongholdKeyPair(stronghold, uint32(addressIndex))

	return &Wallet{KeyPair: keyPair}
}

func initializePlainWallet() *Wallet {

	seedb58, err := config.Store.Seed()

	log.Check(err)

	if seedb58 == "" {
		log.Fatalf("call `init` first")
	}
	seedBytes, err := base58.Decode(seedb58)
	log.Check(err)
	seed := cryptolib.NewSeedFromBytes(seedBytes)
	kp := cryptolib.NewKeyPairFromSeed(seed)

	return &Wallet{KeyPair: kp}
}
