package curve

import (
	"fmt"

	fq "github.com/decentralisedkev/go-jubjub/internal/Fq"
	"github.com/decentralisedkev/go-jubjub/internal/field"
)

// ExtendedPoint represents the affine point `(u/z, v/z)` with
/// `z` nonzero and `t1 * t2 = uv/z`.
type ExtendedPoint struct {
	u, v, z, t1, t2 fq.FieldQ
}

// Equal returns true if a and b are equal
// => (a.u * b.z) = (b.u * a.z) & (a.v * b.z) = (b.v * a.z)
func (e *ExtendedPoint) Equal(a, b ExtendedPoint) bool {

	// XXX: Optimise later on, one is to remove so many mallocs

	var c1, c2, c3, c4 *fq.FieldQ

	c1.Mul(a.u, b.z)
	c2.Mul(b.u, a.z)
	c3.Mul(a.v, b.z)
	c4.Mul(b.v, a.z)

	return field.Equal(c1.Field, c2.Field) && field.Equal(c3.Field, c4.Field)
}

func (e *ExtendedPoint) SetZero() *ExtendedPoint {
	e.u.SetZero()
	e.t1.SetZero()
	e.t2.SetZero()
	e.v.SetOne()
	e.z.SetOne()
	return e
}

// Neg negates the u and t1 value
// returning point (-u, v, z, -t1, t2)
func (e *ExtendedPoint) Neg(a ExtendedPoint) *ExtendedPoint {
	e.u.Neg(a.u)
	e.t1.Neg(a.t1)
	e.t2.Set(a.t2)
	e.v.Set(a.v)
	e.z.Set(a.z)
	return e
}

// Double doubles the point
func (e *ExtendedPoint) Double() *ExtendedPoint {
	// XXX: impl projective points and re-do formula
	// Or do it in Completed

	var uu, vv, zz2, uv2, vvPlusUU, vvMinusUU fq.FieldQ

	uu.Square(e.u)
	vv.Square(e.v)

	zz2.Square(e.z)
	zz2.Double()

	uv2.Add(e.u, e.v)
	uv2.Square(uv2)

	vvPlusUU.Add(vv, uu)
	vvMinusUU.Sub(vv, uu)

	uv2.Sub(uv2, vvPlusUU)
	zz2.Sub(zz2, vvMinusUU)

	cp := CompletedPoint{
		u: uv2,
		v: vvPlusUU,
		z: vvMinusUU,
		t: zz2,
	}

	e.SetCompleted(cp)

	return e
}

// MulCof multiplies the point by the cofactor of 8
func (e *ExtendedPoint) MulCof() *ExtendedPoint {
	e.Double()
	e.Double()
	e.Double()
	return e
}

// Identity returns the identity point
func (e *ExtendedPoint) Identity() *ExtendedPoint {
	e.u.Neg(e.u)
	e.t1.Neg(e.t1)
	return e
}

// SetAffine sets the Affine Point af, to an Extended Point
func (e *ExtendedPoint) SetAffine(af AffinePoint) *ExtendedPoint {

	e.u.Set(af.u)
	e.t1.Set(af.u)

	e.v.Set(af.v)
	e.t2.Set(af.v)

	e.z.SetOne()

	return e
}

// SetCompleted sets the completedPoint c to ExtendedPoint e
func (e *ExtendedPoint) SetCompleted(c CompletedPoint) *ExtendedPoint {
	e.u.Mul(c.u, c.t)
	e.v.Mul(c.v, c.z)
	e.z.Mul(c.t, c.z)

	e.t1.Set(c.u)
	e.t2.Set(c.v)
	return e
}

func (e *ExtendedPoint) Bytes() *ExtendedPoint {
	e.u.Mul(c.u, c.t)
	e.v.Mul(c.v, c.z)
	e.z.Mul(c.t, c.z)

	e.t1.Set(c.u)
	e.t2.Set(c.v)
	return e
}

func (e *ExtendedPoint) isOnCurveVarTime() bool {

	// XXX: Remove once we change everything from pointers to values, to avoid nil dereferencing
	var af AffinePoint
	af.SetExtended(e)

	var s, t12 fq.FieldQ

	s.Mul(af.u, af.v)
	s.Mul(s, e.z)
	t12.Mul(e.t1, e.t2)

	fmt.Println("Not Zero", !e.z.IsZero())
	fmt.Println("affine on curve", af.isOnCurveVarTime())

	return !e.z.IsZero() && af.isOnCurveVarTime() && field.Equal(t12.Field, s.Field)
}

