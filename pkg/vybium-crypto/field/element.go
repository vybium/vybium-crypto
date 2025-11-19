// Package field provides finite field arithmetic operations using Montgomery representation.
//
// This package implements efficient arithmetic over the Goldilocks prime field (p = 2^64 - 2^32 + 1).
// Values are stored in Montgomery form (x * 2^64 mod P) to enable efficient modular multiplication
// without expensive division operations. This representation is commonly used in zero-knowledge
// proof systems for its performance characteristics.
package field

import (
	"encoding/binary"
	"fmt"
	"math/big"
	"math/bits"
)

// P is the prime modulus: 2^64 - 2^32 + 1
const P uint64 = 0xFFFFFFFF00000001

// R2 is 2^128 mod P, used for conversion into Montgomery representation
const R2 uint64 = 0xFFFFFFFE00000001

// MINUS_TWO_INVERSE is the constant -2^(-1) mod P
const MINUS_TWO_INVERSE uint64 = 0x7FFFFFFF80000000

// Element represents a field element in the base field F_p where p = 2^64 - 2^32 + 1.
// The value is stored in Montgomery representation for efficient arithmetic operations.
// All arithmetic operations (Add, Mul, Sub, etc.) work directly on Montgomery-form values.
type Element struct {
	// value in Montgomery form (value * 2^64 mod P)
	value uint64
}

// Constants for common field elements
var (
	Zero = Element{0}
	One  = New(1)
	Max  = New(P - 1)
)

// New creates a new field element from a uint64 value.
// The value is converted to Montgomery form.
//
// Production implementation.
func New(value uint64) Element {
	// Convert to Montgomery form: montyred(value * R2)
	product := mul128(value, R2)
	return Element{value: montyred(product)}
}

// NewFromRaw creates an element directly from Montgomery form.
// This is used for internal operations and deserialization.
//
// Production implementation.
func NewFromRaw(raw uint64) Element {
	return Element{value: raw}
}

// NewFromInt64 creates a new field element from an int64 value.
// Negative values are handled correctly.
func NewFromInt64(value int64) Element {
	if value < 0 {
		// Handle negative values: -x mod P = P - x
		absValue := uint64(-value) % P
		if absValue == 0 {
			return Zero
		}
		return New(P - absValue)
	}
	return New(uint64(value) % P)
}

// NewFromBigInt creates a new field element from a big.Int value.
func NewFromBigInt(value *big.Int) Element {
	// Reduce modulo P
	mod := new(big.Int).SetUint64(P)
	reduced := new(big.Int).Mod(value, mod)

	// Handle negative results
	if reduced.Sign() < 0 {
		reduced.Add(reduced, mod)
	}

	return New(reduced.Uint64())
}

// Value returns the canonical uint64 value of the field element.
// This converts from Montgomery form back to normal form.
//
// Production implementation.
func (e Element) Value() uint64 {
	// Convert from Montgomery form: montyred(value)
	return montyred(uint128{lo: e.value, hi: 0})
}

// RawValue returns the raw Montgomery form value.
// This is useful for serialization and internal operations.
//
// Production implementation.
func (e Element) RawValue() uint64 {
	return e.value
}

// String returns the string representation of the field element.
// Returns the canonical (non-Montgomery) value for readability.
func (e Element) String() string {
	return fmt.Sprintf("%d", e.Value())
}

// Hex returns the lowercase hexadecimal representation of the canonical value.
// Production implementation.
func (e Element) Hex() string {
	return fmt.Sprintf("%x", e.Value())
}

// HexUpper returns the uppercase hexadecimal representation of the canonical value.
// Production implementation.
func (e Element) HexUpper() string {
	return fmt.Sprintf("%X", e.Value())
}

// IsZero returns true if the element is zero.
func (e Element) IsZero() bool {
	return e.value == 0
}

// IsOne returns true if the element is one.
func (e Element) IsOne() bool {
	return e.Equal(One)
}

// Add performs field addition: (a + b) mod P
// Uses the optimized addition from twenty-first.
//
// Production implementation.
func (e Element) Add(other Element) Element {
	// Compute a + b = a - (p - b)
	// This clever trick avoids conditional logic in many cases
	x1, c1 := bits.Sub64(e.value, P-other.value, 0)

	// If there was a borrow, add P back
	if c1 != 0 {
		return Element{value: x1 + P}
	}
	return Element{value: x1}
}

// Sub performs field subtraction: (a - b) mod P
// Uses the optimized subtraction from twenty-first.
//
// Production implementation.
func (e Element) Sub(other Element) Element {
	// Perform subtraction with borrow detection
	x1, c1 := bits.Sub64(e.value, other.value, 0)

	// Adjust for borrow: x1 - ((1 + !P) * borrow)
	// This is equivalent to: if c1 { x1 + P } else { x1 }
	return Element{value: x1 - ((1 + ^P) * c1)}
}

// Mul performs field multiplication: (a * b) mod P
// Uses Montgomery multiplication for efficiency.
//
// Production implementation.
func (e Element) Mul(other Element) Element {
	// Montgomery multiplication: montyred(a * b)
	product := mul128(e.value, other.value)
	return Element{value: montyred(product)}
}

