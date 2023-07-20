package cache

import (
	"encoding/binary"
	"errors"
	"math"
	"sync"

	"github.com/VictoriaMetrics/fastcache"
)

// fastcache instance
var cache *fastcache.Cache

// only needed or handleCounter. Fastcache is thread-safe
var mutex = sync.Mutex{}

// handle counter
var handleCounter uint32

func init() {
	// todo make it parametrizable (probably build a component with params)
	cache = fastcache.New(512 * 1024 * 1024)
}

func buildKey(handle uint32, key []byte) []byte {
	var newKey = make([]byte, 4+len(key))
	binary.LittleEndian.PutUint32(newKey, handle)
	copy(newKey[4:], key)

	return newKey
}

// get from cache
func Get(handle uint32, key []byte) []byte {
	return cache.Get(nil, buildKey(handle, key))
}

// store into cache
func Set(handle uint32, key []byte, value []byte) {
	if value == nil {
		return
	}
	copied := make([]byte, len(value))
	copy(copied, value)

	cache.Set(buildKey(handle, key), copied)
}

// check if key in cache
func Has(handle uint32, key []byte) bool {
	return cache.Has(buildKey(handle, key))
}

func HasGet(handle uint32, key []byte) ([]byte, bool) {
	return cache.HasGet(nil, buildKey(handle, key))
}

// create new handle
func NewHandle() (uint32, error) {
	mutex.Lock()
	defer mutex.Unlock()

	// check if we can create more ...
	if handleCounter == math.MaxUint32 {
		return 0, errors.New("too many cache handles")
	}
	handleCounter++

	return handleCounter, nil
}
