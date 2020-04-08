package bucket

import (
	"bytes"
	"fmt"
	"github.com/Devister/small-cache/entry"
	"log"
	"sync"
	"unsafe"
)

type ArrayBucket struct {
	sync.RWMutex
	entrySize int
	entryCap  int
	entryPtrs []uint64
	pool      *entry.Pool
	overflow  bool
}

func (b *ArrayBucket) Set(hkey uint64, key []byte, value []byte) {
	var err error

	b.Lock()
	defer b.Unlock()

	// get entry and idx
	e, idx := b.getEntry(hkey, key)
	if e != nil {
		// entry exist
		if e.Cap() >= len(key)+len(value) {
			// capacity is big enough
			if err := e.Set(key, value); err != nil {
				log.Fatal("unexpected, entry set key and value failed: ", err.Error())
			}
		} else {
			// capacity is not enough, new entry and replace the old one
			b.pool.RecycleEntry(e)
			e, err = b.pool.GetEntry(key, value)
			if err != nil {
				fmt.Println("[warn] bucket set failed, can not get entry, error: ", err.Error())
				return
			}
			b.entryPtrs[idx] = uint64(uintptr(unsafe.Pointer(e)))
		}
	} else {
		// entry does not exist, new entry and append to entryPtrs
		e, err = b.pool.GetEntry(key, value)
		if err != nil {
			fmt.Println("[warn] bucket set failed, can not get entry, error: ", err.Error())
			return
		}
		ptr := uintptr(unsafe.Pointer(e))
		b.entryPtrs = append(b.entryPtrs, uint64(ptr), hkey)
		if len(b.entryPtrs) >= b.entryCap {
			b.overflow = true
		}
	}
}

func (b *ArrayBucket) Get(hkey uint64, key []byte) []byte {
	b.RLock()
	defer b.RUnlock()

	e, _ := b.getEntry(hkey, key)
	if e == nil {
		return nil
	}
	return e.Value()
}

func (b *ArrayBucket) Delete(hkey uint64, key []byte) {
	b.Lock()
	defer b.Unlock()

	e, idx := b.getEntry(hkey, key)
	if e == nil {
		return
	}

	b.pool.RecycleEntry(e)
	if len(b.entryPtrs) == 2 {
		b.entryPtrs = b.entryPtrs[:0]
	} else {
		lastIdx := len(b.entryPtrs) - 2
		b.entryPtrs[idx] = b.entryPtrs[lastIdx]
		b.entryPtrs[idx+1] = b.entryPtrs[lastIdx+1]
		b.entryPtrs = b.entryPtrs[:lastIdx]
	}
}

func (b *ArrayBucket) getEntry(hkey uint64, key []byte) (*entry.Entry, int) {
	for i := 0; i < len(b.entryPtrs); i += 2 {
		if b.entryPtrs[i+1] == hkey {
			ptr := uintptr(b.entryPtrs[i])
			e := (*entry.Entry)(unsafe.Pointer(ptr))
			if bytes.Equal(e.Key(), key) {
				return e, i
			}
		}
	}
	return nil, -1
}
