// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package clique

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/samber/lo"

	"github.com/iotaledger/hive.go/core/generics/shrinkingmap"
	"github.com/iotaledger/hive.go/core/logger"
	"github.com/iotaledger/wasp/packages/cryptolib"
	"github.com/iotaledger/wasp/packages/gpa"
	"github.com/iotaledger/wasp/packages/hashing"
	"github.com/iotaledger/wasp/packages/peering"
	"github.com/iotaledger/wasp/packages/peering/clique/cliqueDist"
	"github.com/iotaledger/wasp/packages/util/pipe"
)

const (
	defaultTimeout = 5 * time.Second
	distTimeTick   = 1 * time.Second
)

const (
	msgTypeCliqueChecker byte = iota
)

// The parameter of the callback stands for link[src][dst]error
type Callback func(link map[cryptolib.PublicKeyKey]map[cryptolib.PublicKeyKey]error)

// Checks the peering connectivity between the specified nodes.
// For the name: a clique -- is a fully-connected subgraph.
type Clique interface {
	// The timeout for this operation is taken from the Context deadline.
	Check(ctx context.Context, nodes []*cryptolib.PublicKey, callback Callback)
}

type cliqueImpl struct {
	dist         gpa.GPA
	requestPipe  pipe.Pipe[*checkRequest]
	requests     *shrinkingmap.ShrinkingMap[hashing.HashValue, *checkRequest]
	trustedPipe  pipe.Pipe[[]*peering.TrustedPeer]
	netRecvPipe  pipe.Pipe[*peering.PeerMessageIn]
	netPeeringID peering.PeeringID
	netPeerPubs  map[gpa.NodeID]*cryptolib.PublicKey
	net          peering.NetworkProvider
	log          *logger.Logger
}

var _ Clique = &cliqueImpl{}

func New(
	ctx context.Context,
	nodeIdentity *cryptolib.KeyPair,
	net peering.NetworkProvider,
	tns peering.TrustedNetworkManager,
	log *logger.Logger,
) Clique {
	// there is only one Clique checker per Wasp node, so the identifier is a constant.
	netPeeringID := peering.HashPeeringIDFromBytes([]byte("CliqueChecker"))
	ci := &cliqueImpl{
		dist:         cliqueDist.New(nodeIdentity, time.Now(), log).AsGPA(),
		requestPipe:  pipe.NewInfinitePipe[*checkRequest](),
		requests:     shrinkingmap.New[hashing.HashValue, *checkRequest](),
		trustedPipe:  pipe.NewInfinitePipe[[]*peering.TrustedPeer](),
		netRecvPipe:  pipe.NewInfinitePipe[*peering.PeerMessageIn](),
		netPeeringID: netPeeringID,
		netPeerPubs:  map[gpa.NodeID]*cryptolib.PublicKey{},
		net:          net,
		log:          log,
	}
	netRecvPipeInCh := ci.netRecvPipe.In()
	netAttachID := net.Attach(&netPeeringID, peering.PeerMessageReceiverClique, func(recv *peering.PeerMessageIn) {
		if recv.MsgType != msgTypeCliqueChecker {
			ci.log.Warnf("Unexpected message, type=%v", recv.MsgType)
			return
		}
		netRecvPipeInCh <- recv
	})
	trustedCancel := tns.TrustedPeersListener(func(trustedPeers []*peering.TrustedPeer) {
		ci.trustedPipe.In() <- trustedPeers
	})
	go ci.run(ctx, netAttachID, trustedCancel)
	return ci
}

func (ci *cliqueImpl) Check(ctx context.Context, nodes []*cryptolib.PublicKey, callback Callback) {
	var timeout time.Duration
	now := time.Now()
	if deadline, ok := ctx.Deadline(); ok {
		timeout = deadline.Sub(now)
	} else {
		timeout = defaultTimeout
	}
	ci.requestPipe.In() <- newCheckRequest(ctx, nodes, now, timeout, callback)
}

func (ci *cliqueImpl) run(ctx context.Context, netAttachID interface{}, trustedCancel context.CancelFunc) {
	requestPipeOutCh := ci.requestPipe.Out()
	trustedPipeOutCh := ci.trustedPipe.Out()
	netRecvPipeOutCh := ci.netRecvPipe.Out()
	distTimeTicker := time.NewTicker(distTimeTick)
	for {
		select {
		case req, ok := <-requestPipeOutCh:
			if !ok {
				requestPipeOutCh = nil
				continue
			}
			ci.handleRequest(req)
		case recv, ok := <-trustedPipeOutCh:
			if !ok {
				trustedPipeOutCh = nil
				continue
			}
			ci.handleTrusted(recv)
		case recv, ok := <-netRecvPipeOutCh:
			if !ok {
				netRecvPipeOutCh = nil
				continue
			}
			ci.handleNetMessage(recv)
		case timestamp := <-distTimeTicker.C:
			ci.handleDistTimeTick(timestamp)
		case <-ctx.Done():
			ci.net.Detach(netAttachID)
			trustedCancel()
			return
		}
	}
}

