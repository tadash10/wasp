package kv

import "github.com/iotaledger/wasp/packages/cache"

type cachedKVStoreReader struct {
	KVStoreReader
	partition []byte
}

// NewCachedKVStoreReader creates a KVStoreReader with an prefix to a fastcache instance.
// IMPORTANT: there is no logic for cache invalidation, so make sure that the
// underlying KVStoreReader is never mutated.
func NewCachedKVStoreReader(r KVStoreReader) KVStoreReader {
	newPartition, err := cache.NewPartition()
	if err != nil {
		panic(err)
	}
	return &cachedKVStoreReader{
		KVStoreReader: r,
		partition:     newPartition,
	}
}

func (c *cachedKVStoreReader) Get(key Key) []byte {
	if v := cache.Get(c.partition, []byte(key)); v != nil {
		return v
	}
	v := c.KVStoreReader.Get(key)
	cache.Set(c.partition, []byte(key), v)
	return v
}

func (c *cachedKVStoreReader) Has(key Key) bool {
	return cache.Has(c.partition, []byte(key))
}
