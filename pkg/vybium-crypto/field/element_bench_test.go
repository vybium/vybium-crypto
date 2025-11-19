package field

import (
	"testing"
)

// Benchmarks for field operations
// Target performance (compared to twenty-first in Rust):
// - Add: < 5 ns (Rust: ~2 ns)
// - Mul: < 10 ns (Rust: ~5 ns)
// - Inv: < 500 ns (Rust: ~300 ns)

func BenchmarkElementNew(b *testing.B) {
	var result Element
	for i := 0; i < b.N; i++ {
		result = New(uint64(i))
	}
	_ = result
}

func BenchmarkElementAdd(b *testing.B) {
	a := New(123456789)
	c := New(987654321)
	var result Element

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result = a.Add(c)
	}
	_ = result
}

func BenchmarkElementSub(b *testing.B) {
	a := New(987654321)
	c := New(123456789)
	var result Element

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result = a.Sub(c)
	}
	_ = result
}

func BenchmarkElementMul(b *testing.B) {
	a := New(123456789)
	c := New(987654321)
	var result Element

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result = a.Mul(c)
	}
	_ = result
}

func BenchmarkElementSquare(b *testing.B) {
	a := New(123456789)
	var result Element

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result = a.Square()
	}
	_ = result
}

func BenchmarkElementInverse(b *testing.B) {
	a := New(123456789)
	var result Element

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result = a.Inverse()
	}
	_ = result
}

func BenchmarkElementDiv(b *testing.B) {
	a := New(987654321)
	c := New(123456789)
	var result Element

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result = a.Div(c)
	}
	_ = result
}

func BenchmarkElementModPow(b *testing.B) {
	a := New(123456789)
	exp := uint64(12345)
	var result Element

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result = a.ModPow(exp)
	}
	_ = result
}

func BenchmarkElementNeg(b *testing.B) {
	a := New(123456789)
	var result Element

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result = a.Neg()
	}
	_ = result
}

func BenchmarkElementValue(b *testing.B) {
	a := New(123456789)
	var result uint64

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result = a.Value()
	}
	_ = result
}

// Benchmark Montgomery reduction directly
func BenchmarkMontyred(b *testing.B) {
	x := uint128{lo: 0x123456789ABCDEF0, hi: 0xFEDCBA9876543210}
	var result uint64

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result = montyred(x)
	}
	_ = result
}

// Benchmark 128-bit multiplication
func BenchmarkMul128(b *testing.B) {
	a := uint64(0x123456789ABCDEF0)
	c := uint64(0xFEDCBA9876543210)
	var result uint128

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result = mul128(a, c)
	}
	_ = result
}

// Benchmark batch operations
func BenchmarkElementBatchAdd(b *testing.B) {
	// Simulate batch addition (common in polynomial operations)
	size := 1024
	elements := make([]Element, size)
	for i := 0; i < size; i++ {
		elements[i] = New(uint64(i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sum := Zero
		for j := 0; j < size; j++ {
			sum = sum.Add(elements[j])
		}
	}
}

func BenchmarkElementBatchMul(b *testing.B) {
	// Simulate batch multiplication
	size := 1024
	elements := make([]Element, size)
	for i := 0; i < size; i++ {
		elements[i] = New(uint64(i + 1))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		product := One
		for j := 0; j < size; j++ {
			product = product.Mul(elements[j])
		}
	}
}

// Benchmark serialization
func BenchmarkElementMarshalBinary(b *testing.B) {
	a := New(123456789)
	var result []byte
	var err error

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result, err = a.MarshalBinary()
		if err != nil {
			b.Fatal(err)
		}
	}
	_ = result
}

func BenchmarkElementUnmarshalBinary(b *testing.B) {
	a := New(123456789)
	data, _ := a.MarshalBinary()
	var result Element

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := result.UnmarshalBinary(data)
		if err != nil {
			b.Fatal(err)
		}
	}
	_ = result
}

// Benchmark common patterns
func BenchmarkElementLinearCombination(b *testing.B) {
	// Compute a*x + b*y + c*z (common in polynomial evaluation)
	a := New(123)
	x := New(456)
	c := New(789)
	y := New(101112)
	d := New(131415)
	z := New(161718)
	var result Element

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result = a.Mul(x).Add(c.Mul(y)).Add(d.Mul(z))
	}
	_ = result
}

func BenchmarkElementHornerEvaluation(b *testing.B) {
	// Horner's method for polynomial evaluation: a0 + x(a1 + x(a2 + x*a3))
	coeffs := []Element{
		New(123),
		New(456),
		New(789),
		New(101112),
		New(131415),
	}
	x := New(42)
	var result Element

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result = coeffs[len(coeffs)-1]
		for j := len(coeffs) - 2; j >= 0; j-- {
			result = result.Mul(x).Add(coeffs[j])
		}
	}
	_ = result
}
