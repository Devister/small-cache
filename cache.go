package small_cache

import (
	"crypto/sha1"
	"fmt"
	"github.com/Devister/small-cache/bucket"
	"github.com/Devister/small-cache/entry"
	xxhash "github.com/cespare/xxhash/v2"
	"hash"
	"sync"
)

type sCache struct {
	bucketNum   uint64
	memoryLimit int
	buckets     []bucket.Bucket // interface brings 20% overhead
	// hash func with lock takes 200% overhead, using xxhash instead
	hashFunc hash.Hash
	hashLock sync.Mutex
}

func NewCache(cfg *CacheConfig) *sCache {
	if err := parseConfig(cfg); err != nil {
		fmt.Println("[error] new cache failed, parse config error: ", err.Error())
		return nil
	}
	c := &sCache{
		bucketNum:   cfg.BucketNum,
		memoryLimit: cfg.HardMemoryLimit,
		buckets:     make([]bucket.Bucket, 0, cfg.BucketNum),
		hashFunc:    sha1.New(),
	}
	pool := entry.NewPool()
	for i := uint64(0); i < cfg.BucketNum; i++ {
		b := bucket.NewArrayBucket(pool)
		c.buckets = append(c.buckets, b)
	}
	return c
}

func (c *sCache) Set(key, value []byte) error {
	hkey := c.getHashKey(key)
	b := c.getBucket(hkey)
	return b.Set(hkey, key, value)
}

func (c *sCache) Get(key []byte) ([]byte, error) {
	hkey := c.getHashKey(key)
	b := c.getBucket(hkey)
	return b.Get(hkey, key)
}

func (c *sCache) Delete(key []byte) error {
	hkey := c.getHashKey(key)
	b := c.getBucket(hkey)
	return b.Delete(hkey, key)
}

func (c *sCache) getBucket(hkey uint64) bucket.Bucket {
	bucketId := hkey & (c.bucketNum - 1)
	return c.buckets[bucketId]
}

func (c *sCache) getHashKey(key []byte) uint64 {
	//c.hashLock.Lock()
	//defer c.hashLock.Unlock()
	//c.hashFunc.Reset()
	//c.hashFunc.Write(key)
	//return binary.BigEndian.Uint64(c.hashFunc.Sum([]byte{}))
	return xxhash.Sum64(key)
}
