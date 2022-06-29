package store

import (
	"github.com/iotaledger/wasp/tools/wasp-cli/config"
	"github.com/iotaledger/wasp/tools/wasp-cli/log"
	"github.com/spf13/cobra"
)

var dumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "Dump the store",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		dump, err := config.Store.Dump()
		var dumpFormatted [][]string

		for k, v := range dump {
			item := []string{k, v}

			dumpFormatted = append(dumpFormatted, item)
		}

		log.Check(err)
		log.PrintTable([]string{"Key", "Value"}, dumpFormatted)
	},
}

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset store",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		err := config.Store.Reset()

		log.Check(err)

		log.Printf("Store has been reset!")
	},
}
