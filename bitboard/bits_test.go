package bitboard

import (
	"strconv"
	"testing"
)

func TestPrecompute(t *testing.T) {
	c := Precompute(5)
	if c.B != (1<<5)-1 {
		t.Error("c.b(5):", strconv.FormatUint(c.B, 2))
	}
	if c.T != ((1<<5)-1)<<(4*5) {
		t.Error("c.t(5):", strconv.FormatUint(c.T, 2))
	}
	if c.R != 0x0108421 {
		t.Error("c.r(5):", strconv.FormatUint(c.R, 2))
	}
	if c.L != 0x1084210 {
		t.Error("c.l(5):", strconv.FormatUint(c.L, 2))
	}
	if c.Mask != 0x1ffffff {
		t.Error("c.mask(5):", strconv.FormatUint(c.Mask, 2))
	}

	c = Precompute(8)
	if c.B != (1<<8)-1 {
		t.Error("c.b(8):", strconv.FormatUint(c.B, 2))
	}
	if c.T != ((1<<8)-1)<<(7*8) {
		t.Error("c.t(8):", strconv.FormatUint(c.T, 2))
	}
	if c.R != 0x101010101010101 {
		t.Error("c.r(8):", strconv.FormatUint(c.R, 2))
	}
	if c.L != 0x8080808080808080 {
		t.Error("c.l(8):", strconv.FormatUint(c.L, 2))
	}
	if c.Mask != ^uint64(0) {
		t.Error("c.mask(8):", strconv.FormatUint(c.Mask, 2))
	}
}
