// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package jsonrpc

import (
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/iotaledger/wasp/packages/cryptolib"
	"github.com/iotaledger/wasp/packages/iscp"
	"github.com/iotaledger/wasp/packages/kv/dict"
)

type ChainBackend interface {
	EstimateGasOnLedger(scName string, funName string, transfer *iscp.FungibleTokens, args dict.Dict) (uint64, *iscp.FungibleTokens, error)
	PostOnLedgerRequest(scName string, funName string, transfer *iscp.FungibleTokens, args dict.Dict, gasBudget uint64) error
	EstimateGasOffLedger(scName string, funName string, args dict.Dict) (uint64, *iscp.FungibleTokens, error)
	PostOffLedgerRequest(scName string, funName string, args dict.Dict, gasBudget uint64) error
	CallView(scName string, funName string, args dict.Dict) (dict.Dict, error)
	Signer() cryptolib.VariantKeyPair
	EVMSendTransaction(tx *types.Transaction) error
	EVMEstimateGas(callMsg ethereum.CallMsg) (uint64, error)
	ISCCallView(scName string, funName string, args dict.Dict) (dict.Dict, error)
}
