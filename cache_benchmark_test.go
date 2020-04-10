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
	cfg := &CacheConfig{BucketNum: 0x10000}
	cache := NewCache(cfg)
	fmt.Println("[debug] cache initialized")
	rand.Seed(time.Now().Unix())

	//b.N = 0x10000
	b.RunParallel(func(pb *testing.PB) {
		id := rand.Int()
		counter := 0

		b.ReportAllocs()
		for pb.Next() {
			cache.Set([]byte(fmt.Sprintf("key-%d-%d", id, counter)), message)
			counter = counter + 1
		}
		fmt.Println("set number: ", counter)
	})
	//counter := 0
	//for {
	//	id := rand.Int() % 10
	//	counter = rand.Int() % 1000
	//	key := fmt.Sprintf("key-%d-%d", id, counter)
	//	//fmt.Println(key)
	//	cache.Set([]byte(key), message)
	//}
}

func blob(char byte, len int) []byte {
	return bytes.Repeat([]byte{char}, len)
}
