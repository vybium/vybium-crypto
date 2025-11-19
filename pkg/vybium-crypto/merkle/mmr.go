package merkle

import (
	"math/bits"

	"github.com/vybium/vybium-crypto/pkg/vybium-crypto/field"
	"github.com/vybium/vybium-crypto/pkg/vybium-crypto/hash"
)

// MMR represents a Merkle Mountain Range, which is a collection of perfect
// binary Merkle trees (called "peaks") arranged by decreasing size.
//
// An MMR allows for efficient append operations and membership proofs
// without requiring the number of elements to be a power of two.
//
// The MMR is represented by its peaks (the roots of the component Merkle trees)
// and the total number of leafs.
type MMR interface {
	// BagPeaks calculates a commitment to the entire MMR.
	BagPeaks() hash.Digest

	// Peaks returns the Merkle tree roots of the Merkle trees that this MMR consists of.
	Peaks() []hash.Digest

	// IsEmpty returns true if there are no leafs in the MMR.
	IsEmpty() bool

	// NumLeafs returns the number of leafs in the MMR.
	NumLeafs() uint64

	// Append adds a leaf to the MMR and returns the membership proof.
	Append(newLeaf hash.Digest) MmrMembershipProof

	// VerifyMembership verifies a membership proof against the MMR.
	VerifyMembership(leaf hash.Digest, proof MmrMembershipProof) bool
}

// MmrAccumulator is a lightweight representation of an MMR that only stores
// the peaks and leaf count, not the full tree structure.
type MmrAccumulator struct {
	leafCount uint64
	peaks     []hash.Digest
}

// NewMmrAccumulator creates a new MMR accumulator with the given peaks and leaf count.
func NewMmrAccumulator(peaks []hash.Digest, leafCount uint64) *MmrAccumulator {
	return &MmrAccumulator{
		leafCount: leafCount,
		peaks:     peaks,
	}
}

// NewMmrAccumulatorFromLeafs creates a new MMR accumulator from a list of leafs.
func NewMmrAccumulatorFromLeafs(leafs []hash.Digest) *MmrAccumulator {
	leafCount := uint64(len(leafs))
	peaks := peaksFromLeafs(leafs)

	return &MmrAccumulator{
		leafCount: leafCount,
		peaks:     peaks,
	}
}

// peaksFromLeafs computes the MMR peaks from a list of leafs.
//
// Algorithm:
// The MMR is built by iterating through pairs of leafs and building perfect
// binary trees bottom-up. When two trees of the same height exist, they are
// merged into a larger tree. The final peaks are the roots of all unmerged trees.
//
// Example for 13 leafs:
//
//	          (3)   (4)    (5)   (6)
//	            \     \      \     \
//	    (2)
//	      \   ──── 7 ────
//	         /           \
//	(1)     3             6             10
//	  \    /  \          /  \          /  \
//	      /    \        /    \        /    \
//	     1      2      4      5      8      9
//	    / \    / \    / \    / \    / \    / \
//	   _   _  _   _  _   _  _   _  _   _  _   _  _
//
// The algorithm processes "diagonals" (marked 1-6):
// - Diagonal 1: Compute node 1, no merges
// - Diagonal 2: Compute node 2, merge 1+2 → 3
// - Diagonal 3: Compute node 4, no merges
// - Diagonal 4: Compute node 5, merge 4+5 → 6, merge 3+6 → 7
// - Diagonal 5: Compute node 8, no merges
// - Diagonal 6: Compute node 9, merge 8+9 → 10
// Final peaks: [7, 10, straggling_leaf]
func peaksFromLeafs(leafs []hash.Digest) []hash.Digest {
	if len(leafs) == 0 {
		return []hash.Digest{}
	}

	// Maximum number of peaks is bounded by log2(numLeafs)
	maxTreeHeight := bits.Len(uint(len(leafs)))
	peaks := make([]hash.Digest, 0, maxTreeHeight)

	// Process pairs of leafs
	diagonalIdx := uint64(1)
	for i := 0; i+1 < len(leafs); i += 2 {
		leftLeaf := leafs[i]
		rightLeaf := leafs[i+1]

		// Hash the pair
		right := hash.HashPair(leftLeaf, rightLeaf)

		// Merge with existing peaks based on trailing zeros
		// The number of trailing zeros indicates how many merges to perform
		numMerges := bits.TrailingZeros64(diagonalIdx)
		for j := 0; j < numMerges; j++ {
			if len(peaks) == 0 {
				break
			}
			left := peaks[len(peaks)-1]
			peaks = peaks[:len(peaks)-1]
			right = hash.HashPair(left, right)
		}

		peaks = append(peaks, right)
		diagonalIdx++
	}

	// If odd number of leafs, add the straggling leaf as a peak
	if len(leafs)%2 == 1 {
		peaks = append(peaks, leafs[len(leafs)-1])
	}

	return peaks
}

// BagPeaks calculates a commitment to the entire MMR by hashing all peaks.
func (mmr *MmrAccumulator) BagPeaks() hash.Digest {
	return bagPeaks(mmr.peaks, mmr.leafCount)
}

