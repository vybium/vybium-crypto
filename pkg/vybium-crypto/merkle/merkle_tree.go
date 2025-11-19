package merkle

import (
	"fmt"
	"math/bits"

	"github.com/vybium/vybium-crypto/pkg/vybium-crypto/hash"
)

// MerkleTreeNodeIndex indexes internal nodes of a MerkleTree.
// Convention:
//   - Nothing lives at index 0
//   - Index 1 points to the root
//   - Indices 2 and 3 contain the two children of the root
//   - Indices 4 and 5 contain the two children of node 2
//   - And so on...
type MerkleTreeNodeIndex = uint64

// MerkleTreeLeafIndex indexes the leafs of a Merkle tree, left to right,
// starting with zero.
type MerkleTreeLeafIndex = uint64

// MerkleTreeHeight counts the number of layers in the Merkle tree,
// not including the root.
type MerkleTreeHeight = uint32

const (
	// RootIndex is the index of the root node.
	RootIndex MerkleTreeNodeIndex = 1
)

// MerkleTree is a binary tree of digests used to efficiently prove
// the inclusion of items in a set.
// The tree can hold at most 2^62 leafs (height up to 62).
// The hash function used is Tip5.
type MerkleTree struct {
	nodes []hash.Digest
}

// New builds a MerkleTree with the given leafs.
// Returns an error if:
// - the number of leafs is zero
// - the number of leafs is not a power of two
func New(leafs []hash.Digest) (*MerkleTree, error) {
	nodes, err := initializeMerkleTreeNodes(leafs)
	if err != nil {
		return nil, err
	}

	numRemainingNodes := len(leafs)
	return sequentiallyFillTree(nodes, numRemainingNodes)
}

// initializeMerkleTreeNodes validates the input and initializes the node array.
func initializeMerkleTreeNodes(leafs []hash.Digest) ([]hash.Digest, error) {
	numLeafs := len(leafs)

	if numLeafs == 0 {
		return nil, fmt.Errorf("cannot create Merkle tree with zero leafs")
	}

	if !isPowerOfTwo(uint32(numLeafs)) {
		return nil, fmt.Errorf("number of leafs must be a power of two, got %d", numLeafs)
	}

	// Tree needs space for: 1 (unused) + numLeafs (internal nodes) + numLeafs (leafs) - 1
	// = 2 * numLeafs
	nodes := make([]hash.Digest, 2*numLeafs)

	// Copy leafs to the second half of the array
	copy(nodes[numLeafs:], leafs)

	return nodes, nil
}

// sequentiallyFillTree fills the tree by hashing pairs of nodes bottom-up.
func sequentiallyFillTree(nodes []hash.Digest, numRemainingNodes int) (*MerkleTree, error) {
	for numRemainingNodes > 1 {
		for i := 0; i < numRemainingNodes; i += 2 {
			left := nodes[numRemainingNodes+i]
			right := nodes[numRemainingNodes+i+1]
			nodes[numRemainingNodes/2+i/2] = hash.HashPair(left, right)
		}
		numRemainingNodes /= 2
	}

	return &MerkleTree{nodes: nodes}, nil
}

// Root returns the root of the Merkle tree.
func (mt *MerkleTree) Root() hash.Digest {
	if len(mt.nodes) == 0 {
		return hash.ZeroDigest()
	}
	return mt.nodes[RootIndex]
}

// Height returns the height of the Merkle tree.
func (mt *MerkleTree) Height() MerkleTreeHeight {
	if len(mt.nodes) <= 1 {
		return 0
	}
	// nodes = 2 * numLeafs, so numLeafs = len(nodes) / 2
	// height = log2(numLeafs)
	numLeafs := len(mt.nodes) / 2
	return uint32(bits.Len(uint(numLeafs)) - 1)
}

// NumLeafs returns the number of leafs in the tree.
func (mt *MerkleTree) NumLeafs() uint64 {
	if len(mt.nodes) <= 1 {
		return 0
	}
	return uint64(len(mt.nodes) / 2)
}

// Size returns the total number of nodes in the tree (including internal nodes and leafs).
// Production implementation.
// The tree has 2 * numLeafs nodes total (1 unused + numLeafs internal + numLeafs leafs - 1).
func (mt *MerkleTree) Size() int {
	if len(mt.nodes) <= 1 {
		return 0
	}
	return len(mt.nodes)
}

// GetLeaf returns the leaf at the specified index.
func (mt *MerkleTree) GetLeaf(index MerkleTreeLeafIndex) (hash.Digest, error) {
	numLeafs := mt.NumLeafs()
	if index >= numLeafs {
		return hash.ZeroDigest(), fmt.Errorf("leaf index %d out of range [0, %d)", index, numLeafs)
	}

	// Leafs are stored in the second half of the nodes array
	leafNodeIndex := numLeafs + index
	return mt.nodes[leafNodeIndex], nil
}

