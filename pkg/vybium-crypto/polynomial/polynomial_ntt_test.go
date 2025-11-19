package polynomial

import (
	"testing"

	"github.com/vybium/vybium-crypto/pkg/vybium-crypto/field"
)

// TestMulNTT tests NTT-based polynomial multiplication
func TestMulNTT(t *testing.T) {
	p1 := New([]field.Element{field.New(1), field.New(2), field.New(3)})
	p2 := New([]field.Element{field.New(4), field.New(5)})

	// Test NTT multiplication
	resultNTT := p1.MulNTT(p2)
	resultNaive := p1.Mul(p2)

	if !resultNTT.Equal(resultNaive) {
		t.Error("MulNTT result differs from naive multiplication")
	}
}

// TestEvaluateNTT tests NTT-based polynomial evaluation
func TestEvaluateNTT(t *testing.T) {
	// Create a polynomial
	coeffs := []field.Element{field.New(1), field.New(2), field.New(3), field.New(4)}
	p := New(coeffs)

	// EvaluateNTT evaluates at powers of a primitive root of unity
	domainSize := 8
	resultsNTT := p.EvaluateNTT(domainSize)

	if len(resultsNTT) != domainSize {
		t.Fatalf("EvaluateNTT returned %d results, expected %d", len(resultsNTT), domainSize)
	}

	// Verify by interpolating back
	pReconstructed := InterpolateNTT(resultsNTT)

	// Check that we get the same coefficients (up to original degree)
	for i := 0; i < len(coeffs); i++ {
		if i < len(pReconstructed.coefficients) {
			if !pReconstructed.coefficients[i].Equal(coeffs[i]) {
				t.Errorf("Coefficient mismatch at index %d", i)
			}
		}
	}
}

// TestInterpolateNTT tests NTT-based polynomial interpolation
func TestInterpolateNTT(t *testing.T) {
	// Start with a known polynomial
	originalCoeffs := []field.Element{field.New(1), field.New(2), field.New(3)}
	original := New(originalCoeffs)

	// Evaluate at powers of primitive root (NTT domain)
	domainSize := 8
	evaluations := original.EvaluateNTT(domainSize)

	// Interpolate back using NTT
	reconstructed := InterpolateNTT(evaluations)

	// Verify we get the original coefficients back
	for i := 0; i < len(originalCoeffs); i++ {
		if i < len(reconstructed.coefficients) {
			if !reconstructed.coefficients[i].Equal(originalCoeffs[i]) {
				t.Errorf("Coefficient mismatch at index %d: got %v, want %v",
					i, reconstructed.coefficients[i], originalCoeffs[i])
			}
		} else {
			t.Errorf("Reconstructed polynomial missing coefficient at index %d", i)
		}
	}
}

// TestDivideNTT tests NTT-based polynomial division
func TestDivideNTT(t *testing.T) {
	// dividend = x^3 + 2x^2 + 3x + 4
	dividend := New([]field.Element{field.New(4), field.New(3), field.New(2), field.New(1)})
	// divisor = x + 1
	divisor := New([]field.Element{field.New(1), field.New(1)})

	quotientNTT, remainderNTT := dividend.DivideNTT(divisor)
	quotientNaive, remainderNaive := dividend.Divide(divisor)

	if !quotientNTT.Equal(quotientNaive) {
		t.Error("DivideNTT quotient differs from naive division")
	}

	if !remainderNTT.Equal(remainderNaive) {
		t.Error("DivideNTT remainder differs from naive division")
	}

	// Verify: dividend = quotient * divisor + remainder
	result := quotientNTT.Mul(divisor).Add(remainderNTT)
	if !result.Equal(dividend) {
		t.Error("Division property violated: dividend != quotient * divisor + remainder")
	}
}

// TestMod tests polynomial modular reduction
func TestMod(t *testing.T) {
	// poly = x^3 + 2x^2 + 3x + 4
	poly := New([]field.Element{field.New(4), field.New(3), field.New(2), field.New(1)})
	// modulus = x^2 + 1
	modulus := New([]field.Element{field.New(1), field.Zero, field.New(1)})

	remainder := poly.Mod(modulus)

	// remainder should have degree < modulus degree
	if remainder.Degree() >= modulus.Degree() {
		t.Errorf("Mod remainder degree %d >= modulus degree %d", remainder.Degree(), modulus.Degree())
	}

	// Verify: poly = q * modulus + remainder for some q
	_, rem := poly.Divide(modulus)
	if !remainder.Equal(rem) {
		t.Error("Mod result differs from Divide remainder")
	}
}

