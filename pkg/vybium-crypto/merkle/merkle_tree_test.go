package merkle

import (
	"testing"

	"github.com/vybium/vybium-crypto/pkg/vybium-crypto/field"
	"github.com/vybium/vybium-crypto/pkg/vybium-crypto/hash"
)

// createTestLeafs creates test leafs for Merkle tree testing.
func createTestLeafs(count int) []hash.Digest {
	leafs := make([]hash.Digest, count)
	for i := 0; i < count; i++ {
		// Create a digest with varying values
		leafs[i] = hash.NewDigest([hash.DigestLen]field.Element{
			field.New(uint64(i)),
			field.New(uint64(i * 2)),
			field.New(uint64(i * 3)),
			field.New(uint64(i * 4)),
			field.New(uint64(i * 5)),
		})
	}
	return leafs
}

func TestMerkleTreeCreation(t *testing.T) {
	tests := []struct {
		name      string
		numLeafs  int
		shouldErr bool
	}{
		{"2 leafs", 2, false},
		{"4 leafs", 4, false},
		{"8 leafs", 8, false},
		{"16 leafs", 16, false},
		{"1 leaf", 1, false},
		{"0 leafs", 0, true},
		{"3 leafs (not power of 2)", 3, true},
		{"5 leafs (not power of 2)", 5, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			leafs := createTestLeafs(tt.numLeafs)
			tree, err := New(leafs)

			if tt.shouldErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if tree == nil {
					t.Error("Expected tree but got nil")
				}
			}
		})
	}
}

func TestMerkleTreeRoot(t *testing.T) {
	leafs := createTestLeafs(4)
	tree, err := New(leafs)
	if err != nil {
		t.Fatalf("Failed to create tree: %v", err)
	}

	root := tree.Root()
	if root.IsZero() {
		t.Error("Root should not be zero")
	}

	// Create another tree with same leafs - should have same root
	tree2, _ := New(leafs)
	root2 := tree2.Root()
	if !root.Equal(root2) {
		t.Error("Same leafs should produce same root")
	}

	// Different leafs should produce different root
	leafs2 := createTestLeafs(4)
	leafs2[0] = hash.NewDigest([hash.DigestLen]field.Element{
		field.New(999),
		field.New(999),
		field.New(999),
		field.New(999),
		field.New(999),
	})
	tree3, _ := New(leafs2)
	root3 := tree3.Root()
	if root.Equal(root3) {
		t.Error("Different leafs should produce different roots")
	}
}

func TestMerkleTreeHeight(t *testing.T) {
	tests := []struct {
		numLeafs       int
		expectedHeight uint32
	}{
		{1, 0},
		{2, 1},
		{4, 2},
		{8, 3},
		{16, 4},
		{32, 5},
		{64, 6},
		{128, 7},
	}

	for _, tt := range tests {
		leafs := createTestLeafs(tt.numLeafs)
		tree, err := New(leafs)
		if err != nil {
			t.Fatalf("Failed to create tree with %d leafs: %v", tt.numLeafs, err)
		}

		height := tree.Height()
		if height != tt.expectedHeight {
			t.Errorf("Expected height %d for %d leafs, got %d",
				tt.expectedHeight, tt.numLeafs, height)
		}
	}
}

func TestMerkleTreeGetLeaf(t *testing.T) {
	leafs := createTestLeafs(8)
	tree, err := New(leafs)
	if err != nil {
		t.Fatalf("Failed to create tree: %v", err)
	}

	// Test valid indices
	for i := uint64(0); i < 8; i++ {
		leaf, err := tree.GetLeaf(i)
		if err != nil {
			t.Errorf("Failed to get leaf %d: %v", i, err)
		}
		if !leaf.Equal(leafs[i]) {
			t.Errorf("Leaf %d doesn't match original", i)
		}
	}

	// Test invalid index
	_, err = tree.GetLeaf(8)
	if err == nil {
		t.Error("Expected error for out-of-range index")
	}
}

func TestMerkleTreeAuthenticationPath(t *testing.T) {
	leafs := createTestLeafs(8)
	tree, err := New(leafs)
	if err != nil {
		t.Fatalf("Failed to create tree: %v", err)
	}

	root := tree.Root()
	height := tree.Height()

	// Test authentication paths for all leafs
	for i := uint64(0); i < 8; i++ {
		authPath, err := tree.AuthenticationPath(i)
		if err != nil {
			t.Errorf("Failed to get auth path for leaf %d: %v", i, err)
		}

		if len(authPath) != int(height) {
			t.Errorf("Auth path length should be %d, got %d", height, len(authPath))
		}

		// Verify the authentication path
		leaf, _ := tree.GetLeaf(i)
		if !VerifyInclusionProof(root, i, leaf, authPath) {
			t.Errorf("Auth path verification failed for leaf %d", i)
		}
	}
}

