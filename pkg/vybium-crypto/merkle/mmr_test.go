package merkle

import (
	"testing"

	"github.com/vybium/vybium-crypto/pkg/vybium-crypto/field"
	"github.com/vybium/vybium-crypto/pkg/vybium-crypto/hash"
)

func TestMmrCreation(t *testing.T) {
	leafs := createTestLeafs(13)
	mmr := NewMmrAccumulatorFromLeafs(leafs)

	if mmr.IsEmpty() {
		t.Error("MMR should not be empty")
	}

	if mmr.NumLeafs() != 13 {
		t.Errorf("Expected 13 leafs, got %d", mmr.NumLeafs())
	}

	// Check consistency: number of peaks should equal number of 1-bits in leaf count
	// 13 = 0b1101 = 3 bits set, so we should have 3 peaks
	expectedPeaks := 3
	if len(mmr.Peaks()) != expectedPeaks {
		t.Errorf("Expected %d peaks, got %d", expectedPeaks, len(mmr.Peaks()))
	}

	if !mmr.IsConsistent() {
		t.Error("MMR should be self-consistent")
	}
}

func TestMmrEmptyCase(t *testing.T) {
	mmr := NewMmrAccumulatorFromLeafs([]hash.Digest{})

	if !mmr.IsEmpty() {
		t.Error("Empty MMR should be empty")
	}

	if mmr.NumLeafs() != 0 {
		t.Error("Empty MMR should have 0 leafs")
	}

	if len(mmr.Peaks()) != 0 {
		t.Error("Empty MMR should have 0 peaks")
	}
}

func TestMmrPowerOfTwo(t *testing.T) {
	// When leaf count is a power of 2, MMR should have exactly 1 peak
	tests := []int{1, 2, 4, 8, 16, 32, 64}

	for _, numLeafs := range tests {
		leafs := createTestLeafs(numLeafs)
		mmr := NewMmrAccumulatorFromLeafs(leafs)

		if len(mmr.Peaks()) != 1 {
			t.Errorf("MMR with %d leafs (power of 2) should have 1 peak, got %d",
				numLeafs, len(mmr.Peaks()))
		}
	}
}

func TestMmrPeaksFromLeafs(t *testing.T) {
	tests := []struct {
		numLeafs      int
		expectedPeaks int
	}{
		{1, 1},   // 0b1
		{2, 1},   // 0b10
		{3, 2},   // 0b11
		{4, 1},   // 0b100
		{5, 2},   // 0b101
		{6, 2},   // 0b110
		{7, 3},   // 0b111
		{8, 1},   // 0b1000
		{13, 3},  // 0b1101
		{15, 4},  // 0b1111
		{16, 1},  // 0b10000
		{31, 5},  // 0b11111
		{32, 1},  // 0b100000
		{100, 3}, // 0b1100100
	}

	for _, tt := range tests {
		leafs := createTestLeafs(tt.numLeafs)
		peaks := peaksFromLeafs(leafs)

		if len(peaks) != tt.expectedPeaks {
			t.Errorf("With %d leafs, expected %d peaks, got %d",
				tt.numLeafs, tt.expectedPeaks, len(peaks))
		}
	}
}

func TestMmrAppend(t *testing.T) {
	mmr := NewMmrAccumulatorFromLeafs([]hash.Digest{})

	// Append leafs one by one
	for i := 0; i < 10; i++ {
		newLeaf := hash.NewDigest([hash.DigestLen]field.Element{
			field.New(uint64(i)),
			field.New(uint64(i * 2)),
			field.New(uint64(i * 3)),
			field.New(uint64(i * 4)),
			field.New(uint64(i * 5)),
		})

		proof := mmr.Append(newLeaf)

		// Verify the proof
		if !mmr.VerifyMembership(newLeaf, proof) {
			t.Errorf("Membership proof failed for leaf %d", i)
		}

		// Check leaf count is updated
		expectedLeafCount := uint64(i + 1)
		if mmr.NumLeafs() != expectedLeafCount {
			t.Errorf("After appending leaf %d, expected %d leafs, got %d",
				i, expectedLeafCount, mmr.NumLeafs())
		}

		// Check consistency
		if !mmr.IsConsistent() {
			t.Errorf("MMR inconsistent after appending leaf %d", i)
		}
	}
}

