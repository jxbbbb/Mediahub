package utils

import (
	"math"
	"math/rand/v2"
	"testing"
)

func TestToBase62(t *testing.T) {
	for i := 0; i < 1000; i++ {
		d := rand.Int64N(math.MaxInt64)
		str := ToBase62(d)
		d1 := ToBase10(str)
		if d1 != d {
			t.Errorf("ToBase62(%d)=%d, want %d", d, d1, d)
		}

	}
}
