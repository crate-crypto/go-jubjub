package fq

import (
	"bufio"
	"strings"
	"testing"
)

// var LARGEST = Fq{
// 	0xffffffff00000000,
// 	0x53bda402fffe5bfe,
// 	0x3339d80809a1d805,
// 	0x73eda753299d7d48,
// }

// func TestEquality(t *testing.T) {
// 	z := Fq{}
// 	z.SetZero()
// 	assert.Equal(t, zero, z)

// 	one := Fq{}
// 	one.SetOne()
// 	assert.Equal(t, one, R)
// }

// func TestAddition(t *testing.T) {
// 	var temp = LARGEST
// 	z := temp.Add(&temp, &Fq{1, 0, 0, 0})
// 	assert.Equal(t, true, z.Equal(&zero, z))

// 	temp = LARGEST
// 	temp.Add(&temp, &temp)

// 	expected := &Fq{
// 		0xfffffffeffffffff,
// 		0x53bda402fffe5bfe,
// 		0x3339d80809a1d805,
// 		0x73eda753299d7d48,
// 	}
// 	assert.Equal(t, true, expected.Equal(expected, &temp))

// 	file, err := os.Open("testdata/fqadd_test.txt")
// 	assert.Equal(t, nil, err)
// 	defer file.Close()

// 	r := bufio.NewReader(file)
// 	line, e := readLn(r)
// 	for e == nil {
// 		s := strings.Split(line, ",")
// 		a1, err := strconv.ParseUint(s[0], 10, 64)
// 		assert.Equal(t, nil, err)
// 		a2, err := strconv.ParseUint(s[1], 10, 64)
// 		assert.Equal(t, nil, err)
// 		a3, err := strconv.ParseUint(s[2], 10, 64)
// 		assert.Equal(t, nil, err)
// 		a4, err := strconv.ParseUint(s[3], 10, 64)
// 		assert.Equal(t, nil, err)
// 		b1, err := strconv.ParseUint(s[4], 10, 64)
// 		assert.Equal(t, nil, err)
// 		b2, err := strconv.ParseUint(s[5], 10, 64)
// 		assert.Equal(t, nil, err)
// 		b3, err := strconv.ParseUint(s[6], 10, 64)
// 		assert.Equal(t, nil, err)
// 		b4, err := strconv.ParseUint(s[7], 10, 64)
// 		assert.Equal(t, nil, err)
// 		expectedAVec := s[8]
// 		expectedBVec := s[9]
// 		expectedRes := s[10]

// 		a := &Fq{a1, a2, a3, a4}
// 		aRev := reverse(a.intoBytes())
// 		assert.Equal(t, expectedAVec, hex.EncodeToString(aRev))

// 		b := &Fq{b1, b2, b3, b4}
// 		bRev := reverse(b.intoBytes())
// 		assert.Equal(t, expectedBVec, hex.EncodeToString(bRev))

// 		res := a.Add(b, a)
// 		resRev := reverse(res.intoBytes())
// 		assert.Equal(t, expectedRes, hex.EncodeToString(resRev))

// 		line, e = readLn(r)
// 	}

// }

// func TestSub(t *testing.T) {
// 	var temp = LARGEST
// 	temp.Sub(&temp, &LARGEST)
// 	assert.Equal(t, zero, temp)

// 	temp = zero
// 	temp.Sub(&temp, &LARGEST)

// 	var temp2 = q
// 	temp2.Sub(&temp2, &LARGEST)

// 	assert.Equal(t, true, temp.Equal(&temp, &temp2))

// 	file, err := os.Open("testdata/fqsub_test.txt")
// 	assert.Equal(t, nil, err)
// 	defer file.Close()

