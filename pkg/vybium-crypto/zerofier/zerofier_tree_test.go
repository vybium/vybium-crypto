package zerofier

import (
	"testing"

	"github.com/vybium/vybium-crypto/pkg/vybium-crypto/field"
)

func TestNewZerofierTree(t *testing.T) {
	tests := []struct {
		name   string
		domain []field.Element
	}{
		{
			name:   "Empty domain",
			domain: []field.Element{},
		},
		{
			name:   "Single point",
			domain: []field.Element{field.One},
		},
		{
			name:   "Two points",
			domain: []field.Element{field.One, field.New(2)},
		},
		{
			name:   "Small domain",
			domain: []field.Element{field.One, field.New(2), field.New(3), field.New(4)},
		},
		{
			name:   "Medium domain",
			domain: []field.Element{field.New(1), field.New(2), field.New(3), field.New(4), field.New(5), field.New(6), field.New(7), field.New(8)},
		},
		{
			name:   "Large domain",
			domain: make([]field.Element, 50),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize test data
			for i := range tt.domain {
				tt.domain[i] = field.New(uint64(i + 1))
			}

			tree := NewZerofierTree(tt.domain)

			// Verify tree is not nil
			if tree == nil {
				t.Error("NewZerofierTree() returned nil")
			}

			// Verify tree is valid
			if err := tree.Validate(); err != nil {
				t.Errorf("NewZerofierTree() validation failed: %v", err)
			}

			// Verify point count
			expectedPoints := len(tt.domain)
			actualPoints := tree.PointCount()
			if actualPoints != expectedPoints {
				t.Errorf("PointCount() = %d, want %d", actualPoints, expectedPoints)
			}
		})
	}
}

func TestNewLeaf(t *testing.T) {
	tests := []struct {
		name   string
		points []field.Element
	}{
		{
			name:   "Single point",
			points: []field.Element{field.One},
		},
		{
			name:   "Multiple points",
			points: []field.Element{field.One, field.New(2), field.New(3)},
		},
		{
			name:   "Many points",
			points: make([]field.Element, 20),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize test data
			for i := range tt.points {
				tt.points[i] = field.New(uint64(i + 1))
			}

			leaf := NewLeaf(tt.points)

			// Verify it's a leaf node
			if !leaf.IsLeaf() {
				t.Error("NewLeaf() did not create a leaf node")
			}

			// Verify points
			actualPoints := leaf.GetPoints()
			if len(actualPoints) != len(tt.points) {
				t.Errorf("GetPoints() length = %d, want %d", len(actualPoints), len(tt.points))
			}

			for i, point := range actualPoints {
				if !point.Equal(tt.points[i]) {
					t.Errorf("GetPoints()[%d] = %v, want %v", i, point, tt.points[i])
				}
			}

			// Verify zerofier is not nil
			zerofier := leaf.GetZerofier()
			if zerofier == nil {
				t.Error("GetZerofier() returned nil")
			}

			// Verify zerofier degree
			expectedDegree := len(tt.points)
			if zerofier.Degree() != expectedDegree {
				t.Errorf("Zerofier degree = %d, want %d", zerofier.Degree(), expectedDegree)
			}
		})
	}
}

func TestNewBranch(t *testing.T) {
	// Create two leaf nodes
	leftPoints := []field.Element{field.One, field.New(2)}
	rightPoints := []field.Element{field.New(3), field.New(4)}
	left := NewLeaf(leftPoints)
	right := NewLeaf(rightPoints)

	// Create branch
	branch := NewBranch(left, right)

	// Verify it's a branch node
	if !branch.IsBranch() {
		t.Error("NewBranch() did not create a branch node")
	}

	// Verify children
	if branch.GetLeft() != left {
		t.Error("GetLeft() returned wrong node")
	}
	if branch.GetRight() != right {
		t.Error("GetRight() returned wrong node")
	}

	// Verify zerofier is not nil
	zerofier := branch.GetZerofier()
	if zerofier == nil {
		t.Error("GetZerofier() returned nil")
	}

	// Verify zerofier is the product of children
	leftZerofier := left.GetZerofier()
	rightZerofier := right.GetZerofier()
	expectedZerofier := leftZerofier.Mul(rightZerofier)
	if !zerofier.Equal(expectedZerofier) {
		t.Error("Branch zerofier is not the product of children")
	}
}

