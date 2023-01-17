// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package cliqueDist

import (
	"fmt"
	"time"

	"github.com/samber/lo"

	"github.com/iotaledger/hive.go/core/logger"
	"github.com/iotaledger/wasp/packages/cryptolib"
	"github.com/iotaledger/wasp/packages/gpa"
	"github.com/iotaledger/wasp/packages/hashing"
)

type session struct {
	id              hashing.HashValue    // SessionID
	me              gpa.NodeID           // This node.
	initiatorPubKey *cryptolib.PublicKey // nil on the initiator node.
	initiatorNodeID gpa.NodeID           // me, on the initiator node.
	started         time.Time
	timeout         time.Duration
	callback        Callback
	nodePubs        []*cryptolib.PublicKey
	nodeIDs         map[gpa.NodeID]*cryptolib.PublicKey
	required        map[gpa.NodeID]interface{}
	collected       map[gpa.NodeID]map[gpa.NodeID]*LinkStatus // responses[SRC][DST]
	log             *logger.Logger
}

func newSession(
	id hashing.HashValue,
	me gpa.NodeID,
	myKeyPair *cryptolib.KeyPair,
	now time.Time,
	timeout time.Duration,
	initiatorPubKey *cryptolib.PublicKey,
	initiatorNodeID gpa.NodeID,
	callback Callback,
	nodePubs []*cryptolib.PublicKey,
	trusted map[gpa.NodeID]*cryptolib.PublicKey,
	log *logger.Logger,
) *session {
	nodeIDs := map[gpa.NodeID]*cryptolib.PublicKey{}
	for _, nodePub := range nodePubs {
		nodeIDs[gpa.NodeIDFromPublicKey(nodePub)] = nodePub
	}
	required := map[gpa.NodeID]interface{}{}
	collected := map[gpa.NodeID]map[gpa.NodeID]*LinkStatus{
		me: {},
	}
	for nodeID, nodePub := range nodeIDs {
		if nodeID.Equals(me) || nodeID.Equals(initiatorNodeID) {
			continue
		}
		if _, ok := trusted[nodeID]; !ok {
			collected[me][nodeID] = newLinkStatusFailed(id, myKeyPair, nodePub, fmt.Errorf("non-trusted"))
			continue
		}
		required[nodeID] = nil
	}
	if timeout > maxStepTimeout {
		timeout = maxStepTimeout
	}
	return &session{
		id:              id,
		me:              me,
		initiatorPubKey: initiatorPubKey,
		initiatorNodeID: initiatorNodeID,
		started:         now,
		timeout:         timeout,
		callback:        callback,
		nodePubs:        nodePubs,
		nodeIDs:         nodeIDs,
		required:        required,
		collected:       collected,
		log:             log,
	}
}

func (s *session) ID() hashing.HashValue {
	return s.id
}

func (s *session) InitiatorNodeID() gpa.NodeID {
	return s.initiatorNodeID
}

func (s *session) InitiatorPubKey() *cryptolib.PublicKey {
	return s.initiatorPubKey
}

func (s *session) StatusString() string {
	return fmt.Sprintf("{session, sub=%v, |required|=%v}", s.initiatorPubKey != nil, len(s.required))
}

func (s *session) MakeReqMsgs() gpa.OutMessages {
	msgs := gpa.NoMessages()
	sub := lo.If(s.initiatorPubKey != nil, []*cryptolib.PublicKey{}).Else(s.nodePubs)
	for nodeID := range s.required {
		msgs.Add(newMsgQuery(nodeID, s.id, s.nodeIDs[s.me], sub, s.timeout/2))
	}
	return msgs
}

func (s *session) AddResponse(msg *msgResponse) {
	senderPub, ok := s.nodeIDs[msg.Sender()]
	if !ok {
		s.log.Warnf("unexpected sender")
		return
	}
	if _, ok := s.required[msg.Sender()]; !ok {
		s.log.Warnf("we haven't asked the response")
		return
	}
	//
	// Validate and record the direct answer (me -> responseSender).
	if !s.nodeIDs[s.me].Equals(msg.response.srcPubKey) {
		s.log.Warnf("src pub key invalid.")
		return
	}
	if !senderPub.Equals(msg.response.dstPubKey) {
		s.log.Warnf("dst pub key invalid, expect=%v, got=%v", senderPub, msg.response.dstPubKey)
		return
	}
	if !msg.response.Validate(s.id) {
		s.log.Warnf("response invalid")
		return
	}
	delete(s.required, msg.Sender())
	s.addLinkStatus(msg.response, s.me, msg.Sender())
	//
	// Validate and record the sub-answers.
	for _, subRes := range msg.subResponses {
		if !senderPub.Equals(subRes.srcPubKey) {
			continue
		}
		subDstNodeID := gpa.NodeIDFromPublicKey(subRes.dstPubKey)
		if _, ok := s.nodeIDs[subDstNodeID]; !ok {
			continue
		}
		if !subRes.Validate(s.id) {
			continue
		}
		s.addLinkStatus(subRes, msg.Sender(), subDstNodeID)
	}
}

func (s *session) addLinkStatus(status *LinkStatus, src, dst gpa.NodeID) {
	sub, ok := s.collected[src]
	if !ok || sub == nil {
		sub = map[gpa.NodeID]*LinkStatus{}
		s.collected[src] = sub
	}
	if ls, ok := sub[dst]; ok && ls != nil {
		return // Already have it.
	}
	sub[dst] = status
}

func (s *session) HaveAllResponses() (bool, []*LinkStatus) {
	if len(s.required) > 0 {
		return false, nil
	}
	res := []*LinkStatus{}
	for _, srcLS := range s.collected {
		for _, ls := range srcLS {
			res = append(res, ls)
		}
	}
	if s.callback != nil {
		s.callback(s.id, res)
		return true, nil
	}
	return true, res
}

func (s *session) HaveTimeout(now time.Time, myKeyPair *cryptolib.KeyPair) (bool, []*LinkStatus) {
	if s.started.Add(s.timeout).After(now) {
		// No timeout yet.
		return false, nil
	}
	res := []*LinkStatus{}
	for dst := range s.required {
		res = append(res, newLinkStatusFailed(s.id, myKeyPair, s.nodeIDs[dst], fmt.Errorf("timeout")))
	}
	for _, srcLS := range s.collected {
		for _, ls := range srcLS {
			res = append(res, ls)
		}
	}
	if s.callback != nil {
		s.callback(s.id, res)
		return true, nil
	}
	return true, res
}