// Div performs field division: (a / b) mod P
func (e Element) Div(other Element) Element {
	// Division is multiplication by inverse: a / b = a * b^(-1)
	return e.Mul(other.Inverse())
}

// Square computes e^2 mod P
func (e Element) Square() Element {
	return e.Mul(e)
}

// Inverse computes the multiplicative inverse: a^(-1) mod P
// Uses the optimized inversion chain from twenty-first.
//
// Production implementation.
func (e Element) Inverse() Element {
	if e.IsZero() {
		panic("attempted to find the multiplicative inverse of zero")
	}

	// Helper function for repeated squaring
	exp := func(base Element, exponent uint64) Element {
		result := base
		for i := uint64(0); i < exponent; i++ {
			result = result.Square()
		}
		return result
	}

	// Optimized inversion chain from twenty-first
	// This computes a^(P-2) mod P using an addition chain
	x := e
	bin2Ones := x.Square().Mul(x)                  // a^3
	bin3Ones := bin2Ones.Square().Mul(x)           // a^7
	bin6Ones := exp(bin3Ones, 3).Mul(bin3Ones)     // a^63
	bin12Ones := exp(bin6Ones, 6).Mul(bin6Ones)    // a^(2^12 - 1)
	bin24Ones := exp(bin12Ones, 12).Mul(bin12Ones) // a^(2^24 - 1)
	bin30Ones := exp(bin24Ones, 6).Mul(bin6Ones)   // a^(2^30 - 1)
	bin31Ones := bin30Ones.Square().Mul(x)         // a^(2^31 - 1)
	bin31Ones1Zero := bin31Ones.Square()           // a^(2^32 - 2)
	bin32Ones := bin31Ones.Square().Mul(x)         // a^(2^32 - 1)

	return exp(bin31Ones1Zero, 32).Mul(bin32Ones)
}

// ModPow computes modular exponentiation: a^exp mod P
// Uses binary exponentiation in Montgomery form.
//
// Production implementation.
func (e Element) ModPow(exp uint64) Element {
	if exp == 0 {
		return One
	}

	// Binary exponentiation
	acc := One
	bitLength := bits.Len64(exp)

	for i := 0; i < bitLength; i++ {
		acc = acc.Square()
		// Check bit from most significant to least
		if exp&(1<<(bitLength-1-i)) != 0 {
			acc = acc.Mul(e)
		}
	}

	return acc
}

// Neg returns the additive inverse: -a mod P
func (e Element) Neg() Element {
	if e.IsZero() {
		return Zero
	}
	return Element{value: P - e.value}
}

// Equal returns true if two elements are equal.
func (e Element) Equal(other Element) bool {
	return e.value == other.value
}

// Less returns true if this element's canonical representation is less than the other.
func (e Element) Less(other Element) bool {
	return e.Value() < other.Value()
}

// Greater returns true if this element's canonical representation is greater than the other.
func (e Element) Greater(other Element) bool {
	return e.Value() > other.Value()
}

// ToBigInt converts the field element to a big.Int.
func (e Element) ToBigInt() *big.Int {
	return new(big.Int).SetUint64(e.Value())
}

// ToBytes returns the little-endian byte representation.
func (e Element) ToBytes() [8]byte {
	var bytes [8]byte
	binary.LittleEndian.PutUint64(bytes[:], e.value)
	return bytes
}

// FromBytes creates an element from little-endian bytes.
func FromBytes(bytes [8]byte) Element {
	raw := binary.LittleEndian.Uint64(bytes[:])
	return NewFromRaw(raw)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (e Element) MarshalBinary() ([]byte, error) {
	bytes := e.ToBytes()
	return bytes[:], nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (e *Element) UnmarshalBinary(data []byte) error {
	if len(data) != 8 {
		return fmt.Errorf("invalid data length: expected 8 bytes, got %d", len(data))
	}

	var bytes [8]byte
	copy(bytes[:], data)
	*e = FromBytes(bytes)
	return nil
}

// Generator returns a generator for the entire field.
// The generator for this field is 7.
func Generator() Element {
	return New(7)
}

// uint128 represents a 128-bit unsigned integer.
type uint128 struct {
	lo, hi uint64
}

// mul128 performs 64-bit × 64-bit → 128-bit multiplication.
func mul128(a, b uint64) uint128 {
	hi, lo := bits.Mul64(a, b)
	return uint128{lo: lo, hi: hi}
}

// montyred performs Montgomery reduction: reduces a 128-bit value modulo P.
// This is the core operation for efficient modular arithmetic.
//
// This implements Montgomery reduction for efficient modular arithmetic.
func montyred(x uint128) uint64 {
	xl := x.lo
	xh := x.hi

	// a = xl + (xl << 32), with overflow detection
	a, e := bits.Add64(xl, xl<<32, 0)

	// b = a - (a >> 32) - overflow_flag
	b := a - (a >> 32) - e

	// r = xh - b, with borrow detection
	r, c := bits.Sub64(xh, b, 0)

	// Final adjustment
	// This is equivalent to: if c { r + P } else { r }
	return r - ((1 + ^P) * c)
}
