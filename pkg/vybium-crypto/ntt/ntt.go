// Package ntt provides Number Theoretic Transform operations.
//
// NTT is the finite field equivalent of FFT, used for fast polynomial multiplication.
// It enables O(n log n) polynomial multiplication by converting between coefficient
// and evaluation representations. This is essential for efficient STARK proof generation.
// Reference: https://eprint.iacr.org/2016/504.pdf (Longa and Naehrig)
package ntt

import (
	"fmt"
	"math/bits"
	"sync"

	"github.com/vybium/vybium-crypto/pkg/vybium-crypto/field"
)

// Cache for twiddle factors to avoid recomputation
var (
	nttTwiddleCache  = make(map[uint32][][]field.Element)
	inttTwiddleCache = make(map[uint32][][]field.Element)
	swapIndicesCache = make(map[uint32][]int)
	cacheMutex       sync.RWMutex
)

// NTT performs an in-place Number Theoretic Transform.
// The input slice length must be a power of 2.
//
// Production implementation.
//
// Panics if:
// - len(x) is not a power of 2
// - len(x) > 2^31
func NTT(x []field.Element) {
	n := len(x)
	if n == 0 {
		return
	}

	// Validate length
	if n&(n-1) != 0 {
		panic(fmt.Sprintf("NTT requires power-of-2 length, got %d", n))
	}
	if n > (1 << 31) {
		panic(fmt.Sprintf("NTT length too large: %d", n))
	}

	// Get or compute twiddle factors
	twiddles := getTwiddleFactors(uint32(n), false)

	// Perform NTT
	nttUnchecked(x, twiddles)
}

// INTT performs an in-place Inverse Number Theoretic Transform.
// The input slice length must be a power of 2.
//
// Production implementation.
//
// Panics if:
// - len(x) is not a power of 2
// - len(x) > 2^31
func INTT(x []field.Element) {
	n := len(x)
	if n == 0 {
		return
	}

	// Validate length
	if n&(n-1) != 0 {
		panic(fmt.Sprintf("INTT requires power-of-2 length, got %d", n))
	}
	if n > (1 << 31) {
		panic(fmt.Sprintf("INTT length too large: %d", n))
	}

	// Get or compute twiddle factors (using inverse root of unity)
	twiddles := getTwiddleFactors(uint32(n), true)

	// Perform INTT
	nttUnchecked(x, twiddles)

	// Unscale by multiplying by 1/n
	unscale(x)
}

// nttUnchecked performs the core NTT algorithm.
// Assumes:
// - len(x) is a power of 2
// - twiddle factors are correct for len(x)
//
// Production implementation.
func nttUnchecked(x []field.Element, twiddles [][]field.Element) {
	n := uint32(len(x))
	if n <= 1 {
		return
	}

	// Bit-reverse permutation
	swapIndices := getSwapIndices(n)
	for i, revI := range swapIndices {
		if revI > 0 {
			x[i], x[revI] = x[revI], x[i]
		}
	}

	// Cooley-Tukey butterfly operations
	m := uint32(1)
	for _, twiddleRow := range twiddles {
		k := uint32(0)
		for k < n {
			for j := uint32(0); j < m; j++ {
				idx1 := k + j
				idx2 := k + j + m

				u := x[idx1]
				v := x[idx2].Mul(twiddleRow[j])

				x[idx1] = u.Add(v)
				x[idx2] = u.Sub(v)
			}
			k += 2 * m
		}
		m *= 2
	}
}

// unscale multiplies every element by 1/n.
// Used after INTT to complete the inverse transform.
func unscale(x []field.Element) {
	if len(x) == 0 {
		return
	}

	nInv := field.New(uint64(len(x))).Inverse()
	for i := range x {
		x[i] = x[i].Mul(nInv)
	}
}

// getTwiddleFactors returns the twiddle factors for NTT/INTT of given length.
// Results are cached for performance.
//
// If inverse is true, uses the inverse root of unity (for INTT).
func getTwiddleFactors(n uint32, inverse bool) [][]field.Element {
	cacheMutex.RLock()
	cache := nttTwiddleCache
	if inverse {
		cache = inttTwiddleCache
	}
	if twiddles, ok := cache[n]; ok {
		cacheMutex.RUnlock()
		return twiddles
	}
	cacheMutex.RUnlock()

	// Compute twiddle factors
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	// Double-check after acquiring write lock
	if twiddles, ok := cache[n]; ok {
		return twiddles
	}

	// Get primitive root of unity for this domain size
	omega := field.PrimitiveRootOfUnity(uint64(n))
	if omega.IsZero() {
		panic(fmt.Sprintf("no primitive root of unity for n=%d", n))
	}

	if inverse {
		omega = omega.Inverse()
	}

	// Compute twiddle factors
	log2N := bits.Len32(n) - 1
	twiddles := make([][]field.Element, log2N)

	for i := uint32(0); i < uint32(log2N); i++ {
		m := uint32(1) << i
		exponent := n / (2 * m)
		wm := omega.ModPow(uint64(exponent))

		twiddleRow := make([]field.Element, m)
		twiddleRow[0] = field.One
		for j := uint32(1); j < m; j++ {
			twiddleRow[j] = twiddleRow[j-1].Mul(wm)
		}

		twiddles[i] = twiddleRow
	}

	cache[n] = twiddles
	return twiddles
}

// getSwapIndices returns the bit-reverse permutation indices.
// For index i, if swapIndices[i] > 0, then i should be swapped with swapIndices[i].
// If swapIndices[i] == 0 or swapIndices[i] == i, no swap is needed.
//
// Production implementation.
func getSwapIndices(n uint32) []int {
	cacheMutex.RLock()
	if indices, ok := swapIndicesCache[n]; ok {
		cacheMutex.RUnlock()
		return indices
	}
	cacheMutex.RUnlock()

	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	// Double-check after acquiring write lock
	if indices, ok := swapIndicesCache[n]; ok {
		return indices
	}

	log2N := bits.Len32(n) - 1
	indices := make([]int, n)

	for k := uint32(0); k < n; k++ {
		revK := bitReverse(k, uint32(log2N))

		// Only store swap if k < revK to avoid double-swapping
		if k < revK {
			indices[k] = int(revK)
		} else {
			indices[k] = 0 // No swap needed
		}
	}

	swapIndicesCache[n] = indices
	return indices
}

// bitReverse reverses the bits of k using log2N bits.
// Production implementation.
func bitReverse(k, log2N uint32) uint32 {
	// Standard bit-reversal algorithm
	k = ((k & 0x55555555) << 1) | ((k & 0xaaaaaaaa) >> 1)
	k = ((k & 0x33333333) << 2) | ((k & 0xcccccccc) >> 2)
	k = ((k & 0x0f0f0f0f) << 4) | ((k & 0xf0f0f0f0) >> 4)
	k = ((k & 0x00ff00ff) << 8) | ((k & 0xff00ff00) >> 8)
	k = bits.RotateLeft32(k, 16)
	return k >> (32 - log2N)
}

// NextPowerOfTwo returns the next power of 2 >= n.
func NextPowerOfTwo(n int) int {
	if n <= 0 {
		return 1
	}
	if n&(n-1) == 0 {
		return n
	}
	return 1 << bits.Len(uint(n))
}

// IsPowerOfTwo returns true if n is a power of 2.
func IsPowerOfTwo(n int) bool {
	return n > 0 && (n&(n-1)) == 0
}
