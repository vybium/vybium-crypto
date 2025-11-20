// Package bfieldcodec provides serialization and deserialization for BFieldElement sequences.
//
// BFieldCodec is a serialization format that encodes data as sequences of BFieldElement values.
// It's used for STARK proofs and other cryptographic data that needs to be serialized
// in a format compatible with zero-knowledge proof systems.
//
// Key Features:
// - Length-prefixed encoding for dynamic-size types
// - Static-length encoding for fixed-size types
// - Support for nested structures and collections
// - Error handling for malformed sequences
package bfieldcodec

import (
	"fmt"
	"math/big"

	"github.com/vybium/vybium-crypto/pkg/vybium-crypto/field"
	"github.com/vybium/vybium-crypto/pkg/vybium-crypto/xfield"
)

// BFieldCodecError represents errors that can occur during encoding/decoding.
type BFieldCodecError struct {
	Type    ErrorType
	Message string
}

type ErrorType int

const (
	ErrorEmptySequence ErrorType = iota
	ErrorSequenceTooShort
	ErrorSequenceTooLong
	ErrorElementOutOfRange
	ErrorMissingLengthIndicator
	ErrorInvalidLengthIndicator
	ErrorInnerDecodingFailure
	ErrorUnsupportedType
)

func (e BFieldCodecError) Error() string {
	return fmt.Sprintf("BFieldCodec error [%d]: %s", e.Type, e.Message)
}

// BFieldCodec defines the interface for types that can be encoded/decoded to/from BFieldElement sequences.
// This mirrors twenty-first's BFieldCodec trait from bfield_codec.rs
//
// Note: Go's type system doesn't support the same level of static dispatch as Rust's trait system,
// so we use interface methods and runtime type checking where necessary.
type BFieldCodec interface {
	// Encode converts the value to a sequence of BFieldElement values.
	Encode() []field.Element

	// Decode creates a value from a sequence of BFieldElement values.
	// Returns the decoded value (as BFieldCodec interface) and an error if the sequence is malformed.
	// Callers should type-assert the result to the expected concrete type.
	Decode(sequence []field.Element) (BFieldCodec, error)

	// StaticLength returns the fixed length in BFieldElements if known at compile time.
	// Returns nil for dynamic-length types (those with length prefixes).
	StaticLength() *int
}

// EncodeBFieldElement encodes a single BFieldElement.
func EncodeBFieldElement(element field.Element) []field.Element {
	return []field.Element{element}
}

// DecodeBFieldElement decodes a single BFieldElement from a sequence.
func DecodeBFieldElement(sequence []field.Element) (field.Element, error) {
	if len(sequence) == 0 {
		return field.Zero, BFieldCodecError{ErrorEmptySequence, "empty sequence"}
	}
	if len(sequence) > 1 {
		return field.Zero, BFieldCodecError{ErrorSequenceTooLong, "sequence too long for single element"}
	}
	return sequence[0], nil
}

// EncodeUint64 encodes a uint64 as two BFieldElement values (32 bits each).
func EncodeUint64(value uint64) []field.Element {
	// Split into two 32-bit parts
	low := uint32(value & 0xFFFFFFFF)
	high := uint32((value >> 32) & 0xFFFFFFFF)
	return []field.Element{
		field.New(uint64(low)),
		field.New(uint64(high)),
	}
}

// DecodeUint64 decodes a uint64 from two BFieldElement values.
func DecodeUint64(sequence []field.Element) (uint64, error) {
	if len(sequence) < 2 {
		return 0, BFieldCodecError{ErrorSequenceTooShort, "need at least 2 elements for uint64"}
	}
	if len(sequence) > 2 {
		return 0, BFieldCodecError{ErrorSequenceTooLong, "too many elements for uint64"}
	}

	low := sequence[0].Value()
	high := sequence[1].Value()

	// Validate that values fit in 32 bits
	if low > 0xFFFFFFFF || high > 0xFFFFFFFF {
		return 0, BFieldCodecError{ErrorElementOutOfRange, "element out of range for uint64"}
	}

	return (uint64(high) << 32) | uint64(low), nil
}

// EncodeUint32 encodes a uint32 as a single BFieldElement.
func EncodeUint32(value uint32) []field.Element {
	return []field.Element{field.New(uint64(value))}
}

