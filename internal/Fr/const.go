package fr

var zero = Fr{0, 0, 0, 0}

// INV = -(r^{-1} mod 2^64) mod 2^64
const INV uint64 = 0x1ba3a358ef788ef9

// R = 2^256 mod r
// the montgomoery modulus
var R = Fr{0x25f80bb3b99607d9, 0xf315d62f66b6e750, 0x932514eeeb8814f4, 0x09a6fc6f479155c6}

// R2 (Rsquared) = 2^512 mod r
var R2 = Fr{0x67719aa495e57731, 0x51b0cef09ce3fc26, 0x69dab7fac026e9a5, 0x04f6547b8d127688}

// r is modulus in Fr
// r = 0x0e7db4ea6533afa906673b0101343b00a6682093ccc81082d0970e5ed6f72cb7
var r = Fr{0xd0970e5ed6f72cb7, 0xa6682093ccc81082, 0x06673b0101343b00, 0x0e7db4ea6533afa9}

// NEG1 = -R = -(2^256 mod r) mod r
var NEG1 = Fr{0xaa9f02ab1d6124de, 0xb3524a6466112932, 0x7342261215ac260b, 0x4d6b87b1da259e2}
