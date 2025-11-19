package hash

import (
	"testing"

	"github.com/vybium/vybium-crypto/pkg/vybium-crypto/field"
)

func TestTip5BasicOperations(t *testing.T) {
	// Test basic Tip5 operations
	tip5 := New(VariableLength)

	// Test initial state
	for i := 0; i < StateSize; i++ {
		if !tip5.state[i].Equal(field.Zero) {
			t.Errorf("Initial state should be zeros, got %v", tip5.state[i])
		}
	}

	// Test fixed length domain
	tip5Fixed := New(FixedLength)
	for i := 0; i < Rate; i++ {
		if !tip5Fixed.state[i].Equal(field.Zero) {
			t.Errorf("Rate part should be zeros, got %v", tip5Fixed.state[i])
		}
	}
	for i := Rate; i < StateSize; i++ {
		if !tip5Fixed.state[i].Equal(field.One) {
			t.Errorf("Capacity part should be ones, got %v", tip5Fixed.state[i])
		}
	}
}

func TestTip5Hash10(t *testing.T) {
	// Test Hash10 function
	input := [Rate]field.Element{
		field.New(1), field.New(2), field.New(3), field.New(4), field.New(5),
		field.New(6), field.New(7), field.New(8), field.New(9), field.New(10),
	}

	digest := Hash10(input)

	// Check that digest has correct length
	if len(digest) != DigestLen {
		t.Errorf("Expected digest length %d, got %d", DigestLen, len(digest))
	}

	// Check that digest is not all zeros
	allZero := true
	for _, elem := range digest {
		if !elem.Equal(field.Zero) {
			allZero = false
			break
		}
	}
	if allZero {
		t.Error("Digest should not be all zeros")
	}
}

func TestTip5HashPair(t *testing.T) {
	// Test HashPair function
	left := [DigestLen]field.Element{
		field.New(1), field.New(2), field.New(3), field.New(4), field.New(5),
	}
	right := [DigestLen]field.Element{
		field.New(6), field.New(7), field.New(8), field.New(9), field.New(10),
	}

	digest := HashPair(left, right)

	// Check that digest has correct length
	if len(digest) != DigestLen {
		t.Errorf("Expected digest length %d, got %d", DigestLen, len(digest))
	}

	// Check that digest is not all zeros
	allZero := true
	for _, elem := range digest {
		if !elem.Equal(field.Zero) {
			allZero = false
			break
		}
	}
	if allZero {
		t.Error("Digest should not be all zeros")
	}
}

func TestTip5HashVarlen(t *testing.T) {
	// Test HashVarlen function
	input := []field.Element{
		field.New(1), field.New(2), field.New(3), field.New(4), field.New(5),
	}

	digest := HashVarlen(input)

	// Check that digest has correct length
	if len(digest) != DigestLen {
		t.Errorf("Expected digest length %d, got %d", DigestLen, len(digest))
	}

	// Check that digest is not all zeros
	allZero := true
	for _, elem := range digest {
		if !elem.Equal(field.Zero) {
			allZero = false
			break
		}
	}
	if allZero {
		t.Error("Digest should not be all zeros")
	}
}

func TestTip5Permutation(t *testing.T) {
	// Test permutation function
	tip5 := New(VariableLength)

	// Set some initial state
	for i := 0; i < StateSize; i++ {
		tip5.state[i] = field.New(uint64(i + 1))
	}

	// Store initial state
	initialState := tip5.state

	// Apply permutation
	tip5.Permutation()

	// Check that state changed
	changed := false
	for i := 0; i < StateSize; i++ {
		if !tip5.state[i].Equal(initialState[i]) {
			changed = true
			break
		}
	}
	if !changed {
		t.Error("Permutation should change the state")
	}
}