// DecodeUint32 decodes a uint32 from a single BFieldElement.
func DecodeUint32(sequence []field.Element) (uint32, error) {
	if len(sequence) == 0 {
		return 0, BFieldCodecError{ErrorEmptySequence, "empty sequence"}
	}
	if len(sequence) > 1 {
		return 0, BFieldCodecError{ErrorSequenceTooLong, "sequence too long for uint32"}
	}

	value := sequence[0].Value()
	if value > 0xFFFFFFFF {
		return 0, BFieldCodecError{ErrorElementOutOfRange, "element out of range for uint32"}
	}

	return uint32(value), nil
}

// EncodeUint16 encodes a uint16 as a single BFieldElement.
func EncodeUint16(value uint16) []field.Element {
	return []field.Element{field.New(uint64(value))}
}

// DecodeUint16 decodes a uint16 from a single BFieldElement.
func DecodeUint16(sequence []field.Element) (uint16, error) {
	if len(sequence) == 0 {
		return 0, BFieldCodecError{ErrorEmptySequence, "empty sequence"}
	}
	if len(sequence) > 1 {
		return 0, BFieldCodecError{ErrorSequenceTooLong, "sequence too long for uint16"}
	}

	value := sequence[0].Value()
	if value > 0xFFFF {
		return 0, BFieldCodecError{ErrorElementOutOfRange, "element out of range for uint16"}
	}

	return uint16(value), nil
}

// EncodeUint8 encodes a uint8 as a single BFieldElement.
func EncodeUint8(value uint8) []field.Element {
	return []field.Element{field.New(uint64(value))}
}

// DecodeUint8 decodes a uint8 from a single BFieldElement.
func DecodeUint8(sequence []field.Element) (uint8, error) {
	if len(sequence) == 0 {
		return 0, BFieldCodecError{ErrorEmptySequence, "empty sequence"}
	}
	if len(sequence) > 1 {
		return 0, BFieldCodecError{ErrorSequenceTooLong, "sequence too long for uint8"}
	}

	value := sequence[0].Value()
	if value > 0xFF {
		return 0, BFieldCodecError{ErrorElementOutOfRange, "element out of range for uint8"}
	}

	return uint8(value), nil
}

// EncodeBool encodes a boolean as a single BFieldElement (0 or 1).
func EncodeBool(value bool) []field.Element {
	if value {
		return []field.Element{field.One}
	}
	return []field.Element{field.Zero}
}

// DecodeBool decodes a boolean from a single BFieldElement.
func DecodeBool(sequence []field.Element) (bool, error) {
	if len(sequence) == 0 {
		return false, BFieldCodecError{ErrorEmptySequence, "empty sequence"}
	}
	if len(sequence) > 1 {
		return false, BFieldCodecError{ErrorSequenceTooLong, "sequence too long for bool"}
	}

	value := sequence[0].Value()
	if value == 0 {
		return false, nil
	} else if value == 1 {
		return true, nil
	} else {
		return false, BFieldCodecError{ErrorElementOutOfRange, "element out of range for bool"}
	}
}

// EncodeXFieldElement encodes an XFieldElement as three BFieldElement values.
func EncodeXFieldElement(element xfield.XFieldElement) []field.Element {
	return []field.Element{
		element.Coefficients[0],
		element.Coefficients[1],
		element.Coefficients[2],
	}
}

// DecodeXFieldElement decodes an XFieldElement from three BFieldElement values.
func DecodeXFieldElement(sequence []field.Element) (xfield.XFieldElement, error) {
	if len(sequence) < 3 {
		return xfield.Zero, BFieldCodecError{ErrorSequenceTooShort, "need at least 3 elements for XFieldElement"}
	}
	if len(sequence) > 3 {
		return xfield.Zero, BFieldCodecError{ErrorSequenceTooLong, "too many elements for XFieldElement"}
	}

	return xfield.New([3]field.Element{
		sequence[0],
		sequence[1],
		sequence[2],
	}), nil
}

// EncodeSlice encodes a slice of BFieldCodec values with length prefix.
func EncodeSlice[T BFieldCodec](slice []T) []field.Element {
	if len(slice) == 0 {
		return []field.Element{field.Zero}
	}

	// Start with length prefix
	result := []field.Element{field.New(uint64(len(slice)))}

	// Encode each element
	for _, item := range slice {
		encoded := item.Encode()
		result = append(result, encoded...)
	}

	return result
}

// DecodeSlice decodes a slice of BFieldCodec values from a sequence with length prefix.

