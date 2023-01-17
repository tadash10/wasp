// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package cliqueDist

import (
	"fmt"
	"time"

	"github.com/iotaledger/hive.go/core/marshalutil"
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
}

var _ gpa.Message = &msgQuery{}

func newMsgQuery(
	recipient gpa.NodeID,
	session hashing.HashValue,
	senderPubKey *cryptolib.PublicKey,
	subQuery []*cryptolib.PublicKey,
	timeout time.Duration,
) gpa.Message {
	return &msgQuery{
		BasicMessage: gpa.NewBasicMessage(recipient),
		session:      session,
		senderPubKey: senderPubKey,
		subQuery:     subQuery,
		timeout:      timeout,
	}
}

func (m *msgQuery) Validate() error {
	if !gpa.NodeIDFromPublicKey(m.senderPubKey).Equals(m.Sender()) {
		return fmt.Errorf("invalid senderPubKey")
	}
	return nil
}

func (m *msgQuery) MarshalBinary() ([]byte, error) {
	mu := marshalutil.New()
	mu.WriteByte(msgTypeQuery)
	mu.Write(m.session)
	mu.Write(m.senderPubKey)
	mu.WriteUint16(uint16(len(m.subQuery)))
	for _, sq := range m.subQuery {
		mu.Write(sq)
	}
	mu.WriteInt64(m.timeout.Milliseconds())
	return mu.Bytes(), nil
}

func (m *msgQuery) UnmarshalBinary(data []byte) error {
	var err error
	var mt byte
	mu := marshalutil.New(data)
	if mt, err = mu.ReadByte(); err != nil {
		return err
	}
	if mt != msgTypeQuery {
		return fmt.Errorf("unexpected message type: %v", mt)
	}
	//
	// m.session
	if m.session, err = hashing.HashValueFromMarshalUtil(mu); err != nil {
		return err
	}
	//
	// m.senderPubKey
	if m.senderPubKey, err = cryptolib.NewPublicKeyFromMarshalUtil(mu); err != nil {
		return err
	}
	//
	// m.subQuery
	subQueryLen, err := mu.ReadUint16()
	if err != nil {
		return err
	}
	m.subQuery = make([]*cryptolib.PublicKey, subQueryLen)
	for i := range m.subQuery {
		if m.subQuery[i], err = cryptolib.NewPublicKeyFromMarshalUtil(mu); err != nil {
			return err
		}
	}
	//
	// m.timeout
	timeoutMs, err := mu.ReadInt64()
	if err != nil {
		return nil
	}
	m.timeout = time.Duration(timeoutMs) * time.Millisecond
	return nil
}

func (m *msgQuery) String() string {
	return fmt.Sprintf("{msgQuery, sender=%v, recipient=%v, |subQuery|=%v}",
		m.Sender().ShortString(),
		m.Recipient().ShortString(),
		len(m.subQuery),
	)
}
