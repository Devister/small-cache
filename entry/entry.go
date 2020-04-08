package entry

import "unsafe"

const (
	entryHeaderSize = uint16(unsafe.Sizeof(EntryHeader{}))
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
	return e.data[entryHeaderSize:][:e.header().KeyLen]
}

func (e *Entry) Value() []byte {
	return e.data[entryHeaderSize:][e.header().KeyLen:e.header().ValueLen]
}

func (e *Entry) Cap() int {
	return len(e.data)
}

func (e *Entry) Set(key, value []byte) error {
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
