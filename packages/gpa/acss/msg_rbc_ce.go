// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package acss

import (
	"io"

	"go.dedis.ch/kyber/v3/suites"

	"github.com/iotaledger/wasp/packages/gpa"
	"github.com/iotaledger/wasp/packages/util/rwutil"
)

// This message is used as a payload of the RBC:
//
// > RBC(C||E)
type msgRBCCEPayload struct {
	gpa.BasicMessage
	suite suites.Suite
	data  []byte
}

func (msg *msgRBCCEPayload) MarshalBinary() ([]byte, error) {
	return rwutil.MarshalBinary(msg)
}

func (msg *msgRBCCEPayload) UnmarshalBinary(data []byte) error {
	return rwutil.UnmarshalBinary(data, msg)
}

func (msg *msgRBCCEPayload) Read(r io.Reader) error {
	rr := rwutil.NewReader(r)
	rr.ReadMessageTypeAndVerify(msgTypeRBCCEPayload)
	msg.data = rr.ReadBytes()
	return rr.Err
}

func (msg *msgRBCCEPayload) Write(w io.Writer) error {
	ww := rwutil.NewWriter(w)
	ww.WriteMessageType(msgTypeRBCCEPayload)
	ww.WriteBytes(msg.data)
	return ww.Err
}
