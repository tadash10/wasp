// Code generated by schema tool; DO NOT EDIT.

// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

import * as wasmtypes from '../wasmtypes';

export const ScName        = 'accounts';
export const ScDescription = 'Chain account ledger contract';
export const HScName       = new wasmtypes.ScHname(0x3c4b5e02);

export const ParamAgentID                = 'a';
export const ParamCollection             = 'C';
export const ParamDestroyTokens          = 'y';
export const ParamForceMinimumBaseTokens = 'f';
export const ParamFoundrySN              = 's';
export const ParamGasReserve             = 'g';
export const ParamNftID                  = 'z';
export const ParamSupplyDeltaAbs         = 'd';
export const ParamTokenID                = 'N';
export const ParamTokenScheme            = 't';

export const ResultAccountNonce     = 'n';
export const ResultAllAccounts      = 'this';
export const ResultAmount           = 'A';
export const ResultAssets           = 'this';
export const ResultBalance          = 'B';
export const ResultBalances         = 'this';
export const ResultFoundries        = 'this';
export const ResultFoundryOutputBin = 'b';
export const ResultFoundrySN        = 's';
export const ResultMapping          = 'this';
export const ResultNftData          = 'e';
export const ResultNftIDs           = 'i';
export const ResultTokens           = 'B';

export const FuncDeposit                      = 'deposit';
export const FuncFoundryCreateNew             = 'foundryCreateNew';
export const FuncFoundryDestroy               = 'foundryDestroy';
export const FuncFoundryModifySupply          = 'foundryModifySupply';
export const FuncHarvest                      = 'harvest';
export const FuncTransferAccountToChain       = 'transferAccountToChain';
export const FuncTransferAllowanceTo          = 'transferAllowanceTo';
export const FuncWithdraw                     = 'withdraw';
export const ViewAccountFoundries             = 'accountFoundries';
export const ViewAccountNFTAmount             = 'accountNFTAmount';
export const ViewAccountNFTAmountInCollection = 'accountNFTAmountInCollection';
export const ViewAccountNFTs                  = 'accountNFTs';
export const ViewAccountNFTsInCollection      = 'accountNFTsInCollection';
export const ViewAccounts                     = 'accounts';
export const ViewBalance                      = 'balance';
export const ViewBalanceBaseToken             = 'balanceBaseToken';
export const ViewBalanceNativeToken           = 'balanceNativeToken';
export const ViewFoundryOutput                = 'foundryOutput';
export const ViewGetAccountNonce              = 'getAccountNonce';
export const ViewGetNativeTokenIDRegistry     = 'getNativeTokenIDRegistry';
export const ViewNftData                      = 'nftData';
export const ViewTotalAssets                  = 'totalAssets';

export const HFuncDeposit                      = new wasmtypes.ScHname(0xbdc9102d);
export const HFuncFoundryCreateNew             = new wasmtypes.ScHname(0x41822f5f);
export const HFuncFoundryDestroy               = new wasmtypes.ScHname(0x85e4c893);
export const HFuncFoundryModifySupply          = new wasmtypes.ScHname(0x76a5868b);
export const HFuncHarvest                      = new wasmtypes.ScHname(0x7b40efbd);
export const HFuncTransferAccountToChain       = new wasmtypes.ScHname(0x07005c45);
export const HFuncTransferAllowanceTo          = new wasmtypes.ScHname(0x23f4e3a1);
export const HFuncWithdraw                     = new wasmtypes.ScHname(0x9dcc0f41);
export const HViewAccountFoundries             = new wasmtypes.ScHname(0xdc3a0c38);
export const HViewAccountNFTAmount             = new wasmtypes.ScHname(0xabefd5b5);
export const HViewAccountNFTAmountInCollection = new wasmtypes.ScHname(0xd7028e1b);
export const HViewAccountNFTs                  = new wasmtypes.ScHname(0x27422359);
export const HViewAccountNFTsInCollection      = new wasmtypes.ScHname(0xa37fb50f);
export const HViewAccounts                     = new wasmtypes.ScHname(0x3c4b5e02);
export const HViewBalance                      = new wasmtypes.ScHname(0x84168cb4);
export const HViewBalanceBaseToken             = new wasmtypes.ScHname(0x4c8ccd0f);
export const HViewBalanceNativeToken           = new wasmtypes.ScHname(0x1fea3104);
export const HViewFoundryOutput                = new wasmtypes.ScHname(0xd9647be3);
export const HViewGetAccountNonce              = new wasmtypes.ScHname(0x529d7df9);
export const HViewGetNativeTokenIDRegistry     = new wasmtypes.ScHname(0x2ad8a59f);
export const HViewNftData                      = new wasmtypes.ScHname(0x83c5c4da);
export const HViewTotalAssets                  = new wasmtypes.ScHname(0xfab0f8d2);
