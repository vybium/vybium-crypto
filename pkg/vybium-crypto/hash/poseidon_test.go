package hash

import (
	"testing"

	"github.com/vybium/vybium-crypto/pkg/vybium-crypto/field"
)

func TestPoseidonCreation(t *testing.T) {
	// Test with default parameters (128-bit security)
	poseidon, err := NewPoseidon(nil)
	if err != nil {
		t.Fatalf("Failed to create Poseidon with default parameters: %v", err)
	}

	if poseidon.securityLevel != 128 {
		t.Errorf("Expected security level 128, got %d", poseidon.securityLevel)
	}

	if poseidon.width != 4 {
		t.Errorf("Expected width 4, got %d", poseidon.width)
	}

	if poseidon.rate != 3 {
		t.Errorf("Expected rate 3, got %d", poseidon.rate)
	}
}

func TestPoseidon128BitSecurity(t *testing.T) {
	params := GetDefaultPoseidonParameters(128)
	poseidon, err := NewPoseidon(params)
	if err != nil {
		t.Fatalf("Failed to create Poseidon: %v", err)
	}

	if poseidon.roundsFull != 8 {
		t.Errorf("Expected 8 full rounds, got %d", poseidon.roundsFull)
	}

	if poseidon.roundsPartial != 84 {
		t.Errorf("Expected 84 partial rounds, got %d", poseidon.roundsPartial)
	}
}

func TestPoseidon256BitSecurity(t *testing.T) {
	params := GetDefaultPoseidonParameters(256)
	poseidon, err := NewPoseidon(params)
	if err != nil {
		t.Fatalf("Failed to create Poseidon: %v", err)
	}

	if poseidon.roundsFull != 8 {
		t.Errorf("Expected 8 full rounds, got %d", poseidon.roundsFull)
	}

	if poseidon.roundsPartial != 170 {
		t.Errorf("Expected 170 partial rounds, got %d", poseidon.roundsPartial)
	}
}

func TestPoseidonHashEmptyInput(t *testing.T) {
	poseidon, err := NewPoseidon(nil)
	if err != nil {
		t.Fatalf("Failed to create Poseidon: %v", err)
	}

	result := poseidon.Hash([]field.Element{})
	if !result.IsZero() {
		t.Error("Hash of empty input should be zero")
	}
}

func TestPoseidonHashSingleElement(t *testing.T) {
	poseidon, err := NewPoseidon(nil)
	if err != nil {
		t.Fatalf("Failed to create Poseidon: %v", err)
	}

	input := []field.Element{field.New(42)}
	result := poseidon.Hash(input)

	if result.IsZero() {
		t.Error("Hash of non-empty input should not be zero")
	}
}

func TestPoseidonHashMultipleElements(t *testing.T) {
	poseidon, err := NewPoseidon(nil)
	if err != nil {
		t.Fatalf("Failed to create Poseidon: %v", err)
	}

	inputs := []field.Element{
		field.New(1),
		field.New(2),
		field.New(3),
		field.New(4),
		field.New(5),
	}

	result := poseidon.Hash(inputs)
	if result.IsZero() {
		t.Error("Hash should not be zero")
	}
}

func TestPoseidonHashDeterminism(t *testing.T) {
	poseidon, err := NewPoseidon(nil)
	if err != nil {
		t.Fatalf("Failed to create Poseidon: %v", err)
	}

	inputs := []field.Element{
		field.New(1),
		field.New(2),
		field.New(3),
	}

	result1 := poseidon.Hash(inputs)
	result2 := poseidon.Hash(inputs)

	if !result1.Equal(result2) {
		t.Error("Hash should be deterministic")
	}
}

func TestPoseidonHashDifferentInputs(t *testing.T) {
	poseidon, err := NewPoseidon(nil)
	if err != nil {
		t.Fatalf("Failed to create Poseidon: %v", err)
	}

	inputs1 := []field.Element{field.New(1), field.New(2)}
	inputs2 := []field.Element{field.New(2), field.New(1)}

	result1 := poseidon.Hash(inputs1)
	result2 := poseidon.Hash(inputs2)

	if result1.Equal(result2) {
		t.Error("Different inputs should produce different hashes")
	}
}

func TestPoseidonHashTwo(t *testing.T) {
	poseidon, err := NewPoseidon(nil)
	if err != nil {
		t.Fatalf("Failed to create Poseidon: %v", err)
	}

	left := field.New(123)
	right := field.New(456)

	result := poseidon.HashTwo(left, right)
	if result.IsZero() {
		t.Error("HashTwo result should not be zero")
	}

	// Should be same as hashing array
	result2 := poseidon.Hash([]field.Element{left, right})
	if !result.Equal(result2) {
		t.Error("HashTwo should equal Hash of array")
	}
}

func TestPoseidonSpongeCreation(t *testing.T) {
	sponge, err := NewPoseidonSponge(nil)
	if err != nil {
		t.Fatalf("Failed to create PoseidonSponge: %v", err)
	}

	if sponge.poseidon == nil {
		t.Error("Sponge should have non-nil poseidon instance")
	}

	if len(sponge.state) != sponge.poseidon.width {
		t.Errorf("Sponge state length should be %d, got %d",
			sponge.poseidon.width, len(sponge.state))
	}
}

