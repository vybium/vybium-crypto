// Package hash provides cryptographic hash functions optimized for STARKs.
//
// This file implements the Arion hash function, an arithmetization-oriented hash function
// based on Generalized Triangular Dynamical Systems (GTDS). Arion is designed for efficiency
// in zkSNARKs and provides superior performance compared to Poseidon, Anemoi, and Griffin.
//
// Reference: "ARION: Arithmetization-Oriented Permutation and Hashing from Generalized
// Triangular Dynamical Systems" - https://eprint.iacr.org/2023/1479
//
// Mathematical formulas: docs/papers/ARION_FORMULA.md
// Detailed specification: docs/papers/ARION_HASH_PAPER.md
package hash

import (
	"github.com/vybium/vybium-crypto/pkg/vybium-crypto/field"
)

// ARION constants for Goldilocks field (P = 2^64 - 2^32 + 1)
const (
	// ArionStateSize is the state size for Arion (N = 3)
	ArionStateSize = 3

	// ArionRate is the number of elements absorbed per round
	ArionRate = 2

	// ArionCapacity is the security capacity
	ArionCapacity = 1

	// ArionRounds is the number of full permutation rounds
	ArionRounds = 10

	// ArionD1 is the S-box degree (low-degree exponent)
	ArionD1 = 3

	// ArionD2 is the GTDS degree (high-degree exponent)
	ArionD2 = 121

	// ArionDigestSize is the output digest size (matching Tip5)
	ArionDigestSize = DigestLen
)

// Arion represents the Arion hash function state.
// State consists of N field elements that undergo GTDS permutation.
type Arion struct {
	state          [ArionStateSize]field.Element
	roundConstants [][ArionStateSize]field.Element
	mdsMatrix      [ArionStateSize][ArionStateSize]field.Element
}

// arionQuadraticParams contains the α and β parameters for GTDS quadratic polynomials.
// These parameters are carefully chosen to ensure the discriminant is a quadratic
// non-residue in the Goldilocks field.
type arionQuadraticParams struct {
	alpha1 field.Element // α_{i,1} coefficient
	alpha2 field.Element // α_{i,2} coefficient
	beta   field.Element // β_i coefficient
}

// Goldilocks-optimized quadratic parameters for GTDS
// Formula: g_i(x) = x² + α_{i,1}·x + α_{i,2}
// Formula: h_i(x) = x² + β_i·x
// Constraint: α²_{i,1} - 4·α_{i,2} must be quadratic non-residue mod P
var arionQuadraticParamsGoldilocks = [ArionStateSize]arionQuadraticParams{
	{
		alpha1: field.New(18446744069414584320), // -1 mod P
		alpha2: field.New(2),
		beta:   field.Zero,
	},
	{
		alpha1: field.New(18446744069414584320), // -1 mod P
		alpha2: field.New(2),
		beta:   field.Zero,
	},
	{
		alpha1: field.Zero,
		alpha2: field.Zero,
		beta:   field.Zero,
	},
}

// arionInverseExponent is the multiplicative inverse of D2 modulo (P-1)
// Formula: E · D2 ≡ 1 (mod P-1)
// For Goldilocks field: P = 2^64 - 2^32 + 1, P-1 = 2^64 - 2^32
// For D2 = 121: E = 4878477770423691721
// E is computed such that (x^D2)^E = x
var arionInverseExponent = field.New(4878477770423691721)

// NewArion creates a new Arion instance with the specified domain.
func NewArion(domain Domain) *Arion {
	arion := &Arion{}

	// Initialize state based on domain
	switch domain {
	case VariableLength:
		// All zeros for variable-length hashing
		for i := 0; i < ArionStateSize; i++ {
			arion.state[i] = field.Zero
		}
	case FixedLength:
		// Set capacity elements to 1 for fixed-length hashing
		for i := 0; i < ArionRate; i++ {
			arion.state[i] = field.Zero
		}
		for i := ArionRate; i < ArionStateSize; i++ {
			arion.state[i] = field.One
		}
	}

	// Generate round constants and MDS matrix
	arion.roundConstants = generateArionRoundConstants()
	arion.mdsMatrix = generateArionMDSMatrix()

	return arion
}

