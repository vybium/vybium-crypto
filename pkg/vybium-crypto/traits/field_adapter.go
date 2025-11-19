package traits

import (
	"encoding/json"
	"math/big"

	"github.com/vybium/vybium-crypto/pkg/vybium-crypto/field"
	"github.com/vybium/vybium-crypto/pkg/vybium-crypto/xfield"
)

// BFieldElementAdapter adapts field.Element to implement FiniteField interface
type BFieldElementAdapter struct {
	field.Element
}

// NewBFieldElementAdapter creates a new BFieldElementAdapter
func NewBFieldElementAdapter(elem field.Element) *BFieldElementAdapter {
	return &BFieldElementAdapter{Element: elem}
}

// FiniteField interface implementation
func (b *BFieldElementAdapter) Add(other FiniteField) FiniteField {
	if otherB, ok := other.(*BFieldElementAdapter); ok {
		return &BFieldElementAdapter{Element: b.Element.Add(otherB.Element)}
	}
	panic("type mismatch in Add")
}

func (b *BFieldElementAdapter) Sub(other FiniteField) FiniteField {
	if otherB, ok := other.(*BFieldElementAdapter); ok {
		return &BFieldElementAdapter{Element: b.Element.Sub(otherB.Element)}
	}
	panic("type mismatch in Sub")
}

func (b *BFieldElementAdapter) Mul(other FiniteField) FiniteField {
	if otherB, ok := other.(*BFieldElementAdapter); ok {
		return &BFieldElementAdapter{Element: b.Element.Mul(otherB.Element)}
	}
	panic("type mismatch in Mul")
}

func (b *BFieldElementAdapter) Div(other FiniteField) FiniteField {
	if otherB, ok := other.(*BFieldElementAdapter); ok {
		return &BFieldElementAdapter{Element: b.Element.Div(otherB.Element)}
	}
	panic("type mismatch in Div")
}

func (b *BFieldElementAdapter) Neg() FiniteField {
	return &BFieldElementAdapter{Element: b.Element.Neg()}
}

func (b *BFieldElementAdapter) Equal(other FiniteField) bool {
	if otherB, ok := other.(*BFieldElementAdapter); ok {
		return b.Element.Equal(otherB.Element)
	}
	return false
}

func (b *BFieldElementAdapter) IsZero() bool {
	return b.Element.IsZero()
}

func (b *BFieldElementAdapter) IsOne() bool {
	return b.Element.Equal(field.One)
}

func (b *BFieldElementAdapter) Inverse() FiniteField {
	return &BFieldElementAdapter{Element: b.Element.Inverse()}
}

func (b *BFieldElementAdapter) Square() FiniteField {
	return &BFieldElementAdapter{Element: b.Element.Mul(b.Element)}
}

func (b *BFieldElementAdapter) Pow(exp uint64) FiniteField {
	return &BFieldElementAdapter{Element: b.Element.ModPow(exp)}
}

func (b *BFieldElementAdapter) ToBigInt() *big.Int {
	return b.Element.ToBigInt()
}

func (b *BFieldElementAdapter) FromBigInt(val *big.Int) FiniteField {
	return &BFieldElementAdapter{Element: field.NewFromBigInt(val)}
}

func (b *BFieldElementAdapter) ToUint64() uint64 {
	return b.Element.Value()
}

func (b *BFieldElementAdapter) FromUint64(val uint64) FiniteField {
	return &BFieldElementAdapter{Element: field.New(val)}
}

func (b *BFieldElementAdapter) String() string {
	return b.Element.String()
}

func (b *BFieldElementAdapter) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.Element)
}

func (b *BFieldElementAdapter) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &b.Element)
}

// Inverse interface implementation
func (b *BFieldElementAdapter) InverseOrZero() FiniteField {
	if b.IsZero() {
		return &BFieldElementAdapter{Element: field.Zero}
	}
	return b.Inverse()
}

// PrimitiveRootOfUnity interface implementation
func (b *BFieldElementAdapter) PrimitiveRootOfUnity(n uint64) (FiniteField, bool) {
	root := field.PrimitiveRootOfUnity(n)
	if root.IsZero() {
		return nil, false
	}
	return &BFieldElementAdapter{Element: root}, true
}

// CyclicGroupGenerator interface implementation
func (b *BFieldElementAdapter) GetCyclicGroupElements(max *uint64) []FiniteField {
	// This is a simplified implementation
	// In practice, you'd implement proper cyclic group generation
	var elements []FiniteField
	limit := uint64(100) // Default limit
	if max != nil {
		limit = *max
	}

	for i := uint64(1); i <= limit; i++ {
		elements = append(elements, &BFieldElementAdapter{Element: field.New(i)})
	}

	return elements
}

// ModPowU32 interface implementation
func (b *BFieldElementAdapter) ModPowU32(exp uint32) FiniteField {
	return &BFieldElementAdapter{Element: b.Element.ModPow(uint64(exp))}
}

// ModPowU64 interface implementation
func (b *BFieldElementAdapter) ModPowU64(exp uint64) FiniteField {
	return &BFieldElementAdapter{Element: b.Element.ModPow(exp)}
}

// Helper functions to create adapters
func NewBFieldElement(val uint64) FiniteField {
	return NewBFieldElementAdapter(field.New(val))
}

func NewXFieldElement(coeffs [3]field.Element) FiniteField {
	return NewXFieldElementAdapter(xfield.New(coeffs))
}

// Type checking functions
func (b *BFieldElementAdapter) BFieldElement() bool {
	return true
}
