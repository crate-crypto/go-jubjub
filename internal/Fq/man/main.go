package main

import (
	"fmt"
	"encoding/hex"
	jubjub "github.com/decentralisedkev/go-jubjub"
	
)

func main() {

	// p := jubjub.Point{}

	// s := jubjub.Scalar{}
	// s.Rand()
	// fmt.Println(s.String())

	// newPoint := p.ScalarMultBase(s)
	// fmt.Println(newPoint)

	var a, b jubjub.Scalar
	a.Rand() // random element

	byt, _ := hex.DecodeString(a.String())
	
	var buf [32]byte
	copy(buf[:], byt)

	b.SetBytes(&buf)

	aRes, _ := hex.DecodeString(a.String())
	bRes, _ := hex.DecodeString(b.String())

	fmt.Println(aRes)
	fmt.Println(bRes)

}