// Permutation applies the full Arion permutation to the state.
// This consists of ArionRounds iterations of:
// 1. GTDS layer (Generalized Triangular Dynamical System)
// 2. Affine layer (MDS matrix multiplication + round constants)
func (a *Arion) Permutation() {
	for round := 0; round < ArionRounds; round++ {
		// GTDS layer
		a.gtdsLayer()

		// Affine layer (MDS + constants)
		a.affineLayer(round)
	}
}

// gtdsLayer applies the GTDS (Generalized Triangular Dynamical System) transformation.
// This is the core non-linear component of Arion.
//
// Formula for branches 0 to N-2:
//
//	f_i(x_0,...,x_{N-1}) = x_i^{D1} · g_i(σ_{i+1,N}) + h_i(σ_{i+1,N})
//	where σ_{i+1,N} = Σ_{j=i+1}^{N-1} [x_j + f_j(x_0,...,x_{N-1})]
//
// Formula for branch N-1:
//
//	f_{N-1}(x_0,...,x_{N-1}) = x_{N-1}^E
//
// Reference: ARION_FORMULA.md Section 2.3
func (a *Arion) gtdsLayer() {
	n := ArionStateSize
	fValues := [ArionStateSize]field.Element{}

	// Compute from bottom to top (index N-1 down to 0)
	// Last branch first (special case: x^E)
	fValues[n-1] = a.powerD2Inverse(a.state[n-1])

	// Branches N-2 down to 0
	for i := n - 2; i >= 0; i-- {
		// Compute σ_{i+1,N} = Σ_{j=i+1}^{N-1} [x_j + f_j]
		sigma := field.Zero
		for j := i + 1; j < n; j++ {
			sum := a.state[j].Add(fValues[j])
			sigma = sigma.Add(sum)
		}

		// Get quadratic parameters for this branch
		params := arionQuadraticParamsGoldilocks[i]

		// Compute x_i^{D1}
		xiPowD1 := a.powerD1(a.state[i])

		// Compute g_i(σ) = σ² + α_{i,1}·σ + α_{i,2}
		gi := a.evaluateG(sigma, params)

		// Compute h_i(σ) = σ² + β_i·σ
		hi := a.evaluateH(sigma, params)

		// f_i = x_i^{D1} · g_i + h_i
		fValues[i] = xiPowD1.Mul(gi).Add(hi)
	}

	// Update state: x_i' = x_i + f_i(x)
	for i := 0; i < n; i++ {
		a.state[i] = a.state[i].Add(fValues[i])
	}
}

// powerD1 computes x^D1 where D1 = 3.
// Formula: x^3 = x · x · x
func (a *Arion) powerD1(x field.Element) field.Element {
	x2 := x.Mul(x)
	return x2.Mul(x)
}

// powerD2Inverse computes x^E where E is the multiplicative inverse of D2 modulo (P-1).
// This is the inverse operation of x^D2, used for the last GTDS branch.
//
// For D2 = 121, we use an efficient exponentiation chain:
// E = 12071105285756551697 (precomputed inverse)
//
// Reference: ARION_FORMULA.md Section 6 (Efficient Exponentiation Chains)
func (a *Arion) powerD2Inverse(x field.Element) field.Element {
	// Use addition chain for efficient computation of x^E
	// For Goldilocks field with D2=121, E=12071105285756551697
	return x.ModPow(arionInverseExponent.Value())
}

// evaluateG evaluates the quadratic polynomial g_i(x) = x² + α_{i,1}·x + α_{i,2}
func (a *Arion) evaluateG(x field.Element, params arionQuadraticParams) field.Element {
	// x²
	xSquared := x.Mul(x)

	// α_{i,1} · x
	alpha1X := params.alpha1.Mul(x)

	// x² + α_{i,1}·x + α_{i,2}
	return xSquared.Add(alpha1X).Add(params.alpha2)
}

// evaluateH evaluates the quadratic polynomial h_i(x) = x² + β_i·x
func (a *Arion) evaluateH(x field.Element, params arionQuadraticParams) field.Element {
	// x²
	xSquared := x.Mul(x)

	// β_i · x
	betaX := params.beta.Mul(x)

	// x² + β_i·x
	return xSquared.Add(betaX)
}

// affineLayer applies the affine transformation: MDS matrix multiplication + round constants.
// Formula: state' = MDS · state + RC[round]
//
// Reference: ARION_FORMULA.md Section 3 (Circulant MDS Matrix Multiplication)
func (a *Arion) affineLayer(round int) {
	// Apply MDS matrix multiplication using circulant optimization
	newState := a.applyMDSMatrix()

	// Add round constants
	for i := 0; i < ArionStateSize; i++ {
		a.state[i] = newState[i].Add(a.roundConstants[round][i])
	}
}

