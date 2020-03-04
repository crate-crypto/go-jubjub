package jubjub

import (
	"crypto/sha512"

	curve "github.com/decentralisedkev/go-jubjub/internal"
)

// Point represents a point on the JubJub curve
type Point curve.ExtendedPoint

var basePoint = base()

func base() *curve.ExtendedPoint {
	var p = &Point{}
	p.HashToPoint([]byte("jubjub"))
	return p.ep()
}

func (p *Point) ScalarMultBase(s Scalar) *Point {
	var buf [32]byte
	s.BytesInto(&buf)

	p.SetBase()
	p.ep().MulScalar(*basePoint, buf)

	return p
}

func (p *Point) HashToPoint(d []byte) *Point {
	byt := sha512.Sum512(d)
	p.ep().FromBytes(byt)
	return p
}

func (p *Point) SetBase() *Point {
	p.HashToPoint([]byte("jubjub"))
	return p
}

func (p *Point) ep() *curve.ExtendedPoint {
	return (*curve.ExtendedPoint)(p)
}

func (p *Point) Bytes() []byte {
	return  p.ep().
}
