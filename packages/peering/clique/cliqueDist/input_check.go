// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package cliqueDist

import (
	"time"

	"github.com/iotaledger/wasp/packages/cryptolib"
	"github.com/iotaledger/wasp/packages/gpa"
	"github.com/iotaledger/wasp/packages/hashing"
)

type inputCheck struct {
	sessionID hashing.HashValue
	nodePubs  []*cryptolib.PublicKey
	timeout   time.Duration
	callback  Callback
}

func NewInputCheck(sessionID hashing.HashValue, nodes []*cryptolib.PublicKey, timeout time.Duration, callback Callback) gpa.Input {
	return &inputCheck{
		sessionID: sessionID,
		nodePubs:  nodes,
		timeout:   timeout,
		callback:  callback,
	}
}
