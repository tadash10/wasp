// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package cliqueDist

import (
	"fmt"

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
	panic("implement") // TODO: Implement.
}

func (m *msgResponse) UnmarshalBinary(data []byte) error {
	panic("implement") // TODO: Implement.
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