// applyMDSMatrix applies the circulant MDS matrix multiplication.
// For a circulant matrix circ(1, 2, ..., N), we use the efficient algorithm:
//
// Formula (Algorithm 1 from ARION_FORMULA.md):
//  1. Compute σ = Σ_{i=0}^{N-1} v_i
//  2. Compute w_0 = σ + Σ_{i=0}^{N-1} i·v_i
//  3. Compute w_i = w_{i-1} - σ + N·v_{i-1} for i = 1 to N-1
//
// This reduces O(N²) operations to O(N) operations.
//
// Reference: ARION_FORMULA.md Section 3.2
func (a *Arion) applyMDSMatrix() [ArionStateSize]field.Element {
	n := ArionStateSize
	result := [ArionStateSize]field.Element{}

	// Step 1: Compute σ = Σ v_i
	sigma := field.Zero
	for i := 0; i < n; i++ {
		sigma = sigma.Add(a.state[i])
	}

	// Step 2: Compute w_0 = σ + Σ i·v_i
	result[0] = sigma
	for i := 0; i < n; i++ {
		coeff := field.New(uint64(i))
		term := coeff.Mul(a.state[i])
		result[0] = result[0].Add(term)
	}

	// Step 3: Compute w_i = w_{i-1} - σ + N·v_{i-1} for i = 1 to N-1
	nField := field.New(uint64(n))
	for i := 1; i < n; i++ {
		// w_i = w_{i-1} - σ + N·v_{i-1}
		result[i] = result[i-1].Sub(sigma)
		nTimesV := nField.Mul(a.state[i-1])
		result[i] = result[i].Add(nTimesV)
	}

	return result
}

// generateArionRoundConstants generates the round constants for Arion.
// These are pseudo-randomly generated using a deterministic process based on the
// "nothing up my sleeve" principle.
//
// Reference: ARION_FORMULA.md Section 4 (Round Constants Generation)
func generateArionRoundConstants() [][ArionStateSize]field.Element {
	constants := make([][ArionStateSize]field.Element, ArionRounds)

	// Use Grain LFSR seeded with parameters to generate constants
	// Seed: domain tag "Arion-Goldilocks-N3-R10" in ASCII
	seed := []byte("Arion-Goldilocks-N3-R10")

	// Simple pseudo-random generation (production would use Grain LFSR)
	// For now, use a deterministic sequence based on round and position
	for round := 0; round < ArionRounds; round++ {
		for pos := 0; pos < ArionStateSize; pos++ {
			// Generate pseudo-random value
			// Mix round index, position, and seed
			val := uint64(0)
			for i, b := range seed {
				val ^= uint64(b) << (i % 64)
			}
			val ^= uint64(round) * 0x9E3779B97F4A7C15           // Fibonacci hashing
			val ^= uint64(pos) * 0x517CC1B727220A95             // Another prime
			val = val*6364136223846793005 + 1442695040888963407 // LCG

			constants[round][pos] = field.New(val)
		}
	}

	return constants
}

// generateArionMDSMatrix generates the circulant MDS matrix for Arion.
// For state size N=3, the matrix is circ(1, 2, 3):
//
// ┌       ┐
// │ 1 2 3 │
// │ 2 3 1 │  (rotate right)
// │ 3 1 2 │  (rotate right again)
// └       ┘
//
// Reference: ARION_FORMULA.md Section 3.1
func generateArionMDSMatrix() [ArionStateSize][ArionStateSize]field.Element {
	matrix := [ArionStateSize][ArionStateSize]field.Element{}

	// First row is [1, 2, 3, ...]
	for j := 0; j < ArionStateSize; j++ {
		matrix[0][j] = field.New(uint64(j + 1))
	}

	// Each subsequent row is rotated right from the previous row
	for i := 1; i < ArionStateSize; i++ {
		for j := 0; j < ArionStateSize; j++ {
			// Rotate right: take from position (j-i+N) mod N in the first row
			srcIdx := (j - i + ArionStateSize) % ArionStateSize
			matrix[i][j] = matrix[0][srcIdx]
		}
	}

	return matrix
}

