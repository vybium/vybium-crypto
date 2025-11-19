// Package polynomial provides univariate polynomial operations over finite fields.
//
// This package implements efficient polynomial arithmetic including addition, multiplication,
// evaluation, interpolation, and division. Polynomials are represented as coefficient vectors
// with coefficients stored in order of increasing degree. The implementation supports operations
// commonly needed in zero-knowledge proof systems such as vanishing polynomials and quotient computation.
package polynomial

import (
	"fmt"

	"github.com/vybium/vybium-crypto/pkg/vybium-crypto/field"
)

// Polynomial represents a univariate polynomial with coefficients in F_p.
// Coefficients are stored in order of increasing degree (coefficients[0] is the constant term).
// The zero polynomial is represented as an empty coefficient slice.
type Polynomial struct {
	// coefficients in order of increasing degree
	coefficients []field.Element
}

// New creates a new polynomial from coefficients.
// Coefficients are in order of increasing degree: [c0, c1, c2, ...] represents c0 + c1*x + c2*x^2 + ...
func New(coefficients []field.Element) *Polynomial {
	p := &Polynomial{
		coefficients: make([]field.Element, len(coefficients)),
	}
	copy(p.coefficients, coefficients)
	p.normalize()
	return p
}

// Zero returns the zero polynomial.
func Zero() *Polynomial {
	return &Polynomial{coefficients: []field.Element{}}
}

// One returns the constant polynomial 1.
func One() *Polynomial {
	return &Polynomial{coefficients: []field.Element{field.One}}
}

// X returns the polynomial x (identity polynomial).
func X() *Polynomial {
	return &Polynomial{coefficients: []field.Element{field.Zero, field.One}}
}

// XToThe returns x^n.
func XToThe(n int) *Polynomial {
	if n < 0 {
		panic("negative exponent")
	}
	if n == 0 {
		return One()
	}
	coeffs := make([]field.Element, n+1)
	for i := range coeffs {
		coeffs[i] = field.Zero
	}
	coeffs[n] = field.One
	return &Polynomial{coefficients: coeffs}
}

// Degree returns the degree of the polynomial.
// Returns -1 for the zero polynomial.
func (p *Polynomial) Degree() int {
	deg := len(p.coefficients) - 1
	for deg >= 0 && p.coefficients[deg].IsZero() {
		deg--
	}
	return deg
}

// Coefficients returns the polynomial's coefficients in order of increasing degree.
// The leading coefficient is guaranteed to be non-zero (except for the zero polynomial).
func (p *Polynomial) Coefficients() []field.Element {
	deg := p.Degree()
	if deg < 0 {
		return []field.Element{}
	}
	return p.coefficients[:deg+1]
}

// LeadingCoefficient returns the leading coefficient (coefficient of highest degree term).
// Returns Zero for the zero polynomial.
func (p *Polynomial) LeadingCoefficient() field.Element {
	deg := p.Degree()
	if deg < 0 {
		return field.Zero
	}
	return p.coefficients[deg]
}

// IsZero returns true if this is the zero polynomial.
func (p *Polynomial) IsZero() bool {
	return p.Degree() < 0
}

// IsOne returns true if this is the constant polynomial 1.
func (p *Polynomial) IsOne() bool {
	return p.Degree() == 0 && p.coefficients[0].IsOne()
}

// IsX returns true if this is the polynomial x.
func (p *Polynomial) IsX() bool {
	return p.Degree() == 1 && p.coefficients[0].IsZero() && p.coefficients[1].IsOne()
}

// Equal returns true if two polynomials are equal.
func (p *Polynomial) Equal(other *Polynomial) bool {
	if p.Degree() != other.Degree() {
		return false
	}

	for i := 0; i <= p.Degree(); i++ {
		if !p.coefficients[i].Equal(other.coefficients[i]) {
			return false
		}
	}
	return true
}

// Clone creates a deep copy of the polynomial.
func (p *Polynomial) Clone() *Polynomial {
	coeffs := make([]field.Element, len(p.coefficients))
	copy(coeffs, p.coefficients)
	return &Polynomial{coefficients: coeffs}
}

// normalize removes leading zero coefficients.
func (p *Polynomial) normalize() {
	for len(p.coefficients) > 0 && p.coefficients[len(p.coefficients)-1].IsZero() {
		p.coefficients = p.coefficients[:len(p.coefficients)-1]
	}
}

// Add adds two polynomials.
func (p *Polynomial) Add(other *Polynomial) *Polynomial {
	maxLen := len(p.coefficients)
	if len(other.coefficients) > maxLen {
		maxLen = len(other.coefficients)
	}

	coeffs := make([]field.Element, maxLen)
	for i := range coeffs {
		var a, b field.Element
		if i < len(p.coefficients) {
			a = p.coefficients[i]
		} else {
			a = field.Zero
		}
		if i < len(other.coefficients) {
			b = other.coefficients[i]
		} else {
			b = field.Zero
		}
		coeffs[i] = a.Add(b)
	}

	return New(coeffs)
}

