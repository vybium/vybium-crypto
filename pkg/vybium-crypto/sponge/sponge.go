// Package sponge provides cryptographic sponge construction for hash functions.
//
// A cryptographic sponge is a construction that can absorb arbitrary-length input
// and squeeze arbitrary-length output using a fixed-width permutation function.
// It's the foundation for hash functions like Tip5 and Poseidon.
//
// Key Features:
// - Absorb arbitrary-length input in fixed-size chunks
// - Squeeze arbitrary-length output in fixed-size chunks
// - Padding scheme to handle variable-length inputs
// - Domain separation for different use cases
package sponge

import (
	"fmt"

	"github.com/vybium/vybium-crypto/pkg/vybium-crypto/field"
	"github.com/vybium/vybium-crypto/pkg/vybium-crypto/hash"
)

const (
	// Rate is the number of field elements that can be absorbed in one permutation.
	// This Production implementation.
	Rate = 10
)

// Domain represents the hashing domain for collision prevention.
// Different domains ensure that hashing different types of data
// produces different outputs even if the raw data is identical.
type Domain int

const (
	// VariableLength domain is used for hashing objects that potentially
	// serialize to more than RATE field elements.
	VariableLength Domain = iota

	// FixedLength domain is used for hashing objects that always fit
	// within RATE field elements, e.g. a pair of Digests.
	FixedLength
)

func (d Domain) String() string {
	switch d {
	case VariableLength:
		return "VariableLength"
	case FixedLength:
		return "FixedLength"
	default:
		return "Unknown"
	}
}

// Sponge defines the interface for cryptographic sponge constructions.
// A sponge can absorb arbitrary-length input and squeeze arbitrary-length output
// using a fixed-width permutation function.
//
// Production implementation.
type Sponge interface {
	// Init creates a new sponge instance.
	Init() Sponge

	// Absorb absorbs a chunk of RATE field elements into the sponge state.
	Absorb(input [Rate]field.Element)

	// Squeeze squeezes a chunk of RATE field elements from the sponge state.
	Squeeze() [Rate]field.Element

	// PadAndAbsorbAll absorbs arbitrary-length input with proper padding.
	// This is the main method for hashing variable-length data.
	PadAndAbsorbAll(input []field.Element)

	// Clone creates a copy of the sponge state.
	Clone() Sponge

	// Reset resets the sponge to its initial state.
	Reset()
}

// Tip5Sponge implements the Sponge interface using the Tip5 permutation.
// This is the primary sponge implementation used in STARK proofs.
type Tip5Sponge struct {
	state  [Rate]field.Element
	domain Domain
}

// NewTip5Sponge creates a new Tip5 sponge with the specified domain.
func NewTip5Sponge(domain Domain) *Tip5Sponge {
	return &Tip5Sponge{
		state:  [Rate]field.Element{},
		domain: domain,
	}
}

// Init creates a new Tip5 sponge instance.
func (s *Tip5Sponge) Init() Sponge {
	return NewTip5Sponge(s.domain)
}

// Absorb absorbs a chunk of RATE field elements into the sponge state.
// This uses the Tip5 permutation to update the internal state.
func (s *Tip5Sponge) Absorb(input [Rate]field.Element) {
	// XOR the input with the current state
	for i := 0; i < Rate; i++ {
		s.state[i] = s.state[i].Add(input[i])
	}

	// Apply Tip5 permutation to the state
	s.applyTip5Permutation()
}

// Squeeze squeezes a chunk of RATE field elements from the sponge state.
// This extracts output from the current state and applies permutation.
func (s *Tip5Sponge) Squeeze() [Rate]field.Element {
	// Extract the current state as output
	output := s.state

	// Apply permutation to prepare for next squeeze
	s.applyTip5Permutation()

	return output
}

// PadAndAbsorbAll absorbs arbitrary-length input with proper padding.
// This is the main method for hashing variable-length data.
func (s *Tip5Sponge) PadAndAbsorbAll(input []field.Element) {
	// Process input in chunks of RATE
	for i := 0; i < len(input); i += Rate {
		var chunk [Rate]field.Element

		// Copy available elements
		copyCount := Rate
		if i+Rate > len(input) {
			copyCount = len(input) - i
		}

		for j := 0; j < copyCount; j++ {
			chunk[j] = input[i+j]
		}

		// If this is the last chunk and it's not full, apply padding
		if i+Rate >= len(input) && copyCount < Rate {
			// Pad with [1, 0, 0, ...] - at least one element of padding
			chunk[copyCount] = field.One
			for j := copyCount + 1; j < Rate; j++ {
				chunk[j] = field.Zero
			}
		}

		s.Absorb(chunk)
	}
}

// Clone creates a copy of the sponge state.
func (s *Tip5Sponge) Clone() Sponge {
	clone := &Tip5Sponge{
		state:  s.state,
		domain: s.domain,
	}
	return clone
}

// Reset resets the sponge to its initial state.
func (s *Tip5Sponge) Reset() {
	s.state = [Rate]field.Element{}
}

// applyTip5Permutation applies the Tip5 permutation to the sponge state.
func (s *Tip5Sponge) applyTip5Permutation() {
	// Convert state to digest format for Tip5
	var digest [5]field.Element
	copy(digest[:], s.state[:])

	// Apply Tip5 permutation
	permuted := hash.Tip5Permutation(digest)

	// Copy result back to state
	copy(s.state[:], permuted[:])
}

// PoseidonSponge implements the Sponge interface using the Poseidon permutation.
// This provides an alternative sponge construction using Poseidon hash.
type PoseidonSponge struct {
	state  [Rate]field.Element
	domain Domain
}

