package polynomial

import (
	"testing"

	"github.com/vybium/vybium-crypto/pkg/vybium-crypto/field"
)

func TestPolynomialCreation(t *testing.T) {
	// Test creating polynomials
	coeffs := []field.Element{
		field.New(1), field.New(2), field.New(3),
	}
	p := New(coeffs)

	if p.Degree() != 2 {
		t.Errorf("Expected degree 2, got %d", p.Degree())
	}

	if len(p.Coefficients()) != 3 {
		t.Errorf("Expected 3 coefficients, got %d", len(p.Coefficients()))
	}
}

func TestZeroPolynomial(t *testing.T) {
	p := Zero()

	if !p.IsZero() {
		t.Error("Zero() should create zero polynomial")
	}

	if p.Degree() != -1 {
		t.Errorf("Zero polynomial should have degree -1, got %d", p.Degree())
	}

	if len(p.Coefficients()) != 0 {
		t.Error("Zero polynomial should have empty coefficients")
	}
}

func TestOnePolynomial(t *testing.T) {
	p := One()

	if !p.IsOne() {
		t.Error("One() should create constant polynomial 1")
	}

	if p.Degree() != 0 {
		t.Errorf("Constant polynomial 1 should have degree 0, got %d", p.Degree())
	}
}

func TestXPolynomial(t *testing.T) {
	p := X()

	if !p.IsX() {
		t.Error("X() should create polynomial x")
	}

	if p.Degree() != 1 {
		t.Errorf("X polynomial should have degree 1, got %d", p.Degree())
	}
}

func TestXToThe(t *testing.T) {
	tests := []struct {
		n      int
		degree int
	}{
		{0, 0},
		{1, 1},
		{2, 2},
		{10, 10},
	}

	for _, tt := range tests {
		p := XToThe(tt.n)
		if p.Degree() != tt.degree {
			t.Errorf("XToThe(%d) should have degree %d, got %d",
				tt.n, tt.degree, p.Degree())
		}
	}
}

func TestPolynomialAddition(t *testing.T) {
	// (1 + 2x) + (3 + 4x) = 4 + 6x
	p1 := New([]field.Element{field.New(1), field.New(2)})
	p2 := New([]field.Element{field.New(3), field.New(4)})
	result := p1.Add(p2)

	expected := New([]field.Element{field.New(4), field.New(6)})
	if !result.Equal(expected) {
		t.Error("Polynomial addition failed")
	}
}

func TestPolynomialSubtraction(t *testing.T) {
	// (5 + 3x) - (2 + x) = 3 + 2x
	p1 := New([]field.Element{field.New(5), field.New(3)})
	p2 := New([]field.Element{field.New(2), field.New(1)})
	result := p1.Sub(p2)

	expected := New([]field.Element{field.New(3), field.New(2)})
	if !result.Equal(expected) {
		t.Error("Polynomial subtraction failed")
	}
}

func TestPolynomialNegation(t *testing.T) {
	// -(1 + 2x + 3x^2) = -1 - 2x - 3x^2
	p := New([]field.Element{field.New(1), field.New(2), field.New(3)})
	negP := p.Neg()

	// Check that p + (-p) = 0
	sum := p.Add(negP)
	if !sum.IsZero() {
		t.Error("p + (-p) should be zero")
	}
}

func TestPolynomialMultiplication(t *testing.T) {
	// (1 + x) * (1 + x) = 1 + 2x + x^2
	p := New([]field.Element{field.New(1), field.New(1)})
	result := p.Mul(p)

	expected := New([]field.Element{field.New(1), field.New(2), field.New(1)})
	if !result.Equal(expected) {
		t.Errorf("Polynomial multiplication failed: got %v", result.Coefficients())
	}
}

func TestPolynomialScalarMultiplication(t *testing.T) {
	// 3 * (1 + 2x) = 3 + 6x
	p := New([]field.Element{field.New(1), field.New(2)})
	result := p.ScalarMul(field.New(3))

	expected := New([]field.Element{field.New(3), field.New(6)})
	if !result.Equal(expected) {
		t.Error("Scalar multiplication failed")
	}
}

func TestPolynomialEvaluation(t *testing.T) {
	// p(x) = 1 + 2x + 3x^2
	// p(0) = 1, p(1) = 6, p(2) = 17
	p := New([]field.Element{field.New(1), field.New(2), field.New(3)})

	tests := []struct {
		x        uint64
		expected uint64
	}{
		{0, 1},
		{1, 6},
		{2, 17},
	}

	for _, tt := range tests {
		result := p.Evaluate(field.New(tt.x))
		expected := field.New(tt.expected)
		if !result.Equal(expected) {
			t.Errorf("p(%d) = %d, expected %d", tt.x, result.Value(), tt.expected)
		}
	}
}

func TestPolynomialFormalDerivative(t *testing.T) {
	// p(x) = 1 + 2x + 3x^2 + 4x^3
	// p'(x) = 2 + 6x + 12x^2
	p := New([]field.Element{field.New(1), field.New(2), field.New(3), field.New(4)})
	derivative := p.FormalDerivative()

	expected := New([]field.Element{field.New(2), field.New(6), field.New(12)})
	if !derivative.Equal(expected) {
		t.Error("Formal derivative failed")
	}
}

