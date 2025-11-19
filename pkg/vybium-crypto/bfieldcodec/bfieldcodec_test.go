package bfieldcodec

import (
	"math/big"
	"testing"

	"github.com/vybium/vybium-crypto/pkg/vybium-crypto/field"
	"github.com/vybium/vybium-crypto/pkg/vybium-crypto/xfield"
)

func TestEncodeDecodeBFieldElement(t *testing.T) {
	tests := []struct {
		name     string
		value    field.Element
		expected []field.Element
	}{
		{
			name:     "Zero",
			value:    field.Zero,
			expected: []field.Element{field.Zero},
		},
		{
			name:     "One",
			value:    field.One,
			expected: []field.Element{field.One},
		},
		{
			name:     "Large value",
			value:    field.New(12345),
			expected: []field.Element{field.New(12345)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test encoding
			encoded := EncodeBFieldElement(tt.value)
			if len(encoded) != 1 {
				t.Errorf("EncodeBFieldElement() length = %d, want 1", len(encoded))
			}
			if !encoded[0].Equal(tt.expected[0]) {
				t.Errorf("EncodeBFieldElement() = %v, want %v", encoded[0], tt.expected[0])
			}

			// Test decoding
			decoded, err := DecodeBFieldElement(encoded)
			if err != nil {
				t.Errorf("DecodeBFieldElement() error = %v", err)
			}
			if !decoded.Equal(tt.value) {
				t.Errorf("DecodeBFieldElement() = %v, want %v", decoded, tt.value)
			}
		})
	}
}

func TestDecodeBFieldElementErrors(t *testing.T) {
	tests := []struct {
		name      string
		sequence  []field.Element
		wantError ErrorType
	}{
		{
			name:      "Empty sequence",
			sequence:  []field.Element{},
			wantError: ErrorEmptySequence,
		},
		{
			name:      "Too long sequence",
			sequence:  []field.Element{field.One, field.New(2)},
			wantError: ErrorSequenceTooLong,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := DecodeBFieldElement(tt.sequence)
			if err == nil {
				t.Errorf("DecodeBFieldElement() expected error, got nil")
				return
			}

			bfcErr, ok := err.(BFieldCodecError)
			if !ok {
				t.Errorf("DecodeBFieldElement() error type = %T, want BFieldCodecError", err)
				return
			}

			if bfcErr.Type != tt.wantError {
				t.Errorf("DecodeBFieldElement() error type = %v, want %v", bfcErr.Type, tt.wantError)
			}
		})
	}
}

func TestEncodeDecodeUint64(t *testing.T) {
	tests := []struct {
		name     string
		value    uint64
		expected []field.Element
	}{
		{
			name:     "Zero",
			value:    0,
			expected: []field.Element{field.Zero, field.Zero},
		},
		{
			name:     "Small value",
			value:    12345,
			expected: []field.Element{field.New(12345), field.Zero},
		},
		{
			name:     "Large value",
			value:    0x123456789ABCDEF0,
			expected: []field.Element{field.New(0x9ABCDEF0), field.New(0x12345678)},
		},
		{
			name:     "Max uint64",
			value:    0xFFFFFFFFFFFFFFFF,
			expected: []field.Element{field.New(0xFFFFFFFF), field.New(0xFFFFFFFF)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test encoding
			encoded := EncodeUint64(tt.value)
			if len(encoded) != 2 {
				t.Errorf("EncodeUint64() length = %d, want 2", len(encoded))
			}
			for i, expected := range tt.expected {
				if !encoded[i].Equal(expected) {
					t.Errorf("EncodeUint64()[%d] = %v, want %v", i, encoded[i], expected)
				}
			}

			// Test decoding
			decoded, err := DecodeUint64(encoded)
			if err != nil {
				t.Errorf("DecodeUint64() error = %v", err)
			}
			if decoded != tt.value {
				t.Errorf("DecodeUint64() = %d, want %d", decoded, tt.value)
			}
		})
	}
}

