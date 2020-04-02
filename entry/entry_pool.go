package entry

import ()

type entryCache struct {
	entries []*Entry
}

func (c *entryCache) Get() *Entry {

}

func (c *entryCache) Put(e *Entry) {

}

type Pool struct {
	entryCaches []entryCache
}

func (p *Pool) GetEntry(key, value []byte) *Entry {
	keyLen := uint16(len(key))
	valueLen := uint16(len(value))
	size := keyLen + valueLen
	sizeLevel := (size-1)/8 + 1
	e := p.entryCaches[sizeLevel].Get()
	if e != nil {
		return e
	}
	capacity := sizeLevel * 8
	return &Entry{
		header: EntryHeader{
			Size:     entryHeaderSize + capacity,
			Cap:      capacity,
			KeyLen:   keyLen,
			ValueLen: valueLen,
		},
		data: make([]byte, capacity),
	}
}

func (p *Pool) RecycleEntry(e *Entry) {
	size := e.header.Size
	sizeLevel := (size-1)/8 + 1
	p.entryCaches[sizeLevel].Put(e)
}
