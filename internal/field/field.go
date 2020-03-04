package field

import (
	"encoding/binary"

	"github.com/decentralisedkev/go-jubjub/internal/futil"
)

// Field is a generic field element
type Field [4]uint64

// SetZero sets f to the zero element
func (f *Field) SetZero() *Field {
	copy(f[:], Zero[:])
	return f
}

// IsZero returns true if the field element is equal to the the zero element
func (f *Field) IsZero() bool {
	return f[0]|f[1]|f[2]|f[3] == 0
}

// IsOdd checks whether the element is Odd or even
// Please note that since we are in LE format the LSB will be the first element
// This would be isNegative for Ed25519
// returns 1 if odd and 0 if even
func (f *Field) IsOdd() uint64 {
	if f[0]&1 == 1 {
		return 1
	}
	return 0
}

// Set sets f to the the field element passed as an argument
func (f *Field) Set(a Field) *Field {
	copy(f[:], a[:])
	// *f = a // XXX: check which is faster
	return f
}

// FromU64 converts a uint64 into a field element
func (f *Field) FromU64(a, INV uint64, R2, modulus Field) *Field {

	f[0] = a
	f[1] = 0
	f[2] = 0
	f[3] = 0

	return f.Mul(*f, R2, INV, modulus)
}

// CondSel Conditionally selects and sets a if c = 1 or b if c =  0
func (f *Field) CondSel(a, b Field, c uint64) *Field {

	f[0] = ^(c-1)&a[0] | (c-1)&b[0]
	f[1] = ^(c-1)&a[1] | (c-1)&b[1]
	f[2] = ^(c-1)&a[2] | (c-1)&b[2]
	f[3] = ^(c-1)&a[3] | (c-1)&b[3]

	return f
}

// CondSet Conditionally sets f to a if b == 1
// Taken from Go-ristretto (Bas)
func (f *Field) CondSet(a Field, b uint64) *Field {

	b = -b // b == 0b11111111111111111111111111111111 or 0.
	f[0] ^= b & (f[0] ^ a[0])
	f[1] ^= b & (f[1] ^ a[1])
	f[2] ^= b & (f[2] ^ a[2])
	f[3] ^= b & (f[3] ^ a[3])

	return f
}

// FromBytes takes a 64 byte array and returns
// a point representation in Fq
func (f *Field) FromBytes(byt [64]byte, INV uint64, modulus, R2, R3 Field) *Field {

	var d0, d1 *Field
	d0 = &Field{0, 0, 0, 0}
	d1 = &Field{0, 0, 0, 0}

	d0[0] = binary.LittleEndian.Uint64(byt[0:8])
	d0[1] = binary.LittleEndian.Uint64(byt[8:16])
	d0[2] = binary.LittleEndian.Uint64(byt[16:24])
	d0[3] = binary.LittleEndian.Uint64(byt[24:32])

	d1[0] = binary.LittleEndian.Uint64(byt[32:40])
	d1[1] = binary.LittleEndian.Uint64(byt[40:48])
	d1[2] = binary.LittleEndian.Uint64(byt[48:56])
	d1[3] = binary.LittleEndian.Uint64(byt[56:64])

	d0.Sub(*d0, modulus, modulus)
	d1.Sub(*d1, modulus, modulus)

	// Convert to Montgomery form
	d0.Mul(*d0, R2, INV, modulus)
	d1.Mul(*d1, R3, INV, modulus)

	f.Add(*d0, *d1, modulus)

	return f
}

func Equal(a, b Field) bool {
	return a[0] == b[0] && a[1] == b[1] && a[2] == b[2] && a[3] == b[3]
}

// Cmp compares a and b
// if a > b return 1
// if a==b return 0
// if a < b return -1
func Cmp(a, b Field) int8 {
	for i := 3; i >= 0; i-- { // We are in little endian format
		if a[i] > b[i] {
			return 1
		} else if b[i] > a[i] {
			return -1
		}
	}
	return 0
}

func (f *Field) Sub(a, b, modulus Field) *Field {

	var borrow, carry uint64
	d0, borrow := futil.Sbb(a[0], b[0], borrow)
	d1, borrow := futil.Sbb(a[1], b[1], borrow)
	d2, borrow := futil.Sbb(a[2], b[2], borrow)
	d3, borrow := futil.Sbb(a[3], b[3], borrow)

	// If underflow occurred on the final limb, borrow = 0xfff...fff, otherwise
	// borrow = 0x000...000. Thus, we use it as a mask to conditionally add the modulus.
	d0, carry = futil.Adc(d0, modulus[0]&borrow, carry)
	d1, carry = futil.Adc(d1, modulus[1]&borrow, carry)
	d2, carry = futil.Adc(d2, modulus[2]&borrow, carry)
	d3, carry = futil.Adc(d3, modulus[3]&borrow, carry)

	f[0] = d0
	f[1] = d1
	f[2] = d2
	f[3] = d3

	return f
}

