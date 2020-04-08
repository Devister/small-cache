package entry

import (
	"fmt"
	"testing"
	"unsafe"
)

func TestCalcSize(t *testing.T) {
	lvls := make([]uint16, 2048)
	sizes := make([]uint16, 2048)
	lastLvl := uint16(0)
	for i := uint16(1); i <= 2048; i++ {
		lvl := calcuSizeLevel(i)
		if lvl < lastLvl {
			t.Fatal("level failed: ", i)
		}
		lvls[i-1] = lvl
		lastLvl = lvl
		size := sizeLevel2Size(lvl)
		if size < i {
			t.Fatal("size failed: ", i)
		}
		sizes[i-1] = size
	}

	for _, lvl := range lvls {
		fmt.Print(lvl)
	}
	fmt.Println()
	for _, size := range sizes {
		fmt.Print(size)
	}
	fmt.Println()
}

type SA struct {
	f1 uint16
	f2 uint16
	f3 uint32
	f4 uint64
}

func TestByte2Struct(t *testing.T) {
	data := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	var sa = *(**SA)(unsafe.Pointer(&data))
	fmt.Println(sa.f1, sa.f2, sa.f3, sa.f4)
}

func TestMaxSizeLevel(t *testing.T) {
	fmt.Println(calcuSizeLevel(65535))
}
