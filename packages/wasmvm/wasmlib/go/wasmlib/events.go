// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package wasmlib

import (
	"github.com/iotaledger/wasp/packages/wasmvm/wasmlib/go/wasmlib/wasmtypes"
)

type IEventHandlers interface {
	CallHandler(topic string, dec *wasmtypes.WasmDecoder)
	ID() uint32
}

var nextID = uint32(0)

func EventHandlersGenerateID() uint32 {
	nextID++
	return nextID
}

func NewEventEncoder() *wasmtypes.WasmEncoder {
	enc := wasmtypes.NewWasmEncoder()
	wasmtypes.Uint64Encode(enc, ScFuncContext{}.Timestamp())
	return enc
}

func EventEmit(topic string, enc *wasmtypes.WasmEncoder) {
	ScFuncContext{}.Event(topic + "|" + wasmtypes.HexEncode(enc.Buf()))
}