// TestScale tests polynomial scaling
func TestScale(t *testing.T) {
	// p(x) = 1 + 2x + 3x^2
	poly := New([]field.Element{field.New(1), field.New(2), field.New(3)})
	alpha := field.New(2)

	// p_scaled(x) = p(alpha * x)
	scaled := poly.Scale(alpha)

	// Test at multiple points: p(alpha * x) should equal scaled(x)
	testPoints := []uint64{0, 1, 2, 3, 5, 10}
	for _, xVal := range testPoints {
		x := field.New(xVal)
		scaledEval := scaled.Evaluate(x)
		originalEval := poly.Evaluate(alpha.Mul(x))

		if !originalEval.Equal(scaledEval) {
			t.Errorf("Scale property violated at x=%v: p(alpha*x)=%v != scaled(x)=%v",
				x, originalEval, scaledEval)
		}
	}
}

// TestPolynomialString tests string representation
func TestPolynomialString(t *testing.T) {
	tests := []struct {
		name string
		poly *Polynomial
	}{
		{"Zero", Zero()},
		{"One", One()},
		{"X", X()},
		{"Constant", New([]field.Element{field.New(42)})},
		{"Linear", New([]field.Element{field.New(1), field.New(2)})},
		{"Quadratic", New([]field.Element{field.New(1), field.New(2), field.New(3)})},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			str := tt.poly.String()
			if str == "" {
				t.Error("String() returned empty string")
			}
		})
	}
}

// TestNTTEdgeCases tests edge cases for NTT operations
func TestNTTEdgeCases(t *testing.T) {
	t.Run("MulNTT with zero", func(t *testing.T) {
		p := New([]field.Element{field.New(1), field.New(2)})
		zero := Zero()
		result := p.MulNTT(zero)
		if !result.IsZero() {
			t.Error("MulNTT with zero should return zero")
		}
	})

	t.Run("MulNTT with one", func(t *testing.T) {
		p := New([]field.Element{field.New(1), field.New(2), field.New(3)})
		one := One()
		result := p.MulNTT(one)
		if !result.Equal(p) {
			t.Error("MulNTT with one should return original polynomial")
		}
	})

	t.Run("DivideNTT by one", func(t *testing.T) {
		p := New([]field.Element{field.New(1), field.New(2), field.New(3)})
		one := One()
		quotient, remainder := p.DivideNTT(one)
		if !quotient.Equal(p) || !remainder.IsZero() {
			t.Error("Division by one should return original polynomial with zero remainder")
		}
	})

	t.Run("Mod by large polynomial", func(t *testing.T) {
		small := New([]field.Element{field.New(1), field.New(2)})
		large := New([]field.Element{field.New(1), field.New(2), field.New(3), field.New(4)})
		remainder := small.Mod(large)
		if !remainder.Equal(small) {
			t.Error("Mod by larger polynomial should return original")
		}
	})

	t.Run("Scale by zero", func(t *testing.T) {
		p := New([]field.Element{field.New(1), field.New(2), field.New(3)})
		scaled := p.Scale(field.Zero)
		// Scaling by zero returns zero polynomial
		if !scaled.IsZero() {
			t.Error("Scale by zero should return zero polynomial")
		}
	})

	t.Run("Scale by one", func(t *testing.T) {
		p := New([]field.Element{field.New(1), field.New(2), field.New(3)})
		scaled := p.Scale(field.One)
		if !scaled.Equal(p) {
			t.Error("Scale by one should return original polynomial")
		}
	})
}

// BenchmarkMulNTT benchmarks NTT multiplication
func BenchmarkMulNTT(b *testing.B) {
	coeffs1 := make([]field.Element, 256)
	coeffs2 := make([]field.Element, 256)
	for i := 0; i < 256; i++ {
		coeffs1[i] = field.New(uint64(i))
		coeffs2[i] = field.New(uint64(i * 2))
	}
	p1 := New(coeffs1)
	p2 := New(coeffs2)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = p1.MulNTT(p2)
	}
}

// BenchmarkEvaluateNTT benchmarks NTT evaluation
func BenchmarkEvaluateNTT(b *testing.B) {
	coeffs := make([]field.Element, 256)
	for i := 0; i < 256; i++ {
		coeffs[i] = field.New(uint64(i))
	}
	p := New(coeffs)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = p.EvaluateNTT(128)
	}
}
