package wallet

import (
	"github.com/spf13/viper"

	wasp_wallet_sdk "github.com/iotaledger/wasp-wallet-sdk"
	"github.com/iotaledger/wasp/tools/wasp-cli/cli/wallet/providers"
	"github.com/iotaledger/wasp/tools/wasp-cli/cli/wallet/wallets"
	"github.com/iotaledger/wasp/tools/wasp-cli/log"
)

var AddressIndex uint32

const (
	SchemeInMemory   = "in_memory"
	SchemeLedger     = "sdk_ledger"
	SchemeStronghold = "sdk_stronghold"
)

func GetWalletScheme() string {
	scheme := viper.GetString("wallet.scheme")

	switch scheme {
	case SchemeLedger, SchemeInMemory, SchemeStronghold:
		return scheme
	default:
		log.Fatalf("invalid wallet scheme configured")
	}
	return ""
}

func getIotaSDK() *wasp_wallet_sdk.IOTASDK {
	sdk, err := wasp_wallet_sdk.NewIotaSDK("/home/luke/dev/iota-sdk/target/release/libiota_sdk_go.so")
	log.Check(err)
	return sdk
}

func Load() wallets.Wallet {
	walletScheme := GetWalletScheme()

	log.Printf("Scheme: %v\n", walletScheme)
	switch walletScheme {
	case SchemeInMemory:
		return providers.LoadInMemory(AddressIndex)
	case SchemeLedger:
		return providers.LoadLedgerWallet(getIotaSDK(), AddressIndex)
	case SchemeStronghold:
		return providers.LoadStrongholdWallet(getIotaSDK(), AddressIndex)
	}

	return nil
}

func InitWallet() {
	walletScheme := GetWalletScheme()

	switch walletScheme {
	case SchemeInMemory:
		providers.CreateNewInMemory()
	case SchemeLedger:
		log.Printf("Ledger wallet scheme selected, no initialization required")
	case SchemeStronghold:
		providers.CreateNewStrongholdWallet(getIotaSDK())
	}
}