// - First element is the length prefix (number of items)
// - For each item: if static_length is Some, use that; else read length prefix per item
// - Remaining elements are the encoded items
func DecodeSlice[T BFieldCodec](sequence []field.Element, constructor func() T) ([]T, error) {
	if len(sequence) == 0 {
		return nil, BFieldCodecError{ErrorEmptySequence, "empty sequence"}
	}

	// Read length prefix (number of items in the slice)
	numItems := sequence[0].Value()
	sequence = sequence[1:]

	if numItems == 0 {
		return []T{}, nil
	}

	result := make([]T, numItems)

	// Get the static length of the element type
	sampleItem := constructor()
	staticLen := sampleItem.StaticLength()

	for i := 0; i < int(numItems); i++ {
		var itemLength int
		var itemSequence []field.Element

		// Determine item length: either static or from length prefix
		if staticLen != nil {
			// Static length: use it directly
			itemLength = *staticLen
			if len(sequence) < itemLength {
				return nil, BFieldCodecError{
					ErrorSequenceTooShort,
					fmt.Sprintf("sequence too short for item %d (need %d elements)", i, itemLength),
				}
			}
			itemSequence = sequence[:itemLength]
			sequence = sequence[itemLength:]
		} else {
			// Dynamic length: read length prefix for this item
			if len(sequence) == 0 {
				return nil, BFieldCodecError{
					ErrorMissingLengthIndicator,
					fmt.Sprintf("missing length indicator for item %d", i),
				}
			}
			itemLength = int(sequence[0].Value())
			if len(sequence) < 1+itemLength {
				return nil, BFieldCodecError{
					ErrorSequenceTooShort,
					fmt.Sprintf("sequence too short for item %d (need %d elements after prefix)", i, itemLength),
				}
			}
			itemSequence = sequence[1 : 1+itemLength]
			sequence = sequence[1+itemLength:]
		}

		// Decode the item using its Decode method
		item := constructor()
		decoded, err := item.Decode(itemSequence)
		if err != nil {
			return nil, BFieldCodecError{
				ErrorInnerDecodingFailure,
				fmt.Sprintf("failed to decode item %d: %v", i, err),
			}
		}

		// Type-assert the decoded value
		typedItem, ok := decoded.(T)
		if !ok {
			return nil, BFieldCodecError{
				ErrorUnsupportedType,
				fmt.Sprintf("decoded item %d has unexpected type", i),
			}
		}

		result[i] = typedItem
	}

	// Ensure we consumed all the sequence
	if len(sequence) > 0 {
		return nil, BFieldCodecError{ErrorSequenceTooLong, "trailing data after decoding all items"}
	}

	return result, nil
}

// EncodeTuple encodes a tuple of BFieldCodec values.
func EncodeTuple(values ...BFieldCodec) []field.Element {
	var result []field.Element
	for _, value := range values {
		encoded := value.Encode()
		result = append(result, encoded...)
	}
	return result
}

// EncodeOption encodes an optional BFieldCodec value.
func EncodeOption[T BFieldCodec](value *T) []field.Element {
	if value == nil {
		return []field.Element{field.Zero}
	}
	return append([]field.Element{field.One}, (*value).Encode()...)
}

// DecodeOption decodes an optional BFieldCodec value.

// - First element is boolean indicator (0 = None, 1 = Some)
// - If Some, remaining elements are the encoded value
// - If None, sequence must contain only the indicator
func DecodeOption[T BFieldCodec](sequence []field.Element, constructor func() T) (*T, error) {
	if len(sequence) == 0 {
		return nil, BFieldCodecError{ErrorEmptySequence, "empty sequence"}
	}

	// Decode the boolean indicator from first element
	isSome, err := DecodeBool(sequence[0:1])
	if err != nil {
		return nil, BFieldCodecError{ErrorInnerDecodingFailure, fmt.Sprintf("failed to decode option indicator: %v", err)}
	}

	// If None, ensure no additional data
	if !isSome {
		if len(sequence) > 1 {
			return nil, BFieldCodecError{ErrorSequenceTooLong, "None option should not have trailing data"}
		}
		return nil, nil
	}

	// If Some, decode the value from remaining sequence
	value := constructor()

	// Use the Decode method from the BFieldCodec interface
	decoded, err := value.Decode(sequence[1:])
	if err != nil {
		return nil, BFieldCodecError{ErrorInnerDecodingFailure, fmt.Sprintf("failed to decode option value: %v", err)}
	}

	// Convert back to concrete type
	typedValue, ok := decoded.(T)
	if !ok {
		return nil, BFieldCodecError{ErrorUnsupportedType, "decoded value has unexpected type"}
	}

	return &typedValue, nil
}

