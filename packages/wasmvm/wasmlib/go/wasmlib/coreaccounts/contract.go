// Code generated by schema tool; DO NOT EDIT.

// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package coreaccounts

import "github.com/iotaledger/wasp/packages/wasmvm/wasmlib/go/wasmlib"

type DepositCall struct {
	Func *wasmlib.ScFunc
}

type FoundryCreateNewCall struct {
	Func    *wasmlib.ScFunc
	Params  MutableFoundryCreateNewParams
	Results ImmutableFoundryCreateNewResults
}

type FoundryDestroyCall struct {
	Func   *wasmlib.ScFunc
	Params MutableFoundryDestroyParams
}

type FoundryModifySupplyCall struct {
	Func   *wasmlib.ScFunc
	Params MutableFoundryModifySupplyParams
}

type HarvestCall struct {
	Func   *wasmlib.ScFunc
	Params MutableHarvestParams
}

type TransferAccountToChainCall struct {
	Func   *wasmlib.ScFunc
	Params MutableTransferAccountToChainParams
}

type TransferAllowanceToCall struct {
	Func   *wasmlib.ScFunc
	Params MutableTransferAllowanceToParams
}

type WithdrawCall struct {
	Func *wasmlib.ScFunc
}

type AccountFoundriesCall struct {
	Func    *wasmlib.ScView
	Params  MutableAccountFoundriesParams
	Results ImmutableAccountFoundriesResults
}

type AccountNFTAmountCall struct {
	Func    *wasmlib.ScView
	Params  MutableAccountNFTAmountParams
	Results ImmutableAccountNFTAmountResults
}

type AccountNFTAmountInCollectionCall struct {
	Func    *wasmlib.ScView
	Params  MutableAccountNFTAmountInCollectionParams
	Results ImmutableAccountNFTAmountInCollectionResults
}

type AccountNFTsCall struct {
	Func    *wasmlib.ScView
	Params  MutableAccountNFTsParams
	Results ImmutableAccountNFTsResults
}

type AccountNFTsInCollectionCall struct {
	Func    *wasmlib.ScView
	Params  MutableAccountNFTsInCollectionParams
	Results ImmutableAccountNFTsInCollectionResults
}

type AccountsCall struct {
	Func    *wasmlib.ScView
	Results ImmutableAccountsResults
}

type BalanceCall struct {
	Func    *wasmlib.ScView
	Params  MutableBalanceParams
	Results ImmutableBalanceResults
}

type BalanceBaseTokenCall struct {
	Func    *wasmlib.ScView
	Params  MutableBalanceBaseTokenParams
	Results ImmutableBalanceBaseTokenResults
}

type BalanceNativeTokenCall struct {
	Func    *wasmlib.ScView
	Params  MutableBalanceNativeTokenParams
	Results ImmutableBalanceNativeTokenResults
}

type FoundryOutputCall struct {
	Func    *wasmlib.ScView
	Params  MutableFoundryOutputParams
	Results ImmutableFoundryOutputResults
}

type GetAccountNonceCall struct {
	Func    *wasmlib.ScView
	Params  MutableGetAccountNonceParams
	Results ImmutableGetAccountNonceResults
}

type GetNativeTokenIDRegistryCall struct {
	Func    *wasmlib.ScView
	Results ImmutableGetNativeTokenIDRegistryResults
}

type NftDataCall struct {
	Func    *wasmlib.ScView
	Params  MutableNftDataParams
	Results ImmutableNftDataResults
}

type TotalAssetsCall struct {
	Func    *wasmlib.ScView
	Results ImmutableTotalAssetsResults
}

type Funcs struct{}

var ScFuncs Funcs

// A no-op that has the side effect of crediting any transferred tokens to the sender's account.
func (sc Funcs) Deposit(ctx wasmlib.ScFuncCallContext) *DepositCall {
	return &DepositCall{Func: wasmlib.NewScFunc(ctx, HScName, HFuncDeposit)}
}