// 	r := bufio.NewReader(file)
// 	line, e := readLn(r)
// 	for e == nil {
// 		s := strings.Split(line, ",")
// 		a1, err := strconv.ParseUint(s[0], 10, 64)
// 		assert.Equal(t, nil, err)
// 		a2, err := strconv.ParseUint(s[1], 10, 64)
// 		assert.Equal(t, nil, err)
// 		a3, err := strconv.ParseUint(s[2], 10, 64)
// 		assert.Equal(t, nil, err)
// 		a4, err := strconv.ParseUint(s[3], 10, 64)
// 		assert.Equal(t, nil, err)
// 		b1, err := strconv.ParseUint(s[4], 10, 64)
// 		assert.Equal(t, nil, err)
// 		b2, err := strconv.ParseUint(s[5], 10, 64)
// 		assert.Equal(t, nil, err)
// 		b3, err := strconv.ParseUint(s[6], 10, 64)
// 		assert.Equal(t, nil, err)
// 		b4, err := strconv.ParseUint(s[7], 10, 64)
// 		assert.Equal(t, nil, err)
// 		expectedAVec := s[8]
// 		expectedBVec := s[9]
// 		expectedRes := s[10]

// 		a := &Fq{a1, a2, a3, a4}
// 		aRev := reverse(a.intoBytes())
// 		assert.Equal(t, expectedAVec, hex.EncodeToString(aRev))

// 		b := &Fq{b1, b2, b3, b4}
// 		bRev := reverse(b.intoBytes())
// 		assert.Equal(t, expectedBVec, hex.EncodeToString(bRev))

// 		a.Sub(a, b) //a is now the result
// 		resRev := reverse(a.intoBytes())
// 		assert.Equal(t, expectedRes, hex.EncodeToString(resRev))

// 		line, e = readLn(r)
// 	}

// }
// func TestMul(t *testing.T) {

// 	// XXX: calculate expected dynamically

// 	file, err := os.Open("testdata/fqmul_test.txt")
// 	assert.Equal(t, nil, err)
// 	defer file.Close()

// 	r := bufio.NewReader(file)
// 	line, e := readLn(r)
// 	for e == nil {
// 		s := strings.Split(line, ",")
// 		a1, err := strconv.ParseUint(s[0], 10, 64)
// 		assert.Equal(t, nil, err)
// 		a2, err := strconv.ParseUint(s[1], 10, 64)
// 		assert.Equal(t, nil, err)
// 		a3, err := strconv.ParseUint(s[2], 10, 64)
// 		assert.Equal(t, nil, err)
// 		a4, err := strconv.ParseUint(s[3], 10, 64)
// 		assert.Equal(t, nil, err)
// 		b1, err := strconv.ParseUint(s[4], 10, 64)
// 		assert.Equal(t, nil, err)
// 		b2, err := strconv.ParseUint(s[5], 10, 64)
// 		assert.Equal(t, nil, err)
// 		b3, err := strconv.ParseUint(s[6], 10, 64)
// 		assert.Equal(t, nil, err)
// 		b4, err := strconv.ParseUint(s[7], 10, 64)
// 		assert.Equal(t, nil, err)
// 		expectedAVec := s[8]
// 		expectedBVec := s[9]
// 		expectedRes := s[10]

// 		a := &Fq{a1, a2, a3, a4}
// 		aRev := reverse(a.intoBytes())
// 		assert.Equal(t, expectedAVec, hex.EncodeToString(aRev))

// 		b := &Fq{b1, b2, b3, b4}
// 		bRev := reverse(b.intoBytes())
// 		assert.Equal(t, expectedBVec, hex.EncodeToString(bRev))

// 		a.Mul(b, a) // a is now the result
// 		resRev := reverse(a.intoBytes())
// 		assert.Equal(t, expectedRes, hex.EncodeToString(resRev))

// 		line, e = readLn(r)
// 	}

// }
// func TestSquare(t *testing.T) {

// 	file, err := os.Open("testdata/fqsq_test.txt")
// 	assert.Equal(t, nil, err)
// 	defer file.Close()

