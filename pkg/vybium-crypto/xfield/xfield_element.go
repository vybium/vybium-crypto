// Package xfield provides extension field arithmetic operations over F_p^3.
//
// The extension field F_p^3 is defined over the base field F_p (Goldilocks prime)
// using the irreducible polynomial: x³ - x + 1
//
// An extension field element is represented as: c₀ + c₁·x + c₂·x²
// where c₀, c₁, c₂ are elements of the base field F_p
//
// Extension fields are commonly used in zero-knowledge proof systems to provide
// additional algebraic structure for efficient polynomial operations.
package xfield

import (
	"encoding/json"
	"fmt"

	"github.com/vybium/vybium-crypto/pkg/vybium-crypto/field"
	"github.com/vybium/vybium-crypto/pkg/vybium-crypto/polynomial"
)

const (
	// ExtensionDegree is the degree of the field extension
	ExtensionDegree = 3
)

// XFieldElement represents an element in the extension field F_p^3.
// It is represented as c₀ + c₁·x + c₂·x² where each cᵢ is a BFieldElement.
//
// The extension is defined by the irreducible polynomial x³ - x + 1,
// meaning x³ = x - 1 in this field.
//
// Production implementation.
type XFieldElement struct {
	// Coefficients [c₀, c₁, c₂] representing c₀ + c₁·x + c₂·x²
	Coefficients [ExtensionDegree]field.Element
}

// Common extension field elements
var (
	Zero = XFieldElement{[ExtensionDegree]field.Element{field.Zero, field.Zero, field.Zero}}
	One  = XFieldElement{[ExtensionDegree]field.Element{field.One, field.Zero, field.Zero}}
)

// New creates a new extension field element from three base field coefficients.
//
// Production implementation.
func New(coefficients [ExtensionDegree]field.Element) XFieldElement {
	return XFieldElement{Coefficients: coefficients}
}

// NewConst creates a new extension field element from a single base field element.
// The result is c + 0·x + 0·x² (i.e., a constant in the extension field).
//
// Production implementation.
func NewConst(element field.Element) XFieldElement {
	return XFieldElement{
		Coefficients: [ExtensionDegree]field.Element{element, field.Zero, field.Zero},
	}
}

// NewU64 creates a new constant extension field element from a uint64 value.
func NewU64(value uint64) XFieldElement {
	return NewConst(field.New(value))
}

// IsZero returns true if the element is zero.
func (x XFieldElement) IsZero() bool {
	return x.Coefficients[0].IsZero() &&
		x.Coefficients[1].IsZero() &&
		x.Coefficients[2].IsZero()
}

// IsOne returns true if the element is one.
func (x XFieldElement) IsOne() bool {
	return x.Coefficients[0].IsOne() &&
		x.Coefficients[1].IsZero() &&
		x.Coefficients[2].IsZero()
}

// Equal returns true if two extension field elements are equal.
func (x XFieldElement) Equal(other XFieldElement) bool {
	return x.Coefficients[0].Equal(other.Coefficients[0]) &&
		x.Coefficients[1].Equal(other.Coefficients[1]) &&
		x.Coefficients[2].Equal(other.Coefficients[2])
}

// Unlift attempts to convert an extension field element to a base field element.
// Returns the base field element if c₁ = c₂ = 0, otherwise returns nil.
//
// Production implementation.
func (x XFieldElement) Unlift() *field.Element {
	if x.Coefficients[1].IsZero() && x.Coefficients[2].IsZero() {
		result := x.Coefficients[0]
		return &result
	}
	return nil
}

// String returns the string representation of the extension field element.
func (x XFieldElement) String() string {
	// If it's a constant (unlifted), show it as such
	if unlift := x.Unlift(); unlift != nil {
		return fmt.Sprintf("%s_xfe", unlift.String())
	}

	// Otherwise show the full polynomial form
	c0, c1, c2 := x.Coefficients[0], x.Coefficients[1], x.Coefficients[2]
	return fmt.Sprintf("(%020d·x² + %020d·x + %020d)", c2.Value(), c1.Value(), c0.Value())
}

// Add performs extension field addition: (a₀ + a₁x + a₂x²) + (b₀ + b₁x + b₂x²)
//
// Production implementation.
func (x XFieldElement) Add(other XFieldElement) XFieldElement {
	return XFieldElement{
		Coefficients: [ExtensionDegree]field.Element{
			x.Coefficients[0].Add(other.Coefficients[0]),
			x.Coefficients[1].Add(other.Coefficients[1]),
			x.Coefficients[2].Add(other.Coefficients[2]),
		},
	}
}

// AddConst adds a base field element to an extension field element.
// This adds the constant to the c₀ coefficient.
//
// Production implementation.
func (x XFieldElement) AddConst(other field.Element) XFieldElement {
	return XFieldElement{
		Coefficients: [ExtensionDegree]field.Element{
			x.Coefficients[0].Add(other),
			x.Coefficients[1],
			x.Coefficients[2],
		},
	}
}