// GetNode returns the node at the specified node index.
func (mt *MerkleTree) GetNode(nodeIndex MerkleTreeNodeIndex) (hash.Digest, error) {
	if nodeIndex >= uint64(len(mt.nodes)) || nodeIndex == 0 {
		return hash.ZeroDigest(), fmt.Errorf("node index %d out of range [1, %d)", nodeIndex, len(mt.nodes))
	}
	return mt.nodes[nodeIndex], nil
}

// AuthenticationPath returns the authentication path (also called Merkle proof or witness)
// for the leaf at the specified index.
// The authentication path is the list of sibling hashes needed to recompute the root.
func (mt *MerkleTree) AuthenticationPath(leafIndex MerkleTreeLeafIndex) ([]hash.Digest, error) {
	numLeafs := mt.NumLeafs()
	if leafIndex >= numLeafs {
		return nil, fmt.Errorf("leaf index %d out of range [0, %d)", leafIndex, numLeafs)
	}

	height := mt.Height()
	path := make([]hash.Digest, height)

	// Start at the leaf node
	nodeIndex := numLeafs + leafIndex

	// Walk up the tree, collecting sibling hashes
	for i := uint32(0); i < height; i++ {
		// Get sibling index (flip the least significant bit)
		siblingIndex := nodeIndex ^ 1
		path[i] = mt.nodes[siblingIndex]

		// Move to parent
		nodeIndex /= 2
	}

	return path, nil
}

// VerifyInclusionProof verifies that a leaf with the given digest is at the specified
// index in a Merkle tree with the given root, using the provided authentication path.
func VerifyInclusionProof(root hash.Digest, leafIndex MerkleTreeLeafIndex, leaf hash.Digest, authPath []hash.Digest) bool {
	// Recompute the root by hashing up the tree
	currentHash := leaf
	currentIndex := leafIndex

	for _, siblingHash := range authPath {
		// If currentIndex is even, current is left child; otherwise, right child
		if currentIndex%2 == 0 {
			currentHash = hash.HashPair(currentHash, siblingHash)
		} else {
			currentHash = hash.HashPair(siblingHash, currentHash)
		}
		currentIndex /= 2
	}

	return currentHash.Equal(root)
}

// MerkleTreeInclusionProof is a full inclusion proof for multiple leafs.
type MerkleTreeInclusionProof struct {
	// TreeHeight is the stated height of the Merkle tree this proof is relative to.
	TreeHeight MerkleTreeHeight

	// IndexedLeafs contains the leafs the proof is about (leaf index -> digest).
	IndexedLeafs []LeafIndexDigestPair

	// AuthenticationStructure is the proof's witness: de-duplicated authentication
	// structure for the leafs this proof is about.
	AuthenticationStructure []hash.Digest
}

// LeafIndexDigestPair represents a leaf index and its digest.
type LeafIndexDigestPair struct {
	Index  MerkleTreeLeafIndex
	Digest hash.Digest
}

// NewInclusionProof creates an inclusion proof for the specified leaf indices.
func (mt *MerkleTree) NewInclusionProof(leafIndices []MerkleTreeLeafIndex) (*MerkleTreeInclusionProof, error) {
	numLeafs := mt.NumLeafs()
	for _, idx := range leafIndices {
		if idx >= numLeafs {
			return nil, fmt.Errorf("leaf index %d out of range [0, %d)", idx, numLeafs)
		}
	}

	// Build indexed leafs
	indexedLeafs := make([]LeafIndexDigestPair, len(leafIndices))
	for i, idx := range leafIndices {
		leaf, _ := mt.GetLeaf(idx)
		indexedLeafs[i] = LeafIndexDigestPair{Index: idx, Digest: leaf}
	}

	// Build authentication structure (de-duplicated)
	authStructure := mt.buildAuthenticationStructure(leafIndices)

	return &MerkleTreeInclusionProof{
		TreeHeight:              mt.Height(),
		IndexedLeafs:            indexedLeafs,
		AuthenticationStructure: authStructure,
	}, nil
}

