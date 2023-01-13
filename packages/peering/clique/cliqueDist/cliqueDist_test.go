// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package cliqueDist_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/iotaledger/wasp/packages/cryptolib"
	"github.com/iotaledger/wasp/packages/gpa"
	"github.com/iotaledger/wasp/packages/hashing"
	"github.com/iotaledger/wasp/packages/peering/clique/cliqueDist"
	"github.com/iotaledger/wasp/packages/testutil/testlogger"
)

func TestBasic(t *testing.T) {
	log := testlogger.NewLogger(t)
	nodeCount := 4
	nodeKeys := make([]*cryptolib.KeyPair, nodeCount)
	nodePubs := make([]*cryptolib.PublicKey, nodeCount)
	nodeIDs := make([]gpa.NodeID, nodeCount)
	for i := range nodeKeys {
		nodeKeys[i] = cryptolib.NewKeyPair()
		nodePubs[i] = nodeKeys[i].GetPublicKey()
		nodeIDs[i] = gpa.NodeIDFromPublicKey(nodePubs[i])
	}
	now := time.Now()
	nodes := map[gpa.NodeID]gpa.GPA{}
	for i, nodeID := range nodeIDs {
		nodes[nodeID] = cliqueDist.New(nodeKeys[i], now, log.Named(fmt.Sprintf("N%v", i)))
	}
	tc := gpa.NewTestContext(nodes)

	done := false
	tc.WithInput(nodeIDs[0], cliqueDist.NewInputCheck(hashing.HashDataBlake2b([]byte{0}), nodePubs, time.Second, func(sessionID hashing.HashValue, linkStates []*cliqueDist.LinkStatus) {
		for _, ls := range linkStates {
			require.True(t, ls.Validate(sessionID))
			require.NoError(t, ls.FailReason())
		}
		done = true
	}))
	tc.RunAll()
	tc.PrintAllStatusStrings("All done", t.Logf)
	require.True(t, done)
}

func TestSilent(t *testing.T) {
	log := testlogger.NewLogger(t)
	nodeCount := 4
	nodeLast := nodeCount - 1
	nodeKeys := make([]*cryptolib.KeyPair, nodeCount)
	nodePubs := make([]*cryptolib.PublicKey, nodeCount)
	nodeIDs := make([]gpa.NodeID, nodeCount)
	for i := range nodeKeys {
		nodeKeys[i] = cryptolib.NewKeyPair()
		nodePubs[i] = nodeKeys[i].GetPublicKey()
		nodeIDs[i] = gpa.NodeIDFromPublicKey(nodePubs[i])
		t.Logf("NodeID[%v]=%v", i, nodeIDs[i].ShortString())
	}
	now := time.Now()
	nodes := map[gpa.NodeID]gpa.GPA{}
	for i, nodeID := range nodeIDs {
		if i == nodeLast {
			nodes[nodeID] = gpa.MakeTestSilentNode()
		} else {
			nodes[nodeID] = cliqueDist.New(nodeKeys[i], now, log.Named(fmt.Sprintf("N%v", i)))
		}
	}
	tc := gpa.NewTestContext(nodes)

	done := false
	tc.WithInput(nodeIDs[0], cliqueDist.NewInputCheck(hashing.HashDataBlake2b([]byte{0}), nodePubs, 5*time.Second, func(sessionID hashing.HashValue, linkStates []*cliqueDist.LinkStatus) {
		for _, ls := range linkStates {
			require.True(t, ls.Validate(sessionID))
			t.Logf("LS: %v", ls.ShortString())
			if ls.DstPubKey().Equals(nodePubs[nodeLast]) {
				require.ErrorContains(t, ls.FailReason(), "timeout")
			} else {
				require.NoError(t, ls.FailReason())
			}
		}
		done = true
	}))
	tc.RunAll()
	now = now.Add(3 * time.Second)
	for nodeID := range nodes {
		tc.WithInput(nodeID, cliqueDist.NewInputTimeTick(now))
	}
	tc.RunAll()
	now = now.Add(3 * time.Second)
	for nodeID := range nodes {
		tc.WithInput(nodeID, cliqueDist.NewInputTimeTick(now))
	}
	tc.RunAll()
	tc.PrintAllStatusStrings("All done", t.Logf)
	require.True(t, done)
}