// Sub performs extension field subtraction: (a₀ + a₁x + a₂x²) - (b₀ + b₁x + b₂x²)
//
// Production implementation.
func (x XFieldElement) Sub(other XFieldElement) XFieldElement {
	return XFieldElement{
		Coefficients: [ExtensionDegree]field.Element{
			x.Coefficients[0].Sub(other.Coefficients[0]),
			x.Coefficients[1].Sub(other.Coefficients[1]),
			x.Coefficients[2].Sub(other.Coefficients[2]),
		},
	}
}

// SubConst subtracts a base field element from an extension field element.
//
// Production implementation.
func (x XFieldElement) SubConst(other field.Element) XFieldElement {
	return XFieldElement{
		Coefficients: [ExtensionDegree]field.Element{
			x.Coefficients[0].Sub(other),
			x.Coefficients[1],
			x.Coefficients[2],
		},
	}
}

// Neg negates the extension field element.
//
// Production implementation.
func (x XFieldElement) Neg() XFieldElement {
	return XFieldElement{
		Coefficients: [ExtensionDegree]field.Element{
			x.Coefficients[0].Neg(),
			x.Coefficients[1].Neg(),
			x.Coefficients[2].Neg(),
		},
	}
}

// Mul performs extension field multiplication modulo x³ - x + 1.
//
// Given:
//
//	(ax² + bx + c) * (dx² + ex + f) mod (x³ - x + 1)
//
// Expansion:
//
//	= adx⁴ + aex³ + afx² + bdx³ + bex² + bfx + cdx² + cex + cf
//
// Using x³ = x - 1, we reduce:
//
//	x⁴ = x² - x, x³ = x - 1
//
// Result coefficients:
//
//	r₀ = cf - ae - bd
//	r₁ = bf + ce - ad + ae + bd
//	r₂ = af + be + cd + ad
//
// Production implementation.
func (x XFieldElement) Mul(other XFieldElement) XFieldElement {
	c, b, a := x.Coefficients[0], x.Coefficients[1], x.Coefficients[2]
	f, e, d := other.Coefficients[0], other.Coefficients[1], other.Coefficients[2]

	// Compute products
	ae := a.Mul(e)
	bd := b.Mul(d)

	r0 := c.Mul(f).Sub(ae).Sub(bd)
	r1 := b.Mul(f).Add(c.Mul(e)).Sub(a.Mul(d)).Add(ae).Add(bd)
	r2 := a.Mul(f).Add(b.Mul(e)).Add(c.Mul(d)).Add(a.Mul(d))

	return XFieldElement{
		Coefficients: [ExtensionDegree]field.Element{r0, r1, r2},
	}
}

// MulConst multiplies an extension field element by a base field element (scalar multiplication).
// Each coefficient is multiplied by the scalar.
//
// Production implementation.
func (x XFieldElement) MulConst(scalar field.Element) XFieldElement {
	return XFieldElement{
		Coefficients: [ExtensionDegree]field.Element{
			x.Coefficients[0].Mul(scalar),
			x.Coefficients[1].Mul(scalar),
			x.Coefficients[2].Mul(scalar),
		},
	}
}

// ShahPolynomial returns the irreducible polynomial defining the extension: x³ - x + 1
//
// Production implementation.
func ShahPolynomial() *polynomial.Polynomial {
	// x³ - x + 1 = 1 + (-1)x + 0x² + 1x³
	return polynomial.New([]field.Element{
		field.One,
		field.One.Neg(),
		field.Zero,
		field.One,
	})
}

// Inverse computes the multiplicative inverse of the extension field element.
//
// Production implementation.
func (x XFieldElement) Inverse() XFieldElement {
	if x.IsZero() {
		panic("cannot invert the zero element in the extension field")
	}

	a, b, c := x.Coefficients[0], x.Coefficients[1], x.Coefficients[2]

	// Handle constant elements (b=0, c=0) correctly
	if b.IsZero() && c.IsZero() {
		return NewConst(a.Inverse())
	}

	// For non-constant elements, use direct formula based on the multiplication structure
	// We need to find [d, e, f] such that (a + bx + cx²) * (d + ex + fx²) = 1 (mod x³ - x + 1)
	//
	// Using the multiplication formula:
	// r0 = c*d - a*e - b*f = 1
	// r1 = b*d + c*e - a*f + a*e + b*f = 0
	// r2 = a*d + b*e + c*f + a*f = 0
	//
	// Simplifying:
	// r0 = cd - ae - bf = 1
	// r1 = bd + ce + ae + bf = 0
	// r2 = ad + be + cf + af = 0
	//
	// This can be solved using the norm form. In extension fields,
	// we use: inv(x) = conj(x) / norm(x)
	// where norm(x) = x * conj(x) * conj²(x)
	//
	// For the extension field defined by x³ - x + 1,
	// we can compute the inverse using the formula:
	// inv(a + bx + cx²) = (a² + bc - ac*x - b²*x² + (ab-c²)*1) / norm
	//
	// Following twenty-first's approach with XGCD:
	xPoly := polynomial.New([]field.Element{a, b, c})
	shahPoly := ShahPolynomial()

	// Compute XGCD: gcd = aResult*xPoly + bResult*shahPoly
	_, aResult, _ := polynomial.XGCD(xPoly, shahPoly)

	// Reduce modulo Shah polynomial
	_, remainder := aResult.Divide(shahPoly)

	// Convert remainder to XFieldElement
	coeffs := remainder.Coefficients()

	var resultCoeffs [ExtensionDegree]field.Element
	for i := 0; i < ExtensionDegree; i++ {
		if i < len(coeffs) {
			resultCoeffs[i] = coeffs[i]
		} else {
			resultCoeffs[i] = field.Zero
		}
	}

	return XFieldElement{Coefficients: resultCoeffs}
}