func TestZerofierTreeTypes(t *testing.T) {
	// Test leaf node
	leaf := NewLeaf([]field.Element{field.One})
	if !leaf.IsLeaf() {
		t.Error("IsLeaf() should return true for leaf node")
	}
	if leaf.IsBranch() {
		t.Error("IsBranch() should return false for leaf node")
	}
	if leaf.IsPadding() {
		t.Error("IsPadding() should return false for leaf node")
	}

	// Test branch node
	left := NewLeaf([]field.Element{field.One})
	right := NewLeaf([]field.Element{field.New(2)})
	branch := NewBranch(left, right)
	if leaf.IsBranch() {
		t.Error("IsBranch() should return false for leaf node")
	}
	if !branch.IsBranch() {
		t.Error("IsBranch() should return true for branch node")
	}
	if branch.IsLeaf() {
		t.Error("IsLeaf() should return false for branch node")
	}
	if branch.IsPadding() {
		t.Error("IsPadding() should return false for branch node")
	}

	// Test padding node
	padding := &ZerofierTree{Type: Padding}
	if !padding.IsPadding() {
		t.Error("IsPadding() should return true for padding node")
	}
	if padding.IsLeaf() {
		t.Error("IsLeaf() should return false for padding node")
	}
	if padding.IsBranch() {
		t.Error("IsBranch() should return false for padding node")
	}
}

func TestZerofierTreeClone(t *testing.T) {
	// Create a complex tree
	domain := []field.Element{field.New(1), field.New(2), field.New(3), field.New(4), field.New(5)}
	tree := NewZerofierTree(domain)

	// Clone the tree
	clone := tree.Clone()

	// Verify they are equal
	if !tree.Equal(clone) {
		t.Error("Clone() did not create an equal tree")
	}

	// Verify they are different instances
	if tree == clone {
		t.Error("Clone() returned the same instance")
	}

	// Verify clone is valid
	if err := clone.Validate(); err != nil {
		t.Errorf("Clone validation failed: %v", err)
	}
}

func TestZerofierTreeEqual(t *testing.T) {
	domain1 := []field.Element{field.New(1), field.New(2), field.New(3)}
	domain2 := []field.Element{field.New(1), field.New(2), field.New(3)}
	domain3 := []field.Element{field.New(1), field.New(2), field.New(4)}

	tree1 := NewZerofierTree(domain1)
	tree2 := NewZerofierTree(domain2)
	tree3 := NewZerofierTree(domain3)

	// Test equal trees
	if !tree1.Equal(tree2) {
		t.Error("Equal trees should be equal")
	}

	// Test different trees
	if tree1.Equal(tree3) {
		t.Error("Different trees should not be equal")
	}

	// Test nil trees
	if tree1.Equal(nil) {
		t.Error("Tree should not equal nil")
	}
	if !(*ZerofierTree)(nil).Equal(nil) {
		t.Error("Nil trees should be equal")
	}
}

func TestZerofierTreeString(t *testing.T) {
	// Test leaf node
	leaf := NewLeaf([]field.Element{field.One})
	leafStr := leaf.String()
	if leafStr == "" {
		t.Error("String() should not return empty string")
	}

	// Test branch node
	left := NewLeaf([]field.Element{field.One})
	right := NewLeaf([]field.Element{field.New(2)})
	branch := NewBranch(left, right)
	branchStr := branch.String()
	if branchStr == "" {
		t.Error("String() should not return empty string")
	}

	// Test padding node
	padding := &ZerofierTree{Type: Padding}
	paddingStr := padding.String()
	if paddingStr != "Padding" {
		t.Errorf("String() = %s, want 'Padding'", paddingStr)
	}
}

