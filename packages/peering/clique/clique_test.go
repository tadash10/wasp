// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package clique_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/iotaledger/hive.go/core/logger"
	"github.com/iotaledger/wasp/packages/cryptolib"
	"github.com/iotaledger/wasp/packages/peering"
	"github.com/iotaledger/wasp/packages/peering/clique"
	"github.com/iotaledger/wasp/packages/testutil"
	"github.com/iotaledger/wasp/packages/testutil/testlogger"
	"github.com/iotaledger/wasp/packages/testutil/testpeers"
)

func TestBasic(t *testing.T) {
	n := 4
	ctx, ctxCancel := context.WithCancel(context.Background())
	log := testlogger.NewLogger(t)
	defer log.Sync()
	defer ctxCancel()

	peerNetIDs, peerIdentities := testpeers.SetupKeys(uint16(n))
	peerPubKeys := testpeers.PublicKeys(peerIdentities)
	peeringNetwork := testutil.NewPeeringNetwork(
		peerNetIDs, peerIdentities, 10000,
		testutil.NewPeeringNetReliable(log),
		testlogger.WithLevel(log, logger.LevelDebug, false),
	)
	networkProviders := peeringNetwork.NetworkProviders()
	defer peeringNetwork.Close()

	tns := make([]peering.TrustedNetworkManager, n)
	for i := range tns {
		tns[i] = testutil.NewTrustedNetworkManager()
		for _, peerPub := range peerPubKeys {
			tns[i].TrustPeer(peerPub, "some net id")
		}
	}

	cliqueCheckers := make([]clique.Clique, len(peerIdentities))
	for i := range cliqueCheckers {
		cliqueCheckers[i] = clique.New(ctx, peerIdentities[i], networkProviders[i], tns[i], log.Named(fmt.Sprintf("N#%v", i)))
	}
	time.Sleep(20 * time.Millisecond) // Wait until trusted peers are received.

	t.Logf("checker=%v", peerPubKeys[0])
	doneCh := make(chan bool)
	cliqueCheckers[0].Check(ctx, peerPubKeys, func(status map[cryptolib.PublicKeyKey]map[cryptolib.PublicKeyKey]error) {
		for _, src := range peerPubKeys {
			for _, dst := range peerPubKeys {
				t.Logf("Link src=%v, dst=%v: %v", src, dst, status[src.AsKey()][dst.AsKey()])
				require.NoError(t, status[src.AsKey()][dst.AsKey()], "src=%v, dst=%v should not fail", src.String(), dst.String())
			}
		}
		doneCh <- true
	})
	select {
	case <-doneCh:
		// OK
	case <-time.After(10 * time.Second):
		t.Error("timeout waiting for result")
	}
}
