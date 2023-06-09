// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package acss

import (
	"io"

	"github.com/iotaledger/wasp/packages/gpa"
	"github.com/iotaledger/wasp/packages/util/rwutil"
)

type msgVoteKind byte

const (
	msgVoteOK msgVoteKind = iota
	msgVoteREADY
)

// This message is used a vote for the "Bracha-style totality" agreement.
type msgVote struct {
	gpa.BasicMessage
	kind msgVoteKind
}

var _ gpa.Message = new(msgVote)

func (msg *msgVote) MarshalBinary() ([]byte, error) {
	return rwutil.MarshalBinary(msg)
}

func (msg *msgVote) UnmarshalBinary(data []byte) error {
	return rwutil.UnmarshalBinary(data, msg)
}

func (msg *msgVote) Read(r io.Reader) error {
	rr := rwutil.NewReader(r)
	rr.ReadMessageTypeAndVerify(msgTypeVote)
	msg.kind = msgVoteKind(rr.ReadByte())
	return rr.Err
}

func (msg *msgVote) Write(w io.Writer) error {
	ww := rwutil.NewWriter(w)
	ww.WriteMessageType(msgTypeVote)
	ww.WriteByte(byte(msg.kind))
	return ww.Err
}