func TestZerofierTreeDepth(t *testing.T) {
	tests := []struct {
		name     string
		domain   []field.Element
		expected int
	}{
		{
			name:     "Empty domain",
			domain:   []field.Element{},
			expected: 1, // Padding node
		},
		{
			name:     "Single point",
			domain:   []field.Element{field.One},
			expected: 1, // Leaf node
		},
		{
			name:     "Two points",
			domain:   []field.Element{field.One, field.New(2)},
			expected: 1, // Single leaf (within threshold)
		},
		{
			name:     "Many points",
			domain:   make([]field.Element, 50),
			expected: 3, // Tree with multiple levels
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize test data
			for i := range tt.domain {
				tt.domain[i] = field.New(uint64(i + 1))
			}

			tree := NewZerofierTree(tt.domain)
			depth := tree.Depth()

			// Depth should be at least 1
			if depth < 1 {
				t.Errorf("Depth() = %d, want >= 1", depth)
			}

			// For large domains, depth should be reasonable
			if len(tt.domain) > RecursionCutoffThreshold {
				maxExpectedDepth := 10 // Reasonable upper bound
				if depth > maxExpectedDepth {
					t.Errorf("Depth() = %d, want <= %d", depth, maxExpectedDepth)
				}
			}
		})
	}
}

func TestZerofierTreeSize(t *testing.T) {
	tests := []struct {
		name    string
		domain  []field.Element
		minSize int
		maxSize int
	}{
		{
			name:    "Empty domain",
			domain:  []field.Element{},
			minSize: 1, // At least one padding node
			maxSize: 1,
		},
		{
			name:    "Single point",
			domain:  []field.Element{field.One},
			minSize: 1, // One leaf node
			maxSize: 1,
		},
		{
			name:    "Small domain",
			domain:  []field.Element{field.New(1), field.New(2), field.New(3)},
			minSize: 1, // One leaf node (within threshold)
			maxSize: 1,
		},
		{
			name:    "Large domain",
			domain:  make([]field.Element, 50),
			minSize: 1,   // At least one node
			maxSize: 100, // Reasonable upper bound
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize test data
			for i := range tt.domain {
				tt.domain[i] = field.New(uint64(i + 1))
			}

			tree := NewZerofierTree(tt.domain)
			size := tree.Size()

			if size < tt.minSize {
				t.Errorf("Size() = %d, want >= %d", size, tt.minSize)
			}
			if size > tt.maxSize {
				t.Errorf("Size() = %d, want <= %d", size, tt.maxSize)
			}
		})
	}
}

func TestZerofierTreeLeafCount(t *testing.T) {
	tests := []struct {
		name     string
		domain   []field.Element
		expected int
	}{
		{
			name:     "Empty domain",
			domain:   []field.Element{},
			expected: 0, // Only padding nodes
		},
		{
			name:     "Single point",
			domain:   []field.Element{field.One},
			expected: 1, // One leaf node
		},
		{
			name:     "Small domain",
			domain:   []field.Element{field.New(1), field.New(2), field.New(3)},
			expected: 1, // One leaf node (within threshold)
		},
		{
			name:     "Large domain",
			domain:   make([]field.Element, 50),
			expected: 4, // Multiple leaf nodes (50/16 = 4 chunks)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize test data
			for i := range tt.domain {
				tt.domain[i] = field.New(uint64(i + 1))
			}

			tree := NewZerofierTree(tt.domain)
			leafCount := tree.LeafCount()

			if leafCount != tt.expected {
				t.Errorf("LeafCount() = %d, want %d", leafCount, tt.expected)
			}
		})
	}
}

