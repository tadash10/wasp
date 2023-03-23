package verify

import (
	"github.com/iotaledger/wasp/tools/wasp-cli/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)


var (

)

func initVerifyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "verify <command>",
		Short: "Verify code with Blockscout",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			log.Check(cmd.Help())
		},
	}
}

func initContractCmd() *cobra.Command {
	return &cobra.Command{
		Use: "contract",
		Short: "Verify a contract with blockscout",
		Args: cobra.OnlyValidArgs,
		
	}
}

func Init(rootCmd *cobra.Command) {
	verifyCmd := initVerifyCmd()
	rootCmd.AddCommand(verifyCmd)

	contractCmd := initContractCmd()
	contractCmd.Flags().AddFlag(pflag.String())

	verifyCmd.AddCommand(contractCmd)
}
