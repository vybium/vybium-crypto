// Package zerofier provides zerofier tree implementation for polynomial optimization.
//
// A zerofier tree is a balanced binary tree of vanishing polynomials.
// Conceptually, every leaf corresponds to a single point, and the value of
// that leaf is the monic linear polynomial that evaluates to zero there and
// nowhere else. Every non-leaf node is the product of its two children.
//
// This is used for optimizing polynomial operations in STARK proofs,
// particularly for computing polynomials that vanish on specific sets of points.
package zerofier

import (
	"fmt"

	"github.com/vybium/vybium-crypto/pkg/vybium-crypto/field"
	"github.com/vybium/vybium-crypto/pkg/vybium-crypto/polynomial"
)

const (
	// RecursionCutoffThreshold regulates the depth at which the tree is truncated.
	// This controls the number of points contained by each leaf.
	RecursionCutoffThreshold = 16
)

// ZerofierTree represents a balanced binary tree of vanishing polynomials.
// Production implementation.
type ZerofierTree struct {
	// Type indicates the kind of node
	Type NodeType

	// For Leaf nodes
	Points   []field.Element
	Zerofier *polynomial.Polynomial

	// For Branch nodes
	Left  *ZerofierTree
	Right *ZerofierTree

	// For Padding nodes (no additional fields needed)
}

// NodeType represents the type of a zerofier tree node.
type NodeType int

const (
	// Leaf represents a leaf node containing points
	Leaf NodeType = iota
	// Branch represents an internal node with left and right children
	Branch
	// Padding represents a padding node
	Padding
)

func (nt NodeType) String() string {
	switch nt {
	case Leaf:
		return "Leaf"
	case Branch:
		return "Branch"
	case Padding:
		return "Padding"
	default:
		return "Unknown"
	}
}

// NewZerofierTree creates a new zerofier tree from a domain of points.
// Production implementation.
func NewZerofierTree(domain []field.Element) *ZerofierTree {
	if len(domain) == 0 {
		return &ZerofierTree{Type: Padding}
	}

	// Create leaf nodes by chunking the domain
	var nodes []*ZerofierTree
	for i := 0; i < len(domain); i += RecursionCutoffThreshold {
		end := i + RecursionCutoffThreshold
		if end > len(domain) {
			end = len(domain)
		}
		chunk := domain[i:end]
		leaf := NewLeaf(chunk)
		nodes = append(nodes, leaf)
	}

	// Pad to next power of two
	nextPowerOfTwo := nextPowerOfTwo(len(nodes))
	for len(nodes) < nextPowerOfTwo {
		nodes = append(nodes, &ZerofierTree{Type: Padding})
	}

	// Build tree bottom-up
	for len(nodes) > 1 {
		right := nodes[len(nodes)-1]
		left := nodes[len(nodes)-2]
		nodes = nodes[:len(nodes)-2]

		if left.Type == Padding {
			nodes = append([]*ZerofierTree{{Type: Padding}}, nodes...)
		} else {
			branch := NewBranch(left, right)
			nodes = append([]*ZerofierTree{branch}, nodes...)
		}
	}

	if len(nodes) == 0 {
		return &ZerofierTree{Type: Padding}
	}

	return nodes[0]
}

// NewLeaf creates a new leaf node with the given points.
// Production implementation.
func NewLeaf(points []field.Element) *ZerofierTree {
	// Create a copy of the points
	pointsCopy := make([]field.Element, len(points))
	copy(pointsCopy, points)

	// Compute the zerofier polynomial for these points
	zerofier := polynomial.Zerofier(pointsCopy)

	return &ZerofierTree{
		Type:     Leaf,
		Points:   pointsCopy,
		Zerofier: zerofier,
	}
}

// NewBranch creates a new branch node with left and right children.
// Production implementation.
func NewBranch(left, right *ZerofierTree) *ZerofierTree {
	// Compute the zerofier as the product of left and right zerofiers
	leftZerofier := left.GetZerofier()
	rightZerofier := right.GetZerofier()
	zerofier := leftZerofier.Mul(rightZerofier)

	return &ZerofierTree{
		Type:     Branch,
		Left:     left,
		Right:    right,
		Zerofier: zerofier,
	}
}

// GetZerofier returns the zerofier polynomial for this tree node.
// Production implementation.
func (zt *ZerofierTree) GetZerofier() *polynomial.Polynomial {
	switch zt.Type {
	case Leaf, Branch:
		return zt.Zerofier
	case Padding:
		return polynomial.One()
	default:
		panic(fmt.Sprintf("unknown node type: %v", zt.Type))
	}
}

// IsLeaf returns true if this is a leaf node.
func (zt *ZerofierTree) IsLeaf() bool {
	return zt.Type == Leaf
}

// IsBranch returns true if this is a branch node.
func (zt *ZerofierTree) IsBranch() bool {
	return zt.Type == Branch
}

// IsPadding returns true if this is a padding node.
func (zt *ZerofierTree) IsPadding() bool {
	return zt.Type == Padding
}

// GetPoints returns the points contained in this node (only valid for leaf nodes).
func (zt *ZerofierTree) GetPoints() []field.Element {
	if zt.Type != Leaf {
		return nil
	}
	return zt.Points
}

// GetLeft returns the left child (only valid for branch nodes).
func (zt *ZerofierTree) GetLeft() *ZerofierTree {
	if zt.Type != Branch {
		return nil
	}
	return zt.Left
}

// GetRight returns the right child (only valid for branch nodes).
func (zt *ZerofierTree) GetRight() *ZerofierTree {
	if zt.Type != Branch {
		return nil
	}
	return zt.Right
}