func TestDecodeUint64Errors(t *testing.T) {
	tests := []struct {
		name      string
		sequence  []field.Element
		wantError ErrorType
	}{
		{
			name:      "Too short",
			sequence:  []field.Element{field.One},
			wantError: ErrorSequenceTooShort,
		},
		{
			name:      "Too long",
			sequence:  []field.Element{field.One, field.New(2), field.New(3)},
			wantError: ErrorSequenceTooLong,
		},
		{
			name:      "Out of range",
			sequence:  []field.Element{field.New(0x100000000), field.Zero},
			wantError: ErrorElementOutOfRange,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := DecodeUint64(tt.sequence)
			if err == nil {
				t.Errorf("DecodeUint64() expected error, got nil")
				return
			}

			bfcErr, ok := err.(BFieldCodecError)
			if !ok {
				t.Errorf("DecodeUint64() error type = %T, want BFieldCodecError", err)
				return
			}

			if bfcErr.Type != tt.wantError {
				t.Errorf("DecodeUint64() error type = %v, want %v", bfcErr.Type, tt.wantError)
			}
		})
	}
}

func TestEncodeDecodeUint32(t *testing.T) {
	tests := []struct {
		name     string
		value    uint32
		expected field.Element
	}{
		{
			name:     "Zero",
			value:    0,
			expected: field.Zero,
		},
		{
			name:     "Small value",
			value:    12345,
			expected: field.New(12345),
		},
		{
			name:     "Max uint32",
			value:    0xFFFFFFFF,
			expected: field.New(0xFFFFFFFF),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test encoding
			encoded := EncodeUint32(tt.value)
			if len(encoded) != 1 {
				t.Errorf("EncodeUint32() length = %d, want 1", len(encoded))
			}
			if !encoded[0].Equal(tt.expected) {
				t.Errorf("EncodeUint32() = %v, want %v", encoded[0], tt.expected)
			}

			// Test decoding
			decoded, err := DecodeUint32(encoded)
			if err != nil {
				t.Errorf("DecodeUint32() error = %v", err)
			}
			if decoded != tt.value {
				t.Errorf("DecodeUint32() = %d, want %d", decoded, tt.value)
			}
		})
	}
}

func TestEncodeDecodeBool(t *testing.T) {
	tests := []struct {
		name     string
		value    bool
		expected field.Element
	}{
		{
			name:     "False",
			value:    false,
			expected: field.Zero,
		},
		{
			name:     "True",
			value:    true,
			expected: field.One,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test encoding
			encoded := EncodeBool(tt.value)
			if len(encoded) != 1 {
				t.Errorf("EncodeBool() length = %d, want 1", len(encoded))
			}
			if !encoded[0].Equal(tt.expected) {
				t.Errorf("EncodeBool() = %v, want %v", encoded[0], tt.expected)
			}

			// Test decoding
			decoded, err := DecodeBool(encoded)
			if err != nil {
				t.Errorf("DecodeBool() error = %v", err)
			}
			if decoded != tt.value {
				t.Errorf("DecodeBool() = %v, want %v", decoded, tt.value)
			}
		})
	}
}

func TestDecodeBoolErrors(t *testing.T) {
	tests := []struct {
		name      string
		sequence  []field.Element
		wantError ErrorType
	}{
		{
			name:      "Empty sequence",
			sequence:  []field.Element{},
			wantError: ErrorEmptySequence,
		},
		{
			name:      "Too long",
			sequence:  []field.Element{field.One, field.New(2)},
			wantError: ErrorSequenceTooLong,
		},
		{
			name:      "Out of range",
			sequence:  []field.Element{field.New(2)},
			wantError: ErrorElementOutOfRange,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := DecodeBool(tt.sequence)
			if err == nil {
				t.Errorf("DecodeBool() expected error, got nil")
				return
			}

			bfcErr, ok := err.(BFieldCodecError)
			if !ok {
				t.Errorf("DecodeBool() error type = %T, want BFieldCodecError", err)
				return
			}

			if bfcErr.Type != tt.wantError {
				t.Errorf("DecodeBool() error type = %v, want %v", bfcErr.Type, tt.wantError)
			}
		})
	}
}

