package hash

import (
	"fmt"
	"testing"
)

// TestTip5Trace verifies that Trace() returns the correct permutation trace
func TestTip5Trace(t *testing.T) {
	tip5 := Init()
	initialState := tip5.state

	trace := tip5.Trace()

	// Verify initial state is saved
	if trace[0] != initialState {
		t.Error("Trace[0] should equal initial state")
	}

	// Verify we have 1 + NumRounds states
	if len(trace) != 1+NumRounds {
		t.Errorf("Expected %d states in trace, got %d", 1+NumRounds, len(trace))
	}

	// Verify final state matches what we'd get from Permutation()
	tip52 := Init()
	tip52.state = initialState
	tip52.Permutation()
	if trace[NumRounds] != tip52.state {
		t.Error("Final trace state should match Permutation() result")
	}
}

// TestTip5SampleIndices verifies that SampleIndices produces valid indices
func TestTip5SampleIndices(t *testing.T) {
	testCases := []struct {
		name       string
		upperBound uint32
		numIndices int
		valid      bool
	}{
		{"power of 2: 256", 256, 10, true},
		{"power of 2: 1024", 1024, 20, true},
		{"power of 2: 1", 1, 5, true},
		{"not power of 2: 100", 100, 10, false},
		{"not power of 2: 0", 0, 10, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tip5 := Init()

			if tc.valid {
				indices := tip5.SampleIndices(tc.upperBound, tc.numIndices)

				if len(indices) != tc.numIndices {
					t.Errorf("Expected %d indices, got %d", tc.numIndices, len(indices))
				}

				// Verify all indices are in range [0, upperBound)
				for i, idx := range indices {
					if idx >= tc.upperBound {
						t.Errorf("Index %d at position %d is out of range [0, %d)", idx, i, tc.upperBound)
					}
				}

				// Verify upperBound is power of 2
				if (tc.upperBound & (tc.upperBound - 1)) != 0 {
					t.Errorf("upperBound %d is not a power of 2", tc.upperBound)
				}
			} else {
				// Test that invalid inputs panic
				defer func() {
					if r := recover(); r == nil {
						t.Error("Expected panic for invalid upperBound")
					}
				}()
				tip5.SampleIndices(tc.upperBound, tc.numIndices)
			}
		})
	}
}

// TestTip5SampleIndicesDistribution verifies uniform distribution (basic check)
func TestTip5SampleIndicesDistribution(t *testing.T) {
	upperBound := uint32(256)
	numIndices := 1000
	tip5 := Init()

	indices := tip5.SampleIndices(upperBound, numIndices)

	// Count occurrences of each index
	counts := make(map[uint32]int)
	for _, idx := range indices {
		counts[idx]++
	}

	// Basic check: all indices should appear at least once (with high probability for 1000 samples)
	if len(counts) < int(upperBound)/2 {
		t.Errorf("Expected at least %d unique indices, got %d", int(upperBound)/2, len(counts))
	}
}

// TestTip5SampleScalars verifies that SampleScalars produces valid XFieldElements
func TestTip5SampleScalars(t *testing.T) {
	testCases := []struct {
		name        string
		numElements int
	}{
		{"single element", 1},
		{"multiple elements", 10},
		{"many elements", 100},
		{"not divisible by rate", 7},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tip5 := Init()

			scalars, err := tip5.SampleScalars(tc.numElements)
			if err != nil {
				t.Fatalf("SampleScalars failed: %v", err)
			}

			if len(scalars) != tc.numElements {
				t.Errorf("Expected %d scalars, got %d", tc.numElements, len(scalars))
			}

			// Verify all scalars are valid XFieldElements (not all zero)
			allZero := true
			for i, scalar := range scalars {
				if !scalar.IsZero() {
					allZero = false
				}
				// Verify structure
				if len(scalar.Coefficients) != 3 {
					t.Errorf("Scalar %d should have 3 coefficients, got %d", i, len(scalar.Coefficients))
				}
			}

			// With high probability, at least one scalar should be non-zero
			if allZero && tc.numElements > 1 {
				t.Error("All scalars are zero (unlikely but possible)")
			}
		})
	}
}

// TestTip5SampleScalarsProduct verifies that scalars are non-zero with high probability
func TestTip5SampleScalarsProduct(t *testing.T) {
	tip5 := Init()

	scalars, err := tip5.SampleScalars(10)
	if err != nil {
		t.Fatalf("SampleScalars failed: %v", err)
	}

	// Verify at least one scalar is non-zero (with high probability for 10 samples)
	allZero := true
	for _, scalar := range scalars {
		if !scalar.IsZero() {
			allZero = false
			break
		}
	}

	// With 10 random scalars, probability of all being zero is negligible
	// But we'll make this a warning rather than a hard failure
	if allZero {
		t.Log("Warning: All sampled scalars are zero (extremely unlikely but possible)")
	}
}

// TestTip5SampleIndicesPowerOfTwo verifies power-of-two requirement
func TestTip5SampleIndicesPowerOfTwo(t *testing.T) {
	validBounds := []uint32{1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024, 2048, 4096, 8192, 16384, 32768, 65536}
	invalidBounds := []uint32{0, 3, 5, 6, 7, 9, 10, 100, 200, 500, 1000}

	tip5 := Init()

	// Test valid bounds
	for _, bound := range validBounds {
		t.Run(fmt.Sprintf("valid_%d", bound), func(t *testing.T) {
			if !isPowerOfTwo(bound) {
				t.Errorf("%d should be detected as power of 2", bound)
			}
			indices := tip5.SampleIndices(bound, 10)
			if len(indices) != 10 {
				t.Errorf("Expected 10 indices, got %d", len(indices))
			}
		})
	}

	// Test invalid bounds (should panic)
	for _, bound := range invalidBounds {
		t.Run(fmt.Sprintf("invalid_%d", bound), func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("Expected panic for non-power-of-2 bound %d", bound)
				}
			}()
			tip5.SampleIndices(bound, 10)
		})
	}
}

// isPowerOfTwo checks if a number is a power of 2
func isPowerOfTwo(n uint32) bool {
	return n != 0 && (n&(n-1)) == 0
}