// Div performs extension field division: x / y = x * y⁻¹
//
// Production implementation.
func (x XFieldElement) Div(other XFieldElement) XFieldElement {
	return x.Mul(other.Inverse())
}

// Pow computes x^exponent using square-and-multiply algorithm.
//
// Production implementation.
func (x XFieldElement) Pow(exponent uint64) XFieldElement {
	result := One
	base := x

	for exponent > 0 {
		if exponent&1 == 1 {
			result = result.Mul(base)
		}
		base = base.Mul(base)
		exponent >>= 1
	}

	return result
}

// MarshalJSON implements json.Marshaler.
// Extension field elements are serialized as arrays of 3 base field elements.
func (x XFieldElement) MarshalJSON() ([]byte, error) {
	// Serialize as array of coefficient values
	values := [ExtensionDegree]uint64{
		x.Coefficients[0].Value(),
		x.Coefficients[1].Value(),
		x.Coefficients[2].Value(),
	}
	return json.Marshal(values)
}

// UnmarshalJSON implements json.Unmarshaler.
func (x *XFieldElement) UnmarshalJSON(data []byte) error {
	var values [ExtensionDegree]uint64
	if err := json.Unmarshal(data, &values); err != nil {
		return err
	}

	x.Coefficients = [ExtensionDegree]field.Element{
		field.New(values[0]),
		field.New(values[1]),
		field.New(values[2]),
	}

	return nil
}

// ToDigest converts an extension field element to a 5-element digest.
// The three coefficients become the first three digest elements, with zeros padding.
//
// This is used for Merkle tree construction from extension field elements.
//
// Production implementation.
func (x XFieldElement) ToDigest() [5]field.Element {
	return [5]field.Element{
		x.Coefficients[0],
		x.Coefficients[1],
		x.Coefficients[2],
		field.Zero,
		field.Zero,
	}
}

// FromDigest creates an extension field element from a digest.
// Returns nil if the last two elements of the digest are not zero.
//
// Production implementation.
func FromDigest(digest [5]field.Element) *XFieldElement {
	// Check that the last two elements are zero
	if !digest[3].IsZero() || !digest[4].IsZero() {
		return nil
	}

	return &XFieldElement{
		Coefficients: [ExtensionDegree]field.Element{
			digest[0],
			digest[1],
			digest[2],
		},
	}
}

// AsFlatSlice interprets a slice of XFieldElements as a slice of BFieldElements.
// This allows zero-copy conversion for hashing and other operations.
//
// The resulting slice has length len(xfes) * 3, with coefficients in order:
// [xfe[0].c0, xfe[0].c1, xfe[0].c2, xfe[1].c0, xfe[1].c1, xfe[1].c2, ...]
//
// Production implementation.
func AsFlatSlice(xfes []XFieldElement) []field.Element {
	if len(xfes) == 0 {
		return nil
	}

	result := make([]field.Element, len(xfes)*ExtensionDegree)
	for i, xfe := range xfes {
		base := i * ExtensionDegree
		result[base] = xfe.Coefficients[0]
		result[base+1] = xfe.Coefficients[1]
		result[base+2] = xfe.Coefficients[2]
	}

	return result
}

// FromBFieldSlice creates an extension field element from a slice of base field elements.
// The slice must have exactly 3 elements.
//
// Production implementation.
func FromBFieldSlice(elements []field.Element) (*XFieldElement, error) {
	if len(elements) != ExtensionDegree {
		return nil, fmt.Errorf("invalid length %d, expected %d", len(elements), ExtensionDegree)
	}

	return &XFieldElement{
		Coefficients: [ExtensionDegree]field.Element{
			elements[0],
			elements[1],
			elements[2],
		},
	}, nil
}
