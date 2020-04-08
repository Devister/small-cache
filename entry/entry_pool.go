package entry

import (
	"errors"
	"sync"
)

const (
	maxSize      = 65535
	maxSizeLevel = 89 // maxSizeLevel = calcuSizeLevel(maxSize) + 1
)

var (
	SizeLargeError = errors.New("key and value size is large than 65535")
)

type entryCache struct {
	sync.Mutex
	entries []*Entry
}

func (c *entryCache) Get() *Entry {
	c.Lock()
	defer c.Unlock()
	if len(c.entries) > 0 {
		e := c.entries[len(c.entries)-1]
		c.entries = c.entries[:len(c.entries)-1]
		return e
	}
	return nil
}

func (c *entryCache) Put(e *Entry) {
	c.Lock()
	defer c.Unlock()
	if cap(c.entries) > len(c.entries) {
		c.entries = append(c.entries, e)
	}
}

type Pool struct {
	sync.Mutex
	entryCaches []entryCache
}

func NewPool() *Pool {
	return &Pool{
		entryCaches: make([]entryCache, maxSizeLevel),
	}
}

func (p *Pool) GetEntry(key, value []byte) (*Entry, error) {
	keyLen := uint16(len(key))
	valueLen := uint16(len(value))
	size := keyLen + valueLen + uint16(entryHeaderSize)
	if size > maxSize {
		return nil, SizeLargeError
	}
	sizeLevel := calcuSizeLevel(size)

	e := p.entryCaches[sizeLevel].Get()
	if e != nil {
		return e, nil
	}
	capacity := sizeLevel2Size(sizeLevel)
	return &Entry{
		data: make([]byte, capacity, capacity),
	}, nil
}

func (p *Pool) RecycleEntry(e *Entry) {
	size := e.Cap()
	sizeLevel := (size-1)/8 + 1
	p.entryCaches[sizeLevel].Put(e)
}

func calcuSizeLevel(size uint16) uint16 {
	size -= 1
	if size < 128 {
		return size/8 + 1
	}
	n, m := log(size)
	return (size-m)/(m/8) + (n-5)*8 + 1
}

func log(size uint16) (uint16, uint16) {
	count := uint16(0)
	num := uint16(1)
	for size > 1 {
		size /= 2
		count++
		num *= 2
	}
	return count, num
}

func sizeLevel2Size(lvl uint16) uint16 {
	if lvl <= 16 {
		return lvl * 8
	}
	n := lvl / 8
	m := uint16(8 * (1<<n - 1))
	size := m * (8 + lvl - (n * 8))
	return size
}
