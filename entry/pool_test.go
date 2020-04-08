package entry

import (
	"fmt"
	"testing"
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
