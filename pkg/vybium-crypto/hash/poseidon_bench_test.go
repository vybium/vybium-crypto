package hash

import (
	"testing"

	"github.com/vybium/vybium-crypto/pkg/vybium-crypto/field"
)

// BenchmarkPoseidonHash128 benchmarks Poseidon with 128-bit security
func BenchmarkPoseidonHash128(b *testing.B) {
	poseidon, err := NewPoseidon(GetDefaultPoseidonParameters(128))
	if err != nil {
		b.Fatal(err)
	}

	inputs := make([]field.Element, 10)
	for i := 0; i < 10; i++ {
		inputs[i] = field.New(uint64(i + 1))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = poseidon.Hash(inputs)
	}
}

// BenchmarkPoseidonHash256 benchmarks Poseidon with 256-bit security
func BenchmarkPoseidonHash256(b *testing.B) {
	poseidon, err := NewPoseidon(GetDefaultPoseidonParameters(256))
	if err != nil {
		b.Fatal(err)
	}

	inputs := make([]field.Element, 10)
	for i := 0; i < 10; i++ {
		inputs[i] = field.New(uint64(i + 1))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = poseidon.Hash(inputs)
	}
}

// BenchmarkPoseidonHashTwo benchmarks hashing two elements
func BenchmarkPoseidonHashTwo(b *testing.B) {
	poseidon, err := NewPoseidon(nil)
	if err != nil {
		b.Fatal(err)
	}

	left := field.New(123)
	right := field.New(456)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = poseidon.HashTwo(left, right)
	}
}

// BenchmarkPoseidonSponge benchmarks the sponge construction
func BenchmarkPoseidonSponge(b *testing.B) {
	inputs := make([]field.Element, 10)
	for i := 0; i < 10; i++ {
		inputs[i] = field.New(uint64(i + 1))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sponge, _ := NewPoseidonSponge(nil)
		sponge.Absorb(inputs)
		_ = sponge.Squeeze(1)
	}
}

// BenchmarkGrainLFSR benchmarks the Grain LFSR parameter generation
func BenchmarkGrainLFSR(b *testing.B) {
	params := GetDefaultPoseidonParameters(128)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lfsr := NewGrainLFSR(params)
		_ = lfsr.NextFieldElement()
	}
}

// BenchmarkMDSMatrixGeneration benchmarks MDS matrix generation
func BenchmarkMDSMatrixGeneration(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := generatePoseidonMDSMatrix(4)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkRoundConstantsGeneration benchmarks round constants generation
func BenchmarkRoundConstantsGeneration(b *testing.B) {
	params := GetDefaultPoseidonParameters(128)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := generatePoseidonRoundConstants(params)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkPoseidonVaryingInputSizes benchmarks different input sizes
func BenchmarkPoseidonVaryingInputSizes(b *testing.B) {
	poseidon, err := NewPoseidon(nil)
	if err != nil {
		b.Fatal(err)
	}

	sizes := []int{1, 2, 5, 10, 20, 50, 100}

	for _, size := range sizes {
		inputs := make([]field.Element, size)
		for i := 0; i < size; i++ {
			inputs[i] = field.New(uint64(i + 1))
		}

		b.Run("size_"+string(rune('0'+size)), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = poseidon.Hash(inputs)
			}
		})
	}
}

// BenchmarkPoseidonFullRound benchmarks a single full round
func BenchmarkPoseidonFullRound(b *testing.B) {
	poseidon, err := NewPoseidon(nil)
	if err != nil {
		b.Fatal(err)
	}

	state := make([]field.Element, poseidon.width)
	for i := 0; i < poseidon.width; i++ {
		state[i] = field.New(uint64(i + 1))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = poseidon.fullRound(state, 0)
	}
}

// BenchmarkPoseidonPartialRound benchmarks a single partial round
func BenchmarkPoseidonPartialRound(b *testing.B) {
	poseidon, err := NewPoseidon(nil)
	if err != nil {
		b.Fatal(err)
	}

	state := make([]field.Element, poseidon.width)
	for i := 0; i < poseidon.width; i++ {
		state[i] = field.New(uint64(i + 1))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = poseidon.partialRound(state, 0)
	}
}

// BenchmarkSboxComputation benchmarks the S-box operation
func BenchmarkSboxComputation(b *testing.B) {
	poseidon, err := NewPoseidon(nil)
	if err != nil {
		b.Fatal(err)
	}

	x := field.New(42)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = poseidon.sbox(x)
	}
}

// BenchmarkMDSMatrixApplication benchmarks MDS matrix multiplication
func BenchmarkMDSMatrixApplication(b *testing.B) {
	poseidon, err := NewPoseidon(nil)
	if err != nil {
		b.Fatal(err)
	}

	state := make([]field.Element, poseidon.width)
	for i := 0; i < poseidon.width; i++ {
		state[i] = field.New(uint64(i + 1))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = poseidon.applyMDSMatrix(state)
	}
}

// BenchmarkPoseidonInitialization benchmarks hash initialization
func BenchmarkPoseidonInitialization(b *testing.B) {
	params := GetDefaultPoseidonParameters(128)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := NewPoseidon(params)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkPoseidonPermutation benchmarks the full permutation
func BenchmarkPoseidonPermutation(b *testing.B) {
	poseidon, err := NewPoseidon(nil)
	if err != nil {
		b.Fatal(err)
	}

	state := make([]field.Element, poseidon.width)
	for i := 0; i < poseidon.width; i++ {
		state[i] = field.New(uint64(i + 1))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = poseidon.poseidonPermutation(state)
	}
}

// BenchmarkPoseidonConvenience benchmarks the convenience function
func BenchmarkPoseidonConvenience(b *testing.B) {
	inputs := make([]field.Element, 10)
	for i := 0; i < 10; i++ {
		inputs[i] = field.New(uint64(i + 1))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = PoseidonHash(inputs)
	}
}