// EncodeArray encodes an array of BFieldCodec values.
func EncodeArray[T BFieldCodec](array []T) []field.Element {
	var result []field.Element
	for _, item := range array {
		encoded := item.Encode()
		result = append(result, encoded...)
	}
	return result
}

// DecodeArray decodes an array of BFieldCodec values.

// - Each element is decoded sequentially using its StaticLength or Decode method
// - Total number of elements is specified by the caller
func DecodeArray[T BFieldCodec](sequence []field.Element, length int, constructor func() T) ([]T, error) {
	if len(sequence) == 0 && length > 0 {
		return nil, BFieldCodecError{ErrorEmptySequence, "empty sequence"}
	}

	result := make([]T, length)

	// Get the static length of the element type
	sampleItem := constructor()
	staticLen := sampleItem.StaticLength()

	offset := 0
	for i := 0; i < length; i++ {
		var itemLength int
		var itemSequence []field.Element

		// Determine item length: either static or from remaining sequence
		if staticLen != nil {
			itemLength = *staticLen
		} else {
			// For dynamic-length items in an array, we'd need additional context
			// This is a limitation of the current design
			return nil, BFieldCodecError{
				ErrorUnsupportedType,
				"cannot decode arrays of dynamic-length items without length indicators",
			}
		}

		if offset+itemLength > len(sequence) {
			return nil, BFieldCodecError{
				ErrorSequenceTooShort,
				fmt.Sprintf("sequence too short for element %d (need %d elements at offset %d)", i, itemLength, offset),
			}
		}

		itemSequence = sequence[offset : offset+itemLength]
		offset += itemLength

		// Decode the item using its Decode method
		item := constructor()
		decoded, err := item.Decode(itemSequence)
		if err != nil {
			return nil, BFieldCodecError{
				ErrorInnerDecodingFailure,
				fmt.Sprintf("failed to decode element %d: %v", i, err),
			}
		}

		// Type-assert the decoded value
		typedItem, ok := decoded.(T)
		if !ok {
			return nil, BFieldCodecError{
				ErrorUnsupportedType,
				fmt.Sprintf("decoded element %d has unexpected type", i),
			}
		}

		result[i] = typedItem
	}

	// Validate we consumed exactly the right amount
	if offset != len(sequence) {
		return nil, BFieldCodecError{
			ErrorSequenceTooLong,
			fmt.Sprintf("sequence length mismatch: expected %d elements, got %d", offset, len(sequence)),
		}
	}

	return result, nil
}

// Helper functions for common encoding patterns

// EncodeLengthPrefix adds a length prefix to an encoded sequence.
func EncodeLengthPrefix(encoded []field.Element) []field.Element {
	return append([]field.Element{field.New(uint64(len(encoded)))}, encoded...)
}

// DecodeLengthPrefix reads a length prefix from a sequence.
func DecodeLengthPrefix(sequence []field.Element) (length int, remaining []field.Element, error error) {
	if len(sequence) == 0 {
		return 0, nil, BFieldCodecError{ErrorEmptySequence, "empty sequence"}
	}

	length = int(sequence[0].Value())
	if len(sequence) < 1+length {
		return 0, nil, BFieldCodecError{ErrorSequenceTooShort, "sequence too short for indicated length"}
	}

	return length, sequence[1:], nil
}

// ValidateSequenceLength checks if a sequence has the expected length.
func ValidateSequenceLength(sequence []field.Element, expected int) error {
	if len(sequence) < expected {
		return BFieldCodecError{ErrorSequenceTooShort, fmt.Sprintf("need at least %d elements", expected)}
	}
	if len(sequence) > expected {
		return BFieldCodecError{ErrorSequenceTooLong, fmt.Sprintf("too many elements, expected %d", expected)}
	}
	return nil
}

// ConvertToBigInt converts a BFieldElement to big.Int for compatibility.
func ConvertToBigInt(element field.Element) *big.Int {
	return big.NewInt(int64(element.Value()))
}

// ConvertFromBigInt converts a big.Int to BFieldElement.
func ConvertFromBigInt(value *big.Int) field.Element {
	// Take modulo the field prime
	prime := new(big.Int).SetUint64(field.P)
	mod := new(big.Int).Mod(value, prime)
	return field.New(mod.Uint64())
}
