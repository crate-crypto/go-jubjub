package fr

import (
	"encoding/binary"

	"github.com/decentralisedkev/go-jubjub/internal/futil"
)

// Fr for now represents the scalar field in Montgomery form
type Fr [4]uint64

// Zero sets f to the zero element
func (f *Fr) Zero() *Fr {
	copy(f[:], zero[:])
	return f
}

// IsZero returns true if Fr is the zero element
func (f *Fr) IsZero() bool {
	return f[0]|f[1]|f[2]|f[3] == 0
}

// One sets f to the one element
func (f *Fr) One() *Fr {

	copy(f[:], R[:])
	return f
}

// FromU64 converts a uint64 into a field element
func (f *Fr) FromU64(a uint64) *Fr {

	f[0] = a
	f[1] = 0
	f[2] = 0
	f[3] = 0

	return f.Mul(f, &R2)
}

// Double doubles f by adding it to itself
func (f *Fr) Double() *Fr {
	f.Add(f, f)
	return f
}

// Equal returns true, if a ==b
func (f *Fr) Equal(a, b *Fr) bool {
	return a[0] == b[0] && a[1] == b[1] && a[2] == b[2] && a[3] == b[3]
}

// ConditionalSel selects and sets a if c = 1 or b if c =  0
func (f *Fr) ConditionalSel(a, b *Fr, c uint64) *Fr {

	f[0] = ^(c-1)&a[0] | (c-1)&b[0]
	f[1] = ^(c-1)&a[1] | (c-1)&b[1]
	f[2] = ^(c-1)&a[2] | (c-1)&b[2]
	f[3] = ^(c-1)&a[3] | (c-1)&b[3]

	return f
}

// ConditionalSet sets f to a if b == 1
// Taken from Go-ristretto (Bas)
func (f *Fr) ConditionalSet(a *Fr, b uint64) *Fr {

	b = -b // b == 0b11111111111111111111111111111111 or 0.
	f[0] ^= b & (f[0] ^ a[0])
	f[1] ^= b & (f[1] ^ a[1])
	f[2] ^= b & (f[2] ^ a[2])
	f[3] ^= b & (f[3] ^ a[3])

	return f
}