func TestEncodeDecodeXFieldElement(t *testing.T) {
	tests := []struct {
		name     string
		value    xfield.XFieldElement
		expected []field.Element
	}{
		{
			name:     "Zero",
			value:    xfield.Zero,
			expected: []field.Element{field.Zero, field.Zero, field.Zero},
		},
		{
			name:     "One",
			value:    xfield.One,
			expected: []field.Element{field.One, field.Zero, field.Zero},
		},
		{
			name:     "General element",
			value:    xfield.New([3]field.Element{field.New(1), field.New(2), field.New(3)}),
			expected: []field.Element{field.New(1), field.New(2), field.New(3)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test encoding
			encoded := EncodeXFieldElement(tt.value)
			if len(encoded) != 3 {
				t.Errorf("EncodeXFieldElement() length = %d, want 3", len(encoded))
			}
			for i, expected := range tt.expected {
				if !encoded[i].Equal(expected) {
					t.Errorf("EncodeXFieldElement()[%d] = %v, want %v", i, encoded[i], expected)
				}
			}

			// Test decoding
			decoded, err := DecodeXFieldElement(encoded)
			if err != nil {
				t.Errorf("DecodeXFieldElement() error = %v", err)
			}
			if !decoded.Equal(tt.value) {
				t.Errorf("DecodeXFieldElement() = %v, want %v", decoded, tt.value)
			}
		})
	}
}

func TestDecodeXFieldElementErrors(t *testing.T) {
	tests := []struct {
		name      string
		sequence  []field.Element
		wantError ErrorType
	}{
		{
			name:      "Too short",
			sequence:  []field.Element{field.One, field.New(2)},
			wantError: ErrorSequenceTooShort,
		},
		{
			name:      "Too long",
			sequence:  []field.Element{field.One, field.New(2), field.New(3), field.New(4)},
			wantError: ErrorSequenceTooLong,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := DecodeXFieldElement(tt.sequence)
			if err == nil {
				t.Errorf("DecodeXFieldElement() expected error, got nil")
				return
			}

			bfcErr, ok := err.(BFieldCodecError)
			if !ok {
				t.Errorf("DecodeXFieldElement() error type = %T, want BFieldCodecError", err)
				return
			}

			if bfcErr.Type != tt.wantError {
				t.Errorf("DecodeXFieldElement() error type = %v, want %v", bfcErr.Type, tt.wantError)
			}
		})
	}
}

func TestEncodeLengthPrefix(t *testing.T) {
	tests := []struct {
		name     string
		sequence []field.Element
		expected []field.Element
	}{
		{
			name:     "Empty sequence",
			sequence: []field.Element{},
			expected: []field.Element{field.Zero},
		},
		{
			name:     "Single element",
			sequence: []field.Element{field.One},
			expected: []field.Element{field.One, field.One},
		},
		{
			name:     "Multiple elements",
			sequence: []field.Element{field.One, field.New(2), field.New(3)},
			expected: []field.Element{field.New(3), field.One, field.New(2), field.New(3)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded := EncodeLengthPrefix(tt.sequence)
			if len(encoded) != len(tt.expected) {
				t.Errorf("EncodeLengthPrefix() length = %d, want %d", len(encoded), len(tt.expected))
			}
			for i, expected := range tt.expected {
				if !encoded[i].Equal(expected) {
					t.Errorf("EncodeLengthPrefix()[%d] = %v, want %v", i, encoded[i], expected)
				}
			}
		})
	}
}