// buildAuthenticationStructure builds the de-duplicated authentication structure
// for the given leaf indices.
func (mt *MerkleTree) buildAuthenticationStructure(leafIndices []MerkleTreeLeafIndex) []hash.Digest {
	numLeafs := mt.NumLeafs()
	height := mt.Height()

	// Track which nodes are revealed (either as leafs or in the authentication path)
	revealed := make(map[MerkleTreeNodeIndex]bool)

	// Mark all revealed leafs
	for _, idx := range leafIndices {
		nodeIndex := numLeafs + idx
		revealed[nodeIndex] = true
	}

	// Collect authentication path nodes
	var authNodes []hash.Digest

	// For each revealed leaf, walk up and collect siblings not already revealed
	for _, leafIdx := range leafIndices {
		nodeIndex := numLeafs + leafIdx

		for level := uint32(0); level < height; level++ {
			siblingIndex := nodeIndex ^ 1
			parentIndex := nodeIndex / 2

			// If sibling is not revealed and parent is needed
			if !revealed[siblingIndex] {
				authNodes = append(authNodes, mt.nodes[siblingIndex])
				revealed[siblingIndex] = true
			}

			// Mark parent as revealed
			revealed[parentIndex] = true
			nodeIndex = parentIndex
		}
	}

	return authNodes
}

// Verify verifies the inclusion proof.
func (proof *MerkleTreeInclusionProof) Verify(root hash.Digest) bool {
	if len(proof.IndexedLeafs) == 0 {
		return false
	}

	// Build partial tree from the proof
	partialTree := newPartialMerkleTree(proof.TreeHeight, proof.IndexedLeafs, proof.AuthenticationStructure)

	// Compute root from partial tree
	computedRoot := partialTree.computeRoot()

	return computedRoot.Equal(root)
}

// partialMerkleTree is a helper for verifying inclusion proofs.
type partialMerkleTree struct {
	treeHeight  MerkleTreeHeight
	leafIndices []MerkleTreeLeafIndex
	nodes       map[MerkleTreeNodeIndex]hash.Digest
}

// newPartialMerkleTree creates a partial Merkle tree from the proof data.
func newPartialMerkleTree(height MerkleTreeHeight, indexedLeafs []LeafIndexDigestPair, authStructure []hash.Digest) *partialMerkleTree {
	nodes := make(map[MerkleTreeNodeIndex]hash.Digest)
	leafIndices := make([]MerkleTreeLeafIndex, len(indexedLeafs))

	numLeafs := uint64(1) << height

	// Add leafs
	for i, pair := range indexedLeafs {
		nodeIndex := numLeafs + pair.Index
		nodes[nodeIndex] = pair.Digest
		leafIndices[i] = pair.Index
	}

	// Add authentication structure nodes
	authIdx := 0
	for _, leafIdx := range leafIndices {
		nodeIndex := numLeafs + leafIdx

		for level := uint32(0); level < height; level++ {
			siblingIndex := nodeIndex ^ 1

			// If sibling is not in nodes, take it from authStructure
			if _, exists := nodes[siblingIndex]; !exists && authIdx < len(authStructure) {
				nodes[siblingIndex] = authStructure[authIdx]
				authIdx++
			}

			nodeIndex /= 2
		}
	}

	return &partialMerkleTree{
		treeHeight:  height,
		leafIndices: leafIndices,
		nodes:       nodes,
	}
}

// computeRoot computes the root from the partial tree.
func (pt *partialMerkleTree) computeRoot() hash.Digest {
	numLeafs := uint64(1) << pt.treeHeight

	// Build up the tree level by level
	for level := pt.treeHeight; level > 0; level-- {
		levelStart := uint64(1) << level

		for nodeIdx := levelStart; nodeIdx < 2*levelStart; nodeIdx += 2 {
			left, leftExists := pt.nodes[nodeIdx]
			right, rightExists := pt.nodes[nodeIdx+1]

			if leftExists && rightExists {
				parentIdx := nodeIdx / 2
				pt.nodes[parentIdx] = hash.HashPair(left, right)
			}
		}
	}

	// Return root
	if root, exists := pt.nodes[RootIndex]; exists {
		return root
	}

	// Fallback: reconstruct from leafs
	for _, leafIdx := range pt.leafIndices {
		nodeIndex := numLeafs + leafIdx
		currentHash := pt.nodes[nodeIndex]
		currentIndex := leafIdx

		for level := uint32(0); level < pt.treeHeight; level++ {
			siblingIndex := nodeIndex ^ 1
			siblingHash, exists := pt.nodes[siblingIndex]
			if !exists {
				return hash.ZeroDigest()
			}

			if currentIndex%2 == 0 {
				currentHash = hash.HashPair(currentHash, siblingHash)
			} else {
				currentHash = hash.HashPair(siblingHash, currentHash)
			}

			nodeIndex /= 2
			currentIndex /= 2
		}

		return currentHash
	}

	return hash.ZeroDigest()
}

// isPowerOfTwo checks if a number is a power of two.
func isPowerOfTwo(n uint32) bool {
	return n > 0 && (n&(n-1) == 0)
}