func TestMmrBagPeaks(t *testing.T) {
	leafs := createTestLeafs(7)
	mmr := NewMmrAccumulatorFromLeafs(leafs)

	bag1 := mmr.BagPeaks()
	if bag1.IsZero() {
		t.Error("Bagged peaks should not be zero")
	}

	// Same MMR should produce same bag
	mmr2 := NewMmrAccumulatorFromLeafs(leafs)
	bag2 := mmr2.BagPeaks()
	if !bag1.Equal(bag2) {
		t.Error("Same MMR should produce same bagged peaks")
	}

	// Different MMR should produce different bag
	leafs2 := createTestLeafs(8)
	mmr3 := NewMmrAccumulatorFromLeafs(leafs2)
	bag3 := mmr3.BagPeaks()
	if bag1.Equal(bag3) {
		t.Error("Different MMRs should produce different bagged peaks")
	}
}

func TestMmrMembershipProof(t *testing.T) {
	// Build MMR with initial leafs
	initialLeafs := createTestLeafs(5)
	mmr := NewMmrAccumulatorFromLeafs(initialLeafs)

	// Append a new leaf and get its proof
	newLeaf := hash.NewDigest([hash.DigestLen]field.Element{
		field.New(999),
		field.New(888),
		field.New(777),
		field.New(666),
		field.New(555),
	})
	proof := mmr.Append(newLeaf)

	// Verify the proof
	if !mmr.VerifyMembership(newLeaf, proof) {
		t.Error("Valid membership proof should verify")
	}

	// Wrong leaf should fail
	wrongLeaf := hash.NewDigest([hash.DigestLen]field.Element{
		field.New(111),
		field.New(222),
		field.New(333),
		field.New(444),
		field.New(555),
	})
	if mmr.VerifyMembership(wrongLeaf, proof) {
		t.Error("Wrong leaf should not verify")
	}

	// Modified proof should fail
	modifiedProof := proof
	if len(modifiedProof.AuthPath) > 0 {
		modifiedProof.AuthPath[0] = wrongLeaf
		if mmr.VerifyMembership(newLeaf, modifiedProof) {
			t.Error("Modified proof should not verify")
		}
	}
}

func TestMmrConsistency(t *testing.T) {
	tests := []struct {
		name       string
		leafCount  uint64
		numPeaks   int
		consistent bool
	}{
		{"1 leaf, 1 peak", 1, 1, true},
		{"2 leafs, 1 peak", 2, 1, true},
		{"3 leafs, 2 peaks", 3, 2, true},
		{"7 leafs, 3 peaks", 7, 3, true},
		{"1 leaf, 2 peaks (inconsistent)", 1, 2, false},
		{"7 leafs, 2 peaks (inconsistent)", 7, 2, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			peaks := make([]hash.Digest, tt.numPeaks)
			for i := 0; i < tt.numPeaks; i++ {
				peaks[i] = hash.NewDigest([hash.DigestLen]field.Element{
					field.New(uint64(i)),
					field.Zero,
					field.Zero,
					field.Zero,
					field.Zero,
				})
			}

			mmr := NewMmrAccumulator(peaks, tt.leafCount)
			consistent := mmr.IsConsistent()

			if consistent != tt.consistent {
				t.Errorf("Expected consistency %v, got %v", tt.consistent, consistent)
			}
		})
	}
}

