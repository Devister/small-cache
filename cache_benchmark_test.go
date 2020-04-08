package small_cache

import (
	"bytes"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

var message = blob('a', 256)

func BenchmarkWriteOnCache(b *testing.B) {
	writeOnCache(b)
}

func writeOnCache(b *testing.B) {
	cfg := &CacheConfig{}
	cache := NewCache(cfg)
	rand.Seed(time.Now().Unix())

	//b.RunParallel(func(pb *testing.PB) {
	//	id := rand.Int()
	//	counter := 0
	//
	//	b.ReportAllocs()
	//	for pb.Next() {
	//		cache.Set([]byte(fmt.Sprintf("key-%d-%d", id, counter)), message)
	//		counter = counter + 1
	//	}
	//})
	counter := 0
	for {
		id := rand.Int() % 10
		cache.Set([]byte(fmt.Sprintf("key-%d-%d", id, counter)), message)
		counter = rand.Int() % 1000
	}
}

func blob(char byte, len int) []byte {
	return bytes.Repeat([]byte{char}, len)
}
