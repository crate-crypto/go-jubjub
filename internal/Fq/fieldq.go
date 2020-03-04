package fq

import (
	"crypto/rand"
	"encoding/hex"
	"errors"

	"github.com/decentralisedkev/go-jubjub/internal/field"
)

type FieldQ struct {
	field.Field
}

// // INV = -(q^{-1} mod 2^64) mod 2^64
// const INV uint64 = 0xfffffffeffffffff

var (

	// R = 2^256 mod q
	montR = field.Field{0x00000001fffffffe, 0x5884b7fa00034802, 0x998c4fefecbc4ff5, 0x1824b159acc5056f}

	// R2 = 2^512 mod q
	montR2 = field.Field{0xc999e990f3f29c6d, 0x2b6cedcb87925c23, 0x05d314967254398f, 0x0748d9d99f59ff11}

	// R3 = R2 * 2^256 mod q = 2^768 mod q
	montR3 = field.Field{0xc62c1807439b73af, 0x1b3e0d188cf06990, 0x73d13c71c7b5f418, 0x6e2a5bb9c8db33e9}

	// ROOTOFUNITY GENERATOR^t where t * 2^s + 1 = q with t odd.
	rootOfUnity = field.Field{0xb9b58d8c5f0e466a, 0x5b1b4c801819d7ec, 0x0af53ae352a31e64, 0x5bf3adda19e9b27b}

	// D = -(10240/10241)
	d = field.Field{0x2a522455b974f6b0, 0xfc6cc9ef0d9acab3, 0x7a08fb94c27628d1, 0x57f8f6a8fe0e262e}

	// D2 = 2 * d
	d2 = field.Field{0x54a448ac72e9ed5f, 0xa51befdb1b373967, 0xc0d81f217b4a799e, 0x3c0445fed27ecf14}
)

func (f *FieldQ) IntoBytes() []byte {
	return f.Field.IntoBytes(INV, qMod)
}

func (f *FieldQ) BytesInto(buf *[32]byte) {
	f.Field.BytesInto(buf, qMod, INV)
}

func (f *FieldQ) Set(a FieldQ) *FieldQ {
	f.Field.Set(a.Field)
	return f
}
func (f *FieldQ) SetOne() *FieldQ {
	f.Field.Set(montR)
	return f
}
func (f *FieldQ) SetD() *FieldQ {
	f.Field.Set(d)
	return f
}
func (f *FieldQ) SetD2() *FieldQ {
	f.Field.Set(d2)
	return f
}

func (f *FieldQ) FromBytes(byt [64]byte) *FieldQ {
	f.Field.FromBytes(byt, INV, qMod, r2, r3)
	return f
}

func (f *FieldQ) PowVarTime(b [4]uint64) *FieldQ {
	f.Field.PowVarTime(b, montR, INV, qMod)
	return f
}

func (f *FieldQ) LegendreSymbolVarTime() *FieldQ {
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

func (f *FieldQ) SqrtVarTime() (*FieldQ, error) {
	one := &FieldQ{field.Field{0, 0, 0, 0}}
	one.Set(FieldQ{montR})
	zero := &FieldQ{field.Field{0, 0, 0, 0}}
	tmp := &FieldQ{field.Field{0, 0, 0, 0}}

	*tmp = *f

	lgs := tmp.LegendreSymbolVarTime()

	if field.Equal(lgs.Field, zero.Field) {
		return f, nil
	}
	if !field.Equal(lgs.Field, one.Field) {
		return f, errors.New("legendre symbol does not equal one, sqrt not found") // XXX: We should bubble up an error for this
	}

	*tmp = *f
	r := *tmp.PowVarTime([4]uint64{0x7fff2dff80000000, 0x04d0ec02a9ded201, 0x94cebea4199cec04, 0x0000000039f6d3a9})

	*tmp = *f
	t := *tmp.PowVarTime([4]uint64{0xfffe5bfeffffffff, 0x09a1d80553bda402, 0x299d7d483339d808, 0x0000000073eda753})

	c := FieldQ{rootOfUnity}
	m := S

	for !field.Equal(t.Field, one.Field) {

		var i = uint32(1)

		t2i := FieldQ{field.Field{0, 0, 0, 0}}
		t2i.Square(t)

		for !field.Equal(t2i.Field, one.Field) {
			t2i.Square(t2i)
			i++
		}

		for k := uint32(0); k < m-i-1; k++ {
			c.Square(c)
		}

		r.Mul(r, c)
		c.Square(c)
		t.Mul(t, c)
		m = i

	}

	*f = r

	return f, nil
}

// Rand returns a random field element
func (f *FieldQ) Rand() *FieldQ {
	var buf [64]byte
	rand.Read(buf[:])
	f.FromBytes(buf)
	return f
}

// FromU64 converts a uint64 into a field element
func (f *FieldQ) FromU64(a uint64) *FieldQ {
	f.Field.FromU64(a, INV, r2, qMod)
	return f
}

// Sub Subtracts one field from another
func (f *FieldQ) Sub(a, b FieldQ) *FieldQ {
	f.Field.Sub(a.Field, b.Field, qMod)
	return f
}

// Add adds one field from another
func (f *FieldQ) Add(a, b FieldQ) *FieldQ {
	f.Field.Add(a.Field, b.Field, qMod)
	return f
}

// Double doubles the field element
func (f *FieldQ) Double() *FieldQ {
	f.Field.Double(qMod)
	return f
}

// Mul multiplies one field from another
func (f *FieldQ) Mul(a, b FieldQ) *FieldQ {
	f.Field.Mul(a.Field, b.Field, INV, qMod)
	return f
}

// Square squares a field.
func (f *FieldQ) Square(a FieldQ) *FieldQ {
	f.Field.Square(a.Field, INV, qMod)
	return f
}

// Neg negates a field.
func (f *FieldQ) Neg(a FieldQ) *FieldQ {
	f.Field.Neg(a.Field, qMod)
	return f
}

func (f *FieldQ) String() string {
	var s [32]byte
	f.BytesInto(&s)

	// reverse bytes
	for i, j := 0, len(s)-1; i <= j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}

	return "0x" + hex.EncodeToString(s[:])
}