func (ci *cliqueImpl) handleRequest(req *checkRequest) {
	var sessionID hashing.HashValue
	_, err := rand.Read(sessionID[:])
	if err != nil {
		panic(fmt.Errorf("cannot read random data: %w", err))
	}
	for ci.requests.Has(sessionID) {
		// Just to be sure we generated a unique session id. At least at this node.
		sessionID = hashing.HashDataBlake2b(sessionID.Bytes())
	}
	ci.requests.Set(sessionID, req)

	distCallback := func(sessionID hashing.HashValue, linkStates []*cliqueDist.LinkStatus) {
		linkStateToStr := func(ls *cliqueDist.LinkStatus, _ int) string {
			return ls.ShortString()
		}
		ci.log.Debugf("received response for sessionID=%v linkStates=%v", sessionID, lo.Map(linkStates, linkStateToStr))
		if cr, ok := ci.requests.Get(sessionID); ok {
			ci.requests.Delete(sessionID)
			req.callback(cr.collectLinkStates(linkStates))
		}
	}
	ci.sendMessages(ci.dist.Input(cliqueDist.NewInputCheck(sessionID, req.nodes, req.timeout, distCallback)))
}

func (ci *cliqueImpl) handleTrusted(recv []*peering.TrustedPeer) {
	trustedMap := map[gpa.NodeID]*cryptolib.PublicKey{}
	for _, tp := range recv {
		trustedMap[gpa.NodeIDFromPublicKey(tp.PubKey())] = tp.PubKey()
	}
	ci.netPeerPubs = trustedMap
	ci.sendMessages(ci.dist.Input(cliqueDist.NewInputTrustedPeers(trustedMap)))
}

func (ci *cliqueImpl) handleDistTimeTick(timestamp time.Time) {
	ci.sendMessages(ci.dist.Input(cliqueDist.NewInputTimeTick(timestamp)))
	ci.requests.ForEach(func(sessionID hashing.HashValue, cr *checkRequest) bool {
		if cr.ctx.Err() == nil {
			return true
		}
		cr.callback(cr.collectLinkStates([]*cliqueDist.LinkStatus{}))
		ci.requests.Delete(sessionID)
		return true
	})
}

func (ci *cliqueImpl) handleNetMessage(recv *peering.PeerMessageIn) {
	msg, err := ci.dist.UnmarshalMessage(recv.MsgData)
	if err != nil {
		ci.log.Warnf("cannot parse message: %v", err)
		return
	}
	msg.SetSender(gpa.NodeIDFromPublicKey(recv.SenderPubKey))
	outMsgs := ci.dist.Message(msg) // Output is handled via callbacks in this case.
	ci.sendMessages(outMsgs)
}

func (ci *cliqueImpl) sendMessages(outMsgs gpa.OutMessages) {
	if outMsgs == nil {
		return
	}
	outMsgs.MustIterate(func(m gpa.Message) {
		msgData, err := m.MarshalBinary()
		if err != nil {
			ci.log.Warnf("Failed to send a message: %v", err)
			return
		}
		pm := &peering.PeerMessageData{
			PeeringID:   ci.netPeeringID,
			MsgReceiver: peering.PeerMessageReceiverClique,
			MsgType:     msgTypeCliqueChecker,
			MsgData:     msgData,
		}
		if pubKey, ok := ci.netPeerPubs[m.Recipient()]; ok {
			ci.net.SendMsgByPubKey(pubKey, pm)
		} else {
			ci.log.Warnf("Dropping out message, pub key not know: %v", m.Recipient())
		}
	})
}

////////////////////////////////////////////////////////////////////////////////
// checkRequest

type checkRequest struct {
	ctx      context.Context
	nodes    []*cryptolib.PublicKey
	started  time.Time
	timeout  time.Duration
	callback Callback
}

func newCheckRequest(
	ctx context.Context,
	nodes []*cryptolib.PublicKey,
	started time.Time,
	timeout time.Duration,
	callback Callback,
) *checkRequest {
	return &checkRequest{ctx: ctx, nodes: nodes, started: started, timeout: timeout, callback: callback}
}

func (cr *checkRequest) collectLinkStates(linkStates []*cliqueDist.LinkStatus) map[cryptolib.PublicKeyKey]map[cryptolib.PublicKeyKey]error {
	res := map[cryptolib.PublicKeyKey]map[cryptolib.PublicKeyKey]error{}
	//
	// Collect received info.
	for _, ls := range linkStates {
		src := ls.SrcPubKey().AsKey()
		dst := ls.DstPubKey().AsKey()
		resForSrc, ok := res[src]
		if !ok {
			resForSrc = map[cryptolib.PublicKeyKey]error{}
			res[src] = resForSrc
		}
		resForDst, ok := res[dst]
		if !ok {
			resForDst = map[cryptolib.PublicKeyKey]error{}
			res[dst] = resForDst
		}
		if ls.FailReason() == nil { // Link is alive both ways.
			resForDst[src] = nil
			resForSrc[dst] = nil
		}
		resForSrc[dst] = ls.FailReason()
	}
	//
	// Fill with unknown links.
	for _, src := range cr.nodes {
		for _, dst := range cr.nodes {
			resForSrc, hadSrc := res[src.AsKey()]
			if !hadSrc {
				resForSrc = map[cryptolib.PublicKeyKey]error{}
				res[src.AsKey()] = resForSrc
			}
			if src.Equals(dst) && hadSrc {
				resForSrc[dst.AsKey()] = nil // OK, we have something from it.
				continue
			}
			if _, ok := resForSrc[dst.AsKey()]; !ok {
				resForSrc[dst.AsKey()] = fmt.Errorf("unknown")
			}
		}
	}
	return res
}
