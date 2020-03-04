package curve

import (
	fq "github.com/decentralisedkev/go-jubjub/internal/Fq"
	"github.com/decentralisedkev/go-jubjub/internal/field"
)

// AffinePoint represents an affine point `(u, v)` on the
/// curve `-u^2 + v^2 = 1 + d.u^2.v^2` over `Fq` with
/// `d = -(10240/10241)`
type AffinePoint struct {
	u, v fq.FieldQ
}

// Neg negates the u value in (u,v)
// returning point (-u, v)
func (af *AffinePoint) Neg() *AffinePoint {
	af.u.Neg(af.u)
	return af
}

// Identity returns (0,1)
func (af *AffinePoint) Identity() *AffinePoint {

	af.u.SetZero()
	af.v.SetOne()

	return af
}

// Equal returns whether two affine points
//b and c are equal.
func (af *AffinePoint) Equal(b, c *AffinePoint) bool {
	return (field.Equal(af.u.Field, b.u.Field) && field.Equal(af.v.Field, b.v.Field))
}

// SetExtended sets the Extended Point e, to an Affine Point
func (af *AffinePoint) SetExtended(e *ExtendedPoint) *AffinePoint {
	// XXX: before doing inverses, we should have a check for infinity
	e.z.Inverse(e.z) // z = 1/z

	af.u.Mul(e.u, e.z)
	af.v.Mul(e.v, e.z)

	return af
}

// IntoBytes converts the af element into its little-endian
// byte representation
func (af *AffinePoint) IntoBytes() []byte {

	var tmp, u [32]byte
	af.v.BytesInto(&tmp)
	af.u.BytesInto(&u)

	// Encode the sign of the u-coordinate in the most
	// significant bit.
	tmp[31] |= u[0] << 7

	return tmp[:]
}

type AffineNielsPoint struct {
	vPlusU, VminusU, t2d fq.FieldQ
}

// Identity sets the afn to the identity point
func (afn *AffineNielsPoint) Identity() *AffineNielsPoint {
	afn.VminusU.SetOne()
	afn.vPlusU.SetOne()
	afn.t2d.SetZero()
	return afn
}

// Zero sets AffineNielsPoint to Zero
func (afn *AffineNielsPoint) Zero() *AffineNielsPoint {
	afn.vPlusU.SetOne()
	afn.VminusU.SetOne()
	afn.t2d.SetZero()
	return afn
}

// SetAffine sets the AffineNielsPoint from an AffinePoint
func (afn *AffineNielsPoint) SetAffine(af AffinePoint) *AffineNielsPoint {

	var d2 fq.FieldQ
	d2.SetD2()

	afn.vPlusU.Add(af.u, af.v)
	afn.VminusU.Sub(af.v, af.u)
	afn.t2d.Mul(af.u, af.v)
	afn.t2d.Mul(afn.t2d, d2)
	return afn
}

func (af *AffinePoint) isOnCurveVarTime() bool {

	var u2, v2, lhs, rhs, one, d fq.FieldQ
	d.SetD()

	one.SetOne()

	u2.Square(af.u) // u^2
	v2.Square(af.v) // v^2

	lhs.Sub(v2, u2) // v^2 - u^2

	rhs.Mul(u2, v2)   // u^2.v^2
	rhs.Mul(rhs, d)   // du^2.v^2
	rhs.Add(rhs, one) // 1 + du^2.v^2

	return field.Equal(lhs.Field, rhs.Field) // 1 + du^2.v^2 == v^2 - u^2
}
