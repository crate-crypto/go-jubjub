package futil

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdc(t *testing.T) {

	file, err := os.Open("testdata/adc_test.txt")
	assert.Equal(t, nil, err)
	defer file.Close()

	r := bufio.NewReader(file)
	line, e := readLn(r)
	for e == nil {
		s := strings.Split(line, ",")
		a, err := strconv.ParseUint(s[0], 10, 64)
		assert.Equal(t, nil, err)
		b, err := strconv.ParseUint(s[1], 10, 64)
		assert.Equal(t, nil, err)
		carry, err := strconv.ParseUint(s[2], 10, 64)
		assert.Equal(t, nil, err)
		expectedres, err := strconv.ParseUint(s[3], 10, 64)
		assert.Equal(t, nil, err)
		expectedCarry, err := strconv.ParseUint(s[4], 10, 64)
		assert.Equal(t, nil, err)

		haveRes, haveCarry := Adc(a, b, carry)
		assert.Equal(t, expectedres, haveRes)
		assert.Equal(t, expectedCarry, haveCarry)
		line, e = readLn(r)
	}
}

func TestAdc2(t *testing.T) {

}

func TestSbb(t *testing.T) {
	file, err := os.Open("testdata/sbb_test.txt")
	assert.Equal(t, nil, err)
	defer file.Close()

	r := bufio.NewReader(file)
	line, e := readLn(r)
	for e == nil {
		s := strings.Split(line, ",")
		a, err := strconv.ParseUint(s[0], 10, 64)
		assert.Equal(t, nil, err)
		b, err := strconv.ParseUint(s[1], 10, 64)
		assert.Equal(t, nil, err)
		borrow, err := strconv.ParseUint(s[2], 10, 64)
		assert.Equal(t, nil, err)
		expectedres, err := strconv.ParseUint(s[3], 10, 64)
		assert.Equal(t, nil, err)
		expectedBorrow, err := strconv.ParseUint(s[4], 10, 64)
		assert.Equal(t, nil, err)

		haveRes, haveBorrow := Sbb(a, b, borrow)
		assert.Equal(t, expectedres, haveRes)
		assert.Equal(t, expectedBorrow, haveBorrow)
		line, e = readLn(r)
	}
}

func TestMac(t *testing.T) {
	file, err := os.Open("testdata/mac_test.txt")
	assert.Equal(t, nil, err)
	defer file.Close()

	r := bufio.NewReader(file)
	line, e := readLn(r)
	for e == nil {
		s := strings.Split(line, ",")
		a, err := strconv.ParseUint(s[0], 10, 64)
		assert.Equal(t, nil, err)
		b, err := strconv.ParseUint(s[1], 10, 64)
		assert.Equal(t, nil, err)
		c, err := strconv.ParseUint(s[2], 10, 64)
		assert.Equal(t, nil, err)
		carry, err := strconv.ParseUint(s[3], 10, 64)
		assert.Equal(t, nil, err)
		expectedres, err := strconv.ParseUint(s[4], 10, 64)
		assert.Equal(t, nil, err)
		expectedCarry, err := strconv.ParseUint(s[5], 10, 64)
		assert.Equal(t, nil, err)

		haveRes, haveCarry := Mac(a, b, c, carry)
		assert.Equal(t, expectedres, haveRes)
		assert.Equal(t, expectedCarry, haveCarry)
		line, e = readLn(r)
	}
}

// helper method
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
