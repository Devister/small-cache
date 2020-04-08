package entry

import (
	"errors"
	"fmt"
	"unsafe"
)

const (
	entryHeaderSize = unsafe.Sizeof(EntryHeader{})
)

var (
	sizeLargeError = errors.New("size is large than capacity")
)

// caller should check if KeyLen + ValueLen <= Cap;
type EntryHeader struct {
	KeyLen   uint16
	ValueLen uint16
}

type Entry struct {
	/* data consists of: EntryHeader(32bit), key(keyLen bit), value(valueLen bit) */
	/* len(data) == cap(data) */
	data []byte
}

func (e *Entry) Key() []byte {
	if e.Cap() < int(entryHeaderSize) {
		fmt.Println("capacity of entry: ", e.Cap())
		return nil
	}
	if h := e.header(); h == nil {
		fmt.Println("header of entry is nil")
		return nil
	} else {
		if h.KeyLen+uint16(entryHeaderSize) > uint16(e.Cap()) {
			fmt.Println("key length over entry capacity")
			return nil
		}
	}

	return e.data[entryHeaderSize:][:e.header().KeyLen]
}

func (e *Entry) Value() []byte {
	return e.data[entryHeaderSize:][e.header().KeyLen : e.header().KeyLen+e.header().ValueLen]
}

func (e *Entry) Cap() int {
	return cap(e.data)
}

func (e *Entry) Set(key, value []byte) error {
	if EntryLen(key, value) > e.Cap() {
		return sizeLargeError
	}
	e.header().KeyLen = uint16(len(key))
	e.header().ValueLen = uint16(len(value))
	copy(e.data[entryHeaderSize:], key)
	copy(e.data[entryHeaderSize:][e.header().KeyLen:], value)
	return nil
}

func (e *Entry) header() *EntryHeader {
	header := *(**EntryHeader)(unsafe.Pointer(&e.data))
	return header
}

func EntryLen(key, value []byte) int {
	return len(key) + len(value) + int(entryHeaderSize)
}
