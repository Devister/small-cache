package small_cache

import "github.com/Devister/small-cache/bucket"

type sCache struct {
	bucketNum   int
	memoryLimit int64
	buckets     []bucket.Bucket
}

func NewCache(memoryLimit int64) *sCache {
	return &sCache{}
}
