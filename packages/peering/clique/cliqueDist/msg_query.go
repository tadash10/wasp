// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package cliqueDist

import (
	"fmt"
	"time"

	"github.com/iotaledger/wasp/packages/cryptolib"
	"github.com/iotaledger/wasp/packages/gpa"
	"github.com/iotaledger/wasp/packages/hashing"
)

type msgQuery struct {
	gpa.BasicMessage
	session      hashing.HashValue
	senderPubKey *cryptolib.PublicKey
	subQuery     []*cryptolib.PublicKey
	timeout      time.Duration
	ttl          byte
}

var _ gpa.Message = &msgQuery{}

func newMsgQuery(
	recipient gpa.NodeID,
	session hashing.HashValue,
	senderPubKey *cryptolib.PublicKey,
	subQuery []*cryptolib.PublicKey,
	timeout time.Duration,
	ttl byte,
) gpa.Message {
	return &msgQuery{
		BasicMessage: gpa.NewBasicMessage(recipient),
		session:      session,
		senderPubKey: senderPubKey,
		subQuery:     subQuery,
		timeout:      timeout,
		ttl:          ttl,
	}
}

func (m *msgQuery) Validate() error {
	if !gpa.NodeIDFromPublicKey(m.senderPubKey).Equals(m.Sender()) {
		return fmt.Errorf("invalid senderPubKey")
	}
	return nil
}

func (m *msgQuery) MarshalBinary() ([]byte, error) {
	panic("implement") // TODO: Implement.
}

func (m *msgQuery) UnmarshalBinary(data []byte) error {
	panic("implement") // TODO: Implement.
}

func (m *msgQuery) String() string {
	return fmt.Sprintf("{msgQuery, sender=%v, recipient=%v, ttl=%v, |subQuery|=%v}", m.Sender().ShortString(), m.Recipient().ShortString(), m.ttl, len(m.subQuery))
}
