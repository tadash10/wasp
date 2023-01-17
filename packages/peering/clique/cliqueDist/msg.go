// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package cliqueDist

import (
	"errors"
	"fmt"

	"github.com/iotaledger/wasp/packages/gpa"
	"github.com/iotaledger/wasp/packages/util"
)

const (
	msgTypeQuery byte = iota
	msgTypeResponse
)

func (cd *cliqueDist) UnmarshalMessage(data []byte) (gpa.Message, error) {
	if len(data) < 1 {
		return nil, errors.New("cliqueDist::UnmarshalMessage: data too short")
	}
	switch data[0] {
	case msgTypeQuery:
		m := &msgQuery{}
		if err := m.UnmarshalBinary(data); err != nil {
			return nil, fmt.Errorf("cannot unmarshal cliqueDist.msgQuery: %w", err)
		}
		return m, nil
	case msgTypeResponse:
		m := &msgResponse{}
		if err := m.UnmarshalBinary(data); err != nil {
			return nil, fmt.Errorf("cannot unmarshal cliqueDist.msgResponse: %w", err)
		}
		return m, nil
	default:
		return nil, fmt.Errorf("cliqueDist::UnmarshalMessage: cannot parse message starting with: %v", util.PrefixHex(data, 20))
	}
}
