// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package v1

import (
	"time"

	"github.com/iotaledger/wasp/packages/webapi/v1/admapi"
	"github.com/iotaledger/wasp/packages/webapi/v1/evm"
	"github.com/iotaledger/wasp/packages/webapi/v1/info"
	"github.com/iotaledger/wasp/packages/webapi/v1/reqstatus"
	"github.com/iotaledger/wasp/packages/webapi/v1/request"
	"github.com/iotaledger/wasp/packages/webapi/v1/state"

	"github.com/pangpanglabs/echoswagger/v2"

	loggerpkg "github.com/iotaledger/hive.go/core/logger"
	"github.com/iotaledger/wasp/packages/authentication"
	"github.com/iotaledger/wasp/packages/chain/chainutil"
	"github.com/iotaledger/wasp/packages/chains"
	"github.com/iotaledger/wasp/packages/dkg"
	metricspkg "github.com/iotaledger/wasp/packages/metrics"
	"github.com/iotaledger/wasp/packages/peering"
	"github.com/iotaledger/wasp/packages/registry"
	"github.com/iotaledger/wasp/packages/wal"
)

var log *loggerpkg.Logger

func Init(
	logger *loggerpkg.Logger,
	server echoswagger.ApiRoot,
	network peering.NetworkProvider,
	tnm peering.TrustedNetworkManager,
	registryProvider registry.Provider,
	chainsProvider chains.Provider,
	nodeProvider dkg.NodeProvider,
	shutdown admapi.ShutdownFunc,
	metrics *metricspkg.Metrics,
	w *wal.WAL,
	authConfig authentication.AuthConfiguration,
	nodeOwnerAddresses []string,
	apiCacheTTL time.Duration,
	publisherPort int,
) {
	log = logger

	pub := server.Group("public", "v1").SetDescription("Public endpoints")
	addWebSocketEndpoint(pub, log)

	info.AddEndpoints(pub, network, publisherPort)
	reqstatus.AddEndpoints(pub, chainsProvider.ChainProvider())
	state.AddEndpoints(pub, chainsProvider)
	evm.AddEndpoints(pub, chainsProvider, network.Self().PubKey)
	request.AddEndpoints(
		pub,
		chainsProvider.ChainProvider(),
		chainutil.GetAccountBalance,
		chainutil.HasRequestBeenProcessed,
		chainutil.CheckNonce,
		network.Self().PubKey(),
		apiCacheTTL,
		log,
	)

	adm := server.Group("admin", "v1").SetDescription("Admin endpoints")

	admapi.AddEndpoints(
		logger.Named("webapi/adm"),
		adm,
		network,
		tnm,
		registryProvider,
		chainsProvider,
		nodeProvider,
		shutdown,
		metrics,
		w,
		authConfig,
		nodeOwnerAddresses,
	)
	log.Infof("added web api endpoints")
}
