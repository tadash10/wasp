// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package cliqueDist

import (
	"fmt"
	"time"

	"github.com/iotaledger/hive.go/core/generics/shrinkingmap"
	"github.com/iotaledger/hive.go/core/logger"
	"github.com/iotaledger/wasp/packages/cryptolib"
	"github.com/iotaledger/wasp/packages/gpa"
	"github.com/iotaledger/wasp/packages/hashing"
)

const (
	maxStepTimeout = 30 * time.Second
)

type Callback func(sessionID hashing.HashValue, linkStates []*LinkStatus)

type cliqueDist struct {
	me        gpa.NodeID
	now       time.Time
	myKeyPair *cryptolib.KeyPair
	sessions  *shrinkingmap.ShrinkingMap[hashing.HashValue, *session]
	log       *logger.Logger
}

var _ gpa.GPA = &cliqueDist{}

func New(myKeyPair *cryptolib.KeyPair, now time.Time, log *logger.Logger) gpa.GPA {
	return &cliqueDist{
		me:        gpa.NodeIDFromPublicKey(myKeyPair.GetPublicKey()),
		now:       now,
		myKeyPair: myKeyPair,
		sessions:  shrinkingmap.New[hashing.HashValue, *session](),
		log:       log,
	}
}

func (cd *cliqueDist) Input(input gpa.Input) gpa.OutMessages {
	switch input := input.(type) {
	case *inputCheck:
		return cd.handleInputCheck(input)
	case *inputTimeTick:
		return cd.handleInputTimeTick(input)
	default:
		panic(fmt.Errorf("unexpected input %T: %+v", input, input))
	}
}

func (cd *cliqueDist) Message(msg gpa.Message) gpa.OutMessages {
	switch msg := msg.(type) {
	case *msgQuery:
		return cd.handleMsgQuery(msg)
	case *msgResponse:
		return cd.handleMsgResponse(msg)
	default:
		panic(fmt.Errorf("unexpected message %T: %+v", msg, msg))
	}
}

func (cd *cliqueDist) Output() gpa.Output {
	return nil // No output.
}

func (cd *cliqueDist) StatusString() string {
	sessionsStr := ""
	cd.sessions.ForEach(func(_ hashing.HashValue, s *session) bool {
		sessionsStr += s.StatusString()
		return true
	})
	return fmt.Sprintf("{Clique, sessions=%v}", sessionsStr)
}

func (cd *cliqueDist) handleInputCheck(input *inputCheck) gpa.OutMessages {
	sessionID := input.sessionID
	for {
		if !cd.sessions.Has(sessionID) {
			break
		}
		sessionID = hashing.HashDataBlake2b(sessionID[:])
	}
	s := newSession(sessionID, 1, cd.now, input.timeout, cd.myKeyPair.GetPublicKey(), input.callback, cd.me, cd.me, input.nodePubs, cd.log)
	cd.sessions.Set(sessionID, s)
	return s.MakeReqMsgs()
}

func (cd *cliqueDist) handleInputTimeTick(input *inputTimeTick) gpa.OutMessages {
	cd.now = input.timestamp
	msgs := gpa.NoMessages()
	cd.sessions.ForEach(func(_ hashing.HashValue, s *session) bool {
		if timeout, subRes := s.HaveTimeout(cd.now, cd.myKeyPair); timeout {
			cd.sessions.Delete(s.id)
			if subRes != nil {
				msgs.Add(newMsgResponse(
					s.FromNodeID(),
					s.ID(),
					NewLinkStatusOK(s.id, s.FromPubKey(), cd.myKeyPair),
					subRes,
				))
			}
		}
		return true
	})
	return msgs
}

func (cd *cliqueDist) handleMsgQuery(msg *msgQuery) gpa.OutMessages {
	cd.log.Debugf("handleMsgQuery, err=%+v", msg)
	if err := msg.Validate(); err != nil {
		cd.log.Warnf("failed to validate the received msgQuery: %v", err)
		return nil
	}
	if msg.ttl == 0 || len(msg.subQuery) == 0 {
		// Always respond to leaf queries.
		return gpa.NoMessages().Add(newMsgResponse(
			msg.Sender(),
			msg.session,
			NewLinkStatusOK(msg.session, msg.senderPubKey, cd.myKeyPair),
			[]*LinkStatus{},
		))
	}
	if cd.sessions.Has(msg.session) {
		// Duplicate request, ignore it.
		return nil
	}
	msgS := newSession(msg.session, 0, cd.now, msg.timeout, msg.senderPubKey, nil, cd.me, msg.Sender(), msg.subQuery, cd.log)
	cd.sessions.Set(msg.session, msgS)
	return msgS.MakeReqMsgs()
}

func (cd *cliqueDist) handleMsgResponse(msg *msgResponse) gpa.OutMessages {
	cd.log.Debugf("handleMsgResponse, err=%+v", msg)
	s, ok := cd.sessions.Get(msg.session)
	if !ok {
		// Ignore random or outdated responses.
		cd.log.Warnf("session not found, dropping the message.")
		return nil
	}
	s.AddResponse(msg)
	if have, subRes := s.HaveAllResponses(); have {
		cd.sessions.Delete(s.id)
		if subRes != nil {
			return gpa.NoMessages().Add(newMsgResponse(
				s.FromNodeID(),
				s.ID(),
				NewLinkStatusOK(s.id, s.FromPubKey(), cd.myKeyPair),
				subRes,
			))
		}
	}
	return nil
}