// Add Adds one field to another
func (f *Field) Add(lhs, rhs, modulus Field) *Field {

	var carry uint64
	d0, carry := futil.Adc(lhs[0], rhs[0], carry)
	d1, carry := futil.Adc(lhs[1], rhs[1], carry)
	d2, carry := futil.Adc(lhs[2], rhs[2], carry)
	d3, _ := futil.Adc(lhs[3], rhs[3], carry)

	f[0] = d0
	f[1] = d1
	f[2] = d2
	f[3] = d3

	f.Sub(*f, modulus, modulus)

	return f
}

func (f *Field) Double(modulus Field) *Field {
	return f.Add(*f, *f, modulus)
}

func (f *Field) Mul(lhs, rhs Field, INV uint64, modulus Field) *Field {
	// XXX: Optimise later

	r0, carry := futil.Mac(0, lhs[0], rhs[0], 0)
	r1, carry := futil.Mac(0, lhs[0], rhs[1], carry)
	r2, carry := futil.Mac(0, lhs[0], rhs[2], carry)
	r3, r4 := futil.Mac(0, lhs[0], rhs[3], carry)

	r1, carry = futil.Mac(r1, lhs[1], rhs[0], 0)
	r2, carry = futil.Mac(r2, lhs[1], rhs[1], carry)
	r3, carry = futil.Mac(r3, lhs[1], rhs[2], carry)
	r4, r5 := futil.Mac(r4, lhs[1], rhs[3], carry)

	r2, carry = futil.Mac(r2, lhs[2], rhs[0], 0)
	r3, carry = futil.Mac(r3, lhs[2], rhs[1], carry)
	r4, carry = futil.Mac(r4, lhs[2], rhs[2], carry)
	r5, r6 := futil.Mac(r5, lhs[2], rhs[3], carry)

	r3, carry = futil.Mac(r3, lhs[3], rhs[0], 0)
	r4, carry = futil.Mac(r4, lhs[3], rhs[1], carry)
	r5, carry = futil.Mac(r5, lhs[3], rhs[2], carry)
	r6, r7 := futil.Mac(r6, lhs[3], rhs[3], carry)

	*f = *montRed(r0, r1, r2, r3, r4, r5, r6, r7, INV, modulus)

	return f
}

func (f *Field) Square(a Field, INV uint64, modulus Field) *Field {

	r1, carry := futil.Mac(0, a[0], a[1], 0)
	r2, carry := futil.Mac(0, a[0], a[2], carry)
	r3, r4 := futil.Mac(0, a[0], a[3], carry)

	r3, carry = futil.Mac(r3, a[1], a[2], 0)
	r4, r5 := futil.Mac(r4, a[1], a[3], carry)

	r5, r6 := futil.Mac(r5, a[2], a[3], 0)

	r7 := r6 >> 63
	r6 = (r6 << 1) | (r5 >> 63)
	r5 = (r5 << 1) | (r4 >> 63)
	r4 = (r4 << 1) | (r3 >> 63)
	r3 = (r3 << 1) | (r2 >> 63)
	r2 = (r2 << 1) | (r1 >> 63)
	r1 = r1 << 1

	r0, carry := futil.Mac(0, a[0], a[0], 0)
	r1, carry = futil.Adc(0, r1, carry)
	r2, carry = futil.Mac(r2, a[1], a[1], carry)
	r3, carry = futil.Adc(0, r3, carry)

	r4, carry = futil.Mac(r4, a[2], a[2], carry)

	r5, carry = futil.Adc(0, r5, carry)

	r6, carry = futil.Mac(r6, a[3], a[3], carry)
	r7, _ = futil.Adc(0, r7, carry)

	*f = *montRed(r0, r1, r2, r3, r4, r5, r6, r7, INV, modulus)

	return f
}

// Neg negates a Fq
func (f *Field) Neg(a, modulus Field) *Field {

	var borrow uint64
	d0, borrow := futil.Sbb(modulus[0], a[0], borrow)
	d1, borrow := futil.Sbb(modulus[1], a[1], borrow)
	d2, borrow := futil.Sbb(modulus[2], a[2], borrow)
	d3, _ := futil.Sbb(modulus[3], a[3], borrow)

	msk := a[0]|a[1]|a[2]|a[3] == 0

	var mask uint64
	if !msk {
		mask-- // uint64 max
	}

	// `tmp` could be `MODULUS` if `self` was zero. Create a mask that is
	// zero if `self` was zero, and `u64::max_value()` if self was nonzero.

	f[0] = d0 & mask
	f[1] = d1 & mask
	f[2] = d2 & mask
	f[3] = d3 & mask

	return f
}

