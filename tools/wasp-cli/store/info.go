package store

import (
	"github.com/iotaledger/wasp/tools/wasp-cli/config"
	"github.com/iotaledger/wasp/tools/wasp-cli/log"
	"github.com/spf13/cobra"
)

var resetCmd = &cobra.Command{
	Use:   "show",
	Short: "Show the store",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

		err := config.Store.Reset()

		if err != nil {
			log.Fatalf("Failed to reset store: %v", err)
		} else {
			log.Printf("Store was reset")
		}
	},
}

var initCmd = &cobra.Command{
	Use:   "show",
	Short: "Show the store",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

		/*
			wallet := Load()
			address := wallet.Address()

			outs, err := config.L1Client().OutputMap(address)
			log.Check(err)

			log.Printf("Address index %d\n", addressIndex)
			log.Printf("  Address: %s\n", address.Bech32(config.L1NetworkPrefix()))
			log.Printf("  Balance:\n")
			if log.VerboseFlag {
				printOutputsByOutputID(outs)
			} else {
				printOutputsByAsset(outs)
			}*/
	},
}
