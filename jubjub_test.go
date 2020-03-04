package jubjub

import (
	"fmt"
	"testing"
)

func TestHashToPoint(t *testing.T) {
	var p Point

	hashed := p.HashToPoint([]byte("hello world"))

	fmt.Println(hashed)
}
