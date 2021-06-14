package request

import (
	"net/http"
	"testing"

	"github.com/iotaledger/wasp/packages/chain"
	"github.com/iotaledger/wasp/packages/coretypes"
	"github.com/iotaledger/wasp/packages/coretypes/request"
	"github.com/iotaledger/wasp/packages/coretypes/requestargs"
	"github.com/iotaledger/wasp/packages/kv/dict"
	"github.com/iotaledger/wasp/packages/testutil/testchain"
	"github.com/iotaledger/wasp/packages/testutil/testlogger"
	"github.com/iotaledger/wasp/packages/webapi/routes"
	"github.com/iotaledger/wasp/packages/webapi/testutil"
	"github.com/stretchr/testify/require"
)

func createMockedGetChain(t *testing.T) getChainFn {
	return func(chainID *coretypes.ChainID) chain.ChainCore {
		return testchain.NewMockedChainCore(t, *chainID, testlogger.NewLogger(t))
	}
}

const foo = "foo"

func dummyOffledgerRequest() *request.RequestOffLedger {
	contract := coretypes.Hn("somecontract")
	entrypoint := coretypes.Hn("someentrypoint")
	args := requestargs.New(
		dict.Dict{foo: []byte("bar")},
	)
	return request.NewRequestOffLedger(contract, entrypoint, args)
}

func TestNewRequest(t *testing.T) {
	instance := &offLedgerReqAPI{
		getChain: createMockedGetChain(t),
	}

	var res error

	testutil.CallWebAPIRequestHandler(
		t,
		instance.handleNewRequest,
		http.MethodPost,
		routes.NewRequest(":chainID"),
		map[string]string{"chainID": coretypes.RandomChainID().Base58()},
		map[string]string{"request": dummyOffledgerRequest().Base64()},
		nil,
		http.StatusAccepted,
	)
	require.NoError(t, res)
}