func TestDecodeLengthPrefix(t *testing.T) {
	tests := []struct {
		name          string
		sequence      []field.Element
		wantLength    int
		wantRemaining []field.Element
		wantError     bool
	}{
		{
			name:      "Empty sequence",
			sequence:  []field.Element{},
			wantError: true,
		},
		{
			name:          "Zero length",
			sequence:      []field.Element{field.Zero},
			wantLength:    0,
			wantRemaining: []field.Element{},
			wantError:     false,
		},
		{
			name:          "Valid sequence",
			sequence:      []field.Element{field.New(3), field.One, field.New(2), field.New(3)},
			wantLength:    3,
			wantRemaining: []field.Element{field.One, field.New(2), field.New(3)},
			wantError:     false,
		},
		{
			name:      "Too short",
			sequence:  []field.Element{field.New(3), field.One},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			length, remaining, err := DecodeLengthPrefix(tt.sequence)
			if tt.wantError {
				if err == nil {
					t.Errorf("DecodeLengthPrefix() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("DecodeLengthPrefix() error = %v", err)
				return
			}

			if length != tt.wantLength {
				t.Errorf("DecodeLengthPrefix() length = %d, want %d", length, tt.wantLength)
			}

			if len(remaining) != len(tt.wantRemaining) {
				t.Errorf("DecodeLengthPrefix() remaining length = %d, want %d", len(remaining), len(tt.wantRemaining))
			}

			for i, expected := range tt.wantRemaining {
				if !remaining[i].Equal(expected) {
					t.Errorf("DecodeLengthPrefix() remaining[%d] = %v, want %v", i, remaining[i], expected)
				}
			}
		})
	}
}

func TestValidateSequenceLength(t *testing.T) {
	tests := []struct {
		name      string
		sequence  []field.Element
		expected  int
		wantError bool
	}{
		{
			name:      "Correct length",
			sequence:  []field.Element{field.One, field.New(2), field.New(3)},
			expected:  3,
			wantError: false,
		},
		{
			name:      "Too short",
			sequence:  []field.Element{field.One, field.New(2)},
			expected:  3,
			wantError: true,
		},
		{
			name:      "Too long",
			sequence:  []field.Element{field.One, field.New(2), field.New(3), field.New(4)},
			expected:  3,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSequenceLength(tt.sequence, tt.expected)
			if tt.wantError {
				if err == nil {
					t.Errorf("ValidateSequenceLength() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("ValidateSequenceLength() error = %v", err)
			}
		})
	}
}

func TestConvertToFromBigInt(t *testing.T) {
	tests := []struct {
		name     string
		value    field.Element
		bigValue *big.Int
	}{
		{
			name:     "Zero",
			value:    field.Zero,
			bigValue: big.NewInt(0),
		},
		{
			name:     "One",
			value:    field.One,
			bigValue: big.NewInt(1),
		},
		{
			name:     "Large value",
			value:    field.New(12345),
			bigValue: big.NewInt(12345),
		},
		{
			name:     "Very large value",
			value:    field.New(0x123456789ABCDEF0),
			bigValue: big.NewInt(0x123456789ABCDEF0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test conversion to big.Int
			converted := ConvertToBigInt(tt.value)
			if converted.Cmp(tt.bigValue) != 0 {
				t.Errorf("ConvertToBigInt() = %v, want %v", converted, tt.bigValue)
			}

			// Test conversion from big.Int
			back := ConvertFromBigInt(tt.bigValue)
			if !back.Equal(tt.value) {
				t.Errorf("ConvertFromBigInt() = %v, want %v", back, tt.value)
			}
		})
	}
}

func TestConvertFromBigIntModulo(t *testing.T) {
	// Test that values larger than the field prime are reduced modulo P
	largeValue := new(big.Int)
	largeValue.SetString("18446744069414584322", 10) // P + 1

	converted := ConvertFromBigInt(largeValue)
	expected := field.One // (P + 1) mod P = 1

	if !converted.Equal(expected) {
		t.Errorf("ConvertFromBigInt(large) = %v, want %v", converted, expected)
	}
}

// Benchmark tests
func BenchmarkEncodeBFieldElement(b *testing.B) {
	element := field.New(12345)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = EncodeBFieldElement(element)
	}
}

func BenchmarkDecodeBFieldElement(b *testing.B) {
	sequence := []field.Element{field.New(12345)}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = DecodeBFieldElement(sequence)
	}
}

func BenchmarkEncodeUint64(b *testing.B) {
	value := uint64(0x123456789ABCDEF0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = EncodeUint64(value)
	}
}

func BenchmarkDecodeUint64(b *testing.B) {
	sequence := []field.Element{field.New(0x9ABCDEF0), field.New(0x12345678)}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = DecodeUint64(sequence)
	}
}

func BenchmarkEncodeXFieldElement(b *testing.B) {
	element := xfield.New([3]field.Element{field.New(1), field.New(2), field.New(3)})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = EncodeXFieldElement(element)
	}
}

func BenchmarkDecodeXFieldElement(b *testing.B) {
	sequence := []field.Element{field.New(1), field.New(2), field.New(3)}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = DecodeXFieldElement(sequence)
	}
}