func (f *Field) IntoBytes(INV uint64, modulus Field) []byte {

	// Turn into canonical form by computing (a.R) / R = a
	tmp := montRed(f[0], f[1], f[2], f[3], 0, 0, 0, 0, INV, modulus)

	res := make([]byte, 32, 32)

	binary.LittleEndian.PutUint64(res[0:8], tmp[0])
	binary.LittleEndian.PutUint64(res[8:16], tmp[1])
	binary.LittleEndian.PutUint64(res[16:24], tmp[2])
	binary.LittleEndian.PutUint64(res[24:32], tmp[3])

	return res
}

// BytesInto  converts f into a little endian byte slice
// XXX: same as IntoBytes,
func (f *Field) BytesInto(buf *[32]byte, modulus Field, INV uint64) *Field {

	// Turn into canonical form by computing (a.R) / R = a
	tmp := montRed(f[0], f[1], f[2], f[3], 0, 0, 0, 0, INV, modulus)

	buf[0] = uint8(tmp[0])
	buf[1] = uint8(tmp[0] >> 8)
	buf[2] = uint8(tmp[0] >> 16)
	buf[3] = uint8(tmp[0] >> 24)
	buf[4] = uint8(tmp[0] >> 32)
	buf[5] = uint8(tmp[0] >> 40)
	buf[6] = uint8(tmp[0] >> 48)
	buf[7] = uint8(tmp[0] >> 56)
	buf[8] = uint8(tmp[1])
	buf[9] = uint8(tmp[1] >> 8)
	buf[10] = uint8(tmp[1] >> 16)
	buf[11] = uint8(tmp[1] >> 24)
	buf[12] = uint8(tmp[1] >> 32)
	buf[13] = uint8(tmp[1] >> 40)
	buf[14] = uint8(tmp[1] >> 48)
	buf[15] = uint8(tmp[1] >> 56)
	buf[16] = uint8(tmp[2])
	buf[17] = uint8(tmp[2] >> 8)
	buf[18] = uint8(tmp[2] >> 16)
	buf[19] = uint8(tmp[2] >> 24)
	buf[20] = uint8(tmp[2] >> 32)
	buf[21] = uint8(tmp[2] >> 40)
	buf[22] = uint8(tmp[2] >> 48)
	buf[23] = uint8(tmp[2] >> 56)
	buf[24] = uint8(tmp[3])
	buf[25] = uint8(tmp[3] >> 8)
	buf[26] = uint8(tmp[3] >> 16)
	buf[27] = uint8(tmp[3] >> 24)
	buf[28] = uint8(tmp[3] >> 32)
	buf[29] = uint8(tmp[3] >> 40)
	buf[30] = uint8(tmp[3] >> 48)
	buf[31] = uint8(tmp[3] >> 56)
	return f
}

func montRed(r0, r1, r2, r3, r4, r5, r6, r7, INV uint64, modulus Field) *Field {

	k := r0 * INV
	_, carry := futil.Mac(r0, k, modulus[0], 0)
	r1, carry = futil.Mac(r1, k, modulus[1], carry)
	r2, carry = futil.Mac(r2, k, modulus[2], carry)
	r3, carry = futil.Mac(r3, k, modulus[3], carry)
	r4, carry2 := futil.Adc(r4, 0, carry)

	k = r1 * INV
	_, carry = futil.Mac(r1, k, modulus[0], 0)
	r2, carry = futil.Mac(r2, k, modulus[1], carry)
	r3, carry = futil.Mac(r3, k, modulus[2], carry)
	r4, carry = futil.Mac(r4, k, modulus[3], carry)
	r5, carry2 = futil.Adc(r5, carry2, carry)

	k = r2 * INV
	_, carry = futil.Mac(r2, k, modulus[0], 0)
	r3, carry = futil.Mac(r3, k, modulus[1], carry)
	r4, carry = futil.Mac(r4, k, modulus[2], carry)
	r5, carry = futil.Mac(r5, k, modulus[3], carry)
	r6, carry2 = futil.Adc(r6, carry2, carry)

	k = r3 * INV
	_, carry = futil.Mac(r3, k, modulus[0], 0)
	r4, carry = futil.Mac(r4, k, modulus[1], carry)
	r5, carry = futil.Mac(r5, k, modulus[2], carry)
	r6, carry = futil.Mac(r6, k, modulus[3], carry)
	r7, carry2 = futil.Adc(r7, carry2, carry)

	f := &Field{r4, r5, r6, r7}

	f.Sub(*f, modulus, modulus)

	return f
}

func (f *Field) PowVarTime(b [4]uint64, R Field, INV uint64, modulus Field) *Field {

	res := Field{1, 0, 0, 0}
	res.Set(R)

	for j := range b {

		e := b[len(b)-1-j] // reversed
		for i := 63; i >= 0; i-- {

			res.Square(res, INV, modulus)

			if ((e >> uint64(i)) & 1) == 1 {
				res.Mul(res, *f, INV, modulus)
			}

		}

	}
	*f = res
	return f
}
