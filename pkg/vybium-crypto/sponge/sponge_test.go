package sponge

import (
	"testing"

	"github.com/vybium/vybium-crypto/pkg/vybium-crypto/field"
)

func TestTip5SpongeInit(t *testing.T) {
	sponge := NewTip5Sponge(VariableLength)

	// Test that Init returns a new instance
	newSponge := sponge.Init()
	if newSponge == nil {
		t.Error("Init() returned nil")
	}

	// Test that it's a different instance
	if newSponge == sponge {
		t.Error("Init() returned the same instance")
	}
}

func TestTip5SpongeAbsorb(t *testing.T) {
	sponge := NewTip5Sponge(VariableLength)

	// Create test input
	input := [Rate]field.Element{}
	for i := 0; i < Rate; i++ {
		input[i] = field.New(uint64(i + 1))
	}

	// Test absorb
	sponge.Absorb(input)

	// Verify state changed (we can't easily test the exact state
	// without knowing the Tip5 permutation implementation)
	// For now, just ensure no panic occurs
}

func TestTip5SpongeSqueeze(t *testing.T) {
	sponge := NewTip5Sponge(VariableLength)

	// Absorb some data first
	input := [Rate]field.Element{}
	for i := 0; i < Rate; i++ {
		input[i] = field.New(uint64(i + 1))
	}
	sponge.Absorb(input)

	// Test squeeze
	output := sponge.Squeeze()

	// Verify output has correct length
	if len(output) != Rate {
		t.Errorf("Squeeze() returned %d elements, want %d", len(output), Rate)
	}

	// Verify all elements are valid field elements
	for i, element := range output {
		if element.Value() >= field.P {
			t.Errorf("Squeeze()[%d] = %d, invalid field element", i, element.Value())
		}
	}
}

func TestTip5SpongePadAndAbsorbAll(t *testing.T) {
	tests := []struct {
		name  string
		input []field.Element
	}{
		{
			name:  "Empty input",
			input: []field.Element{},
		},
		{
			name:  "Single element",
			input: []field.Element{field.One},
		},
		{
			name:  "Exact rate",
			input: []field.Element{field.New(1), field.New(2), field.New(3), field.New(4), field.New(5), field.New(6), field.New(7), field.New(8), field.New(9), field.New(10)},
		},
		{
			name:  "Multiple chunks",
			input: []field.Element{field.New(1), field.New(2), field.New(3), field.New(4), field.New(5), field.New(6), field.New(7), field.New(8), field.New(9), field.New(10), field.New(11), field.New(12)},
		},
		{
			name:  "Large input",
			input: make([]field.Element, Rate*3+5),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sponge := NewTip5Sponge(VariableLength)

			// Initialize test data
			for i := range tt.input {
				tt.input[i] = field.New(uint64(i + 1))
			}

			// Test PadAndAbsorbAll
			sponge.PadAndAbsorbAll(tt.input)

			// Verify no panic occurred
			// The actual state verification would require knowing the Tip5 permutation
		})
	}
}

func TestTip5SpongeClone(t *testing.T) {
	sponge := NewTip5Sponge(VariableLength)

	// Absorb some data
	input := [Rate]field.Element{}
	for i := 0; i < Rate; i++ {
		input[i] = field.New(uint64(i + 1))
	}
	sponge.Absorb(input)

	// Clone the sponge
	clone := sponge.Clone()

	// Verify it's a different instance
	if clone == sponge {
		t.Error("Clone() returned the same instance")
	}

	// Verify it's a Tip5Sponge
	if _, ok := clone.(*Tip5Sponge); !ok {
		t.Error("Clone() returned wrong type")
	}
}

func TestTip5SpongeReset(t *testing.T) {
	sponge := NewTip5Sponge(VariableLength)

	// Absorb some data
	input := [Rate]field.Element{}
	for i := 0; i < Rate; i++ {
		input[i] = field.New(uint64(i + 1))
	}
	sponge.Absorb(input)

	// Reset the sponge
	sponge.Reset()

	// Verify state is reset (all zeros)
	// This is a simplified test - in practice, we'd need to expose the state
	// or test through the behavior
}