func TestTip5SboxLayer(t *testing.T) {
	// Test S-box layer
	tip5 := New(VariableLength)

	// Set some initial state
	for i := 0; i < StateSize; i++ {
		tip5.state[i] = field.New(uint64(i + 1))
	}

	// Store initial state
	initialState := tip5.state

	// Apply S-box layer
	tip5.sboxLayer()

	// Check that state changed
	changed := false
	for i := 0; i < StateSize; i++ {
		if !tip5.state[i].Equal(initialState[i]) {
			changed = true
			break
		}
	}
	if !changed {
		t.Error("S-box layer should change the state")
	}
}

func TestTip5SplitAndLookup(t *testing.T) {
	// Test split-and-lookup function
	// Test with a known value
	element := field.New(0x123456789ABCDEF0)
	original := element

	splitAndLookup(&element)

	// Check that element changed
	if element.Equal(original) {
		t.Error("Split-and-lookup should change the element")
	}
}

func TestTip5MdsLayer(t *testing.T) {
	// Test MDS layer
	tip5 := New(VariableLength)

	// Set some initial state
	for i := 0; i < StateSize; i++ {
		tip5.state[i] = field.New(uint64(i + 1))
	}

	// Store initial state
	initialState := tip5.state

	// Apply MDS layer
	tip5.mdsGenerated()

	// Check that state changed
	changed := false
	for i := 0; i < StateSize; i++ {
		if !tip5.state[i].Equal(initialState[i]) {
			changed = true
			break
		}
	}
	if !changed {
		t.Error("MDS layer should change the state")
	}
}

func TestTip5Consistency(t *testing.T) {
	// Test that the same input produces the same output
	input := [Rate]field.Element{
		field.New(1), field.New(2), field.New(3), field.New(4), field.New(5),
		field.New(6), field.New(7), field.New(8), field.New(9), field.New(10),
	}

	digest1 := Hash10(input)
	digest2 := Hash10(input)

	// Check that both digests are equal
	for i := 0; i < DigestLen; i++ {
		if !digest1[i].Equal(digest2[i]) {
			t.Errorf("Hash10 should be deterministic: digest1[%d] = %v, digest2[%d] = %v",
				i, digest1[i], i, digest2[i])
		}
	}
}

func TestTip5DifferentInputs(t *testing.T) {
	// Test that different inputs produce different outputs
	input1 := [Rate]field.Element{
		field.New(1), field.New(2), field.New(3), field.New(4), field.New(5),
		field.New(6), field.New(7), field.New(8), field.New(9), field.New(10),
	}
	input2 := [Rate]field.Element{
		field.New(2), field.New(3), field.New(4), field.New(5), field.New(6),
		field.New(7), field.New(8), field.New(9), field.New(10), field.New(11),
	}

	digest1 := Hash10(input1)
	digest2 := Hash10(input2)

	// Check that digests are different
	equal := true
	for i := 0; i < DigestLen; i++ {
		if !digest1[i].Equal(digest2[i]) {
			equal = false
			break
		}
	}
	if equal {
		t.Error("Different inputs should produce different outputs")
	}
}

// Benchmark tests
func BenchmarkTip5Hash10(b *testing.B) {
	input := [Rate]field.Element{
		field.New(1), field.New(2), field.New(3), field.New(4), field.New(5),
		field.New(6), field.New(7), field.New(8), field.New(9), field.New(10),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Hash10(input)
	}
}

func BenchmarkTip5HashPair(b *testing.B) {
	left := [DigestLen]field.Element{
		field.New(1), field.New(2), field.New(3), field.New(4), field.New(5),
	}
	right := [DigestLen]field.Element{
		field.New(6), field.New(7), field.New(8), field.New(9), field.New(10),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = HashPair(left, right)
	}
}

func BenchmarkTip5HashVarlen(b *testing.B) {
	input := []field.Element{
		field.New(1), field.New(2), field.New(3), field.New(4), field.New(5),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = HashVarlen(input)
	}
}

func BenchmarkTip5Permutation(b *testing.B) {
	tip5 := New(VariableLength)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tip5.Permutation()
	}
}
