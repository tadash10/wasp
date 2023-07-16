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

// only needed or partitionCounter. Fastcache is thread-safe
var mutex sync.Mutex = sync.Mutex{}

// partition counter
var partitionCounter uint32 = 0

func init() {
	// todo make it parametrizable (probably build a component with params)
	cache = fastcache.New(32 * 1024 * 1024)
}

// get from cache
func Get(partition []byte, key []byte) []byte {
	return cache.Get(nil, append(partition, key...))
}

// store into cache
func Set(partition []byte, key []byte, value []byte) {
	cache.Set(append(partition, key...), value)
}

// check if key in cache
func Has(partition []byte, key []byte) bool {
	return cache.Has(append(partition, key...))
}

// create new partition
func NewPartition() ([]byte, error) {
	mutex.Lock()
	defer mutex.Unlock()

	// check if we can create more ...
	if partitionCounter == math.MaxUint32 {
		return nil, errors.New("too many cache partitions")
	}
	partitionCounter++

	var partitionBytes []byte
	binary.LittleEndian.PutUint32(partitionBytes[:], partitionCounter)

	return partitionBytes, nil
}
