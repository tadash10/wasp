package store

import (
	"github.com/iotaledger/wasp/tools/wasp-cli/log"
	"github.com/spf13/cobra"
)

var storeCmd = &cobra.Command{
	Use:   "store <command>",
	Short: "Interact with the secure store",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		log.Check(cmd.Help())
	},
}

func Init(rootCmd *cobra.Command) {
	rootCmd.AddCommand(storeCmd)

	storeCmd.AddCommand(dumpCmd)
	storeCmd.AddCommand(resetCmd)
}
