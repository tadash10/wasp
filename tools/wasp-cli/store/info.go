package store

import (
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show the store",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

	},
}
