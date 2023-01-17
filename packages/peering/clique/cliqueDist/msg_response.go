// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package cliqueDist

import (
	"fmt"

	"github.com/iotaledger/hive.go/core/marshalutil"
	"github.com/iotaledger/wasp/packages/gpa"
	"github.com/iotaledger/wasp/packages/hashing"
)

type msgResponse struct {
	gpa.BasicMessage
	session      hashing.HashValue
	response     *LinkStatus   // Response from the sender.
	subResponses []*LinkStatus // Responses the sender has collected.
}

func newMsgResponse(recipient gpa.NodeID, session hashing.HashValue, response *LinkStatus, subResponses []*LinkStatus) gpa.Message {
	return &msgResponse{
		BasicMessage: gpa.NewBasicMessage(recipient),
		session:      session,
		response:     response,
		subResponses: subResponses,
	}
}

func (m *msgResponse) MarshalBinary() ([]byte, error) {
	mu := marshalutil.New()
	mu.WriteByte(msgTypeResponse)
	mu.Write(m.session)
	mu.Write(m.response)
	mu.WriteUint16(uint16(len(m.subResponses)))
	for _, sr := range m.subResponses {
		mu.Write(sr)
	}
	return mu.Bytes(), nil
}

func (m *msgResponse) UnmarshalBinary(data []byte) error {
	var err error
	var mt byte
	mu := marshalutil.New(data)
	//
	// MsgType
	if mt, err = mu.ReadByte(); err != nil {
		return fmt.Errorf("cannot unmarshal messageType: %w", err)
	}
	if mt != msgTypeResponse {
		return fmt.Errorf("unexpected message type: %v", mt)
	}
	//
	// m.session
	if m.session, err = hashing.HashValueFromMarshalUtil(mu); err != nil {
		return fmt.Errorf("cannot unmarshal session: %w", err)
	}
	//
	// m.response
	if m.response, err = NewLinkStatusFromMarshalUtil(mu); err != nil {
		return fmt.Errorf("cannot unmarshal response: %w", err)
	}
	//
	// m.subResponses
	var subResponsesLen uint16
	if subResponsesLen, err = mu.ReadUint16(); err != nil {
		return fmt.Errorf("cannot unmarshal subResponse count: %w", err)
	}
	m.subResponses = make([]*LinkStatus, subResponsesLen)
	for i := range m.subResponses {
		if m.subResponses[i], err = NewLinkStatusFromMarshalUtil(mu); err != nil {
			return fmt.Errorf("cannot unmarshal subResponse[%v]: %w", i, err)
		}
	}
	return nil
}

func (m *msgResponse) String() string {
	return fmt.Sprintf(
		"{msgResponse, sender=%v, recipient=%v, response=%v, |subResponses|=%v}",
		m.Sender().ShortString(),
		m.Recipient().ShortString(),
		m.response.ShortString(),
		len(m.subResponses),
	)
}