func TestPolynomialMonic(t *testing.T) {
	// 2 + 4x + 6x^2 -> (1/6)(2 + 4x + 6x^2)
	p := New([]field.Element{field.New(2), field.New(4), field.New(6)})
	monic := p.Monic()

	if !monic.LeadingCoefficient().IsOne() {
		t.Error("Monic polynomial should have leading coefficient 1")
	}

	if monic.Degree() != p.Degree() {
		t.Error("Monic polynomial should have same degree")
	}
}

func TestInterpolation(t *testing.T) {
	// Interpolate through (0,1), (1,3), (2,7)
	// Unique polynomial of degree ≤ 2: 1 + x + x^2
	points := [][2]field.Element{
		{field.New(0), field.New(1)},
		{field.New(1), field.New(3)},
		{field.New(2), field.New(7)},
	}

	p := Interpolate(points)

	// Check that polynomial passes through all points
	for _, point := range points {
		x, y := point[0], point[1]
		result := p.Evaluate(x)
		if !result.Equal(y) {
			t.Errorf("Interpolated polynomial doesn't pass through (%d, %d): got %d",
				x.Value(), y.Value(), result.Value())
		}
	}
}

func TestZerofier(t *testing.T) {
	// Zerofier of {1, 2, 3} should be (x-1)(x-2)(x-3)
	points := []field.Element{field.New(1), field.New(2), field.New(3)}
	z := Zerofier(points)

	// Check that it evaluates to zero at all points
	for _, point := range points {
		result := z.Evaluate(point)
		if !result.IsZero() {
			t.Errorf("Zerofier should be zero at %d, got %d",
				point.Value(), result.Value())
		}
	}

	// Check degree
	if z.Degree() != len(points) {
		t.Errorf("Zerofier degree should be %d, got %d", len(points), z.Degree())
	}
}

func TestPolynomialDivision(t *testing.T) {
	// (x^2 + 2x + 1) / (x + 1) = (x + 1) with remainder 0
	dividend := New([]field.Element{field.New(1), field.New(2), field.New(1)})
	divisor := New([]field.Element{field.New(1), field.New(1)})

	quotient, remainder := dividend.Divide(divisor)

	// Check: dividend = quotient * divisor + remainder
	check := quotient.Mul(divisor).Add(remainder)
	if !check.Equal(dividend) {
		t.Error("Division check failed: dividend != quotient * divisor + remainder")
	}

	// For this specific case, remainder should be zero
	if !remainder.IsZero() {
		t.Error("Expected zero remainder")
	}
}

func TestPolynomialDivisionWithRemainder(t *testing.T) {
	// (x^2 + 1) / (x + 1) should have non-zero remainder
	dividend := New([]field.Element{field.New(1), field.New(0), field.New(1)})
	divisor := New([]field.Element{field.New(1), field.New(1)})

	quotient, remainder := dividend.Divide(divisor)

	// Check: dividend = quotient * divisor + remainder
	check := quotient.Mul(divisor).Add(remainder)
	if !check.Equal(dividend) {
		t.Error("Division check failed")
	}
}

func TestPolynomialNormalization(t *testing.T) {
	// Polynomial with trailing zeros should be normalized
	coeffs := []field.Element{
		field.New(1), field.New(2), field.Zero, field.Zero,
	}
	p := New(coeffs)

	if p.Degree() != 1 {
		t.Errorf("Polynomial with trailing zeros should have degree 1, got %d", p.Degree())
	}
}

func TestPolynomialClone(t *testing.T) {
	original := New([]field.Element{field.New(1), field.New(2), field.New(3)})
	clone := original.Clone()

	// Modify clone
	clone.coefficients[0] = field.New(999)

	// Original should be unchanged
	if original.coefficients[0].Value() == 999 {
		t.Error("Clone modification affected original")
	}
}

func TestPolynomialEqual(t *testing.T) {
	p1 := New([]field.Element{field.New(1), field.New(2)})
	p2 := New([]field.Element{field.New(1), field.New(2)})
	p3 := New([]field.Element{field.New(1), field.New(3)})

	if !p1.Equal(p2) {
		t.Error("Equal polynomials not detected as equal")
	}

	if p1.Equal(p3) {
		t.Error("Different polynomials detected as equal")
	}
}

func TestPolynomialBatchEvaluate(t *testing.T) {
	p := New([]field.Element{field.New(1), field.New(2), field.New(3)})
	points := []field.Element{field.New(0), field.New(1), field.New(2)}

	results := p.BatchEvaluate(points)

	if len(results) != len(points) {
		t.Errorf("Expected %d results, got %d", len(points), len(results))
	}

	// Check each result individually
	for i, point := range points {
		expected := p.Evaluate(point)
		if !results[i].Equal(expected) {
			t.Errorf("Batch evaluate mismatch at index %d", i)
		}
	}
}

// Benchmarks
func BenchmarkPolynomialMultiply(b *testing.B) {
	p1 := New([]field.Element{field.New(1), field.New(2), field.New(3)})
	p2 := New([]field.Element{field.New(4), field.New(5), field.New(6)})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = p1.Mul(p2)
	}
}

func BenchmarkPolynomialEvaluate(b *testing.B) {
	p := New([]field.Element{
		field.New(1), field.New(2), field.New(3), field.New(4), field.New(5),
	})
	x := field.New(42)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = p.Evaluate(x)
	}
}

func BenchmarkPolynomialAdd(b *testing.B) {
	p1 := New([]field.Element{field.New(1), field.New(2), field.New(3)})
	p2 := New([]field.Element{field.New(4), field.New(5), field.New(6)})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = p1.Add(p2)
	}
}
