package vm0poc

import (
	"testing"

	"github.com/iotaledger/wasp/packages/solo"
	"github.com/iotaledger/wasp/packages/util"
	"github.com/stretchr/testify/require"
)

func newChain(t *testing.T) *solo.Chain {
	env := solo.New(t)
	ch, _, _ := env.NewChainExt(nil, 0, "ch1", solo.InitChainOptions{
		VMRunner:         NewVMRunner(),
		BypassStardustVM: true,
	})
	return ch
}

func TestBasic(t *testing.T) {
	t.Run("create chain", func(t *testing.T) {
		ch := newChain(t)
		require.False(t, ch.RawState().MustHas(CounterStateVar))
	})
	t.Run("send 1 request", func(t *testing.T) {
		ch := newChain(t)
		req := solo.NewCallParams("dummy", "dummy", ParamDeltaInt64, int64(10))
		_, err := ch.PostRequestOffLedger(req, nil)
		require.NoError(t, err)

		require.True(t, ch.RawState().MustHas(CounterStateVar))
		v := ch.RawState().MustGet(CounterStateVar)
		val, err := util.Int64From8Bytes(v)
		require.NoError(t, err)
		require.EqualValues(t, int64(10), val)
	})
	t.Run("send 5 requests", func(t *testing.T) {
		ch := newChain(t)
		req := solo.NewCallParams("dummy", "dummy", ParamDeltaInt64, int64(1))
		for i := 0; i < 5; i++ {
			_, err := ch.PostRequestOffLedger(req, nil)
			require.NoError(t, err)
			require.True(t, ch.RawState().MustHas(CounterStateVar))
			v := ch.RawState().MustGet(CounterStateVar)
			val, err := util.Int64From8Bytes(v)
			require.NoError(t, err)
			require.EqualValues(t, int64(i+1), val)
		}
	})
	t.Run("send 2 requests", func(t *testing.T) {
		ch := newChain(t)
		req := solo.NewCallParams("dummy", "dummy", ParamDeltaInt64, int64(10))
		_, err := ch.PostRequestOffLedger(req, nil)
		require.NoError(t, err)

		require.True(t, ch.RawState().MustHas(CounterStateVar))
		v := ch.RawState().MustGet(CounterStateVar)
		val, err := util.Int64From8Bytes(v)
		require.NoError(t, err)
		require.EqualValues(t, int64(10), val)

		req = solo.NewCallParams("dummy", "dummy", ParamDeltaInt64, int64(-5))
		_, err = ch.PostRequestOffLedger(req, nil)
		require.NoError(t, err)

		require.True(t, ch.RawState().MustHas(CounterStateVar))
		v = ch.RawState().MustGet(CounterStateVar)
		val, err = util.Int64From8Bytes(v)
		require.NoError(t, err)
		require.EqualValues(t, int64(5), val)
	})
}