func TestMmrClone(t *testing.T) {
	leafs := createTestLeafs(5)
	mmr1 := NewMmrAccumulatorFromLeafs(leafs)
	mmr2 := mmr1.Clone()

	// Should be equal
	if mmr1.NumLeafs() != mmr2.NumLeafs() {
		t.Error("Cloned MMR should have same leaf count")
	}

	if len(mmr1.Peaks()) != len(mmr2.Peaks()) {
		t.Error("Cloned MMR should have same number of peaks")
	}

	// Modify clone
	newLeaf := hash.NewDigest([hash.DigestLen]field.Element{
		field.New(999),
		field.Zero,
		field.Zero,
		field.Zero,
		field.Zero,
	})
	mmr2.Append(newLeaf)

	// Original should be unchanged
	if mmr1.NumLeafs() == mmr2.NumLeafs() {
		t.Error("Modifying clone should not affect original")
	}
}

func TestMmrDeterminism(t *testing.T) {
	// Same leafs should produce same MMR
	leafs := createTestLeafs(13)

	mmr1 := NewMmrAccumulatorFromLeafs(leafs)
	mmr2 := NewMmrAccumulatorFromLeafs(leafs)

	if mmr1.NumLeafs() != mmr2.NumLeafs() {
		t.Error("Determinism: same leafs should produce same leaf count")
	}

	peaks1 := mmr1.Peaks()
	peaks2 := mmr2.Peaks()

	if len(peaks1) != len(peaks2) {
		t.Error("Determinism: same leafs should produce same number of peaks")
	}

	for i := range peaks1 {
		if !peaks1[i].Equal(peaks2[i]) {
			t.Errorf("Determinism: peak %d should match", i)
		}
	}

	bag1 := mmr1.BagPeaks()
	bag2 := mmr2.BagPeaks()
	if !bag1.Equal(bag2) {
		t.Error("Determinism: same leafs should produce same bagged peaks")
	}
}

func TestMmrSequentialAppends(t *testing.T) {
	// Compare sequential appends vs batch creation
	numLeafs := 20
	leafs := createTestLeafs(numLeafs)

	// Build MMR by appending one at a time
	mmr1 := NewMmrAccumulatorFromLeafs([]hash.Digest{})
	for _, leaf := range leafs {
		mmr1.Append(leaf)
	}

	// Build MMR from all leafs at once
	mmr2 := NewMmrAccumulatorFromLeafs(leafs)

	// Should produce same result
	if mmr1.NumLeafs() != mmr2.NumLeafs() {
		t.Error("Sequential and batch should produce same leaf count")
	}

	peaks1 := mmr1.Peaks()
	peaks2 := mmr2.Peaks()

	if len(peaks1) != len(peaks2) {
		t.Error("Sequential and batch should produce same number of peaks")
	}

	for i := range peaks1 {
		if !peaks1[i].Equal(peaks2[i]) {
			t.Errorf("Sequential and batch peak %d should match", i)
		}
	}
}

// Benchmarks
func BenchmarkMmrCreation16(b *testing.B) {
	leafs := createTestLeafs(16)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewMmrAccumulatorFromLeafs(leafs)
	}
}

func BenchmarkMmrCreation256(b *testing.B) {
	leafs := createTestLeafs(256)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewMmrAccumulatorFromLeafs(leafs)
	}
}

func BenchmarkMmrAppend(b *testing.B) {
	mmr := NewMmrAccumulatorFromLeafs(createTestLeafs(100))
	newLeaf := hash.NewDigest([hash.DigestLen]field.Element{
		field.New(999),
		field.Zero,
		field.Zero,
		field.Zero,
		field.Zero,
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mmrCopy := mmr.Clone()
		_ = mmrCopy.Append(newLeaf)
	}
}

func BenchmarkMmrBagPeaks(b *testing.B) {
	mmr := NewMmrAccumulatorFromLeafs(createTestLeafs(100))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mmr.BagPeaks()
	}
}

func BenchmarkMmrVerifyMembership(b *testing.B) {
	mmr := NewMmrAccumulatorFromLeafs(createTestLeafs(100))
	newLeaf := hash.NewDigest([hash.DigestLen]field.Element{
		field.New(999),
		field.Zero,
		field.Zero,
		field.Zero,
		field.Zero,
	})
	proof := mmr.Append(newLeaf)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mmr.VerifyMembership(newLeaf, proof)
	}
}
