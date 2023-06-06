package blocklog

import (
	"bytes"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/iotaledger/wasp/packages/cryptolib"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/kv"
	"github.com/iotaledger/wasp/packages/kv/collections"
	"github.com/iotaledger/wasp/packages/kv/dict"
	"github.com/iotaledger/wasp/packages/vm/gas"
)

func TestSerdeRequestReceipt(t *testing.T) {
	nonce := uint64(time.Now().UnixNano())
	req := isc.NewOffLedgerRequest(isc.RandomChainID(), isc.Hn("0"), isc.Hn("0"), nil, nonce, gas.LimitsDefault.MaxGasPerRequest)
	signedReq := req.Sign(cryptolib.NewKeyPair())
	rec := &RequestReceipt{
		Request: signedReq,
	}
	forward := rec.Bytes()
	back, err := RequestReceiptFromBytes(forward)
	require.NoError(t, err)
	require.EqualValues(t, forward, back.Bytes())
}

func createEventLookupKeys(blocks uint32) []byte {
	keys := make([]byte, 0)

	for blockIndex := uint32(0); blockIndex < blocks; blockIndex++ {
		for reqIndex := uint16(0); reqIndex < 3; reqIndex++ {
			key := NewEventLookupKey(blockIndex, reqIndex, 0).Bytes()

			keys = append(keys, key...)
		}
	}

	return keys
}

func readEventLookupKeys(partition kv.KVStore, contract kv.Key) ([]*EventLookupKey, error) {
	eventLUT := collections.NewMap(partition, prefixSmartContractEventsLookup)
	keyBytes := eventLUT.GetAt([]byte(contract))
	buff := bytes.NewBuffer(keyBytes)

	keys := make([]*EventLookupKey, 0)

	for {
		key, err := EventLookupKeyFromBytes(buff)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return keys, nil
			}

			return nil, err
		}

		keys = append(keys, key)
	}
}

func containsBlockIndex(blockIndex uint32, keys []*EventLookupKey) bool {
	for _, key := range keys {
		if key.BlockIndex() == blockIndex {
			return true
		}
	}

	return false
}

func validatePrunedEventLookupBlock(t *testing.T, partition kv.KVStore, contract kv.Key, prunedBlockIndex uint32, lastBlockIndex uint32) {
	contractKeys, err := readEventLookupKeys(partition, contract)
	require.NoError(t, err)

	for blockIndex := uint32(0); blockIndex < lastBlockIndex; blockIndex++ {
		if blockIndex == prunedBlockIndex {
			require.False(t, containsBlockIndex(blockIndex, contractKeys))
		} else {
			require.True(t, containsBlockIndex(blockIndex, contractKeys))
		}
	}
}

func TestPruneEventLookupTable(t *testing.T) {
	const maxBlocks = 4
	const blockToPrune = 1

	contract0 := kv.Key("0")
	contract1 := kv.Key("1")

	d := dict.Dict{}

	eventLUT := collections.NewMap(d, prefixSmartContractEventsLookup)
	eventLUT.SetAt([]byte(contract0), createEventLookupKeys(maxBlocks))
	eventLUT.SetAt([]byte(contract1), createEventLookupKeys(maxBlocks))

	require.NotPanics(t, func() {
		pruneEventLookupByBlockIndex(d, 1)
	})

	validatePrunedEventLookupBlock(t, d, contract0, blockToPrune, maxBlocks)
	validatePrunedEventLookupBlock(t, d, contract1, blockToPrune, maxBlocks)
}

func createRequestLookupKeys(blocks uint32) []byte {
	keys := make(RequestLookupKeyList, 0)

	for blockIndex := uint32(0); blockIndex < blocks; blockIndex++ {
		for reqIndex := uint16(0); reqIndex < 3; reqIndex++ {
			key := NewRequestLookupKey(blockIndex, reqIndex)

			keys = append(keys, key)
		}
	}

	return keys.Bytes()
}

func validatePrunedRequestIndexLookupBlock(t *testing.T, partition kv.KVStore, contract kv.Key, prunedBlockIndex uint32) {
	requestLookup := collections.NewMap(partition, prefixRequestLookupIndex)
	requestKeys, err := RequestLookupKeyListFromBytes(requestLookup.GetAt([]byte(contract)))
	require.NoError(t, err)

	for _, requestKey := range requestKeys {
		require.False(t, requestKey.BlockIndex() == prunedBlockIndex)
	}
}

func TestPruneRequestIndexLookupTable(t *testing.T) {
	const maxBlocks = 4
	const blockToPrune = 1

	requestIDDigest0 := kv.Key("0")
	requestIDDigest1 := kv.Key("1")

	d := dict.Dict{}

	requestIndexLUT := collections.NewMap(d, prefixRequestLookupIndex)
	requestIndexLUT.SetAt([]byte(requestIDDigest0), createRequestLookupKeys(maxBlocks))
	requestIndexLUT.SetAt([]byte(requestIDDigest1), createRequestLookupKeys(maxBlocks))

	require.NotPanics(t, func() {
		pruneRequestLookupByBlockIndex(d, 1)
	})

	validatePrunedRequestIndexLookupBlock(t, d, requestIDDigest0, blockToPrune)
	validatePrunedRequestIndexLookupBlock(t, d, requestIDDigest1, blockToPrune)
}
