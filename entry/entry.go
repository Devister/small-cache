package entry

import "unsafe"

const (
	entryHeaderSize = uint16(unsafe.Sizeof(EntryHeader{}))
)

// caller should check if KeyLen + ValueLen <= Cap;
type EntryHeader struct {
	Size     uint16
	Cap      uint16
	KeyLen   uint16
	ValueLen uint16
}

type Entry struct {
	header EntryHeader
	data   []byte
}

func (e *Entry) Key() []byte {
	return e.data[:e.header.KeyLen]
}

func (e *Entry) Value() []byte {
	return e.data[e.header.KeyLen:e.header.ValueLen]
}

func (e *Entry) Cap() int {
	return int(e.header.Cap)
}

func (e *Entry) Set(key, value []byte) error {
	e.header.KeyLen = uint16(len(key))
	e.header.ValueLen = uint16(len(value))
	copy(e.data, key)
	copy(e.data[e.header.KeyLen:], value)
	return nil
}
