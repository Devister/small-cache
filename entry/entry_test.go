package entry

import (
	"bytes"
	"testing"
)

func TestEntrySet(t *testing.T) {
	capacity := 16
	e := Entry{data: make([]byte, capacity)}
	if e.Cap() != capacity {
		t.Fatal("capacity not expect")
	}
	key := []byte("key1")
	value2 := []byte("12345678")
	if err := e.Set(key, value2); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(e.Key(), key) {
		t.Fatal("key not correct")
	}
	if !bytes.Equal(e.Value(), value2) {
		t.Fatal("value not correct")
	}
	if e.header().KeyLen != uint16(len(key)) || e.header().ValueLen != uint16(len(value2)) {
		t.Fatal("key or value length not expect")
	}

	value3 := []byte("123456789")
	if err := e.Set(key, value3); err != sizeLargeError {
		t.Fatal("error not expect")
	}
}
