package ntt

import (
	"testing"

	"github.com/vybium/vybium-crypto/pkg/vybium-crypto/field"
)

func TestNTTINTTRoundTrip(t *testing.T) {
	// Test that NTT followed by INTT returns the original values
	sizes := []int{1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024}

	for _, size := range sizes {
		t.Run("size_"+string(rune(size)), func(t *testing.T) {
			// Create random values
			original := make([]field.Element, size)
			for i := range original {
				original[i] = field.New(uint64(i*7 + 13))
			}

			// Copy for transformation
			values := make([]field.Element, size)
			copy(values, original)

			// Apply NTT then INTT
			NTT(values)
			INTT(values)

			// Check that we got back the original
			for i := range original {
				if !original[i].Equal(values[i]) {
					t.Errorf("Round trip failed at index %d: expected %v, got %v",
						i, original[i].Value(), values[i].Value())
				}
			}
		})
	}
}

func TestNTTEmptySlice(t *testing.T) {
	// NTT/INTT on empty slice should not panic
	var empty []field.Element
	NTT(empty)
	INTT(empty)
}

func TestNTTSingleElement(t *testing.T) {
	// NTT/INTT on single element should be identity
	values := []field.Element{field.New(42)}
	original := values[0]

	NTT(values)
	INTT(values)

	if !values[0].Equal(original) {
		t.Errorf("Single element round trip failed: expected %v, got %v",
			original.Value(), values[0].Value())
	}
}

func TestNTTPanicOnNonPowerOfTwo(t *testing.T) {
	// NTT should panic on non-power-of-2 lengths
	defer func() {
		if r := recover(); r == nil {
			t.Error("NTT should panic on non-power-of-2 length")
		}
	}()

	values := make([]field.Element, 7) // Not a power of 2
	NTT(values)
}

func TestINTTPanicOnNonPowerOfTwo(t *testing.T) {
	// INTT should panic on non-power-of-2 lengths
	defer func() {
		if r := recover(); r == nil {
			t.Error("INTT should panic on non-power-of-2 length")
		}
	}()

	values := make([]field.Element, 7) // Not a power of 2
	INTT(values)
}

func TestBitReverse(t *testing.T) {
	tests := []struct {
		k, log2N, expected uint32
	}{
		{0, 3, 0},   // 000 -> 000
		{1, 3, 4},   // 001 -> 100
		{2, 3, 2},   // 010 -> 010
		{3, 3, 6},   // 011 -> 110
		{4, 3, 1},   // 100 -> 001
		{5, 3, 5},   // 101 -> 101
		{6, 3, 3},   // 110 -> 011
		{7, 3, 7},   // 111 -> 111
		{0, 4, 0},   // 0000 -> 0000
		{1, 4, 8},   // 0001 -> 1000
		{8, 4, 1},   // 1000 -> 0001
		{15, 4, 15}, // 1111 -> 1111
	}

	for _, tt := range tests {
		result := bitReverse(tt.k, tt.log2N)
		if result != tt.expected {
			t.Errorf("bitReverse(%d, %d) = %d, expected %d",
				tt.k, tt.log2N, result, tt.expected)
		}
	}
}

func TestNextPowerOfTwo(t *testing.T) {
	tests := []struct {
		n, expected int
	}{
		{0, 1},
		{1, 1},
		{2, 2},
		{3, 4},
		{4, 4},
		{5, 8},
		{7, 8},
		{8, 8},
		{9, 16},
		{15, 16},
		{16, 16},
		{17, 32},
		{100, 128},
		{1000, 1024},
	}

	for _, tt := range tests {
		result := NextPowerOfTwo(tt.n)
		if result != tt.expected {
			t.Errorf("NextPowerOfTwo(%d) = %d, expected %d",
				tt.n, result, tt.expected)
		}
	}
}

func TestIsPowerOfTwo(t *testing.T) {
	tests := []struct {
		n        int
		expected bool
	}{
		{0, false},
		{1, true},
		{2, true},
		{3, false},
		{4, true},
		{5, false},
		{7, false},
		{8, true},
		{16, true},
		{17, false},
		{1024, true},
		{1023, false},
	}

	for _, tt := range tests {
		result := IsPowerOfTwo(tt.n)
		if result != tt.expected {
			t.Errorf("IsPowerOfTwo(%d) = %v, expected %v",
				tt.n, result, tt.expected)
		}
	}
}

func TestNTTCaching(t *testing.T) {
	// Test that twiddle factors are cached
	size := 16
	values1 := make([]field.Element, size)
	values2 := make([]field.Element, size)

	for i := range values1 {
		values1[i] = field.New(uint64(i))
		values2[i] = field.New(uint64(i))
	}

	// First call should compute twiddle factors
	NTT(values1)

	// Second call should use cached twiddle factors
	NTT(values2)

	// Both should produce same result
	for i := range values1 {
		if !values1[i].Equal(values2[i]) {
			t.Errorf("Cached NTT produced different result at index %d", i)
		}
	}
}

func TestNTTINTTWithMaxElement(t *testing.T) {
	// Test with maximum field element value
	size := 16
	values := make([]field.Element, size)
	for i := range values {
		values[i] = field.Max
	}
	original := make([]field.Element, size)
	copy(original, values)

	NTT(values)
	INTT(values)

	for i := range original {
		if !original[i].Equal(values[i]) {
			t.Errorf("Round trip with Max element failed at index %d", i)
		}
	}
}

// Benchmark NTT of various sizes
func BenchmarkNTT16(b *testing.B) {
	values := make([]field.Element, 16)
	for i := range values {
		values[i] = field.New(uint64(i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NTT(values)
	}
}

func BenchmarkNTT256(b *testing.B) {
	values := make([]field.Element, 256)
	for i := range values {
		values[i] = field.New(uint64(i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NTT(values)
	}
}

func BenchmarkNTT1024(b *testing.B) {
	values := make([]field.Element, 1024)
	for i := range values {
		values[i] = field.New(uint64(i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NTT(values)
	}
}

func BenchmarkINTT1024(b *testing.B) {
	values := make([]field.Element, 1024)
	for i := range values {
		values[i] = field.New(uint64(i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		INTT(values)
	}
}

func BenchmarkNTTINTTRoundTrip1024(b *testing.B) {
	values := make([]field.Element, 1024)
	for i := range values {
		values[i] = field.New(uint64(i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NTT(values)
		INTT(values)
	}
}
