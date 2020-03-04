package futil

// util functions for field elements

import "github.com/decentralisedkev/go-jubjub/internal/uint128"

// Adc Computes a + b + carry, returning the result and the new carry over.
func Adc(a, b, carry uint64) (uint64, uint64) {
	res := uint128.FromU64(a)
	res = res.Add(b)
	res = res.Add(carry)

	return res.L, res.H
}

// Sbb Computes a - (b + borrow), returning the result and the new borrow.
func Sbb(a, b, borrow uint64) (uint64, uint64) {

	a128 := uint128.FromU64(a)
	b128 := uint128.FromU64(b)
	borr128 := uint128.FromU64(borrow >> 63)

	bBor := b128.AddU128(borr128)
	res := a128.SubU128(bBor)

	return res.L, res.H

}

// Mac Computes a + (b * c) + carry, returning the result and the new carry over.
func Mac(a, b, c, carry uint64) (uint64, uint64) {

	res := uint128.FromU64(b)
	res = res.MulU64(c)
	res = res.Add(a)
	res = res.Add(carry)

	return res.L, res.H
}

// Load4 interprets a 4-byte unsigned little endian byte-slice as uint64
func Load4(b []byte) uint64 {
	return (uint64(b[0]) | (uint64(b[1]) << 8) |
		(uint64(b[2]) << 16) | (uint64(b[3]) << 24))
}