// Inverse inverts a field element
// If element is zero, it will return nil
func (f *FieldQ) Inverse(a FieldQ) *FieldQ {

	// Check if a is non-zero
	if field.Equal(a.Field, field.Zero) {
		return nil
	}

	var sqrMulti = func(e *FieldQ, n uint64) {
		for i := uint64(0); i < n; i++ {
			e.Square(*e)
		}
	}

	var t0, t1, t2, t3, t4, t5, t6, t7, t8, t9, t10, t11, t12, t13, t14, t15, t16, t17 FieldQ

	t10 = a
	t0.Square(t10)
	t1.Mul(t0, t10)
	t16.Square(t0)
	t6.Square(t16)
	t5.Mul(t6, t0)
	t0.Mul(t6, t16)
	t12.Mul(t5, t16)
	t2.Square(t6)
	t7.Mul(t5, t6)
	t15.Mul(t0, t5)
	t17.Square(t12)
	t1.Mul(t1, t17)
	t3.Mul(t7, t2)
	t8.Mul(t1, t17)
	t4.Mul(t8, t2)
	t9.Mul(t8, t7)
	t7.Mul(t4, t5)
	t11.Mul(t4, t17)
	t5.Mul(t9, t17)
	t14.Mul(t7, t15)
	t13.Mul(t11, t12)
	t12.Mul(t11, t17)
	t15.Mul(t15, t12)
	t16.Mul(t16, t15)
	t3.Mul(t3, t16)
	t17.Mul(t17, t3)
	t0.Mul(t0, t17)
	t6.Mul(t6, t0)
	t2.Mul(t2, t6)
	sqrMulti(&t0, 8)
	t0.Mul(t0, t17)
	sqrMulti(&t0, 9)
	t0.Mul(t0, t16)
	sqrMulti(&t0, 9)
	t0.Mul(t0, t15)
	sqrMulti(&t0, 9)
	t0.Mul(t0, t15)
	sqrMulti(&t0, 7)
	t0.Mul(t0, t14)
	sqrMulti(&t0, 7)
	t0.Mul(t0, t13)
	sqrMulti(&t0, 10)
	t0.Mul(t0, t12)
	sqrMulti(&t0, 9)
	t0.Mul(t0, t11)
	sqrMulti(&t0, 8)
	t0.Mul(t0, t8)
	sqrMulti(&t0, 8)
	t0.Mul(t0, t10)
	sqrMulti(&t0, 14)
	t0.Mul(t0, t9)
	sqrMulti(&t0, 10)
	t0.Mul(t0, t8)
	sqrMulti(&t0, 15)
	t0.Mul(t0, t7)
	sqrMulti(&t0, 10)
	t0.Mul(t0, t6)
	sqrMulti(&t0, 8)
	t0.Mul(t0, t5)
	sqrMulti(&t0, 16)
	t0.Mul(t0, t3)
	sqrMulti(&t0, 8)
	t0.Mul(t0, t2)
	sqrMulti(&t0, 7)
	t0.Mul(t0, t4)
	sqrMulti(&t0, 9)
	t0.Mul(t0, t2)
	sqrMulti(&t0, 8)
	t0.Mul(t0, t3)
	sqrMulti(&t0, 8)
	t0.Mul(t0, t2)
	sqrMulti(&t0, 8)
	t0.Mul(t0, t2)
	sqrMulti(&t0, 8)
	t0.Mul(t0, t2)
	sqrMulti(&t0, 8)
	t0.Mul(t0, t3)
	sqrMulti(&t0, 8)
	t0.Mul(t0, t2)
	sqrMulti(&t0, 8)
	t0.Mul(t0, t2)
	sqrMulti(&t0, 5)
	t0.Mul(t0, t1)
	sqrMulti(&t0, 5)
	t0.Mul(t0, t1)

	*f = t0

	f.Field[0] = t0.Field[0]
	f.Field[1] = t0.Field[1]
	f.Field[2] = t0.Field[2]
	f.Field[3] = t0.Field[3]

	return f
}
