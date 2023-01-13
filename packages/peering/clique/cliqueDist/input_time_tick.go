// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package cliqueDist

import (
	"time"

	"github.com/iotaledger/wasp/packages/gpa"
)

type inputTimeTick struct {
	timestamp time.Time
}

func NewInputTimeTick(timestamp time.Time) gpa.Input {
	return &inputTimeTick{timestamp: timestamp}
}
