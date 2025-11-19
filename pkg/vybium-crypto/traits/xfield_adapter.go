package traits

import (
	"encoding/json"
	"math/big"

	"github.com/vybium/vybium-crypto/pkg/vybium-crypto/field"
	"github.com/vybium/vybium-crypto/pkg/vybium-crypto/xfield"
)

// XFieldElementAdapter adapts xfield.XFieldElement to implement FiniteField interface
type XFieldElementAdapter struct {
	Element xfield.XFieldElement
}

// NewXFieldElementAdapter creates a new XFieldElementAdapter
func NewXFieldElementAdapter(elem xfield.XFieldElement) *XFieldElementAdapter {
	return &XFieldElementAdapter{Element: elem}
}

// FiniteField interface implementation
func (x *XFieldElementAdapter) Add(other FiniteField) FiniteField {
	if otherX, ok := other.(*XFieldElementAdapter); ok {
		return &XFieldElementAdapter{Element: x.Element.Add(otherX.Element)}
	}
	panic("type mismatch in Add")
}

func (x *XFieldElementAdapter) Sub(other FiniteField) FiniteField {
	if otherX, ok := other.(*XFieldElementAdapter); ok {
		return &XFieldElementAdapter{Element: x.Element.Sub(otherX.Element)}
	}
	panic("type mismatch in Sub")
}

func (x *XFieldElementAdapter) Mul(other FiniteField) FiniteField {
	if otherX, ok := other.(*XFieldElementAdapter); ok {
		return &XFieldElementAdapter{Element: x.Element.Mul(otherX.Element)}
	}
	panic("type mismatch in Mul")
}

func (x *XFieldElementAdapter) Div(other FiniteField) FiniteField {
	if otherX, ok := other.(*XFieldElementAdapter); ok {
		return &XFieldElementAdapter{Element: x.Element.Div(otherX.Element)}
	}
	panic("type mismatch in Div")
}

func (x *XFieldElementAdapter) Neg() FiniteField {
	return &XFieldElementAdapter{Element: x.Element.Neg()}
}

func (x *XFieldElementAdapter) Equal(other FiniteField) bool {
	if otherX, ok := other.(*XFieldElementAdapter); ok {
		return x.Element.Equal(otherX.Element)
	}
	return false
}

func (x *XFieldElementAdapter) IsZero() bool {
	return x.Element.IsZero()
}

func (x *XFieldElementAdapter) IsOne() bool {
	one := xfield.New([3]field.Element{field.One, field.Zero, field.Zero})
	return x.Element.Equal(one)
}

func (x *XFieldElementAdapter) Inverse() FiniteField {
	return &XFieldElementAdapter{Element: x.Element.Inverse()}
}

func (x *XFieldElementAdapter) Square() FiniteField {
	return &XFieldElementAdapter{Element: x.Element.Mul(x.Element)}
}

func (x *XFieldElementAdapter) Pow(exp uint64) FiniteField {
	return &XFieldElementAdapter{Element: x.Element.Pow(exp)}
}

func (x *XFieldElementAdapter) ToBigInt() *big.Int {
	// Convert XFieldElement to big.Int by evaluating at x=1
	// For a + bx + cx², this gives a + b + c
	result := big.NewInt(0)
	coeffs := x.Element.Coefficients
	for _, coeff := range coeffs {
		coeffBig := coeff.ToBigInt()
		// Multiply by 1^i = 1, so just add the coefficient
		result.Add(result, coeffBig)
	}
	return result
}

func (x *XFieldElementAdapter) FromBigInt(val *big.Int) FiniteField {
	// Convert big.Int to XFieldElement by treating as constant polynomial
	// This creates a + 0x + 0x² where a is the big.Int value
	coeffs := [3]field.Element{
		field.NewFromBigInt(val),
		field.Zero,
		field.Zero,
	}
	return &XFieldElementAdapter{Element: xfield.New(coeffs)}
}

func (x *XFieldElementAdapter) ToUint64() uint64 {
	// Convert XFieldElement to uint64 by evaluating at x=1
	// For a + bx + cx², this gives a + b + c
	result := uint64(0)
	coeffs := x.Element.Coefficients
	for _, coeff := range coeffs {
		// Multiply by 1^i = 1, so just add the coefficient
		result += coeff.Value()
	}
	return result
}

func (x *XFieldElementAdapter) FromUint64(val uint64) FiniteField {
	if val == 0 {
		zero := xfield.New([3]field.Element{field.Zero, field.Zero, field.Zero})
		return &XFieldElementAdapter{Element: zero}
	}
	if val == 1 {
		one := xfield.New([3]field.Element{field.One, field.Zero, field.Zero})
		return &XFieldElementAdapter{Element: one}
	}
	// For other values, create constant polynomial
	coeffs := [3]field.Element{field.New(val), field.Zero, field.Zero}
	return &XFieldElementAdapter{Element: xfield.New(coeffs)}
}

func (x *XFieldElementAdapter) String() string {
	return x.Element.String()
}

func (x *XFieldElementAdapter) MarshalJSON() ([]byte, error) {
	return json.Marshal(x.Element)
}

func (x *XFieldElementAdapter) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &x.Element)
}

// Inverse interface implementation
func (x *XFieldElementAdapter) InverseOrZero() FiniteField {
	if x.IsZero() {
		zero := xfield.New([3]field.Element{field.Zero, field.Zero, field.Zero})
		return &XFieldElementAdapter{Element: zero}
	}
	return x.Inverse()
}

// PrimitiveRootOfUnity interface implementation
func (x *XFieldElementAdapter) PrimitiveRootOfUnity(n uint64) (FiniteField, bool) {
	// XFieldElement doesn't have primitive root of unity
	// This is a simplified implementation
	return nil, false
}

// CyclicGroupGenerator interface implementation
func (x *XFieldElementAdapter) GetCyclicGroupElements(max *uint64) []FiniteField {
	// This is a simplified implementation
	// In practice, you'd implement proper cyclic group generation
	var elements []FiniteField
	limit := uint64(10) // Default limit for extension fields
	if max != nil {
		limit = *max
	}

	for i := uint64(1); i <= limit; i++ {
		coeffs := [3]field.Element{field.New(i), field.Zero, field.Zero}
		elements = append(elements, &XFieldElementAdapter{Element: xfield.New(coeffs)})
	}

	return elements
}

// ModPowU32 interface implementation
func (x *XFieldElementAdapter) ModPowU32(exp uint32) FiniteField {
	return &XFieldElementAdapter{Element: x.Element.Pow(uint64(exp))}
}

// ModPowU64 interface implementation
func (x *XFieldElementAdapter) ModPowU64(exp uint64) FiniteField {
	return &XFieldElementAdapter{Element: x.Element.Pow(exp)}
}

// Type checking functions
func (x *XFieldElementAdapter) XFieldElement() bool {
	return true
}