// Sub subtracts another polynomial from this one.
func (p *Polynomial) Sub(other *Polynomial) *Polynomial {
	maxLen := len(p.coefficients)
	if len(other.coefficients) > maxLen {
		maxLen = len(other.coefficients)
	}

	coeffs := make([]field.Element, maxLen)
	for i := range coeffs {
		var a, b field.Element
		if i < len(p.coefficients) {
			a = p.coefficients[i]
		} else {
			a = field.Zero
		}
		if i < len(other.coefficients) {
			b = other.coefficients[i]
		} else {
			b = field.Zero
		}
		coeffs[i] = a.Sub(b)
	}

	return New(coeffs)
}

// Neg returns the negation of the polynomial.
func (p *Polynomial) Neg() *Polynomial {
	coeffs := make([]field.Element, len(p.coefficients))
	for i, c := range p.coefficients {
		coeffs[i] = c.Neg()
	}
	return &Polynomial{coefficients: coeffs}
}

// Mul multiplies two polynomials using naive O(n²) algorithm.
// For faster multiplication with NTT, use MulNTT.
func (p *Polynomial) Mul(other *Polynomial) *Polynomial {
	if p.IsZero() || other.IsZero() {
		return Zero()
	}

	degP := p.Degree()
	degQ := other.Degree()
	resultDeg := degP + degQ

	coeffs := make([]field.Element, resultDeg+1)
	for i := range coeffs {
		coeffs[i] = field.Zero
	}

	for i := 0; i <= degP; i++ {
		for j := 0; j <= degQ; j++ {
			product := p.coefficients[i].Mul(other.coefficients[j])
			coeffs[i+j] = coeffs[i+j].Add(product)
		}
	}

	return &Polynomial{coefficients: coeffs}
}

// ScalarMul multiplies the polynomial by a scalar.
func (p *Polynomial) ScalarMul(scalar field.Element) *Polynomial {
	if scalar.IsZero() {
		return Zero()
	}

	coeffs := make([]field.Element, len(p.coefficients))
	for i, c := range p.coefficients {
		coeffs[i] = c.Mul(scalar)
	}
	return &Polynomial{coefficients: coeffs}
}

// Evaluate evaluates the polynomial at a given point using Horner's method.
func (p *Polynomial) Evaluate(x field.Element) field.Element {
	if p.IsZero() {
		return field.Zero
	}

	// Horner's method: a0 + x(a1 + x(a2 + x(...)))
	result := p.coefficients[len(p.coefficients)-1]
	for i := len(p.coefficients) - 2; i >= 0; i-- {
		result = result.Mul(x).Add(p.coefficients[i])
	}
	return result
}

// BatchEvaluate evaluates the polynomial at multiple points.
func (p *Polynomial) BatchEvaluate(points []field.Element) []field.Element {
	results := make([]field.Element, len(points))
	for i, point := range points {
		results[i] = p.Evaluate(point)
	}
	return results
}

// FormalDerivative computes the formal derivative of the polynomial.
// For p(x) = a0 + a1*x + a2*x^2 + ..., returns p'(x) = a1 + 2*a2*x + 3*a3*x^2 + ...
func (p *Polynomial) FormalDerivative() *Polynomial {
	if p.Degree() <= 0 {
		return Zero()
	}

	coeffs := make([]field.Element, len(p.coefficients)-1)
	for i := 1; i < len(p.coefficients); i++ {
		// Coefficient * index
		coeffs[i-1] = p.coefficients[i].Mul(field.New(uint64(i)))
	}
	return New(coeffs)
}

// Shift shifts the polynomial: returns p(x - offset).
func (p *Polynomial) Shift(offset field.Element) *Polynomial {
	if p.IsZero() {
		return Zero()
	}

	// Use Horner's method for polynomial composition
	result := New([]field.Element{p.coefficients[len(p.coefficients)-1]})
	xMinusOffset := New([]field.Element{offset.Neg(), field.One})

	for i := len(p.coefficients) - 2; i >= 0; i-- {
		result = result.Mul(xMinusOffset)
		result = result.Add(New([]field.Element{p.coefficients[i]}))
	}

	return result
}

// Scale scales the polynomial: returns p(alpha * x) for a scalar alpha.
func (p *Polynomial) Scale(alpha field.Element) *Polynomial {
	if p.IsZero() || alpha.IsZero() {
		return Zero()
	}

	coeffs := make([]field.Element, len(p.coefficients))
	alphaPower := field.One
	for i := range p.coefficients {
		coeffs[i] = p.coefficients[i].Mul(alphaPower)
		alphaPower = alphaPower.Mul(alpha)
	}
	return New(coeffs)
}