// To get u
/*
-u^2 + v^2 = 1 + d.u^2.v^2
v^2 - 1 = d.u^2.v^2 + u^2
v^2 - 1 = u^2 (d.v^2 + 1)
(v^2 - 1) / (d.v^2 + 1) = u^2
u^2 = (v^2 - 1) / (d.v^2 + 1)
*/
// Will be replaced with Elligator2: https://elligator.cr.yp.to/elligator-20130828.pdf
// Uses a variation of the try-and-increment method, therefore it is var-time, see 1.1: https://eprint.iacr.org/2009/226.pdf
func (e *ExtendedPoint) FromBytes(byt [64]byte) *ExtendedPoint {

	var num, den, d fq.FieldQ
	d.SetD()

	// z = 1
	e.z.SetOne()

	// v = Reduced hash digest
	e.v.FromBytes(byt)

	for { // XXX: I omitted the security parameter k

		num.Square(e.v) // v^2

		den.Set(num)

		den.Mul(d, den) // v^2 * D

		den.Add(den, e.z) //(v^2 * D + 1) // Only used e.z because it was equal to one

		num.Sub(num, e.z) // v^2 -1 // Only used e.z because it was equal to one

		den.Inverse(den) // 1 / (v^2 * D + 1)

		e.u.Mul(num, den) // v^2 - 1 / (v^2 * D + 1)

		_, err := e.u.SqrtVarTime()
		if err == nil {
			break
		}
		e.v.Add(e.v, e.z) // increment
	}
	// determine sign of u
	if (e.v.Field[3] >> 63 & 1) != (e.u.IsOdd()) {
		e.u.Neg(e.u)
	}

	//  t = u*v
	e.t1.Set(e.u)
	e.t2.Set(e.v)

	return e
}

// MulScalar TODO
// Take in a generic Field element
// XXX: Optimise to use ExtendedNiels
func (e *ExtendedPoint) MulScalar(point ExtendedPoint, scalar [32]byte) *ExtendedPoint {

	var res ExtendedPoint
	res.SetZero()

	var cp CompletedPoint

	fmt.Println("Scalar", scalar)

	for _, byt := range scalar {
		for i := byte(8); i > 0; i-- { // Iterate over bit XXX: We need to double check the endianess

			bit := byt >> (i - 1) & 1

			// Double
			res.Double()

			if bit == 1 {

				// Add
				cp.AddExtended(res, point)
				res.SetCompleted(cp)

			}
		}
	}

	*e = res

	return e
}

type ExtendedNielsPoint struct {
	vPlusU, VminusU, z, t2d fq.FieldQ
}

// Identity sets the afn to the identity point
func (en *ExtendedNielsPoint) Identity() *ExtendedNielsPoint {
	en.VminusU.SetOne()
	en.vPlusU.SetOne()
	en.z.SetOne()
	en.t2d.SetZero()
	return en
}

// SetExtended Sets ExtendedNielsPoint from an ExtendedPoint
// Check
func (en *ExtendedNielsPoint) SetExtended(e ExtendedPoint) *ExtendedNielsPoint {

	var d2 fq.FieldQ
	d2.SetD2()

	en.VminusU.Sub(e.v, e.u)
	en.vPlusU.Add(e.v, e.u)
	en.z.Set(e.z)
	en.t2d.Mul(e.u, e.v)
	en.t2d.Mul(en.t2d, d2)
	return en
}

// Check
func (p *ExtendedPoint) AddExtendedNiels(q *ExtendedNielsPoint) *ExtendedPoint {
	var r CompletedPoint
	var t0 fq.FieldQ

	// FeAdd(&r.X, &p.Y, &p.X)
	r.u.Add(p.v, p.u) // u = v+u
	// FeSub(&r.Y, &p.Y, &p.X)
	r.v.Sub(p.v, p.u)
	// FeMul(&r.Z, &r.X, &q.yPlusX)
	r.z.Mul(r.u, q.vPlusU)
	// FeMul(&r.Y, &r.Y, &q.yMinusX)
	r.v.Mul(r.v, q.VminusU)
	// FeMul(&r.T, &q.T2d, &p.T)
	r.t.Mul(p.t1, q.t2d)
	r.t.Mul(r.t, p.t2)
	// FeMul(&r.X, &p.Z, &q.Z)
	r.u.Mul(p.z, q.z)
	// FeAdd(&t0, &r.X, &r.X)
	t0.Add(r.u, r.u)
	// FeSub(&r.X, &r.Z, &r.Y)
	r.u.Sub(r.z, r.v)
	// FeAdd(&r.Y, &r.Z, &r.Y)
	r.v.Add(r.z, r.v)
	// FeAdd(&r.Z, &t0, &r.T)
	r.z.Add(t0, r.t)
	// FeSub(&r.T, &t0, &r.T)
	r.t.Sub(t0, r.t)

	p.SetCompleted(r)
	return p
}

// Check
func (p *ExtendedPoint) AddAffineNiels(q *AffineNielsPoint) *ExtendedPoint {
	var r CompletedPoint
	var t0 fq.FieldQ

	r.u.Add(p.v, p.u)
	r.v.Sub(p.v, p.u)
	r.z.Mul(r.u, q.vPlusU)
	r.v.Mul(r.v, q.VminusU)
	r.t.Mul(q.t2d, p.t1)
	r.t.Mul(r.t, p.t2)
	t0.Add(p.z, p.z)
	r.u.Sub(r.z, r.v)
	r.v.Add(r.z, r.v)
	r.z.Add(t0, r.t)
	r.t.Sub(t0, r.t)

	p.SetCompleted(r)
	return p
}

// TODO
func (en *ExtendedNielsPoint) Neg() *ExtendedNielsPoint {
	return nil
}

// Zero sets the ExtendedNielsPoint to Zero
// Check
func (en *ExtendedNielsPoint) Zero() *ExtendedNielsPoint {
	en.vPlusU.SetOne()
	en.VminusU.SetOne()
	en.z.SetOne()
	en.t2d.SetZero()
	return en
}
