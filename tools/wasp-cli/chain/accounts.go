package chain

import (
	"context"
	"strings"

	"github.com/spf13/cobra"

	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/wasp/clients/chainclient"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/kv/dict"
	"github.com/iotaledger/wasp/packages/vm/core/accounts"
	"github.com/iotaledger/wasp/tools/wasp-cli/cli/cliclients"
	"github.com/iotaledger/wasp/tools/wasp-cli/cli/config"
	"github.com/iotaledger/wasp/tools/wasp-cli/cli/wallet"
	"github.com/iotaledger/wasp/tools/wasp-cli/log"
	"github.com/iotaledger/wasp/tools/wasp-cli/util"
	"github.com/iotaledger/wasp/tools/wasp-cli/waspcmd"
)

func initListAccountsCmd() *cobra.Command {
	var node string
	var chain string

	cmd := &cobra.Command{
		Use:   "list-accounts",
		Short: "List L2 accounts",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			node = waspcmd.DefaultWaspNodeFallback(node)
			chain = defaultChainFallback(chain)

			client := cliclients.WaspClient(node)
			chainID := config.GetChain(chain)

			accountList, _, err := client.CorecontractsApi.AccountsGetAccounts(context.Background(), chainID.String()).Execute() //nolint:bodyclose // false positive
			log.Check(err)

			log.Printf("Total %d account(s) in chain %s\n", len(accountList.Accounts), config.GetChain(chain).String())

			header := []string{"agentid"}
			rows := make([][]string, len(accountList.Accounts))
			for i, account := range accountList.Accounts {
				rows[i] = []string{account}
			}
			log.PrintTable(header, rows)
		},
	}

	waspcmd.WithWaspNodeFlag(cmd, &node)
	withChainFlag(cmd, &chain)
	return cmd
}

func initBalanceCmd() *cobra.Command {
	var node string
	var chain string
	cmd := &cobra.Command{
		Use:   "balance [<agentid>]",
		Short: "Show the L2 balance of the given account",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			node = waspcmd.DefaultWaspNodeFallback(node)
			chain = defaultChainFallback(chain)

			var agentID isc.AgentID
			if len(args) == 0 {
				agentID = isc.NewAgentID(wallet.Load().Address())
			} else {
				var err error
				agentID, err = isc.NewAgentIDFromString(args[0])
				log.Check(err)
			}

			client := cliclients.WaspClient(node)
			chainID := config.GetChain(chain)
			balance, _, err := client.CorecontractsApi.AccountsGetAccountBalance(context.Background(), chainID.String(), agentID.String()).Execute() //nolint:bodyclose // false positive

			log.Check(err)

			header := []string{"token", "amount"}
			rows := make([][]string, len(balance.NativeTokens)+1)

			rows[0] = []string{"base", balance.BaseTokens}
			for k, v := range balance.NativeTokens {
				rows[k+1] = []string{v.Id, v.Amount}
			}

			log.PrintTable(header, rows)
		},
	}

	waspcmd.WithWaspNodeFlag(cmd, &node)
	withChainFlag(cmd, &chain)
	return cmd
}

func NFTFromNFTOutput(nftOutput *iotago.NFTOutput) *isc.NFT {
	if nftOutput == nil {
		return nil
	}

	return &isc.NFT{
		ID:       nftOutput.NFTID,
		Metadata: nftOutput.ImmutableFeatureSet().MetadataFeature().Data,
		Issuer:   nftOutput.ImmutableFeatureSet().IssuerFeature().Address,
	}
}

func initDepositCmd() *cobra.Command {
	var adjustStorageDeposit bool
	var node string
	var chain string

	cmd := &cobra.Command{
		Use:   "deposit [<agentid>] <token-id>:<amount>, [<token-id>:amount ...], <nft>:nftID, [<nft>:nftID ...]",
		Short: "Deposit L1 funds into the given (default: your) L2 account",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			node = waspcmd.DefaultWaspNodeFallback(node)
			chain = defaultChainFallback(chain)

			chainID := config.GetChain(chain)
			if strings.Contains(args[0], ":") {
				// deposit to own agentID
				tokens := util.ParseAssetArgs(args)
				util.WithSCTransaction(config.GetChain(chain), node, func() (*iotago.Transaction, error) {
					client := cliclients.WaspClient(node)

					return cliclients.SCClient(client, chainID, accounts.Contract.Hname()).PostRequest(
						accounts.FuncDeposit.Name,
						chainclient.PostRequestParams{
							Transfer:                 tokens,
							AutoAdjustStorageDeposit: adjustStorageDeposit,
						},
					)
				})
			} else {
				// deposit to some other agentID
				agentID, err := isc.NewAgentIDFromString(args[0])
				log.Check(err)
				tokensStr := strings.Split(strings.Join(args[1:], ""), ",")
				tokens := util.ParseAssetArgs(tokensStr)
				allowance := tokens.Clone()

				w := wallet.Load()
				outputsSet, err := cliclients.L1Client().OutputMap(w.KeyPair.Address())
				log.Check(err)

				var nftOutput *iotago.NFTOutput

				for _, k := range outputsSet {
					output, ok := k.(*iotago.NFTOutput)

					if ok {
						nftOutput = output
					}
				}

				util.WithSCTransaction(config.GetChain(chain), node, func() (*iotago.Transaction, error) {
					client := cliclients.WaspClient(node)

					return cliclients.SCClient(client, chainID, accounts.Contract.Hname()).PostRequest(
						accounts.FuncTransferAllowanceTo.Name,
						chainclient.PostRequestParams{
							Args: dict.Dict{
								accounts.ParamAgentID: agentID.Bytes(),
							},
							Transfer:                 tokens,
							Allowance:                allowance,
							AutoAdjustStorageDeposit: adjustStorageDeposit,
							NFT:                      NFTFromNFTOutput(nftOutput),
						},
					)
				})
			}
		},
	}

	cmd.Flags().BoolVarP(&adjustStorageDeposit, "adjust-storage-deposit", "s", false, "adjusts the amount of base tokens sent, if it's lower than the min storage deposit required")
	waspcmd.WithWaspNodeFlag(cmd, &node)
	withChainFlag(cmd, &chain)

	return cmd
}