func TestPoseidonSpongeBasicOperations(t *testing.T) {
	sponge := NewPoseidonSponge(FixedLength)

	// Test Init
	newSponge := sponge.Init()
	if newSponge == nil {
		t.Error("Init() returned nil")
	}

	// Test Absorb
	input := [Rate]field.Element{}
	for i := 0; i < Rate; i++ {
		input[i] = field.New(uint64(i + 1))
	}
	sponge.Absorb(input)

	// Test Squeeze
	_ = sponge.Squeeze()

	// Test PadAndAbsorbAll
	testInput := []field.Element{field.One, field.New(2), field.New(3)}
	sponge.PadAndAbsorbAll(testInput)

	// Test Clone
	clone := sponge.Clone()
	if clone == sponge {
		t.Error("Clone() returned the same instance")
	}

	// Test Reset
	sponge.Reset()
}

func TestHashVarlen(t *testing.T) {
	tests := []struct {
		name  string
		input []field.Element
	}{
		{
			name:  "Empty input",
			input: []field.Element{},
		},
		{
			name:  "Single element",
			input: []field.Element{field.One},
		},
		{
			name:  "Multiple elements",
			input: []field.Element{field.One, field.New(2), field.New(3), field.New(4), field.New(5)},
		},
		{
			name:  "Large input",
			input: make([]field.Element, 100),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sponge := NewTip5Sponge(VariableLength)

			// Initialize test data
			for i := range tt.input {
				tt.input[i] = field.New(uint64(i + 1))
			}

			// Test HashVarlen
			output := HashVarlen(sponge, tt.input)

			// Verify output has correct length
			if len(output) != Rate {
				t.Errorf("HashVarlen() returned %d elements, want %d", len(output), Rate)
			}

			// Verify all elements are valid field elements
			for i, element := range output {
				if element.Value() >= field.P {
					t.Errorf("HashVarlen()[%d] = %d, invalid field element", i, element.Value())
				}
			}
		})
	}
}

func TestHashFixed(t *testing.T) {
	tests := []struct {
		name  string
		input []field.Element
	}{
		{
			name:  "Single element",
			input: []field.Element{field.One},
		},
		{
			name:  "Multiple elements",
			input: []field.Element{field.One, field.New(2), field.New(3)},
		},
		{
			name:  "Exact rate",
			input: make([]field.Element, Rate),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sponge := NewTip5Sponge(FixedLength)

			// Initialize test data
			for i := range tt.input {
				tt.input[i] = field.New(uint64(i + 1))
			}

			// Test HashFixed
			output := HashFixed(sponge, tt.input)

			// Verify output has correct length
			if len(output) != Rate {
				t.Errorf("HashFixed() returned %d elements, want %d", len(output), Rate)
			}

			// Verify all elements are valid field elements
			for i, element := range output {
				if element.Value() >= field.P {
					t.Errorf("HashFixed()[%d] = %d, invalid field element", i, element.Value())
				}
			}
		})
	}
}

func TestHashFixedPanic(t *testing.T) {
	sponge := NewTip5Sponge(FixedLength)

	// Test that HashFixed panics with input longer than RATE
	input := make([]field.Element, Rate+1)
	for i := range input {
		input[i] = field.New(uint64(i + 1))
	}

	defer func() {
		if r := recover(); r == nil {
			t.Error("HashFixed() should panic with input longer than RATE")
		}
	}()

	HashFixed(sponge, input)
}

func TestSampleIndices(t *testing.T) {
	tests := []struct {
		name        string
		upperBound  int
		numIndices  int
		expectEmpty bool
	}{
		{
			name:        "Zero upper bound",
			upperBound:  0,
			numIndices:  5,
			expectEmpty: true,
		},
		{
			name:        "Zero num indices",
			upperBound:  10,
			numIndices:  0,
			expectEmpty: true,
		},
		{
			name:        "Normal case",
			upperBound:  10,
			numIndices:  5,
			expectEmpty: false,
		},
		{
			name:        "All indices",
			upperBound:  5,
			numIndices:  5,
			expectEmpty: false,
		},
		{
			name:        "More indices than upper bound",
			upperBound:  3,
			numIndices:  10,
			expectEmpty: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sponge := NewTip5Sponge(VariableLength)

			indices := SampleIndices(sponge, tt.upperBound, tt.numIndices)

			if tt.expectEmpty {
				if len(indices) != 0 {
					t.Errorf("SampleIndices() returned %d indices, want 0", len(indices))
				}
				return
			}

			// Verify correct number of indices
			expectedCount := tt.numIndices
			if tt.upperBound < tt.numIndices {
				expectedCount = tt.upperBound
			}

			if len(indices) != expectedCount {
				t.Errorf("SampleIndices() returned %d indices, want %d", len(indices), expectedCount)
			}

			// Verify all indices are in range
			for _, index := range indices {
				if index < 0 || index >= tt.upperBound {
					t.Errorf("SampleIndices() returned index %d, out of range [0, %d)", index, tt.upperBound)
				}
			}

			// Verify no duplicates
			seen := make(map[int]bool)
			for _, index := range indices {
				if seen[index] {
					t.Errorf("SampleIndices() returned duplicate index %d", index)
				}
				seen[index] = true
			}
		})
	}
}

