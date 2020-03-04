package fq

import (
	"bufio"
	"encoding/hex"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/decentralisedkev/go-jubjub/internal/field"
	"github.com/stretchr/testify/assert"
)

var LARGES = FieldQ{
	field.Field{0xffffffff00000000,
		0x53bda402fffe5bfe,
		0x3339d80809a1d805,
		0x73eda753299d7d48},
}

func TestBytesIntoKO(t *testing.T) {
	buf := [32]byte{}
	// buf2 := [32]byte{}

	f := FieldQ{field.Field{0, 2, 3, 4}}
	by := f.IntoBytes()
	f.BytesInto(&buf)
	assert.Equal(t, by, buf[:])
}

func TestSubKO(t *testing.T) {
	var temp = LARGES
	temp.Sub(temp, LARGES)
	assert.Equal(t, zeroF, temp.Field)

	temp.Field = zeroF
	temp.Sub(temp, LARGES)

	var temp2 = qMod
	temp2.Sub(temp2, LARGES.Field, qMod)

	assert.Equal(t, temp.Field, temp2)

	file, err := os.Open("testdata/fqsub_test.txt")
	assert.Equal(t, nil, err)
	defer file.Close()

	r := bufio.NewReader(file)
	line, e := readLn(r)
	for e == nil {
		s := strings.Split(line, ",")
		a1, err := strconv.ParseUint(s[0], 10, 64)
		assert.Equal(t, nil, err)
		a2, err := strconv.ParseUint(s[1], 10, 64)
		assert.Equal(t, nil, err)
		a3, err := strconv.ParseUint(s[2], 10, 64)
		assert.Equal(t, nil, err)
		a4, err := strconv.ParseUint(s[3], 10, 64)
		assert.Equal(t, nil, err)
		b1, err := strconv.ParseUint(s[4], 10, 64)
		assert.Equal(t, nil, err)
		b2, err := strconv.ParseUint(s[5], 10, 64)
		assert.Equal(t, nil, err)
		b3, err := strconv.ParseUint(s[6], 10, 64)
		assert.Equal(t, nil, err)
		b4, err := strconv.ParseUint(s[7], 10, 64)
		assert.Equal(t, nil, err)
		expectedAVec := s[8]
		expectedBVec := s[9]
		expectedRes := s[10]

		a := FieldQ{field.Field{a1, a2, a3, a4}}
		aRev := reverse(a.IntoBytes())
		assert.Equal(t, expectedAVec, hex.EncodeToString(aRev))

		b := FieldQ{field.Field{b1, b2, b3, b4}}
		bRev := reverse(b.IntoBytes())
		assert.Equal(t, expectedBVec, hex.EncodeToString(bRev))

		a.Sub(a, b)
		resRev := reverse(a.IntoBytes())
		assert.Equal(t, expectedRes, hex.EncodeToString(resRev))

		line, e = readLn(r)
	}

}

func TestAdditionKO(t *testing.T) {
	var temp = LARGES
	temp.Add(temp, FieldQ{field.Field{1, 0, 0, 0}})
	assert.Equal(t, true, field.Equal(field.Zero, temp.Field))

	temp = LARGES
	temp.Add(temp, temp)

	expected := &FieldQ{
		field.Field{
			0xfffffffeffffffff,
			0x53bda402fffe5bfe,
			0x3339d80809a1d805,
			0x73eda753299d7d48,
		},
	}
	assert.Equal(t, true, field.Equal(expected.Field, temp.Field))

	file, err := os.Open("testdata/fqadd_test.txt")
	assert.Equal(t, nil, err)
	defer file.Close()

	r := bufio.NewReader(file)
	line, e := readLn(r)
	for e == nil {
		s := strings.Split(line, ",")
		a1, err := strconv.ParseUint(s[0], 10, 64)
		assert.Equal(t, nil, err)
		a2, err := strconv.ParseUint(s[1], 10, 64)
		assert.Equal(t, nil, err)
		a3, err := strconv.ParseUint(s[2], 10, 64)
		assert.Equal(t, nil, err)
		a4, err := strconv.ParseUint(s[3], 10, 64)
		assert.Equal(t, nil, err)
		b1, err := strconv.ParseUint(s[4], 10, 64)
		assert.Equal(t, nil, err)
		b2, err := strconv.ParseUint(s[5], 10, 64)
		assert.Equal(t, nil, err)
		b3, err := strconv.ParseUint(s[6], 10, 64)
		assert.Equal(t, nil, err)
		b4, err := strconv.ParseUint(s[7], 10, 64)
		assert.Equal(t, nil, err)
		expectedAVec := s[8]
		expectedBVec := s[9]
		expectedRes := s[10]

		a := FieldQ{field.Field{a1, a2, a3, a4}}
		aRev := reverse(a.IntoBytes())
		assert.Equal(t, expectedAVec, hex.EncodeToString(aRev))

		b := FieldQ{field.Field{b1, b2, b3, b4}}
		bRev := reverse(b.IntoBytes())
		assert.Equal(t, expectedBVec, hex.EncodeToString(bRev))

		res := a.Add(b, a)
		resRev := reverse(res.IntoBytes())
		assert.Equal(t, expectedRes, hex.EncodeToString(resRev))

		line, e = readLn(r)
	}

}