// Creates a new foundry with the specified token scheme, and assigns the foundry to the sender.
func (sc Funcs) FoundryCreateNew(ctx wasmlib.ScFuncCallContext) *FoundryCreateNewCall {
	f := &FoundryCreateNewCall{Func: wasmlib.NewScFunc(ctx, HScName, HFuncFoundryCreateNew)}
	f.Params.Proxy = wasmlib.NewCallParamsProxy(&f.Func.ScView)
	wasmlib.NewCallResultsProxy(&f.Func.ScView, &f.Results.Proxy)
	return f
}

// Destroys a given foundry output on L1, reimbursing the storage deposit to the caller.
// The foundry must be owned by the caller.
func (sc Funcs) FoundryDestroy(ctx wasmlib.ScFuncCallContext) *FoundryDestroyCall {
	f := &FoundryDestroyCall{Func: wasmlib.NewScFunc(ctx, HScName, HFuncFoundryDestroy)}
	f.Params.Proxy = wasmlib.NewCallParamsProxy(&f.Func.ScView)
	return f
}

// Mints or destroys tokens for the given foundry, which must be owned by the caller.
func (sc Funcs) FoundryModifySupply(ctx wasmlib.ScFuncCallContext) *FoundryModifySupplyCall {
	f := &FoundryModifySupplyCall{Func: wasmlib.NewScFunc(ctx, HScName, HFuncFoundryModifySupply)}
	f.Params.Proxy = wasmlib.NewCallParamsProxy(&f.Func.ScView)
	return f
}

// Moves all tokens from the chain common account to the sender's L2 account.
// The chain owner is the only one who can call this entry point.
func (sc Funcs) Harvest(ctx wasmlib.ScFuncCallContext) *HarvestCall {
	f := &HarvestCall{Func: wasmlib.NewScFunc(ctx, HScName, HFuncHarvest)}
	f.Params.Proxy = wasmlib.NewCallParamsProxy(&f.Func.ScView)
	return f
}

// Transfers the specified allowance from the sender SC's L2 account on
// the target chain to the sender SC's L2 account on the origin chain.
func (sc Funcs) TransferAccountToChain(ctx wasmlib.ScFuncCallContext) *TransferAccountToChainCall {
	f := &TransferAccountToChainCall{Func: wasmlib.NewScFunc(ctx, HScName, HFuncTransferAccountToChain)}
	f.Params.Proxy = wasmlib.NewCallParamsProxy(&f.Func.ScView)
	return f
}

// Transfers the specified allowance from the sender's L2 account
// to the given L2 account on the chain.
func (sc Funcs) TransferAllowanceTo(ctx wasmlib.ScFuncCallContext) *TransferAllowanceToCall {
	f := &TransferAllowanceToCall{Func: wasmlib.NewScFunc(ctx, HScName, HFuncTransferAllowanceTo)}
	f.Params.Proxy = wasmlib.NewCallParamsProxy(&f.Func.ScView)
	return f
}

// Moves tokens from the caller's on-chain account to the caller's L1 address.
// The number of tokens to be withdrawn must be specified via the allowance of the request.
func (sc Funcs) Withdraw(ctx wasmlib.ScFuncCallContext) *WithdrawCall {
	return &WithdrawCall{Func: wasmlib.NewScFunc(ctx, HScName, HFuncWithdraw)}
}

// Returns a set of all foundries owned by the given account.
func (sc Funcs) AccountFoundries(ctx wasmlib.ScViewCallContext) *AccountFoundriesCall {
	f := &AccountFoundriesCall{Func: wasmlib.NewScView(ctx, HScName, HViewAccountFoundries)}
	f.Params.Proxy = wasmlib.NewCallParamsProxy(f.Func)
	wasmlib.NewCallResultsProxy(f.Func, &f.Results.Proxy)
	return f
}

