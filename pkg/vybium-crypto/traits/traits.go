package traits

import (
	"fmt"
	"math/big"
)

// FiniteField represents a finite field element with all required operations.
// This is the core trait that all field elements must implement.
type FiniteField interface {
	// Basic arithmetic operations
	Add(other FiniteField) FiniteField
	Sub(other FiniteField) FiniteField
	Mul(other FiniteField) FiniteField
	Div(other FiniteField) FiniteField
	Neg() FiniteField

	// Comparison operations
	Equal(other FiniteField) bool
	IsZero() bool
	IsOne() bool

	// Field-specific operations
	Inverse() FiniteField
	Square() FiniteField
	Pow(exp uint64) FiniteField

	// Conversion operations
	ToBigInt() *big.Int
	FromBigInt(val *big.Int) FiniteField
	ToUint64() uint64
	FromUint64(val uint64) FiniteField

	// String representation
	String() string

	// Serialization
	MarshalJSON() ([]byte, error)
	UnmarshalJSON(data []byte) error
}

// Inverse represents elements that can compute multiplicative inverses.
type Inverse interface {
	// Inverse returns the multiplicative inverse: a * a.Inverse() == 1
	// Panics if the element is zero (has no inverse)
	Inverse() FiniteField

	// InverseOrZero returns the inverse or zero if the element is zero
	InverseOrZero() FiniteField
}

// PrimitiveRootOfUnity represents elements that can find primitive roots of unity.
type PrimitiveRootOfUnity interface {
	// PrimitiveRootOfUnity returns a primitive n-th root of unity, if it exists
	PrimitiveRootOfUnity(n uint64) (FiniteField, bool)
}

// CyclicGroupGenerator represents elements that can generate cyclic group elements.
type CyclicGroupGenerator interface {
	// GetCyclicGroupElements returns elements of the cyclic group up to max elements
	GetCyclicGroupElements(max *uint64) []FiniteField
}

// ModPowU32 represents elements that can compute modular exponentiation with u32 exponent.
type ModPowU32 interface {
	// ModPowU32 computes base^exp mod field_order
	ModPowU32(exp uint32) FiniteField
}

// ModPowU64 represents elements that can compute modular exponentiation with u64 exponent.
type ModPowU64 interface {
	// ModPowU64 computes base^exp mod field_order
	ModPowU64(exp uint64) FiniteField
}

// BatchInversion performs Montgomery batch inversion for efficiency.
// This is a key optimization for field arithmetic.
func BatchInversion(elements []FiniteField) ([]FiniteField, error) {
	if len(elements) == 0 {
		return []FiniteField{}, nil
	}

	// Check for zeros
	for i, elem := range elements {
		if elem.IsZero() {
			return nil, fmt.Errorf("cannot perform batch inversion on zero element at index %d", i)
		}
	}

	n := len(elements)
	one := elements[0].FromUint64(1)

	// Scratch space for intermediate products
	scratch := make([]FiniteField, n)
	acc := one
	scratch[0] = elements[0]

	// Forward pass: compute products
	for i := 0; i < n; i++ {
		scratch[i] = acc
		acc = acc.Mul(elements[i])
	}

	// Get the inverse of the total product
	acc = acc.Inverse()

	// Backward pass: compute individual inverses
	result := make([]FiniteField, n)
	copy(result, elements)

	for i := n - 1; i >= 0; i-- {
		tmp := acc.Mul(result[i])
		result[i] = acc.Mul(scratch[i])
		acc = tmp
	}

	return result, nil
}

// Square computes the square of a field element efficiently.
func Square(elem FiniteField) FiniteField {
	return elem.Mul(elem)
}

// Pow computes element^exp efficiently using binary exponentiation.
func Pow(elem FiniteField, exp uint64) FiniteField {
	if exp == 0 {
		// Return 1 (multiplicative identity)
		one := elem.FromUint64(1)
		return one
	}

	if exp == 1 {
		return elem
	}

	// Binary exponentiation
	result := elem.FromUint64(1)
	base := elem

	for exp > 0 {
		if exp&1 == 1 {
			result = result.Mul(base)
		}
		base = base.Mul(base)
		exp >>= 1
	}

	return result
}

// ValidateFiniteField ensures that a FiniteField implementation is correct.
// This is used for testing and validation.
func ValidateFiniteField(elem FiniteField) error {
	// Test basic properties
	zero := elem.FromUint64(0)
	one := elem.FromUint64(1)

	// Test zero properties
	if !elem.IsZero() && elem.Equal(zero) {
		return fmt.Errorf("IsZero() and Equal(zero) are inconsistent")
	}

	// Test one properties
	if !elem.IsOne() && elem.Equal(one) {
		return fmt.Errorf("IsOne() and Equal(one) are inconsistent")
	}

	// Test additive identity: a + 0 = a
	if !elem.Add(zero).Equal(elem) {
		return fmt.Errorf("additive identity failed: a + 0 != a")
	}

	// Test multiplicative identity: a * 1 = a
	if !elem.Mul(one).Equal(elem) {
		return fmt.Errorf("multiplicative identity failed: a * 1 != a")
	}

	// Test additive inverse: a + (-a) = 0
	neg := elem.Neg()
	if !elem.Add(neg).Equal(zero) {
		return fmt.Errorf("additive inverse failed: a + (-a) != 0")
	}

	// Test multiplicative inverse (if not zero)
	if !elem.IsZero() {
		inv := elem.Inverse()
		if !elem.Mul(inv).Equal(one) {
			return fmt.Errorf("multiplicative inverse failed: a * a^-1 != 1")
		}
	}

	return nil
}

// FieldElementType represents the type of field element for type checking.
type FieldElementType int

const (
	BFieldElementType FieldElementType = iota
	XFieldElementType
)

// GetFieldElementType returns the type of a field element.
func GetFieldElementType(elem FiniteField) FieldElementType {
	// This is a simple type check - in practice, you'd use type assertions
	// or a more sophisticated type system
	switch elem.(type) {
	case interface{ BFieldElement() bool }:
		return BFieldElementType
	case interface{ XFieldElement() bool }:
		return XFieldElementType
	default:
		// Try to determine by string representation or other means
		return BFieldElementType // Default fallback
	}
}