// NewPoseidonSponge creates a new Poseidon sponge with the specified domain.
func NewPoseidonSponge(domain Domain) *PoseidonSponge {
	return &PoseidonSponge{
		state:  [Rate]field.Element{},
		domain: domain,
	}
}

// Init creates a new Poseidon sponge instance.
func (s *PoseidonSponge) Init() Sponge {
	return NewPoseidonSponge(s.domain)
}

// Absorb absorbs a chunk of RATE field elements into the sponge state.
func (s *PoseidonSponge) Absorb(input [Rate]field.Element) {
	// XOR the input with the current state
	for i := 0; i < Rate; i++ {
		s.state[i] = s.state[i].Add(input[i])
	}

	// Apply Poseidon permutation to the state
	s.applyPoseidonPermutation()
}

// Squeeze squeezes a chunk of RATE field elements from the sponge state.
func (s *PoseidonSponge) Squeeze() [Rate]field.Element {
	// Extract the current state as output
	output := s.state

	// Apply permutation to prepare for next squeeze
	s.applyPoseidonPermutation()

	return output
}

// PadAndAbsorbAll absorbs arbitrary-length input with proper padding.
func (s *PoseidonSponge) PadAndAbsorbAll(input []field.Element) {
	// Process input in chunks of RATE
	for i := 0; i < len(input); i += Rate {
		var chunk [Rate]field.Element

		// Copy available elements
		copyCount := Rate
		if i+Rate > len(input) {
			copyCount = len(input) - i
		}

		for j := 0; j < copyCount; j++ {
			chunk[j] = input[i+j]
		}

		// If this is the last chunk and it's not full, apply padding
		if i+Rate >= len(input) && copyCount < Rate {
			// Pad with [1, 0, 0, ...] - at least one element of padding
			chunk[copyCount] = field.One
			for j := copyCount + 1; j < Rate; j++ {
				chunk[j] = field.Zero
			}
		}

		s.Absorb(chunk)
	}
}

// Clone creates a copy of the sponge state.
func (s *PoseidonSponge) Clone() Sponge {
	clone := &PoseidonSponge{
		state:  s.state,
		domain: s.domain,
	}
	return clone
}

// Reset resets the sponge to its initial state.
func (s *PoseidonSponge) Reset() {
	s.state = [Rate]field.Element{}
}

// applyPoseidonPermutation applies the Poseidon permutation to the sponge state.
func (s *PoseidonSponge) applyPoseidonPermutation() {
	// Convert state to digest format for Poseidon
	var digest [5]field.Element
	copy(digest[:], s.state[:])

	// Apply Poseidon permutation
	permuted := hash.PoseidonPermutation(digest)

	// Copy result back to state
	copy(s.state[:], permuted[:])
}

// HashVarlen hashes variable-length input using the specified sponge.
// This is a convenience function that handles the full sponge protocol.
func HashVarlen(sponge Sponge, input []field.Element) []field.Element {
	// Reset sponge to initial state
	sponge.Reset()

	// Absorb all input with padding
	sponge.PadAndAbsorbAll(input)

	// Squeeze output (typically one chunk for hash)
	output := sponge.Squeeze()
	return output[:]
}

// HashFixed hashes fixed-length input using the specified sponge.
// This is optimized for inputs that fit within RATE elements.
func HashFixed(sponge Sponge, input []field.Element) []field.Element {
	if len(input) > Rate {
		panic(fmt.Sprintf("input length %d exceeds RATE %d", len(input), Rate))
	}

	// Reset sponge to initial state
	sponge.Reset()

	// Pad input to RATE elements
	var chunk [Rate]field.Element
	copy(chunk[:], input)

	// If input is shorter than RATE, pad with [1, 0, 0, ...]
	if len(input) < Rate {
		chunk[len(input)] = field.One
		for i := len(input) + 1; i < Rate; i++ {
			chunk[i] = field.Zero
		}
	}

	// Absorb the padded chunk
	sponge.Absorb(chunk)

	// Squeeze output
	output := sponge.Squeeze()
	return output[:]
}

// SampleIndices samples random indices from a range using the sponge.
// This is used for FRI queries and other randomized sampling.
func SampleIndices(sponge Sponge, upperBound int, numIndices int) []int {
	if upperBound <= 0 || numIndices <= 0 {
		return []int{}
	}

	// Limit numIndices to upperBound to prevent infinite loops
	if numIndices > upperBound {
		numIndices = upperBound
	}

	indices := make([]int, 0, numIndices)
	used := make(map[int]bool)

	// Generate random field elements and convert to indices
	for len(indices) < numIndices {
		// Squeeze random elements from sponge
		random := sponge.Squeeze()

		for _, element := range random {
			// Convert field element to index using proper modulo
			val := element.Value()
			// Ensure positive result by adding upperBound before modulo
			index := int((val + uint64(upperBound)) % uint64(upperBound))

			// Avoid duplicates
			if !used[index] {
				indices = append(indices, index)
				used[index] = true

				if len(indices) >= numIndices {
					break
				}
			}
		}
	}

	return indices
}

// ValidateSpongeInput validates that input is appropriate for the sponge.
func ValidateSpongeInput(input []field.Element) error {
	if len(input) == 0 {
		return fmt.Errorf("input cannot be empty")
	}

	// Check for reasonable input length (prevent DoS)
	maxLength := 1024 * 1024 // 1MB worth of field elements
	if len(input) > maxLength {
		return fmt.Errorf("input too long: %d elements (max %d)", len(input), maxLength)
	}

	return nil
}

// GetSpongeRate returns the rate constant for sponge operations.
func GetSpongeRate() int {
	return Rate
}

// IsValidDomain checks if a domain value is valid.
func IsValidDomain(domain Domain) bool {
	return domain == VariableLength || domain == FixedLength
}
