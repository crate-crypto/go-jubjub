package jubjub

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"math/big"

	"github.com/decentralisedkev/go-jubjub/internal/field"
	"github.com/decentralisedkev/go-jubjub/internal/futil"
)

// Scalar represents the scalar field of Jubjub
type Scalar struct {
	field.Field
}

// INV = -(r^{-1} mod 2^64) mod 2^64
const INV uint64 = 0x1ba3a358ef788ef9

var (
	// Note: r is the modulus and R is the R related to the montgomery reduce; the mont modulus
	// R = 2^256 mod r
	// the montgomoery modulus
	montR = field.Field{0x25f80bb3b99607d9, 0xf315d62f66b6e750, 0x932514eeeb8814f4, 0x09a6fc6f479155c6}

	// R2 (Rsquared) = 2^512 mod r
	montR2 = field.Field{0x67719aa495e57731, 0x51b0cef09ce3fc26, 0x69dab7fac026e9a5, 0x04f6547b8d127688}

	// R3 (RCubed) = 2^768 mod r
	montR3 = field.Field{0xe0d6c6563d830544, 0x323e3883598d0f85, 0xf0fea3004c2e2ba8, 0x05874f84946737ec}

	// r is modulus in the scalar
	// r = 0x0e7db4ea6533afa906673b0101343b00a6682093ccc81082d0970e5ed6f72cb7
	rMod = field.Field{0xd0970e5ed6f72cb7, 0xa6682093ccc81082, 0x06673b0101343b00, 0x0e7db4ea6533afa9}

	// NEG1 = -R = -(2^256 mod r) mod r
	NEG1 = field.Field{0xaa9f02ab1d6124de, 0xb3524a6466112932, 0x7342261215ac260b, 0x4d6b87b1da259e2}
)

func (s *Scalar) BytesInto(buf *[32]byte) {
	s.Field.BytesInto(buf, rMod, INV)
}

// Rand returns a random scalar in the Scalar field
func (s *Scalar) Rand() *Scalar {
	var buf [64]byte
	rand.Read(buf[:])
	s.FromBytes(buf)
	return s
}

func (s *Scalar) FromBytes(byt [64]byte) *Scalar {
	s.Field.FromBytes(byt, INV, rMod, montR2, montR3)
	return s
}

// String returns the string representation of s
func (s *Scalar) String() string {
	var buf [32]byte
	s.BytesInto(&buf)

	// reverse bytes
	for i, j := 0, len(buf)-1; i <= j; i, j = i+1, j-1 {
		buf[i], buf[j] = buf[j], buf[i]
	}

	return hex.EncodeToString(buf[:])
}

// Add adds two scalars together s.t. s = a+b
func (s *Scalar) Add(a, b Scalar) *Scalar {
	s.Field.Add(a.Field, b.Field, rMod)
	return s
}

func (s *Scalar) Double() *Scalar {
	s.Field.Double(rMod)
	return s
}

func (s *Scalar) BigInt() *big.Int {
	var buf [32]byte
	s.BytesInto(&buf)

	// reverse bytes
	for i, j := 0, len(buf)-1; i <= j; i, j = i+1, j-1 {
		buf[i], buf[j] = buf[j], buf[i]
	}

	bi := big.NewInt(0).SetBytes(buf[:])
	return bi
}

// Sub subtracts two scalars s.t. s = a-b
func (s *Scalar) Sub(a, b Scalar) *Scalar {
	s.Field.Sub(a.Field, b.Field, rMod)
	return s
}

// Mul Multiplies two scalars s.t. s = a * b
func (s *Scalar) Mul(a, b Scalar) *Scalar {
	s.Field.Mul(a.Field, b.Field, INV, rMod)
	return s
}

// // Inverse returns the inverse s.t. s = 1/a
// func (s *Scalar) Inverse(a *Scalar) *Scalar {
// 	return nil
// }

// Neg returns the Negation of a scalar s.t. s = -a
func (s *Scalar) Neg(a Scalar) *Scalar {
	s.Field.Neg(a.Field, rMod)
	return s
}

// SetZero sets s = 0
func (s *Scalar) SetZero() *Scalar {
	s.Field.SetZero()
	return s
}

// SetOne sets s = 1
func (s *Scalar) SetOne() *Scalar {
	s.Field.Set(montR)
	return s
}

// Set sets s to a Field value
func (s *Scalar) SetField(a field.Field) *Scalar {
	s.Field.Set(a)
	return s
}

// Set sets s to scalar `a`
func (s *Scalar) Set(a Scalar) *Scalar {
	s.Field.Set(a.Field)
	return s
}

// HashToScalar hashes the slice d into a scalar returning s mod R
func (s *Scalar) HashToScalar(d []byte) *Scalar {
	byt := sha512.Sum512(d)
	s.Field.FromBytes(byt, INV, rMod, montR2, montR3)
	return s
}

// Square sets s = a * a
func (s *Scalar) Square(a Scalar) *Scalar {
	s.Field.Square(a.Field, INV, rMod)
	return s
}

// SetReduce sets s = a mod r
func (s *Scalar) Reduce(a Scalar) *Scalar {
	if field.Cmp(a.Field, rMod) >= 0 {
		s.Field.Sub(a.Field, rMod, rMod)
	}
	return s
}

// MulAdd multiplies and adds three numbers s.t. s= a*b +c
func (s *Scalar) MulAdd(a, b, c Scalar) *Scalar {
	s.Mul(a, b).Add(*s, c)
	return s
}

// MulSub multiplies and subtracts three numbers s.t. s= a*b -c
func (s *Scalar) MulSub(a, b, c Scalar) *Scalar {
	s.Mul(a, b).Sub(*s, c)
	return s
}

// XXX: Not working, should unmarshall and marshall and it gives the same Scalar onject
func (s *Scalar) SetBytes(buf *[32]byte) *Scalar {
	s.Field[0] = futil.Load4(buf[0:])
	s.Field[1] = futil.Load4(buf[8:])
	s.Field[2] = futil.Load4(buf[16:])
	s.Field[3] = futil.Load4(buf[24:]) & 0x1fffffff
	s.Field.Sub(s.Field, rMod, rMod)
	return s
}