// Clone creates a deep copy of the zerofier tree.
func (zt *ZerofierTree) Clone() *ZerofierTree {
	if zt == nil {
		return nil
	}

	clone := &ZerofierTree{
		Type: zt.Type,
	}

	switch zt.Type {
	case Leaf:
		clone.Points = make([]field.Element, len(zt.Points))
		copy(clone.Points, zt.Points)
		clone.Zerofier = zt.Zerofier.Clone()
	case Branch:
		clone.Left = zt.Left.Clone()
		clone.Right = zt.Right.Clone()
		clone.Zerofier = zt.Zerofier.Clone()
	case Padding:
		// No additional fields for padding nodes
	}

	return clone
}

// Equal checks if two zerofier trees are equal.
func (zt *ZerofierTree) Equal(other *ZerofierTree) bool {
	if zt == nil && other == nil {
		return true
	}
	if zt == nil || other == nil {
		return false
	}

	if zt.Type != other.Type {
		return false
	}

	switch zt.Type {
	case Leaf:
		// Check points
		if len(zt.Points) != len(other.Points) {
			return false
		}
		for i, point := range zt.Points {
			if !point.Equal(other.Points[i]) {
				return false
			}
		}
		// Check zerofier
		return zt.Zerofier.Equal(other.Zerofier)
	case Branch:
		// Check children
		return zt.Left.Equal(other.Left) && zt.Right.Equal(other.Right)
	case Padding:
		return true
	default:
		return false
	}
}

// String returns a string representation of the zerofier tree.
func (zt *ZerofierTree) String() string {
	if zt == nil {
		return "nil"
	}

	switch zt.Type {
	case Leaf:
		return fmt.Sprintf("Leaf(points=%d)", len(zt.Points))
	case Branch:
		return fmt.Sprintf("Branch(left=%s, right=%s)", zt.Left.String(), zt.Right.String())
	case Padding:
		return "Padding"
	default:
		return "Unknown"
	}
}

// Depth returns the depth of the tree.
func (zt *ZerofierTree) Depth() int {
	if zt == nil {
		return 0
	}

	switch zt.Type {
	case Leaf, Padding:
		return 1
	case Branch:
		leftDepth := zt.Left.Depth()
		rightDepth := zt.Right.Depth()
		if leftDepth > rightDepth {
			return leftDepth + 1
		}
		return rightDepth + 1
	default:
		return 0
	}
}

// Size returns the total number of nodes in the tree.
func (zt *ZerofierTree) Size() int {
	if zt == nil {
		return 0
	}

	switch zt.Type {
	case Leaf, Padding:
		return 1
	case Branch:
		return 1 + zt.Left.Size() + zt.Right.Size()
	default:
		return 0
	}
}

// LeafCount returns the number of leaf nodes in the tree.
func (zt *ZerofierTree) LeafCount() int {
	if zt == nil {
		return 0
	}

	switch zt.Type {
	case Leaf:
		return 1
	case Branch:
		return zt.Left.LeafCount() + zt.Right.LeafCount()
	case Padding:
		return 0
	default:
		return 0
	}
}

// PointCount returns the total number of points in the tree.
func (zt *ZerofierTree) PointCount() int {
	if zt == nil {
		return 0
	}

	switch zt.Type {
	case Leaf:
		return len(zt.Points)
	case Branch:
		return zt.Left.PointCount() + zt.Right.PointCount()
	case Padding:
		return 0
	default:
		return 0
	}
}

// Validate checks if the zerofier tree is valid.
func (zt *ZerofierTree) Validate() error {
	if zt == nil {
		return fmt.Errorf("tree is nil")
	}

	switch zt.Type {
	case Leaf:
		if len(zt.Points) == 0 {
			return fmt.Errorf("leaf node has no points")
		}
		if zt.Zerofier == nil {
			return fmt.Errorf("leaf node has nil zerofier")
		}
		// Verify that the zerofier polynomial has the correct degree
		expectedDegree := len(zt.Points)
		if zt.Zerofier.Degree() != expectedDegree {
			return fmt.Errorf("leaf zerofier degree %d, expected %d", zt.Zerofier.Degree(), expectedDegree)
		}
	case Branch:
		if zt.Left == nil {
			return fmt.Errorf("branch node has nil left child")
		}
		if zt.Right == nil {
			return fmt.Errorf("branch node has nil right child")
		}
		if zt.Zerofier == nil {
			return fmt.Errorf("branch node has nil zerofier")
		}
		// Validate children
		if err := zt.Left.Validate(); err != nil {
			return fmt.Errorf("left child validation failed: %w", err)
		}
		if err := zt.Right.Validate(); err != nil {
			return fmt.Errorf("right child validation failed: %w", err)
		}
		// Verify that the zerofier is the product of children
		leftZerofier := zt.Left.GetZerofier()
		rightZerofier := zt.Right.GetZerofier()
		expectedZerofier := leftZerofier.Mul(rightZerofier)
		if !zt.Zerofier.Equal(expectedZerofier) {
			return fmt.Errorf("branch zerofier is not the product of children")
		}
	case Padding:
		// Padding nodes are always valid
	default:
		return fmt.Errorf("unknown node type: %v", zt.Type)
	}

	return nil
}

// Helper functions

// nextPowerOfTwo returns the next power of two greater than or equal to n.
func nextPowerOfTwo(n int) int {
	if n <= 0 {
		return 1
	}
	if n&(n-1) == 0 {
		return n
	}

	// Find the next power of two
	power := 1
	for power < n {
		power <<= 1
	}
	return power
}

// IsPowerOfTwo checks if n is a power of two.
func IsPowerOfTwo(n int) bool {
	return n > 0 && n&(n-1) == 0
}

// GetRecursionCutoffThreshold returns the recursion cutoff threshold.
func GetRecursionCutoffThreshold() int {
	return RecursionCutoffThreshold
}
