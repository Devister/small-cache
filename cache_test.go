package small_cache

import (
	"bytes"
	"testing"
)

func TestCacheSetAndGet(t *testing.T) {
	c := NewCache(&CacheConfig{})
	key := []byte("key1")
	value1 := []byte("value1")
	value2 := []byte("value2")

	if err := c.Set(key, value1); err != nil {
		t.Fatal(err)
	}
	if v, err := c.Get(key); err != nil {
		t.Fatal(err)
	} else {
		if !bytes.Equal(v, value1) {
			t.Fatal("value not expect")
		}
	}

	if err := c.Set(key, value2); err != nil {
		t.Fatal(err)
	}
	if v, err := c.Get(key); err != nil {
		t.Fatal(err)
	} else {
		if !bytes.Equal(v, value2) {
			t.Fatal("value not expect")
		}
	}
}
