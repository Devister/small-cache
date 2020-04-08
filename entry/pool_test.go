package entry

import (
	"bytes"
	"fmt"
	"testing"
)

func TestCalcSize(t *testing.T) {
	boundary := 65535
	max := float64(0)
	maxNum := 0
	//lvls := make([]uint16, boundary + 1)
	//sizes := make([]uint16, boundary + 1)
	lastLvl := uint16(0)
	for i := 1; i <= boundary; i++ {
		lvl := calcuSizeLevel(uint16(i))
		if lvl < lastLvl {
			t.Fatal("level failed: ", i)
		}
		//lvls[i-1] = lvl
		lastLvl = lvl
		size := sizeLevel2Size(lvl)
		if size < uint16(i) {
			t.Fatal("size failed: ", i)
		}
		if i > 56 {
			over := (float64(size) - float64(i)) / float64(size)
			if over > max {
				max = over
				maxNum = i
			}
		}

		//sizes[i-1] = size
	}

	fmt.Println("max memory waste ratio: ", max, " at number ", maxNum)
	//for _, lvl := range lvls {
	//	fmt.Print(lvl)
	//}
	//fmt.Println()
	//for _, size := range sizes {
	//	fmt.Print(size)
	//}
	//fmt.Println()
}

func TestSizeLevel(t *testing.T) {
	lvl := calcuSizeLevel(1)
	fmt.Println("level: ", lvl)
	size := sizeLevel2Size(lvl)
	fmt.Println("size: ", size)
}

type keyValuePair struct {
	key   []byte
	value []byte
	err   error
}

type entryTestCase struct {
	size       int
	lvl        int
	expectSize int
}

func TestPool_GetEntry(t *testing.T) {
	p := NewPool()

	kvPairs := []*keyValuePair{
		{key: []byte{1}, value: []byte{1}},
		{key: []byte{1, 2, 3, 4}, value: []byte{4, 3, 2, 1}},
		{key: []byte{1, 2, 3, 4}, value: bytes.Repeat([]byte{1}, 65527)},
		{key: []byte{1, 2, 3, 4}, value: bytes.Repeat([]byte{1}, 65528), err: SizeLargeError},
	}

	for _, kvPair := range kvPairs {
		poolGetEntry(p, kvPair, t)
	}
}

func TestPool_RecycleEntry(t *testing.T) {
	p := NewPool()

	testCases := []*entryTestCase{
		{size: 16, lvl: 2, expectSize: 1},
		{size: 32, lvl: 4, expectSize: 1},
		{size: 16, lvl: 2, expectSize: 2},
		{size: 16, lvl: 2, expectSize: 3},
		{size: 16, lvl: 2, expectSize: 4},
		{size: 16, lvl: 2, expectSize: 5},
		{size: 16, lvl: 2, expectSize: 5},
		{size: 32, lvl: 4, expectSize: 2},
		{size: 65535, lvl: 88, expectSize: 1},
	}

	for _, tc := range testCases {
		e := &Entry{data: make([]byte, tc.size, tc.size)}
		p.RecycleEntry(e)
		if len(p.entryCaches[tc.lvl].entries) != tc.expectSize {
			t.Fatal("entryCaches length is not expect")
		}
	}

	key := []byte{1, 2, 3, 4}
	value := []byte{1, 2, 3, 4}
	num := 5
	for i := 0; i < 6; i++ {
		_, _ = p.GetEntry(key, value)
		if num--; num < 0 {
			num = 0
		}
		if len(p.entryCaches[2].entries) != num {
			t.Fatal("entryCaches length is not expect")
		}
	}
}

func poolGetEntry(p *Pool, kvPair *keyValuePair, t *testing.T) {
	e, err := p.GetEntry(kvPair.key, kvPair.value)
	if err != kvPair.err {
		t.Fatal("expect error: ", kvPair.err, ", got error: ", err)
	}
	if err == nil {
		if err := e.Set(kvPair.key, kvPair.value); err != kvPair.err {
			t.Fatal(err)
		}
	}
}