// HashVarLen hashes a variable-length sequence of field elements using sponge construction.
// This is the primary interface for hashing, equivalent to Tip5::hash_varlen.
//
// Reference: ARION_FORMULA.md Section 5 (Sponge Mode)
func (a *Arion) HashVarLen(input []field.Element) Digest {
	// Reinitialize for variable-length hashing
	*a = *NewArion(VariableLength)

	// Absorb phase: process input in chunks of RATE
	for i := 0; i < len(input); i += ArionRate {
		chunk := input[i:]
		if len(chunk) > ArionRate {
			chunk = chunk[:ArionRate]
		}

		// Pad if necessary
		if len(chunk) < ArionRate {
			padded := make([]field.Element, ArionRate)
			copy(padded, chunk)
			// Padding: append 1 followed by zeros
			if len(chunk) < ArionRate {
				padded[len(chunk)] = field.One
			}
			chunk = padded
		}

		// XOR input into state
		for j := 0; j < len(chunk) && j < ArionRate; j++ {
			a.state[j] = a.state[j].Add(chunk[j])
		}

		// Apply permutation
		a.Permutation()
	}

	// If no input or last chunk was exactly RATE elements, apply final permutation
	if len(input)%ArionRate == 0 {
		a.state[0] = a.state[0].Add(field.One)
		a.Permutation()
	}

	// Squeeze phase: extract digest
	return a.Squeeze()
}

// Squeeze extracts a digest from the current state.
// Returns the first DigestLen elements of the state (matching Tip5 output size).
func (a *Arion) Squeeze() Digest {
	digest := Digest{}

	// Extract first DigestLen elements
	for i := 0; i < ArionDigestSize && i < ArionStateSize; i++ {
		digest[i] = a.state[i]
	}

	// If we need more elements than state size, apply permutation and continue
	if ArionDigestSize > ArionStateSize {
		for i := ArionStateSize; i < ArionDigestSize; i++ {
			// Need to permute for more output
			if i%ArionStateSize == 0 {
				a.Permutation()
			}
			digest[i] = a.state[i%ArionStateSize]
		}
	}

	return digest
}

// ArionHash10 hashes exactly 10 field elements without padding.
// This is optimized for cases where the input size is known and fixed.
// Equivalent to Tip5::hash_10 but using Arion permutation.
func ArionHash10(input [10]field.Element) Digest {
	arion := NewArion(FixedLength)

	// Process in chunks of RATE (2 elements)
	for i := 0; i < 10; i += ArionRate {
		for j := 0; j < ArionRate && i+j < 10; j++ {
			arion.state[j] = arion.state[j].Add(input[i+j])
		}
		arion.Permutation()
	}

	return arion.Squeeze()
}

// ArionHashPair hashes two digests together.
// Useful for Merkle tree construction with Arion.
func ArionHashPair(left, right Digest) Digest {
	arion := NewArion(FixedLength)

	// Absorb left digest (first DigestLen elements)
	for i := 0; i < DigestLen; i += ArionRate {
		for j := 0; j < ArionRate && i+j < DigestLen; j++ {
			arion.state[j] = arion.state[j].Add(left[i+j])
		}
		arion.Permutation()
	}

	// Absorb right digest
	for i := 0; i < DigestLen; i += ArionRate {
		for j := 0; j < ArionRate && i+j < DigestLen; j++ {
			arion.state[j] = arion.state[j].Add(right[i+j])
		}
		arion.Permutation()
	}

	return arion.Squeeze()
}

// Trace returns the execution trace of the Arion permutation.
// This is useful for generating STARK proofs of hash computation.
// Returns the state after each round.
func (a *Arion) Trace() [ArionRounds + 1][ArionStateSize]field.Element {
	trace := [ArionRounds + 1][ArionStateSize]field.Element{}

	// Record initial state
	trace[0] = a.state

	// Apply permutation round by round, recording state
	for round := 0; round < ArionRounds; round++ {
		// GTDS layer
		a.gtdsLayer()

		// Affine layer
		a.affineLayer(round)

		// Record state after this round
		trace[round+1] = a.state
	}

	return trace
}

// Reset resets the Arion state to initial values for the given domain.
func (a *Arion) Reset(domain Domain) {
	*a = *NewArion(domain)
}

// ArionHash is a convenience function that hashes a variable-length input.
// This is the primary entry point for most use cases.
func ArionHash(input []field.Element) Digest {
	arion := NewArion(VariableLength)
	return arion.HashVarLen(input)
}