// 	r := bufio.NewReader(file)
// 	line, e := readLn(r)
// 	for e == nil {
// 		s := strings.Split(line, ",")
// 		a1, err := strconv.ParseUint(s[0], 10, 64)
// 		assert.Equal(t, nil, err)
// 		a2, err := strconv.ParseUint(s[1], 10, 64)
// 		assert.Equal(t, nil, err)
// 		a3, err := strconv.ParseUint(s[2], 10, 64)
// 		assert.Equal(t, nil, err)
// 		a4, err := strconv.ParseUint(s[3], 10, 64)
// 		assert.Equal(t, nil, err)
// 		expectedAVec := s[4]
// 		expectedRes := s[5]

// 		a := Fq{a1, a2, a3, a4}
// 		aRev := reverse(a.intoBytes())
// 		assert.Equal(t, expectedAVec, hex.EncodeToString(aRev))

// 		a.Square(&a) // a is now the result
// 		resRev := reverse(a.intoBytes())
// 		assert.Equal(t, expectedRes, hex.EncodeToString(resRev))

// 		line, e = readLn(r)
// 	}

// }

func TestSqrtVarTime(t *testing.T) {
	f := Fq{1, 2, 3, 4}
	f.SqrtVarTime()
	t.Fail()
}

// func TestNegation(t *testing.T) {
// 	var temp = LARGEST
// 	temp.Neg(&temp)
// 	assert.Equal(t, true, temp.Equal(&temp, &Fq{1, 0, 0, 0}))

// 	temp = zero
// 	temp.Neg(&temp)
// 	assert.Equal(t, true, temp.Equal(&temp, &zero))

// 	temp = Fq{1, 0, 0, 0}
// 	temp.Neg(&temp)
// 	assert.Equal(t, true, temp.Equal(&temp, &LARGEST))
// }

// func TestInverse(t *testing.T) {

// 	file, err := os.Open("testdata/fqinv_test.txt")
// 	assert.Equal(t, nil, err)
// 	defer file.Close()

// 	r := bufio.NewReader(file)
// 	line, e := readLn(r)
// 	for e == nil {
// 		s := strings.Split(line, ",")
// 		a1, err := strconv.ParseUint(s[0], 10, 64)
// 		assert.Equal(t, nil, err)
// 		a2, err := strconv.ParseUint(s[1], 10, 64)
// 		assert.Equal(t, nil, err)
// 		a3, err := strconv.ParseUint(s[2], 10, 64)
// 		assert.Equal(t, nil, err)
// 		a4, err := strconv.ParseUint(s[3], 10, 64)
// 		assert.Equal(t, nil, err)
// 		expectedAVec := s[4]
// 		expectedRes := s[5]

// 		a := Fq{a1, a2, a3, a4}
// 		aRev := reverse(a.intoBytes())
// 		assert.Equal(t, expectedAVec, hex.EncodeToString(aRev))

// 		a.Inverse(&a) // a is now the result
// 		resRev := reverse(a.intoBytes())
// 		assert.Equal(t, expectedRes, hex.EncodeToString(resRev))

// 		line, e = readLn(r)
// 	}

// }

// // Compute -(q^{-1} mod 2^64) mod 2^64 by exponentiating by totient(2**64) - 1
// func TestInv(t *testing.T) {

// 	var inv uint64 = 1

// 	for i := 0; i < 63; i++ {
// 		inv = inv * inv
// 		inv = inv * q[0]
// 	}

// 	inv = -inv
// 	assert.Equal(t, inv, INV)
// }

// func TestBytesInto(t *testing.T) {
// 	buf := [32]byte{}
// 	// buf2 := [32]byte{}

// 	f := Fq{0, 2, 3, 4}
// 	by := f.intoBytes()
// 	f.BytesInto(&buf)
// 	assert.Equal(t, by, buf[:])
// }

// helper methods
func reverse(s []byte) []byte {
	for i, j := 0, len(s)-1; i <= j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

func readLn(r *bufio.Reader) (string, error) {
	var (
		isPrefix bool  = true
		err      error = nil
		line, ln []byte
	)
	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}
	return strings.TrimSuffix(string(ln), "\n"), err
}
