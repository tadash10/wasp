// Code generated by schema tool; DO NOT EDIT.

// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package coreaccounts

import "github.com/iotaledger/wasp/packages/wasmvm/wasmlib/go/wasmlib/wasmtypes"

const (
	ScName        = "accounts"
	ScDescription = "Chain account ledger contract"
	HScName       = wasmtypes.ScHname(0x3c4b5e02)
)

const (
	ParamAgentID        = "a"
	ParamCollection     = "C"
	ParamDestroyTokens  = "y"
	ParamFoundrySN      = "s"
	ParamGasReserve     = "g"
	ParamNftID          = "z"
	ParamSupplyDeltaAbs = "d"
	ParamTokenID        = "N"
	ParamTokenScheme    = "t"
)

const (
	ResultAccountNonce     = "n"
	ResultAllAccounts      = "this"
	ResultAmount           = "A"
	ResultAssets           = "this"
	ResultBalance          = "B"
	ResultBalances         = "this"
	ResultFoundries        = "this"
	ResultFoundryOutputBin = "b"
	ResultFoundrySN        = "s"
	ResultMapping          = "this"
	ResultNftData          = "e"
	ResultNftIDs           = "i"
	ResultTokens           = "B"
)

const (
	FuncDeposit                      = "deposit"
	FuncFoundryCreateNew             = "foundryCreateNew"
	FuncFoundryDestroy               = "foundryDestroy"
	FuncFoundryModifySupply          = "foundryModifySupply"
	FuncTransferAccountToChain       = "transferAccountToChain"
	FuncTransferAllowanceTo          = "transferAllowanceTo"
	FuncWithdraw                     = "withdraw"
	ViewAccountFoundries             = "accountFoundries"
	ViewAccountNFTAmount             = "accountNFTAmount"
	ViewAccountNFTAmountInCollection = "accountNFTAmountInCollection"
	ViewAccountNFTs                  = "accountNFTs"
	ViewAccountNFTsInCollection      = "accountNFTsInCollection"
	ViewAccounts                     = "accounts"
	ViewBalance                      = "balance"
	ViewBalanceBaseToken             = "balanceBaseToken"
	ViewBalanceNativeToken           = "balanceNativeToken"
	ViewFoundryOutput                = "foundryOutput"
	ViewGetAccountNonce              = "getAccountNonce"
	ViewGetNativeTokenIDRegistry     = "getNativeTokenIDRegistry"
	ViewNftData                      = "nftData"
	ViewTotalAssets                  = "totalAssets"
)

const (
	HFuncDeposit                      = wasmtypes.ScHname(0xbdc9102d)
	HFuncFoundryCreateNew             = wasmtypes.ScHname(0x41822f5f)
	HFuncFoundryDestroy               = wasmtypes.ScHname(0x85e4c893)
	HFuncFoundryModifySupply          = wasmtypes.ScHname(0x76a5868b)
	HFuncTransferAccountToChain       = wasmtypes.ScHname(0x07005c45)
	HFuncTransferAllowanceTo          = wasmtypes.ScHname(0x23f4e3a1)
	HFuncWithdraw                     = wasmtypes.ScHname(0x9dcc0f41)
	HViewAccountFoundries             = wasmtypes.ScHname(0xdc3a0c38)
	HViewAccountNFTAmount             = wasmtypes.ScHname(0xabefd5b5)
	HViewAccountNFTAmountInCollection = wasmtypes.ScHname(0xd7028e1b)
	HViewAccountNFTs                  = wasmtypes.ScHname(0x27422359)
	HViewAccountNFTsInCollection      = wasmtypes.ScHname(0xa37fb50f)
	HViewAccounts                     = wasmtypes.ScHname(0x3c4b5e02)
	HViewBalance                      = wasmtypes.ScHname(0x84168cb4)
	HViewBalanceBaseToken             = wasmtypes.ScHname(0x4c8ccd0f)
	HViewBalanceNativeToken           = wasmtypes.ScHname(0x1fea3104)
	HViewFoundryOutput                = wasmtypes.ScHname(0xd9647be3)
	HViewGetAccountNonce              = wasmtypes.ScHname(0x529d7df9)
	HViewGetNativeTokenIDRegistry     = wasmtypes.ScHname(0x2ad8a59f)
	HViewNftData                      = wasmtypes.ScHname(0x83c5c4da)
	HViewTotalAssets                  = wasmtypes.ScHname(0xfab0f8d2)
)