func TestSquareKO(t *testing.T) {

	file, err := os.Open("testdata/fqsq_test.txt")
	assert.Equal(t, nil, err)
	defer file.Close()

	r := bufio.NewReader(file)
	line, e := readLn(r)
	for e == nil {
		s := strings.Split(line, ",")
		a1, err := strconv.ParseUint(s[0], 10, 64)
		assert.Equal(t, nil, err)
		a2, err := strconv.ParseUint(s[1], 10, 64)
		assert.Equal(t, nil, err)
		a3, err := strconv.ParseUint(s[2], 10, 64)
		assert.Equal(t, nil, err)
		a4, err := strconv.ParseUint(s[3], 10, 64)
		assert.Equal(t, nil, err)
		expectedAVec := s[4]
		expectedRes := s[5]

		a := FieldQ{field.Field{a1, a2, a3, a4}}
		aRev := reverse(a.IntoBytes())
		assert.Equal(t, expectedAVec, hex.EncodeToString(aRev))

		a.Square(a)
		resRev := reverse(a.IntoBytes())
		assert.Equal(t, expectedRes, hex.EncodeToString(resRev))

		line, e = readLn(r)
	}

}

func TestNegationKO(t *testing.T) {
	var temp = LARGES
	temp.Neg(temp)
	assert.Equal(t, true, field.Equal(temp.Field, field.Field{1, 0, 0, 0}))

	temp.SetZero()
	temp.Neg(temp)
	assert.Equal(t, true, field.Equal(temp.Field, field.Zero))

	temp = FieldQ{field.Field{1, 0, 0, 0}}
	temp.Neg(temp)
	assert.Equal(t, true, field.Equal(temp.Field, LARGES.Field))
}

func TestInverseKO(t *testing.T) {

	file, err := os.Open("testdata/fqinv_test.txt")
	assert.Equal(t, nil, err)
	defer file.Close()

	r := bufio.NewReader(file)
	line, e := readLn(r)
	for e == nil {
		s := strings.Split(line, ",")
		a1, err := strconv.ParseUint(s[0], 10, 64)
		assert.Equal(t, nil, err)
		a2, err := strconv.ParseUint(s[1], 10, 64)
		assert.Equal(t, nil, err)
		a3, err := strconv.ParseUint(s[2], 10, 64)
		assert.Equal(t, nil, err)
		a4, err := strconv.ParseUint(s[3], 10, 64)
		assert.Equal(t, nil, err)
		expectedAVec := s[4]
		expectedRes := s[5]

		a := FieldQ{field.Field{a1, a2, a3, a4}}
		aRev := reverse(a.IntoBytes())
		assert.Equal(t, expectedAVec, hex.EncodeToString(aRev))

		a.Inverse(a)
		resRev := reverse(a.IntoBytes())
		assert.Equal(t, expectedRes, hex.EncodeToString(resRev))

		line, e = readLn(r)
	}

}

// Compute -(q^{-1} mod 2^64) mod 2^64 by exponentiating by totient(2**64) - 1
func TestInvKO(t *testing.T) {

	var inv uint64 = 1

	for i := 0; i < 63; i++ {
		inv = inv * inv
		inv = inv * qMod[0]
	}

	inv = -inv
	assert.Equal(t, inv, INV)
}