func TestValidateSpongeInput(t *testing.T) {
	tests := []struct {
		name    string
		input   []field.Element
		wantErr bool
	}{
		{
			name:    "Empty input",
			input:   []field.Element{},
			wantErr: true,
		},
		{
			name:    "Valid input",
			input:   []field.Element{field.One, field.New(2)},
			wantErr: false,
		},
		{
			name:    "Large input",
			input:   make([]field.Element, 1024*1024+1),
			wantErr: true,
		},
		{
			name:    "Maximum valid input",
			input:   make([]field.Element, 1024*1024),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSpongeInput(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateSpongeInput() expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("ValidateSpongeInput() error = %v", err)
				}
			}
		})
	}
}

func TestGetSpongeRate(t *testing.T) {
	rate := GetSpongeRate()
	if rate != Rate {
		t.Errorf("GetSpongeRate() = %d, want %d", rate, Rate)
	}
}

func TestIsValidDomain(t *testing.T) {
	tests := []struct {
		name   string
		domain Domain
		want   bool
	}{
		{
			name:   "VariableLength",
			domain: VariableLength,
			want:   true,
		},
		{
			name:   "FixedLength",
			domain: FixedLength,
			want:   true,
		},
		{
			name:   "Invalid domain",
			domain: Domain(999),
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsValidDomain(tt.domain)
			if got != tt.want {
				t.Errorf("IsValidDomain() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDomainString(t *testing.T) {
	tests := []struct {
		name     string
		domain   Domain
		expected string
	}{
		{
			name:     "VariableLength",
			domain:   VariableLength,
			expected: "VariableLength",
		},
		{
			name:     "FixedLength",
			domain:   FixedLength,
			expected: "FixedLength",
		},
		{
			name:     "Invalid domain",
			domain:   Domain(999),
			expected: "Unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.domain.String()
			if got != tt.expected {
				t.Errorf("Domain.String() = %s, want %s", got, tt.expected)
			}
		})
	}
}

// Benchmark tests
func BenchmarkTip5SpongeAbsorb(b *testing.B) {
	sponge := NewTip5Sponge(VariableLength)
	input := [Rate]field.Element{}
	for i := 0; i < Rate; i++ {
		input[i] = field.New(uint64(i + 1))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sponge.Absorb(input)
	}
}

func BenchmarkTip5SpongeSqueeze(b *testing.B) {
	sponge := NewTip5Sponge(VariableLength)
	input := [Rate]field.Element{}
	for i := 0; i < Rate; i++ {
		input[i] = field.New(uint64(i + 1))
	}
	sponge.Absorb(input)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = sponge.Squeeze()
	}
}

func BenchmarkHashVarlen(b *testing.B) {
	sponge := NewTip5Sponge(VariableLength)
	input := make([]field.Element, 100)
	for i := range input {
		input[i] = field.New(uint64(i + 1))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = HashVarlen(sponge, input)
	}
}

func BenchmarkHashFixed(b *testing.B) {
	sponge := NewTip5Sponge(FixedLength)
	input := make([]field.Element, Rate)
	for i := range input {
		input[i] = field.New(uint64(i + 1))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = HashFixed(sponge, input)
	}
}

func BenchmarkSampleIndices(b *testing.B) {
	sponge := NewTip5Sponge(VariableLength)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = SampleIndices(sponge, 1000, 10)
	}
}