// Returns the amount of NFTs owned by the given account.
func (sc Funcs) AccountNFTAmount(ctx wasmlib.ScViewCallContext) *AccountNFTAmountCall {
	f := &AccountNFTAmountCall{Func: wasmlib.NewScView(ctx, HScName, HViewAccountNFTAmount)}
	f.Params.Proxy = wasmlib.NewCallParamsProxy(f.Func)
	wasmlib.NewCallResultsProxy(f.Func, &f.Results.Proxy)
	return f
}

// Returns the amount of NFTs in the specified collection owned by the given account.
func (sc Funcs) AccountNFTAmountInCollection(ctx wasmlib.ScViewCallContext) *AccountNFTAmountInCollectionCall {
	f := &AccountNFTAmountInCollectionCall{Func: wasmlib.NewScView(ctx, HScName, HViewAccountNFTAmountInCollection)}
	f.Params.Proxy = wasmlib.NewCallParamsProxy(f.Func)
	wasmlib.NewCallResultsProxy(f.Func, &f.Results.Proxy)
	return f
}

// Returns the NFT IDs for all NFTs owned by the given account.
func (sc Funcs) AccountNFTs(ctx wasmlib.ScViewCallContext) *AccountNFTsCall {
	f := &AccountNFTsCall{Func: wasmlib.NewScView(ctx, HScName, HViewAccountNFTs)}
	f.Params.Proxy = wasmlib.NewCallParamsProxy(f.Func)
	wasmlib.NewCallResultsProxy(f.Func, &f.Results.Proxy)
	return f
}

// Returns the NFT IDs for all NFTs in the specified collection owned by the given account.
func (sc Funcs) AccountNFTsInCollection(ctx wasmlib.ScViewCallContext) *AccountNFTsInCollectionCall {
	f := &AccountNFTsInCollectionCall{Func: wasmlib.NewScView(ctx, HScName, HViewAccountNFTsInCollection)}
	f.Params.Proxy = wasmlib.NewCallParamsProxy(f.Func)
	wasmlib.NewCallResultsProxy(f.Func, &f.Results.Proxy)
	return f
}

// Returns a set of all agent IDs that own assets on the chain.
func (sc Funcs) Accounts(ctx wasmlib.ScViewCallContext) *AccountsCall {
	f := &AccountsCall{Func: wasmlib.NewScView(ctx, HScName, HViewAccounts)}
	wasmlib.NewCallResultsProxy(f.Func, &f.Results.Proxy)
	return f
}

// Returns the fungible tokens owned by the given Agent ID on the chain.
func (sc Funcs) Balance(ctx wasmlib.ScViewCallContext) *BalanceCall {
	f := &BalanceCall{Func: wasmlib.NewScView(ctx, HScName, HViewBalance)}
	f.Params.Proxy = wasmlib.NewCallParamsProxy(f.Func)
	wasmlib.NewCallResultsProxy(f.Func, &f.Results.Proxy)
	return f
}

// Returns the amount of base tokens owned by an agent on the chain
func (sc Funcs) BalanceBaseToken(ctx wasmlib.ScViewCallContext) *BalanceBaseTokenCall {
	f := &BalanceBaseTokenCall{Func: wasmlib.NewScView(ctx, HScName, HViewBalanceBaseToken)}
	f.Params.Proxy = wasmlib.NewCallParamsProxy(f.Func)
	wasmlib.NewCallResultsProxy(f.Func, &f.Results.Proxy)
	return f
}

// Returns the amount of specific native tokens owned by an agent on the chain
func (sc Funcs) BalanceNativeToken(ctx wasmlib.ScViewCallContext) *BalanceNativeTokenCall {
	f := &BalanceNativeTokenCall{Func: wasmlib.NewScView(ctx, HScName, HViewBalanceNativeToken)}
	f.Params.Proxy = wasmlib.NewCallParamsProxy(f.Func)
	wasmlib.NewCallResultsProxy(f.Func, &f.Results.Proxy)
	return f
}

