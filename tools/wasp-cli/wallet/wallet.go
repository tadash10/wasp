package wallet

import (
	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/wasp/packages/cryptolib"
	"github.com/iotaledger/wasp/tools/wasp-cli/config"
	"github.com/iotaledger/wasp/tools/wasp-cli/log"
	"github.com/mr-tron/base58"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
			generatePlainWallet()
			log.Printf("Initialized wallet seed, saved as plain string inside the wasp-cli.json configuration.\n")
		} else if config.IsKeyChainScheme() {
			err = config.Store.GenerateAndStorePlainSeed()
			log.Printf("Initialized wallet seed, saved in key chain [IOTA_Foundation].\n")
		}

		log.Check(err)

		log.Printf("\nIMPORTANT: wasp-cli is alpha phase. Do not use this seed to store funds " +
			"in the mainnet!\n")
	},
}

func Load() *Wallet {
	if config.IsPlainScheme() {
		return loadPlainWallet()
	}

	if config.IsKeyChainScheme() {
		return loadKeyChainWallet()
	}

	log.Fatalf("Invalid wallet scheme")
	return nil
}

func loadPlainWallet() *Wallet {
	seedb58 := viper.GetString("wallet.seed")
	if seedb58 == "" {
		log.Fatalf("call `init` first")
	}
	seedBytes, err := base58.Decode(seedb58)
	log.Check(err)
	seed := cryptolib.NewSeedFromBytes(seedBytes)
	kp := cryptolib.NewKeyPairFromSeed(seed.SubSeed(uint64(addressIndex)))

	return &Wallet{KeyPair: kp}
}

func generatePlainWallet() *Wallet {
	seed := cryptolib.NewSeed()
	seedString := base58.Encode(seed[:])
	viper.Set("wallet.seed", seedString)
	log.Check(viper.WriteConfig())

	kp := cryptolib.NewKeyPairFromSeed(seed.SubSeed(uint64(addressIndex)))

	return &Wallet{KeyPair: kp}
}

func loadKeyChainWallet() *Wallet {
	seedb58, err := config.Store.Seed()

	log.Check(err)

	seedEnclave, err := seedb58.Open()
	defer seedEnclave.Destroy()

	if err != nil {
		//nolint:gocritic
		log.Fatalf("call `init` first") // exitAfterDefer happens here, but is no issue in this place.
	}

	seedBytes, err := base58.Decode(seedEnclave.String())
	log.Check(err)
	seed := cryptolib.NewSeedFromBytes(seedBytes)
	kp := cryptolib.NewKeyPairFromSeed(seed)

	return &Wallet{KeyPair: kp}
}
