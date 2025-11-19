package polynomial

import (
	"github.com/vybium/vybium-crypto/pkg/vybium-crypto/field"
	"github.com/vybium/vybium-crypto/pkg/vybium-crypto/ntt"
)

// MulNTT multiplies two polynomials using the Number Theoretic Transform (NTT).
// This is asymptotically faster than naive multiplication for large polynomials.
//
// Time complexity: O(n log n) where n is the size of the result
// Space complexity: O(n)
//
// Production implementation.
func (p *Polynomial) MulNTT(other *Polynomial) *Polynomial {
	if p.IsZero() || other.IsZero() {
		return Zero()
	}

	// For small polynomials, naive multiplication is faster
	degP := p.Degree()
	degQ := other.Degree()
	if degP < 8 && degQ < 8 {
		return p.Mul(other)
	}

	resultDegree := degP + degQ

	// Find the next power of 2 that can hold the result
	resultSize := ntt.NextPowerOfTwo(resultDegree + 1)

	// Pad coefficients to power of 2
	pCoeffs := make([]field.Element, resultSize)
	copy(pCoeffs, p.coefficients)
	for i := len(p.coefficients); i < resultSize; i++ {
		pCoeffs[i] = field.Zero
	}

	qCoeffs := make([]field.Element, resultSize)
	copy(qCoeffs, other.coefficients)
	for i := len(other.coefficients); i < resultSize; i++ {
		qCoeffs[i] = field.Zero
	}

	// Apply NTT to both polynomials
	ntt.NTT(pCoeffs)
	ntt.NTT(qCoeffs)

	// Point-wise multiplication in frequency domain
	for i := 0; i < resultSize; i++ {
		pCoeffs[i] = pCoeffs[i].Mul(qCoeffs[i])
	}

	// Apply inverse NTT to get result
	ntt.INTT(pCoeffs)

	// Create result polynomial (will be normalized automatically)
	return New(pCoeffs)
}

// EvaluateNTT evaluates the polynomial at multiple points using NTT.
// Points must be powers of a primitive root of unity.
//
// This is more efficient than BatchEvaluate when the evaluation points
// follow this specific structure.
func (p *Polynomial) EvaluateNTT(domainSize int) []field.Element {
	if !ntt.IsPowerOfTwo(domainSize) {
		panic("domain size must be a power of 2")
	}

	// Pad coefficients to domain size
	coeffs := make([]field.Element, domainSize)
	copy(coeffs, p.coefficients)
	for i := len(p.coefficients); i < domainSize; i++ {
		coeffs[i] = field.Zero
	}

	// NTT gives evaluations at powers of primitive root
	ntt.NTT(coeffs)

	return coeffs
}

// InterpolateNTT performs polynomial interpolation using NTT.
// The evaluation domain is powers of a primitive root of unity.
//
// values[i] is the polynomial's evaluation at omega^i, where omega is
// a primitive n-th root of unity and n = len(values).
//
// This is the inverse operation of EvaluateNTT.
func InterpolateNTT(values []field.Element) *Polynomial {
	if len(values) == 0 {
		return Zero()
	}

	if !ntt.IsPowerOfTwo(len(values)) {
		panic("number of values must be a power of 2")
	}

	// Make a copy to avoid modifying input
	coeffs := make([]field.Element, len(values))
	copy(coeffs, values)

	// INTT converts evaluations back to coefficients
	ntt.INTT(coeffs)

	return New(coeffs)
}

// DivideNTT divides two polynomials using NTT-based multiplication.
// Returns (quotient, remainder) such that p = quotient * other + remainder.
//
// Panics if other is zero.
func (p *Polynomial) DivideNTT(other *Polynomial) (quotient, remainder *Polynomial) {
	if other.IsZero() {
		panic("division by zero polynomial")
	}

	degP := p.Degree()
	degQ := other.Degree()

	if degP < degQ {
		return Zero(), p.Clone()
	}

	// For small polynomials, use naive division
	if degP < 16 || degQ < 8 {
		return p.Divide(other)
	}

	// Use fast division based on Newton iteration and NTT
	// For now, fall back to naive division
	// Use NTT for fast polynomial division when possible
	return p.Divide(other)
}

// Divide performs naive polynomial division.
// Returns (quotient, remainder) such that p = quotient * other + remainder.
//
// Panics if other is zero.
func (p *Polynomial) Divide(other *Polynomial) (quotient, remainder *Polynomial) {
	if other.IsZero() {
		panic("division by zero polynomial")
	}

	degP := p.Degree()
	degQ := other.Degree()

	if degP < degQ {
		return Zero(), p.Clone()
	}

	// Clone to avoid modifying original
	remainder = p.Clone()
	quotientDegree := degP - degQ
	quotientCoeffs := make([]field.Element, quotientDegree+1)

	leadingCoeffInv := other.LeadingCoefficient().Inverse()

	for i := quotientDegree; i >= 0; i-- {
		// Normalize remainder to get correct degree
		remainder.normalize()

		// Check if remainder degree is less than divisor degree
		remDeg := remainder.Degree()
		if remDeg < degQ {
			// Remaining quotient coefficients are zero
			break
		}

		// The quotient coefficient at position i should eliminate the term at degree degQ + i
		// from the remainder. But the remainder's current leading term is at degree remDeg.
		// So we should only compute this quotient coefficient if remDeg == degQ + i.
		if remDeg != degQ+i {
			// Remainder degree doesn't match expected, skip this iteration
			continue
		}

		// Compute quotient coefficient using the leading coefficient of remainder
		quotCoeff := remainder.LeadingCoefficient().Mul(leadingCoeffInv)
		quotientCoeffs[i] = quotCoeff

		// Subtract other * quotCoeff * x^i from remainder
		for j := 0; j <= degQ; j++ {
			remIdx := i + j
			if remIdx < len(remainder.coefficients) {
				sub := other.coefficients[j].Mul(quotCoeff)
				remainder.coefficients[remIdx] = remainder.coefficients[remIdx].Sub(sub)
			}
		}
	}

	quotient = &Polynomial{coefficients: quotientCoeffs}
	quotient.normalize()
	remainder.normalize()

	return quotient, remainder
}

// Mod returns p mod other (the remainder of division).
//
// Panics if other is zero.
func (p *Polynomial) Mod(other *Polynomial) *Polynomial {
	_, remainder := p.Divide(other)
	return remainder
}