func TestPoseidonSpongeAbsorbSqueeze(t *testing.T) {
	sponge, err := NewPoseidonSponge(nil)
	if err != nil {
		t.Fatalf("Failed to create PoseidonSponge: %v", err)
	}

	inputs := []field.Element{
		field.New(1),
		field.New(2),
		field.New(3),
	}

	sponge.Absorb(inputs)
	outputs := sponge.Squeeze(2)

	if len(outputs) != 2 {
		t.Errorf("Expected 2 outputs, got %d", len(outputs))
	}

	for _, output := range outputs {
		if output.IsZero() {
			t.Error("Squeezed outputs should not all be zero")
		}
	}
}

func TestGrainLFSRInitialization(t *testing.T) {
	params := GetDefaultPoseidonParameters(128)
	lfsr := NewGrainLFSR(params)

	// Check that state is initialized
	hasTrue := false
	hasFalse := false
	for _, bit := range lfsr.state {
		if bit {
			hasTrue = true
		} else {
			hasFalse = true
		}
	}

	if !hasTrue || !hasFalse {
		t.Error("LFSR state should have both true and false bits")
	}
}

func TestGrainLFSRFieldElementGeneration(t *testing.T) {
	params := GetDefaultPoseidonParameters(128)
	lfsr := NewGrainLFSR(params)

	// Generate several field elements
	elements := make([]field.Element, 10)
	for i := 0; i < 10; i++ {
		elements[i] = lfsr.NextFieldElement()
	}

	// Check that elements are different (with high probability)
	allSame := true
	for i := 1; i < len(elements); i++ {
		if !elements[i].Equal(elements[0]) {
			allSame = false
			break
		}
	}

	if allSame {
		t.Error("LFSR should generate different field elements")
	}
}

func TestMDSMatrixGeneration(t *testing.T) {
	for width := 3; width <= 6; width++ {
		t.Run("width_"+string(rune('0'+width)), func(t *testing.T) {
			matrix, err := generatePoseidonMDSMatrix(width)
			if err != nil {
				t.Fatalf("Failed to generate MDS matrix: %v", err)
			}

			// Check dimensions
			if len(matrix) != width {
				t.Errorf("Expected %d rows, got %d", width, len(matrix))
			}

			for i, row := range matrix {
				if len(row) != width {
					t.Errorf("Row %d: expected %d columns, got %d", i, width, len(row))
				}

				// Check that no element is zero (Cauchy matrix property)
				for j, elem := range row {
					if elem.IsZero() {
						t.Errorf("Matrix[%d][%d] should not be zero", i, j)
					}
				}
			}
		})
	}
}

func TestRoundConstantsGeneration(t *testing.T) {
	params := GetDefaultPoseidonParameters(128)
	constants, err := generatePoseidonRoundConstants(params)
	if err != nil {
		t.Fatalf("Failed to generate round constants: %v", err)
	}

	totalRounds := params.RoundsFull + params.RoundsPartial
	if len(constants) != totalRounds {
		t.Errorf("Expected %d rounds of constants, got %d", totalRounds, len(constants))
	}

	for round, roundConsts := range constants {
		if len(roundConsts) != params.Width {
			t.Errorf("Round %d: expected %d constants, got %d",
				round, params.Width, len(roundConsts))
		}
	}
}

func TestPoseidonConvenienceFunctions(t *testing.T) {
	inputs := []field.Element{field.New(1), field.New(2), field.New(3)}

	// Test PoseidonHash convenience function
	result := PoseidonHash(inputs)
	if result.IsZero() {
		t.Error("PoseidonHash should not return zero for non-empty input")
	}

	// Test PoseidonHashTwo convenience function
	left := field.New(123)
	right := field.New(456)
	result2 := PoseidonHashTwo(left, right)
	if result2.IsZero() {
		t.Error("PoseidonHashTwo should not return zero")
	}
}

func TestPoseidonSboxOptimization(t *testing.T) {
	params := GetDefaultPoseidonParameters(128)
	poseidon, err := NewPoseidon(params)
	if err != nil {
		t.Fatalf("Failed to create Poseidon: %v", err)
	}

	x := field.New(7)

	// S-box should compute x^5
	result := poseidon.sbox(x)

	// Manually compute x^5
	x2 := x.Square()
	x4 := x2.Square()
	x5 := x.Mul(x4)

	if !result.Equal(x5) {
		t.Error("S-box should correctly compute x^5")
	}
}

func TestPoseidonLargeInput(t *testing.T) {
	poseidon, err := NewPoseidon(nil)
	if err != nil {
		t.Fatalf("Failed to create Poseidon: %v", err)
	}

	// Test with large input (more than rate)
	inputs := make([]field.Element, 100)
	for i := 0; i < 100; i++ {
		inputs[i] = field.New(uint64(i + 1))
	}

	result := poseidon.Hash(inputs)
	if result.IsZero() {
		t.Error("Hash of large input should not be zero")
	}
}

func TestPoseidonConsistencyAcrossInstances(t *testing.T) {
	// Two instances with same parameters should produce same output
	params := GetDefaultPoseidonParameters(128)

	poseidon1, err := NewPoseidon(params)
	if err != nil {
		t.Fatalf("Failed to create Poseidon 1: %v", err)
	}

	poseidon2, err := NewPoseidon(params)
	if err != nil {
		t.Fatalf("Failed to create Poseidon 2: %v", err)
	}

	inputs := []field.Element{field.New(1), field.New(2), field.New(3)}

	result1 := poseidon1.Hash(inputs)
	result2 := poseidon2.Hash(inputs)

	if !result1.Equal(result2) {
		t.Error("Same parameters should produce same hash")
	}
}
