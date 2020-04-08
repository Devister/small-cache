package small_cache

import (
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"github.com/Devister/small-cache/bucket"
	"hash"
)

type sCache struct {
	bucketNum   uint64
	memoryLimit int
	buckets     []bucket.Bucket
	hashFunc    hash.Hash
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
	for i := uint64(0); i < cfg.BucketNum; i++ {
		b := bucket.NewArrayBucket()
		c.buckets = append(c.buckets, b)
	}
	return c
}

func (c *sCache) Set(key, value []byte) error {
	hkey := binary.BigEndian.Uint64(c.hashFunc.Sum(key))
	b := c.getBucket(hkey)
	return b.Set(hkey, key, value)
}

func (c *sCache) Get(key []byte) ([]byte, error) {
	hkey := binary.BigEndian.Uint64(c.hashFunc.Sum(key))
	b := c.getBucket(hkey)
	return b.Get(hkey, key)
}

func (c *sCache) Delete(key []byte) error {
	hkey := binary.BigEndian.Uint64(c.hashFunc.Sum(key))
	b := c.getBucket(hkey)
	return b.Delete(hkey, key)
}

func (c *sCache) getBucket(hkey uint64) bucket.Bucket {
	bucketId := hkey & (c.bucketNum - 1)
	return c.buckets[bucketId]
}
