package bucket

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/Devister/small-cache/entry"
	"log"
	"sync"
)

var (
	NilErr = errors.New("no such key")
)

const (
	defaultEntryCap = 512
)

type ArrayBucket struct {
	sync.RWMutex
	entryCap  int
	entryPtrs []*entry.Entry
	hashKeys  []uint64
	pool      *entry.Pool
	overflow  bool
}

func NewArrayBucket(pool *entry.Pool) *ArrayBucket {
	return &ArrayBucket{
		entryCap:  defaultEntryCap,
		entryPtrs: make([]*entry.Entry, 0, defaultEntryCap),
		hashKeys:  make([]uint64, 0, defaultEntryCap),
		pool:      pool,
	}
}

func (b *ArrayBucket) Set(hkey uint64, key []byte, value []byte) error {
	b.Lock()
	defer b.Unlock()

	// get entry and idx
	if entryIdx := b.getEntryIdx(hkey, key); entryIdx >= 0 {
		// entry exist
		e := b.entryPtrs[entryIdx]
		if e.Cap() >= entry.EntryLen(key, value) {
			// capacity is big enough
			if err := e.Set(key, value); err != nil {
				log.Fatal("unexpected, entry set key and value failed: ", err.Error())
				return err
			}
		} else {
			// capacity is not enough, new entry and replace the old one
			b.pool.RecycleEntry(e)
			e, err := b.poolGetEntryAndSet(key, value)
			if err != nil {
				return err
			}
			b.entryPtrs[entryIdx] = e
		}
	} else {
		// entry does not exist, new entry and append to entryPtrs
		e, err := b.poolGetEntryAndSet(key, value)
		if err != nil {
			return err
		}
		b.entryPtrs = append(b.entryPtrs, e)
		b.hashKeys = append(b.hashKeys, hkey)
		if len(b.entryPtrs) >= b.entryCap {
			b.overflow = true
			fmt.Println("[debug] bucket overflow, len: ", len(b.entryPtrs))
		}
	}
	return nil
}

func (b *ArrayBucket) Get(hkey uint64, key []byte) ([]byte, error) {
	b.RLock()
	defer b.RUnlock()

	if idx := b.getEntryIdx(hkey, key); idx == -1 {
		return nil, NilErr
	} else {
		return b.entryPtrs[idx].Value(), nil
	}
}

func (b *ArrayBucket) Delete(hkey uint64, key []byte) error {
	b.Lock()
	defer b.Unlock()

	idx := b.getEntryIdx(hkey, key)
	if idx == -1 {
		return NilErr
	}

	b.pool.RecycleEntry(b.entryPtrs[idx])
	if len(b.entryPtrs) == 2 {
		b.entryPtrs = b.entryPtrs[:0]
		b.hashKeys = b.hashKeys[:0]
	} else {
		lastIdx := len(b.entryPtrs) - 1
		b.entryPtrs[idx] = b.entryPtrs[lastIdx]
		b.entryPtrs = b.entryPtrs[:lastIdx]
		b.hashKeys[idx] = b.hashKeys[lastIdx]
		b.hashKeys = b.hashKeys[:lastIdx]
	}
	return nil
}

func (b *ArrayBucket) getEntryIdx(hkey uint64, key []byte) int {
	for i, hk := range b.hashKeys {
		if hk == hkey && bytes.Equal(b.entryPtrs[i].Key(), key) {
			return i
		}
	}
	return -1
}

func (b *ArrayBucket) poolGetEntryAndSet(key, value []byte) (*entry.Entry, error) {
	e, err := b.pool.GetEntry(key, value)
	if err != nil {
		fmt.Println("[warn] bucket set failed, can not get entry, error: ", err.Error())
		return nil, err
	}
	if err := e.Set(key, value); err != nil {
		log.Fatal("unexpected, entry set key and value failed: ", err.Error())
		return nil, err
	}
	return e, nil
}