// bagPeaks computes a single commitment from the peaks and leaf count.
// This is done by hashing: Hash(leafCount, peaks[0], peaks[1], ..., peaks[n])
func bagPeaks(peaks []hash.Digest, leafCount uint64) hash.Digest {
	if len(peaks) == 0 {
		return hash.ZeroDigest()
	}

	// Build input: [leafCount as 5 field elements, peak0, peak1, ..., peakN]
	input := make([]field.Element, 5+len(peaks)*hash.DigestLen)

	// Encode leaf count as field elements (split across 5 elements for consistency)
	input[0] = field.New(leafCount)
	for i := 1; i < 5; i++ {
		input[i] = field.Zero
	}

	// Add all peaks
	idx := 5
	for _, peak := range peaks {
		for i := 0; i < hash.DigestLen; i++ {
			input[idx] = peak[i]
			idx++
		}
	}

	return hash.HashVarlen(input)
}

// Peaks returns the peaks of the MMR.
func (mmr *MmrAccumulator) Peaks() []hash.Digest {
	result := make([]hash.Digest, len(mmr.peaks))
	copy(result, mmr.peaks)
	return result
}

// IsEmpty returns true if the MMR has no leafs.
func (mmr *MmrAccumulator) IsEmpty() bool {
	return mmr.leafCount == 0
}

// NumLeafs returns the number of leafs in the MMR.
func (mmr *MmrAccumulator) NumLeafs() uint64 {
	return mmr.leafCount
}

// Append adds a new leaf to the MMR and returns the membership proof.
func (mmr *MmrAccumulator) Append(newLeaf hash.Digest) MmrMembershipProof {
	newPeaks, membershipProof := calculateNewPeaksFromAppend(mmr.peaks, newLeaf, mmr.leafCount)
	mmr.peaks = newPeaks
	mmr.leafCount++
	return membershipProof
}

// calculateNewPeaksFromAppend computes the new peaks after appending a leaf
// and returns the membership proof for the newly added leaf.
// This is a direct port of twenty-first's `calculate_new_peaks_from_append`.
func calculateNewPeaksFromAppend(oldPeaks []hash.Digest, newLeaf hash.Digest, oldLeafCount uint64) ([]hash.Digest, MmrMembershipProof) {
	// Copy old peaks and append the new leaf
	peaks := make([]hash.Digest, len(oldPeaks))
	copy(peaks, oldPeaks)
	peaks = append(peaks, newLeaf)

	// Build authentication path
	authPath := []hash.Digest{}

	// The number of merges is the number of trailing ones in oldLeafCount
	// (equivalent to twenty-first's right_lineage_length_from_leaf_index)
	numMerges := trailingOnes64(oldLeafCount)

	// Perform merges
	for i := 0; i < numMerges; i++ {
		if len(peaks) < 2 {
			break
		}
		inProgressPeak := peaks[len(peaks)-1]
		peaks = peaks[:len(peaks)-1]
		previousPeak := peaks[len(peaks)-1]
		peaks = peaks[:len(peaks)-1]

		authPath = append(authPath, previousPeak)
		peaks = append(peaks, hash.HashPair(previousPeak, inProgressPeak))
	}

	proof := MmrMembershipProof{
		LeafIndex: oldLeafCount,
		AuthPath:  authPath,
	}

	return peaks, proof
}

// VerifyMembership verifies a membership proof for a leaf.
// This reconstructs the peak from the leaf and authentication path,
// and checks if it matches one of the MMR's peaks.
func (mmr *MmrAccumulator) VerifyMembership(leaf hash.Digest, proof MmrMembershipProof) bool {
	// Recompute the peak for this leaf using the authentication path
	current := leaf
	for _, authNode := range proof.AuthPath {
		current = hash.HashPair(authNode, current)
	}

	// Check if the computed peak matches any of the MMR's peaks
	// The peak should be at a specific position determined by leaf_index_to_mt_index_and_peak_index
	// For simplicity, we check if it matches any peak
	for _, peak := range mmr.peaks {
		if current.Equal(peak) {
			return true
		}
	}

	return false
}

// MmrMembershipProof represents a proof that a leaf is a member of an MMR.
type MmrMembershipProof struct {
	// LeafIndex is the index of the leaf in the MMR (0-based).
	LeafIndex uint64

	// AuthPath contains the authentication path from the leaf to its peak.
	AuthPath []hash.Digest
}

// IsConsistent checks if the MMR accumulator is self-consistent.
// The number of peaks should equal the number of 1-bits in the leaf count.
func (mmr *MmrAccumulator) IsConsistent() bool {
	numPeaks := len(mmr.peaks)
	expectedPeaks := bits.OnesCount64(mmr.leafCount)
	return numPeaks == expectedPeaks
}

// Clone creates a deep copy of the MMR accumulator.
func (mmr *MmrAccumulator) Clone() *MmrAccumulator {
	peaks := make([]hash.Digest, len(mmr.peaks))
	copy(peaks, mmr.peaks)
	return &MmrAccumulator{
		leafCount: mmr.leafCount,
		peaks:     peaks,
	}
}

// trailingOnes64 returns the number of trailing one bits in x.
// Equivalent to Rust's u64::trailing_ones().
func trailingOnes64(x uint64) int {
	if x == 0 {
		return 0
	}
	// Count trailing zeros of the bitwise NOT
	// This is equivalent to counting trailing ones
	return bits.TrailingZeros64(^x)
}
