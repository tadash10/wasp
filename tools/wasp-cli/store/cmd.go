package store

import (
	"github.com/spf13/cobra"
)

func Init(rootCmd *cobra.Command) {
	rootCmd.AddCommand(resetCmd)
	//rootCmd.AddCommand()
}
