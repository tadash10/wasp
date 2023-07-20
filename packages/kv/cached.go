package kv

import "github.com/iotaledger/wasp/packages/cache"

type cachedKVStoreReader struct {
	KVStoreReader
	handle uint32
}

// NewCachedKVStoreReader wraps a KVStoreReader with an prefix to a fastcache instance.
// IMPORTANT: there is no logic for cache invalidation, so make sure that the
// underlying KVStoreReader is never mutated.
func NewCachedKVStoreReader(r KVStoreReader) KVStoreReader {
	newhandle, err := cache.NewHandle()
	if err != nil {
		panic(err)
	}
	return &cachedKVStoreReader{
		KVStoreReader: r,
		handle:        newhandle,
	}
}

func (c *cachedKVStoreReader) Get(key Key) []byte {
	if v, ok := cache.HasGet(c.handle, []byte(key)); ok {
		return v
	}
	v := c.KVStoreReader.Get(key)
	cache.Set(c.handle, []byte(key), v)
	return v
}

func (c *cachedKVStoreReader) Has(key Key) bool {
	return cache.Has(c.handle, []byte(key))
}