// Monic returns a monic version of the polynomial (leading coefficient = 1).
// Panics if the polynomial is zero.
func (p *Polynomial) Monic() *Polynomial {
	if p.IsZero() {
		panic("cannot make zero polynomial monic")
	}

	leadingCoeff := p.LeadingCoefficient()
	if leadingCoeff.IsOne() {
		return p.Clone()
	}

	inv := leadingCoeff.Inverse()
	return p.ScalarMul(inv)
}

// String returns a string representation of the polynomial.
func (p *Polynomial) String() string {
	if p.IsZero() {
		return "0"
	}

	result := ""
	deg := p.Degree()
	for i := deg; i >= 0; i-- {
		coeff := p.coefficients[i]
		if coeff.IsZero() {
			continue
		}

		if result != "" {
			result += " + "
		}

		if !coeff.IsOne() || i == 0 {
			result += fmt.Sprintf("%v", coeff.Value())
		}

		switch i {
		case 0:
			// Just the coefficient
		case 1:
			result += "x"
		default:
			result += fmt.Sprintf("x^%d", i)
		}
	}

	if result == "" {
		return "0"
	}
	return result
}

// Interpolate performs Lagrange interpolation through the given points.
// Points are (x, y) pairs where y = p(x) for some polynomial p.
// Returns the unique polynomial of degree at most n-1 that passes through all n points.
//
// Panics if:
// - points is empty
// - any two points have the same x-coordinate
func Interpolate(points [][2]field.Element) *Polynomial {
	if len(points) == 0 {
		panic("cannot interpolate through zero points")
	}

	// Check for duplicate x-coordinates
	for i := 0; i < len(points); i++ {
		for j := i + 1; j < len(points); j++ {
			if points[i][0].Equal(points[j][0]) {
				panic("duplicate x-coordinates in interpolation points")
			}
		}
	}

	result := Zero()

	// Lagrange interpolation
	for i, point := range points {
		xi, yi := point[0], point[1]

		// Build Lagrange basis polynomial L_i(x)
		basis := One()
		denominator := field.One

		for j, otherPoint := range points {
			if i == j {
				continue
			}
			xj := otherPoint[0]

			// basis *= (x - xj)
			xMinusXj := New([]field.Element{xj.Neg(), field.One})
			basis = basis.Mul(xMinusXj)

			// denominator *= (xi - xj)
			denominator = denominator.Mul(xi.Sub(xj))
		}

		// basis *= yi / denominator
		basis = basis.ScalarMul(yi.Mul(denominator.Inverse()))

		result = result.Add(basis)
	}

	return result
}

// Zerofier returns the polynomial that has zeros at all given points.
// That is, returns (x - points[0]) * (x - points[1]) * ... * (x - points[n-1]).
func Zerofier(points []field.Element) *Polynomial {
	if len(points) == 0 {
		return One()
	}

	result := One()
	for _, point := range points {
		// Multiply by (x - point)
		linear := New([]field.Element{point.Neg(), field.One})
		result = result.Mul(linear)
	}

	return result
}

// XGCD computes the Extended Euclidean Algorithm for polynomials.
// Returns (gcd, a, b) such that: gcd = a*x + b*y
// The gcd is normalized to have leading coefficient 1.
//
// This is used to compute polynomial inverses in quotient rings.
//
// Production implementation.
func XGCD(x, y *Polynomial) (*Polynomial, *Polynomial, *Polynomial) {
	// Make copies to avoid modifying originals
	xCopy := x.Clone()
	yCopy := y.Clone()

	// Initialize Bezout coefficients - following twenty-first exactly
	aFactor, a1 := One(), Zero()
	bFactor, b1 := Zero(), One()

	// Extended Euclidean algorithm - following twenty-first exactly
	for !yCopy.IsZero() {
		quotient, remainder := xCopy.Divide(yCopy)

		// Update Bezout coefficients - following twenty-first exactly
		c := aFactor.Sub(quotient.Mul(a1))
		d := bFactor.Sub(quotient.Mul(b1))

		// Shift values - following twenty-first exactly
		xCopy = yCopy
		yCopy = remainder
		aFactor = a1
		a1 = c
		bFactor = b1
		b1 = d
	}

	// Normalize result to ensure the gcd has leading coefficient 1
	// Following twenty-first exactly
	lc := xCopy.LeadingCoefficient()
	if lc.IsZero() {
		lc = field.One
	}

	lcInv := lc.Inverse()
	gcd := xCopy.ScalarMul(lcInv)
	aResult := aFactor.ScalarMul(lcInv)
	bResult := bFactor.ScalarMul(lcInv)

	return gcd, aResult, bResult
}
