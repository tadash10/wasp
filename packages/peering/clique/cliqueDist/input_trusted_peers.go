// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package cliqueDist

import (
	"github.com/iotaledger/wasp/packages/cryptolib"
	"github.com/iotaledger/wasp/packages/gpa"
)

type inputTrustedPeers struct {
	trusted map[gpa.NodeID]*cryptolib.PublicKey
}

var _ gpa.Input = &inputTrustedPeers{}

func NewInputTrustedPeers(trusted map[gpa.NodeID]*cryptolib.PublicKey) gpa.Input {
	return &inputTrustedPeers{trusted: trusted}
}
