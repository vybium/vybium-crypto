// Merkle tree examples demonstrating authentication structures
package main

import (
	"fmt"
	"log"

	"github.com/vybium/vybium-crypto/pkg/vybium-crypto/field"
	"github.com/vybium/vybium-crypto/pkg/vybium-crypto/hash"
	"github.com/vybium/vybium-crypto/pkg/vybium-crypto/merkle"
)

func main() {
	fmt.Println("🌳 Vybium Crypto Merkle Tree Examples")
	fmt.Println("====================================")

	// Example 1: Basic Merkle Tree
	fmt.Println("\n1. Basic Merkle Tree:")
	demonstrateBasicMerkleTree()

	// Example 2: Merkle Tree with Proofs
	fmt.Println("\n2. Merkle Tree with Proofs:")
	demonstrateMerkleProofs()

	// Example 3: MMR (Merkle Mountain Range)
	fmt.Println("\n3. Merkle Mountain Range:")
	demonstrateMMR()
}

func demonstrateBasicMerkleTree() {
	// Create some test data
	fieldLeaves := []field.Element{
		field.New(1),
		field.New(2),
		field.New(3),
		field.New(4),
	}

	// Convert to digests
	leaves := make([]hash.Digest, len(fieldLeaves))
	for i, elem := range fieldLeaves {
		leaves[i] = hash.NewDigest([5]field.Element{elem, field.Zero, field.Zero, field.Zero, field.Zero})
	}

	// Build Merkle tree
	tree, err := merkle.New(leaves)
	if err != nil {
		log.Fatalf("Error creating Merkle tree: %v", err)
	}
	root := tree.Root()

	fmt.Printf("   Leaves: %v\n", fieldLeaves)
	fmt.Printf("   Merkle Root: %v\n", root)
	fmt.Printf("   Tree Height: %d\n", tree.Height())
	fmt.Printf("   Leaf Count: %d\n", tree.NumLeafs())
}

func demonstrateMerkleProofs() {
	// Create test data
	fieldLeaves := []field.Element{
		field.New(100),
		field.New(200),
		field.New(300),
		field.New(400),
		field.New(500),
	}

	// Convert to digests
	leaves := make([]hash.Digest, len(fieldLeaves))
	for i, elem := range fieldLeaves {
		leaves[i] = hash.NewDigest([5]field.Element{elem, field.Zero, field.Zero, field.Zero, field.Zero})
	}

	// Build Merkle tree
	tree, err := merkle.New(leaves)
	if err != nil {
		log.Fatalf("Error creating Merkle tree: %v", err)
	}
	root := tree.Root()

	fmt.Printf("   Leaves: %v\n", fieldLeaves)
	fmt.Printf("   Merkle Root: %v\n", root)

	// Generate proof for leaf at index 2
	index := uint64(2)
	proof, err := tree.NewInclusionProof([]uint64{index})
	if err != nil {
		log.Fatalf("Error generating proof: %v", err)
	}

	fmt.Printf("   Proof for leaf %d: %v\n", index, proof)
	fmt.Printf("   Proof length: %d\n", len(proof.AuthenticationStructure))

	// Verify the proof
	isValid := proof.Verify(root)
	fmt.Printf("   Proof verification: %t\n", isValid)
}

func demonstrateMMR() {
	// Create MMR accumulator
	mmr := merkle.NewMmrAccumulator([]hash.Digest{}, 0)

	// Add some elements
	fieldElements := []field.Element{
		field.New(1000),
		field.New(2000),
		field.New(3000),
		field.New(4000),
		field.New(5000),
	}

	// Convert to digests and add to MMR, capturing proofs
	proofs := make([]merkle.MmrMembershipProof, len(fieldElements))
	for i, elem := range fieldElements {
		digest := hash.NewDigest([5]field.Element{elem, field.Zero, field.Zero, field.Zero, field.Zero})
		proof := mmr.Append(digest)
		proofs[i] = proof
		fmt.Printf("   Added element %d: %v\n", i+1, elem)
	}

	// Get the root (bag of peaks)
	root := mmr.BagPeaks()
	fmt.Printf("   MMR Root: %v\n", root)

	// Demonstrate proof for element at index 2
	index := uint64(2)
	proof := proofs[index]
	fmt.Printf("   MMR Proof for index %d: %v\n", index, proof)
	fmt.Printf("   MMR Proof length: %d\n", len(proof.AuthPath))

	// Verify the proof
	element := hash.NewDigest([5]field.Element{fieldElements[index], field.Zero, field.Zero, field.Zero, field.Zero})
	isValid := mmr.VerifyMembership(element, proof)
	fmt.Printf("   MMR Proof verification: %t\n", isValid)
}