// Neg negates a Fr
func (f *Fr) Neg(a *Fr) *Fr {

	d0, borrow := futil.Sbb(r[0], a[0], 0)
	d1, borrow := futil.Sbb(r[1], a[1], borrow)
	d2, borrow := futil.Sbb(r[2], a[2], borrow)
	d3, _ := futil.Sbb(r[3], a[3], borrow)

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

// Sub Subtracts one field from another
func (f *Fr) Sub(lhs, rhs *Fr) *Fr {

	d0, borrow := futil.Sbb(lhs[0], rhs[0], 0)
	d1, borrow := futil.Sbb(lhs[1], rhs[1], borrow)
	d2, borrow := futil.Sbb(lhs[2], rhs[2], borrow)
	d3, borrow := futil.Sbb(lhs[3], rhs[3], borrow)

	// If underflow occurred on the final limb, borrow = 0xfff...fff, otherwise
	// // borrow = 0x000...000. Thus, we use it as a mask to conditionally add the modulus.
	d0, carry := futil.Adc(d0, r[0]&borrow, 0)
	d1, carry = futil.Adc(d1, r[1]&borrow, carry)
	d2, carry = futil.Adc(d2, r[2]&borrow, carry)
	d3, carry = futil.Adc(d3, r[3]&borrow, carry)

	f[0] = d0
	f[1] = d1
	f[2] = d2
	f[3] = d3

	return f
}

// SubNoBorrow Subtracts one field from another
// and ignore the borrow
func (f *Fr) SubNoBorrow(lhs, rhs *Fr) *Fr {

	f[0], _ = futil.Sbb(lhs[0], rhs[0], 0)
	f[1], _ = futil.Sbb(lhs[1], rhs[1], 0)
	f[2], _ = futil.Sbb(lhs[2], rhs[2], 0)
	f[3], _ = futil.Sbb(lhs[3], rhs[3], 0)

	return f
}
func (f *Fr) reduce() {
	if f.Cmp(f, &r) >= 0 { // only reduce if f is >= modulus
		f.SubNoBorrow(f, &r)
	}
}

// Add Adds one field to another
func (f *Fr) Add(lhs, rhs *Fr) *Fr {

	d0, carry := futil.Adc(lhs[0], rhs[0], 0)
	d1, carry := futil.Adc(lhs[1], rhs[1], carry)
	d2, carry := futil.Adc(lhs[2], rhs[2], carry)
	d3, _ := futil.Adc(lhs[3], rhs[3], carry)

	f[0] = d0
	f[1] = d1
	f[2] = d2
	f[3] = d3

	// Normalise
	f.Sub(f, &r)

	return f
}

// Cmp compares a and b
// if a > b return 1
// if a==b return 0
// if a < b return -1
func (f *Fr) Cmp(a, b *Fr) int8 {
	for i := 3; i >= 0; i-- {
		if a[i] > b[i] {
			return 1
		} else if b[i] > a[i] {
			return -1
		}
	}
	return 0
}

// Mul mutiplies two field elements together
func (f *Fr) Mul(lhs, rhs *Fr) *Fr {
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

	red := MontRed(r0, r1, r2, r3, r4, r5, r6, r7)
	f[0] = red[0]
	f[1] = red[1]
	f[2] = red[2]
	f[3] = red[3]
	return f
}

// MulAdd multiples two numbers and adds the third such that f = a*b +c
func (f *Fr) MulAdd(a, b, c *Fr) *Fr {
	// We can optimise this later, then have Mul = MulAdd(a,b,zero)
	f.Mul(a, b)
	f.Add(f, c)
	return f
}

// MulSub multiples two numbers and subs the third such that f = a*b -c
func (f *Fr) MulSub(a, b, c *Fr) *Fr {
	f.Mul(a, b)
	f.Sub(f, c)
	return f
}

func (f *Fr) Square(a *Fr) *Fr {

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

	red := MontRed(r0, r1, r2, r3, r4, r5, r6, r7)
	f[0] = red[0]
	f[1] = red[1]
	f[2] = red[2]
	f[3] = red[3]
	return f
}

func (f *Fr) LegendreSymbolVarTime() *Fr {
	// Legendre symbol computed via Euler's criterion:
	// self^((q - 1) // 2)
	f.PowVarTime([4]uint64{
		0x7fffffff80000000,
		0xa9ded2017fff2dff,
		0x199cec0404d0ec02,
		0x39f6d3a994cebea4,
	})
	return f
}

func (f *Fr) PowVarTime(b [4]uint64) *Fr {

	res := f.One()
	for j := range b {
		e := b[len(b)-1-j] // reverse
		for i := uint64(64); i > 0; i++ {
			res.Square(res)

			if ((e >> i) & 1) == 1 {
				res.Mul(res, res)
			}
		}
	}
	return res
}

func (f *Fr) SqrtVarTime() *Fr {
	return nil
}

// Inverse inverts a field element
// If element is zero, it will return nil
// XXX: TODO This is wrong for Fr, redo
func (f *Fr) Inverse(a *Fr) *Fr {

	var zero Fr
	zero.Zero()

	// Check if f is non-zero
	if f.Equal(f, &zero) {
		return nil
	}

	var sqrMulti = func(e *Fr, n uint64) {
		for i := uint64(0); i < n; i++ {
			e.Square(e)
		}
	}

	var t0, t1, t2, t3, t4, t5, t6, t7, t8, t9, t10, t11, t12, t13, t14, t15, t16, t17 Fr

	t10 = *a
	t0.Square(&t10)
	t1.Mul(&t0, &t10)
	t16.Square(&t0)
	t6.Square(&t16)
	t5.Mul(&t6, &t0)
	t0.Mul(&t6, &t16)
	t12.Mul(&t5, &t16)
	t2.Square(&t6)
	t7.Mul(&t5, &t6)
	t15.Mul(&t0, &t5)
	t17.Square(&t12)
	t1.Mul(&t1, &t17)
	t3.Mul(&t7, &t2)
	t8.Mul(&t1, &t17)
	t4.Mul(&t8, &t2)
	t9.Mul(&t8, &t7)
	t7.Mul(&t4, &t5)
	t11.Mul(&t4, &t17)
	t5.Mul(&t9, &t17)
	t14.Mul(&t7, &t15)
	t13.Mul(&t11, &t12)
	t12.Mul(&t11, &t17)
	t15.Mul(&t15, &t12)
	t16.Mul(&t16, &t15)
	t3.Mul(&t3, &t16)
	t17.Mul(&t17, &t3)
	t0.Mul(&t0, &t17)
	t6.Mul(&t6, &t0)
	t2.Mul(&t2, &t6)
	sqrMulti(&t0, 8)
	t0.Mul(&t0, &t17)
	sqrMulti(&t0, 9)
	t0.Mul(&t0, &t16)
	sqrMulti(&t0, 9)
	t0.Mul(&t0, &t15)
	sqrMulti(&t0, 9)
	t0.Mul(&t0, &t15)
	sqrMulti(&t0, 7)
	t0.Mul(&t0, &t14)
	sqrMulti(&t0, 7)
	t0.Mul(&t0, &t13)
	sqrMulti(&t0, 10)
	t0.Mul(&t0, &t12)
	sqrMulti(&t0, 9)
	t0.Mul(&t0, &t11)
	sqrMulti(&t0, 8)
	t0.Mul(&t0, &t8)
	sqrMulti(&t0, 8)
	t0.Mul(&t0, &t10)
	sqrMulti(&t0, 14)
	t0.Mul(&t0, &t9)
	sqrMulti(&t0, 10)
	t0.Mul(&t0, &t8)
	sqrMulti(&t0, 15)
	t0.Mul(&t0, &t7)
	sqrMulti(&t0, 10)
	t0.Mul(&t0, &t6)
	sqrMulti(&t0, 8)
	t0.Mul(&t0, &t5)
	sqrMulti(&t0, 16)
	t0.Mul(&t0, &t3)
	sqrMulti(&t0, 8)
	t0.Mul(&t0, &t2)
	sqrMulti(&t0, 7)
	t0.Mul(&t0, &t4)
	sqrMulti(&t0, 9)
	t0.Mul(&t0, &t2)
	sqrMulti(&t0, 8)
	t0.Mul(&t0, &t3)
	sqrMulti(&t0, 8)
	t0.Mul(&t0, &t2)
	sqrMulti(&t0, 8)
	t0.Mul(&t0, &t2)
	sqrMulti(&t0, 8)
	t0.Mul(&t0, &t2)
	sqrMulti(&t0, 8)
	t0.Mul(&t0, &t3)
	sqrMulti(&t0, 8)
	t0.Mul(&t0, &t2)
	sqrMulti(&t0, 8)
	t0.Mul(&t0, &t2)
	sqrMulti(&t0, 5)
	t0.Mul(&t0, &t1)
	sqrMulti(&t0, 5)
	t0.Mul(&t0, &t1)

	f[0] = t0[0]
	f[1] = t0[1]
	f[2] = t0[2]
	f[3] = t0[3]
	return f
}

// IntoBytes  converts f into a little endian byte slice
func (f *Fr) IntoBytes() []byte {

	// Turn into canonical form by computing (a.R) / R = a
	tmp := MontRed(f[0], f[1], f[2], f[3], 0, 0, 0, 0)

	res := make([]byte, 32, 32)

	binary.LittleEndian.PutUint64(res[0:8], tmp[0])
	binary.LittleEndian.PutUint64(res[8:16], tmp[1])
	binary.LittleEndian.PutUint64(res[16:24], tmp[2])
	binary.LittleEndian.PutUint64(res[24:32], tmp[3])

	return res
}

// FromBytesVarTime takes a byte slice and returns it's canonical
// representation in Fr. Does not check length of slice
func (f *Fr) FromBytesVarTime(b []byte) *Fr {
	f[0] = binary.LittleEndian.Uint64(b[0:8])
	f[1] = binary.LittleEndian.Uint64(b[8:16])
	f[2] = binary.LittleEndian.Uint64(b[16:24])
	f[3] = binary.LittleEndian.Uint64(b[24:32])
	return f
}

// Sub Subtracts one field from another
// taken from : https://github.com/dalek-cryptography/curve25519-dalek/blob/3a5aef7eeab697e86762765072e1a81763087936/src/backend/u64/scalar.rs#L78
func (f *Fr) Reduce(a *[64]byte) *Fr {

	var words = [8]uint64{}
	for i := uint64(0); i < uint64(len(words)); i++ {
		for j := uint64(0); j < uint64(len(words)); j++ {
			words[i] |= uint64((a[(i*8)+j])) << (j * 8)
		}
	}

	var one uint64
	mask := (one << 52) - 1
	var hi, lo Fr

	lo[0] = words[0] & mask
	lo[1] = ((words[0] >> 52) | (words[1] << 12)) & mask
	lo[2] = ((words[1] >> 40) | (words[2] << 24)) & mask
	lo[3] = ((words[2] >> 28) | (words[3] << 36)) & mask
	hi[0] = (words[4] >> 4) & mask
	hi[1] = ((words[4] >> 56) | (words[5] << 8)) & mask
	hi[2] = ((words[5] >> 44) | (words[6] << 20)) & mask
	hi[3] = ((words[6] >> 32) | (words[7] << 32)) & mask

	// then calculate lo = lo * R , then montReduce (lo * R )/ R
	// hi = hi * R2, then montReduce = (hi*R2) / R

	// f.Mul(a, &r) // XXX:Check on usage on R2

	// // Calulate aR
	// s.Mul(&r, s)
	// // montReduce (aR)
	return nil
}

func MontRed(r0, r1, r2, r3, r4, r5, r6, r7 uint64) *Fr {

	k := r0 * INV
	_, carry := futil.Mac(r0, k, r[0], 0)
	r1, carry = futil.Mac(r1, k, r[1], carry)
	r2, carry = futil.Mac(r2, k, r[2], carry)
	r3, carry = futil.Mac(r3, k, r[3], carry)
	r4, carry2 := futil.Adc(r4, 0, carry)

	k = r1 * INV
	_, carry = futil.Mac(r1, k, r[0], 0)
	r2, carry = futil.Mac(r2, k, r[1], carry)
	r3, carry = futil.Mac(r3, k, r[2], carry)
	r4, carry = futil.Mac(r4, k, r[3], carry)
	r5, carry2 = futil.Adc(r5, carry2, carry)

	k = r2 * INV
	_, carry = futil.Mac(r2, k, r[0], 0)
	r3, carry = futil.Mac(r3, k, r[1], carry)
	r4, carry = futil.Mac(r4, k, r[2], carry)
	r5, carry = futil.Mac(r5, k, r[3], carry)
	r6, carry2 = futil.Adc(r6, carry2, carry)

	k = r3 * INV
	_, carry = futil.Mac(r3, k, r[0], 0)
	r4, carry = futil.Mac(r4, k, r[1], carry)
	r5, carry = futil.Mac(r5, k, r[2], carry)
	r6, carry = futil.Mac(r6, k, r[3], carry)
	r7, carry2 = futil.Adc(r7, carry2, carry)

	f := &Fr{r4, r5, r6, r7}

	f.Sub(f, &r)

	return f
}