// Returns specified foundry output in serialized form.
func (sc Funcs) FoundryOutput(ctx wasmlib.ScViewCallContext) *FoundryOutputCall {
	f := &FoundryOutputCall{Func: wasmlib.NewScView(ctx, HScName, HViewFoundryOutput)}
	f.Params.Proxy = wasmlib.NewCallParamsProxy(f.Func)
	wasmlib.NewCallResultsProxy(f.Func, &f.Results.Proxy)
	return f
}

// Returns the current account nonce for an Agent.
// The account nonce is used to issue unique off-ledger requests.
func (sc Funcs) GetAccountNonce(ctx wasmlib.ScViewCallContext) *GetAccountNonceCall {
	f := &GetAccountNonceCall{Func: wasmlib.NewScView(ctx, HScName, HViewGetAccountNonce)}
	f.Params.Proxy = wasmlib.NewCallParamsProxy(f.Func)
	wasmlib.NewCallResultsProxy(f.Func, &f.Results.Proxy)
	return f
}

// Returns a set of all native tokenIDs that are owned by the chain.
func (sc Funcs) GetNativeTokenIDRegistry(ctx wasmlib.ScViewCallContext) *GetNativeTokenIDRegistryCall {
	f := &GetNativeTokenIDRegistryCall{Func: wasmlib.NewScView(ctx, HScName, HViewGetNativeTokenIDRegistry)}
	wasmlib.NewCallResultsProxy(f.Func, &f.Results.Proxy)
	return f
}

// Returns the data for a given NFT that is on the chain.
func (sc Funcs) NftData(ctx wasmlib.ScViewCallContext) *NftDataCall {
	f := &NftDataCall{Func: wasmlib.NewScView(ctx, HScName, HViewNftData)}
	f.Params.Proxy = wasmlib.NewCallParamsProxy(f.Func)
	wasmlib.NewCallResultsProxy(f.Func, &f.Results.Proxy)
	return f
}

// Returns the balances of all fungible tokens controlled by the chain.
func (sc Funcs) TotalAssets(ctx wasmlib.ScViewCallContext) *TotalAssetsCall {
	f := &TotalAssetsCall{Func: wasmlib.NewScView(ctx, HScName, HViewTotalAssets)}
	wasmlib.NewCallResultsProxy(f.Func, &f.Results.Proxy)
	return f
}

var exportMap = wasmlib.ScExportMap{
	Names: []string{
		FuncDeposit,
		FuncFoundryCreateNew,
		FuncFoundryDestroy,
		FuncFoundryModifySupply,
		FuncHarvest,
		FuncTransferAccountToChain,
		FuncTransferAllowanceTo,
		FuncWithdraw,
		ViewAccountFoundries,
		ViewAccountNFTAmount,
		ViewAccountNFTAmountInCollection,
		ViewAccountNFTs,
		ViewAccountNFTsInCollection,
		ViewAccounts,
		ViewBalance,
		ViewBalanceBaseToken,
		ViewBalanceNativeToken,
		ViewFoundryOutput,
		ViewGetAccountNonce,
		ViewGetNativeTokenIDRegistry,
		ViewNftData,
		ViewTotalAssets,
	},
	Funcs: []wasmlib.ScFuncContextFunction{
		wasmlib.FuncError,
		wasmlib.FuncError,
		wasmlib.FuncError,
		wasmlib.FuncError,
		wasmlib.FuncError,
		wasmlib.FuncError,
		wasmlib.FuncError,
		wasmlib.FuncError,
	},
	Views: []wasmlib.ScViewContextFunction{
		wasmlib.ViewError,
		wasmlib.ViewError,
		wasmlib.ViewError,
		wasmlib.ViewError,
		wasmlib.ViewError,
		wasmlib.ViewError,
		wasmlib.ViewError,
		wasmlib.ViewError,
		wasmlib.ViewError,
		wasmlib.ViewError,
		wasmlib.ViewError,
		wasmlib.ViewError,
		wasmlib.ViewError,
		wasmlib.ViewError,
	},
}

func OnDispatch(index int32) *wasmlib.ScExportMap {
	if index < 0 {
		return exportMap.Dispatch(index)
	}

	panic("Calling core contract?")
}