func TestVerifyInclusionProof(t *testing.T) {
	leafs := createTestLeafs(4)
	tree, err := New(leafs)
	if err != nil {
		t.Fatalf("Failed to create tree: %v", err)
	}

	root := tree.Root()

	// Test valid proof
	leaf, _ := tree.GetLeaf(2)
	authPath, _ := tree.AuthenticationPath(2)
	if !VerifyInclusionProof(root, 2, leaf, authPath) {
		t.Error("Valid proof should verify")
	}

	// Test invalid proof (wrong leaf)
	wrongLeaf := hash.NewDigest([hash.DigestLen]field.Element{
		field.New(999),
		field.New(999),
		field.New(999),
		field.New(999),
		field.New(999),
	})
	if VerifyInclusionProof(root, 2, wrongLeaf, authPath) {
		t.Error("Invalid proof (wrong leaf) should not verify")
	}

	// Test invalid proof (wrong index)
	if VerifyInclusionProof(root, 1, leaf, authPath) {
		t.Error("Invalid proof (wrong index) should not verify")
	}

	// Test invalid proof (wrong root)
	wrongRoot := hash.NewDigest([hash.DigestLen]field.Element{
		field.New(123),
		field.New(456),
		field.New(789),
		field.New(101112),
		field.New(131415),
	})
	if VerifyInclusionProof(wrongRoot, 2, leaf, authPath) {
		t.Error("Invalid proof (wrong root) should not verify")
	}
}

func TestMerkleTreeInclusionProof(t *testing.T) {
	leafs := createTestLeafs(8)
	tree, err := New(leafs)
	if err != nil {
		t.Fatalf("Failed to create tree: %v", err)
	}

	root := tree.Root()

	// Test single leaf proof
	t.Run("single leaf", func(t *testing.T) {
		proof, err := tree.NewInclusionProof([]MerkleTreeLeafIndex{3})
		if err != nil {
			t.Fatalf("Failed to create inclusion proof: %v", err)
		}

		if !proof.Verify(root) {
			t.Error("Single leaf proof should verify")
		}

		if proof.TreeHeight != tree.Height() {
			t.Error("Proof height should match tree height")
		}

		if len(proof.IndexedLeafs) != 1 {
			t.Error("Should have exactly 1 indexed leaf")
		}
	})

	// Test multiple leafs proof
	t.Run("multiple leafs", func(t *testing.T) {
		indices := []MerkleTreeLeafIndex{0, 2, 5}
		proof, err := tree.NewInclusionProof(indices)
		if err != nil {
			t.Fatalf("Failed to create inclusion proof: %v", err)
		}

		if !proof.Verify(root) {
			t.Error("Multiple leafs proof should verify")
		}

		if len(proof.IndexedLeafs) != len(indices) {
			t.Error("Number of indexed leafs should match")
		}

		// Verify each indexed leaf is correct
		for i, pair := range proof.IndexedLeafs {
			expectedLeaf, _ := tree.GetLeaf(indices[i])
			if !pair.Digest.Equal(expectedLeaf) {
				t.Errorf("Indexed leaf %d doesn't match", i)
			}
		}
	})

	// Test out-of-range index
	t.Run("out of range index", func(t *testing.T) {
		_, err := tree.NewInclusionProof([]MerkleTreeLeafIndex{10})
		if err == nil {
			t.Error("Expected error for out-of-range index")
		}
	})
}

func TestMerkleTreeDeterminism(t *testing.T) {
	// Same leafs should always produce same tree
	leafs := createTestLeafs(16)

	tree1, _ := New(leafs)
	tree2, _ := New(leafs)

	root1 := tree1.Root()
	root2 := tree2.Root()

	if !root1.Equal(root2) {
		t.Error("Determinism: same leafs should produce same root")
	}

	// Check all internal nodes match
	for i := MerkleTreeNodeIndex(1); i < MerkleTreeNodeIndex(len(tree1.nodes)); i++ {
		if !tree1.nodes[i].Equal(tree2.nodes[i]) {
			t.Errorf("Determinism: node %d should match", i)
		}
	}
}

// Benchmark Merkle tree operations
func BenchmarkMerkleTreeCreation16(b *testing.B) {
	leafs := createTestLeafs(16)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = New(leafs)
	}
}

func BenchmarkMerkleTreeCreation256(b *testing.B) {
	leafs := createTestLeafs(256)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = New(leafs)
	}
}

func BenchmarkMerkleTreeAuthPath(b *testing.B) {
	leafs := createTestLeafs(256)
	tree, _ := New(leafs)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = tree.AuthenticationPath(100)
	}
}

func BenchmarkVerifyInclusionProof(b *testing.B) {
	leafs := createTestLeafs(256)
	tree, _ := New(leafs)
	root := tree.Root()
	leaf, _ := tree.GetLeaf(100)
	authPath, _ := tree.AuthenticationPath(100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = VerifyInclusionProof(root, 100, leaf, authPath)
	}
}

func BenchmarkMerkleInclusionProof(b *testing.B) {
	leafs := createTestLeafs(256)
	tree, _ := New(leafs)
	indices := []MerkleTreeLeafIndex{10, 50, 100, 150, 200}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = tree.NewInclusionProof(indices)
	}
}
