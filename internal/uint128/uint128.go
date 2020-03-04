package uint128

type Uint128 struct{ H, L uint64 }

// FromU64 converts a uint64 into a uint128
func FromU64(a uint64) *Uint128 {
	return &Uint128{0, a}
}

// ToU64 converts a uint128 into a uint64
// by keeping the lower bits
func (u *Uint128) ToU64() uint64 {
	return u.L
}

// Add adds uint64 to u
func (u *Uint128) Add(n uint64) *Uint128 {
	lo := u.L + n
	hi := u.H
	if u.L > lo {
		hi++
	}
	u.L = lo
	u.H = hi
	return u
}

func (u *Uint128) AddU128(o *Uint128) *Uint128 {
	carry := u.L
	u.L += o.L
	u.H += o.H

	if u.L < carry {
		u.H++
	}
	return u
}

func (u *Uint128) SubU128(o *Uint128) *Uint128 {
	borrow := u.L
	u.L -= o.L
	u.H -= o.H

	if u.L > borrow {
		u.H--
	}
	return u
}

// Sub subs a uint64 number from u
func (u *Uint128) Sub(n uint64) *Uint128 {
	lo := u.L - n
	hi := u.H
	if u.L < lo {
		hi--
	}
	u.L = lo
	u.H = hi
	return u
}

// (Adapted from go's math/big)
// z1<<64 + z0 = x*y
// Adapted from Warren, Hacker's Delight, p. 132.
// DavidMinor
func mult(x, y uint64) (z1, z0 uint64) {
	z0 = x * y // lower 64 bits are easy
	// break the multiplication into (x1 << 32 + x0)(y1 << 32 + y0)
	// which is x1*y1 << 64 + (x0*y1 + x1*y0) << 32 + x0*y0
	// so now we can do 64 bit multiplication and addition and
	// shift the results into the right place
	x0, x1 := x&0x00000000ffffffff, x>>32
	y0, y1 := y&0x00000000ffffffff, y>>32
	w0 := x0 * y0
	t := x1*y0 + w0>>32
	w1 := t & 0x00000000ffffffff
	w2 := t >> 32
	w1 += x0 * y1
	z1 = x1*y1 + w2 + w1>>32
	return
}

// Mul multiplies two Uint128 numbers
func (u *Uint128) Mul(a *Uint128) *Uint128 {
	hl := u.H*a.L + u.L*a.H
	u.H, u.L = mult(u.L, a.L)
	u.H += hl
	return u
}

// MulU64 multiplies a uint64 with a uint128
func (u *Uint128) MulU64(n uint64) *Uint128 {
	hl := u.H * n
	u.H, u.L = mult(u.L, n)
	u.H += hl
	return u
}
