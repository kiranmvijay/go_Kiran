//go:build s390x

package bigmod

import (
	"math/bits"
)

// vmslMul2048 is implemented in s390x assembly - File: nat_s390x.s
func vmslMul2048(a, b, c *uint64, num_chunks int, mulTime, accTime *uint64)

// Helpers for Full Montgomery
func pack56(dst []uint64, src []uint) {
	const mask56 uint64 = (1 << 56) - 1
	for i := 0; i < len(dst); i++ {
		bitpos := 56 * i
		k := bitpos >> 6     // divide by 64
		shift := bitpos & 63 // mod 64

		v := uint64(src[k]) >> shift

		if shift > 8 && k+1 < len(src) {
			v |= uint64(src[k+1]) << (64 - shift)
		}

		dst[i] = v & mask56
	}
}

func unpack56(dst []uint, src []uint64) {
	const mask56 uint64 = (1 << 56) - 1
	for i := range dst {
		dst[i] = 0
	}

	for i := 0; i < len(src); i++ {
		bitpos := 56 * i
		k := bitpos >> 6     // divide by 64
		shift := bitpos & 63 // mod 64

		v := src[i] & mask56

		dst[k] |= uint(v << shift)

		if shift > 8 && k+1 < len(dst) {
			dst[k+1] |= uint(v >> (64 - shift))
		}
	}
}

func geq_final(a, b []uint) bool {
	for i := len(a) - 1; i >= 0; i-- {
		if a[i] > b[i] {
			return true
		}
		if a[i] < b[i] {
			return false
		}
	}
	return true
}

func sub_final(a, b []uint) {
	var borrow uint
	for i := 0; i < len(a); i++ {
		a[i], borrow = bits.Sub(a[i], b[i], borrow)
	}
}