func TestZerofierTreePointCount(t *testing.T) {
	tests := []struct {
		name     string
		domain   []field.Element
		expected int
	}{
		{
			name:     "Empty domain",
			domain:   []field.Element{},
			expected: 0,
		},
		{
			name:     "Single point",
			domain:   []field.Element{field.One},
			expected: 1,
		},
		{
			name:     "Multiple points",
			domain:   []field.Element{field.New(1), field.New(2), field.New(3)},
			expected: 3,
		},
		{
			name:     "Large domain",
			domain:   make([]field.Element, 50),
			expected: 50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize test data
			for i := range tt.domain {
				tt.domain[i] = field.New(uint64(i + 1))
			}

			tree := NewZerofierTree(tt.domain)
			pointCount := tree.PointCount()

			if pointCount != tt.expected {
				t.Errorf("PointCount() = %d, want %d", pointCount, tt.expected)
			}
		})
	}
}

func TestZerofierTreeValidate(t *testing.T) {
	tests := []struct {
		name    string
		domain  []field.Element
		wantErr bool
	}{
		{
			name:    "Valid empty domain",
			domain:  []field.Element{},
			wantErr: false,
		},
		{
			name:    "Valid single point",
			domain:  []field.Element{field.One},
			wantErr: false,
		},
		{
			name:    "Valid multiple points",
			domain:  []field.Element{field.New(1), field.New(2), field.New(3)},
			wantErr: false,
		},
		{
			name:    "Valid large domain",
			domain:  make([]field.Element, 50),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize test data
			for i := range tt.domain {
				tt.domain[i] = field.New(uint64(i + 1))
			}

			tree := NewZerofierTree(tt.domain)
			err := tree.Validate()

			if tt.wantErr {
				if err == nil {
					t.Errorf("Validate() expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Validate() error = %v", err)
				}
			}
		})
	}
}

func TestHelperFunctions(t *testing.T) {
	// Test nextPowerOfTwo
	tests := []struct {
		input    int
		expected int
	}{
		{0, 1},
		{1, 1},
		{2, 2},
		{3, 4},
		{4, 4},
		{5, 8},
		{8, 8},
		{9, 16},
		{16, 16},
		{17, 32},
	}

	for _, tt := range tests {
		t.Run("nextPowerOfTwo", func(t *testing.T) {
			got := nextPowerOfTwo(tt.input)
			if got != tt.expected {
				t.Errorf("nextPowerOfTwo(%d) = %d, want %d", tt.input, got, tt.expected)
			}
		})
	}

	// Test IsPowerOfTwo
	powerTests := []struct {
		input    int
		expected bool
	}{
		{0, false},
		{1, true},
		{2, true},
		{3, false},
		{4, true},
		{5, false},
		{8, true},
		{9, false},
		{16, true},
		{17, false},
	}

	for _, tt := range powerTests {
		t.Run("IsPowerOfTwo", func(t *testing.T) {
			got := IsPowerOfTwo(tt.input)
			if got != tt.expected {
				t.Errorf("IsPowerOfTwo(%d) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}

	// Test GetRecursionCutoffThreshold
	threshold := GetRecursionCutoffThreshold()
	if threshold != RecursionCutoffThreshold {
		t.Errorf("GetRecursionCutoffThreshold() = %d, want %d", threshold, RecursionCutoffThreshold)
	}
}

// Benchmark tests
func BenchmarkNewZerofierTree(b *testing.B) {
	domain := make([]field.Element, 1000)
	for i := range domain {
		domain[i] = field.New(uint64(i + 1))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewZerofierTree(domain)
	}
}

func BenchmarkZerofierTreeClone(b *testing.B) {
	domain := make([]field.Element, 100)
	for i := range domain {
		domain[i] = field.New(uint64(i + 1))
	}
	tree := NewZerofierTree(domain)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = tree.Clone()
	}
}

func BenchmarkZerofierTreeValidate(b *testing.B) {
	domain := make([]field.Element, 100)
	for i := range domain {
		domain[i] = field.New(uint64(i + 1))
	}
	tree := NewZerofierTree(domain)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = tree.Validate()
	}
}
