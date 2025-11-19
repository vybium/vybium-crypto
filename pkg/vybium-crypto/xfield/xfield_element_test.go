package xfield

import (
	"encoding/json"
	"testing"

	"github.com/vybium/vybium-crypto/pkg/vybium-crypto/field"
)

func TestXFieldElementCreation(t *testing.T) {
	tests := []struct {
		name   string
		create func() XFieldElement
		want   XFieldElement
	}{
		{
			name:   "Zero",
			create: func() XFieldElement { return Zero },
			want:   New([3]field.Element{field.Zero, field.Zero, field.Zero}),
		},
		{
			name:   "One",
			create: func() XFieldElement { return One },
			want:   New([3]field.Element{field.One, field.Zero, field.Zero}),
		},
		{
			name:   "NewConst",
			create: func() XFieldElement { return NewConst(field.New(42)) },
			want:   New([3]field.Element{field.New(42), field.Zero, field.Zero}),
		},
		{
			name: "New with all coefficients",
			create: func() XFieldElement {
				return New([3]field.Element{field.New(1), field.New(2), field.New(3)})
			},
			want: New([3]field.Element{field.New(1), field.New(2), field.New(3)}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.create()
			if !got.Equal(tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestXFieldElementIsZero(t *testing.T) {
	tests := []struct {
		name string
		xfe  XFieldElement
		want bool
	}{
		{
			name: "Zero is zero",
			xfe:  Zero,
			want: true,
		},
		{
			name: "One is not zero",
			xfe:  One,
			want: false,
		},
		{
			name: "Constant is not zero",
			xfe:  NewConst(field.New(42)),
			want: false,
		},
		{
			name: "Non-constant is not zero",
			xfe:  New([3]field.Element{field.New(1), field.New(2), field.New(3)}),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.xfe.IsZero()
			if got != tt.want {
				t.Errorf("IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestXFieldElementAdd(t *testing.T) {
	tests := []struct {
		name string
		a, b XFieldElement
		want XFieldElement
	}{
		{
			name: "Zero + Zero",
			a:    Zero,
			b:    Zero,
			want: Zero,
		},
		{
			name: "One + Zero",
			a:    One,
			b:    Zero,
			want: One,
		},
		{
			name: "One + One",
			a:    One,
			b:    One,
			want: NewConst(field.New(2)),
		},
		{
			name: "General addition",
			a:    New([3]field.Element{field.New(1), field.New(2), field.New(3)}),
			b:    New([3]field.Element{field.New(4), field.New(5), field.New(6)}),
			want: New([3]field.Element{field.New(5), field.New(7), field.New(9)}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.a.Add(tt.b)
			if !got.Equal(tt.want) {
				t.Errorf("Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestXFieldElementSub(t *testing.T) {
	tests := []struct {
		name string
		a, b XFieldElement
		want XFieldElement
	}{
		{
			name: "Zero - Zero",
			a:    Zero,
			b:    Zero,
			want: Zero,
		},
		{
			name: "One - One",
			a:    One,
			b:    One,
			want: Zero,
		},
		{
			name: "General subtraction",
			a:    New([3]field.Element{field.New(10), field.New(20), field.New(30)}),
			b:    New([3]field.Element{field.New(1), field.New(2), field.New(3)}),
			want: New([3]field.Element{field.New(9), field.New(18), field.New(27)}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.a.Sub(tt.b)
			if !got.Equal(tt.want) {
				t.Errorf("Sub() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestXFieldElementMulConst(t *testing.T) {
	tests := []struct {
		name   string
		xfe    XFieldElement
		scalar field.Element
		want   XFieldElement
	}{
		{
			name:   "Zero * scalar",
			xfe:    Zero,
			scalar: field.New(5),
			want:   Zero,
		},
		{
			name:   "xfe * 0",
			xfe:    New([3]field.Element{field.New(1), field.New(2), field.New(3)}),
			scalar: field.Zero,
			want:   Zero,
		},
		{
			name:   "xfe * 1",
			xfe:    New([3]field.Element{field.New(1), field.New(2), field.New(3)}),
			scalar: field.One,
			want:   New([3]field.Element{field.New(1), field.New(2), field.New(3)}),
		},
		{
			name:   "Scalar multiplication",
			xfe:    New([3]field.Element{field.New(1), field.New(2), field.New(3)}),
			scalar: field.New(5),
			want:   New([3]field.Element{field.New(5), field.New(10), field.New(15)}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.xfe.MulConst(tt.scalar)
			if !got.Equal(tt.want) {
				t.Errorf("MulConst() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestXFieldElementMul(t *testing.T) {
	tests := []struct {
		name string
		a, b XFieldElement
		want XFieldElement
	}{
		{
			name: "Zero * xfe",
			a:    Zero,
			b:    New([3]field.Element{field.New(1), field.New(2), field.New(3)}),
			want: Zero,
		},
		{
			name: "One * xfe",
			a:    One,
			b:    New([3]field.Element{field.New(1), field.New(2), field.New(3)}),
			want: New([3]field.Element{field.New(1), field.New(2), field.New(3)}),
		},
		{
			name: "xfe * xfe (constants)",
			a:    NewConst(field.New(2)),
			b:    NewConst(field.New(3)),
			want: NewConst(field.New(6)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.a.Mul(tt.b)
			if !got.Equal(tt.want) {
				t.Errorf("Mul() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestXFieldElementNeg(t *testing.T) {
	tests := []struct {
		name string
		xfe  XFieldElement
	}{
		{
			name: "Neg of Zero",
			xfe:  Zero,
		},
		{
			name: "Neg of One",
			xfe:  One,
		},
		{
			name: "Neg of general xfe",
			xfe:  New([3]field.Element{field.New(1), field.New(2), field.New(3)}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			neg := tt.xfe.Neg()
			sum := tt.xfe.Add(neg)
			if !sum.IsZero() {
				t.Errorf("xfe + (-xfe) = %v, want Zero", sum)
			}
		})
	}
}

func TestXFieldElementInverse(t *testing.T) {
	tests := []struct {
		name string
		xfe  XFieldElement
	}{
		{
			name: "Inverse of One",
			xfe:  One,
		},
		{
			name: "Inverse of constant",
			xfe:  NewConst(field.New(7)),
		},
		{
			name: "Inverse of x",
			xfe:  New([3]field.Element{field.Zero, field.One, field.Zero}),
		},
		{
			name: "Inverse of x^2",
			xfe:  New([3]field.Element{field.Zero, field.Zero, field.One}),
		},
		{
			name: "Inverse of general element",
			xfe:  New([3]field.Element{field.New(1), field.New(2), field.New(3)}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inv := tt.xfe.Inverse()
			product := tt.xfe.Mul(inv)
			if !product.IsOne() {
				t.Errorf("xfe * inv(xfe) = %v, want One", product)
			}
		})
	}
}

func TestXFieldElementInverseZeroPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Inverse(Zero) did not panic")
		}
	}()
	_ = Zero.Inverse()
}

func TestXFieldElementDiv(t *testing.T) {
	tests := []struct {
		name string
		a, b XFieldElement
	}{
		{
			name: "One / One",
			a:    One,
			b:    One,
		},
		{
			name: "xfe / xfe (same)",
			a:    New([3]field.Element{field.New(1), field.New(2), field.New(3)}),
			b:    New([3]field.Element{field.New(1), field.New(2), field.New(3)}),
		},
		{
			name: "General division",
			a:    New([3]field.Element{field.New(10), field.New(20), field.New(30)}),
			b:    New([3]field.Element{field.New(1), field.New(2), field.New(3)}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			quotient := tt.a.Div(tt.b)
			product := quotient.Mul(tt.b)
			if !product.Equal(tt.a) {
				t.Errorf("(a / b) * b = %v, want %v", product, tt.a)
			}
		})
	}
}

func TestXFieldElementPow(t *testing.T) {
	tests := []struct {
		name     string
		xfe      XFieldElement
		exponent uint64
		want     XFieldElement
	}{
		{
			name:     "xfe^0",
			xfe:      New([3]field.Element{field.New(1), field.New(2), field.New(3)}),
			exponent: 0,
			want:     One,
		},
		{
			name:     "xfe^1",
			xfe:      New([3]field.Element{field.New(1), field.New(2), field.New(3)}),
			exponent: 1,
			want:     New([3]field.Element{field.New(1), field.New(2), field.New(3)}),
		},
		{
			name:     "One^n",
			xfe:      One,
			exponent: 100,
			want:     One,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.xfe.Pow(tt.exponent)
			if !got.Equal(tt.want) {
				t.Errorf("Pow(%d) = %v, want %v", tt.exponent, got, tt.want)
			}
		})
	}
}

func TestXFieldElementUnlift(t *testing.T) {
	tests := []struct {
		name    string
		xfe     XFieldElement
		wantNil bool
		wantVal uint64
	}{
		{
			name:    "Zero unlifts",
			xfe:     Zero,
			wantNil: false,
			wantVal: 0,
		},
		{
			name:    "One unlifts",
			xfe:     One,
			wantNil: false,
			wantVal: 1,
		},
		{
			name:    "Constant unlifts",
			xfe:     NewConst(field.New(42)),
			wantNil: false,
			wantVal: 42,
		},
		{
			name:    "Non-constant does not unlift",
			xfe:     New([3]field.Element{field.Zero, field.One, field.Zero}),
			wantNil: true,
		},
		{
			name:    "Non-constant does not unlift (c2 != 0)",
			xfe:     New([3]field.Element{field.Zero, field.Zero, field.One}),
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.xfe.Unlift()
			if tt.wantNil {
				if got != nil {
					t.Errorf("Unlift() = %v, want nil", got)
				}
			} else {
				if got == nil {
					t.Errorf("Unlift() = nil, want %d", tt.wantVal)
				} else if got.Value() != tt.wantVal {
					t.Errorf("Unlift().Value() = %d, want %d", got.Value(), tt.wantVal)
				}
			}
		})
	}
}

func TestXFieldElementJSONSerialization(t *testing.T) {
	tests := []struct {
		name string
		xfe  XFieldElement
	}{
		{
			name: "Zero",
			xfe:  Zero,
		},
		{
			name: "One",
			xfe:  One,
		},
		{
			name: "General element",
			xfe:  New([3]field.Element{field.New(1), field.New(2), field.New(3)}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal
			data, err := json.Marshal(tt.xfe)
			if err != nil {
				t.Fatalf("Marshal() error = %v", err)
			}

			// Unmarshal
			var got XFieldElement
			err = json.Unmarshal(data, &got)
			if err != nil {
				t.Fatalf("Unmarshal() error = %v", err)
			}

			// Compare
			if !got.Equal(tt.xfe) {
				t.Errorf("After round-trip: got %v, want %v", got, tt.xfe)
			}
		})
	}
}

func TestXFieldElementDigestConversion(t *testing.T) {
	tests := []struct {
		name string
		xfe  XFieldElement
	}{
		{
			name: "Zero",
			xfe:  Zero,
		},
		{
			name: "One",
			xfe:  One,
		},
		{
			name: "General element",
			xfe:  New([3]field.Element{field.New(1), field.New(2), field.New(3)}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Convert to digest
			digest := tt.xfe.ToDigest()

			// Convert back
			got := FromDigest(digest)
			if got == nil {
				t.Fatal("FromDigest() returned nil")
			}

			// Compare
			if !got.Equal(tt.xfe) {
				t.Errorf("After round-trip: got %v, want %v", got, tt.xfe)
			}
		})
	}
}

func TestFromDigestInvalid(t *testing.T) {
	// Digest with non-zero elements in positions 3 or 4 should fail
	invalidDigest := [5]field.Element{
		field.New(1),
		field.New(2),
		field.New(3),
		field.New(4), // Non-zero!
		field.Zero,
	}

	got := FromDigest(invalidDigest)
	if got != nil {
		t.Errorf("FromDigest(invalid) = %v, want nil", got)
	}
}

func TestAsFlatSlice(t *testing.T) {
	xfes := []XFieldElement{
		New([3]field.Element{field.New(1), field.New(2), field.New(3)}),
		New([3]field.Element{field.New(4), field.New(5), field.New(6)}),
		New([3]field.Element{field.New(7), field.New(8), field.New(9)}),
	}

	flat := AsFlatSlice(xfes)

	// Check length
	expectedLen := len(xfes) * ExtensionDegree
	if len(flat) != expectedLen {
		t.Errorf("AsFlatSlice() length = %d, want %d", len(flat), expectedLen)
	}

	// Check values
	expected := []uint64{1, 2, 3, 4, 5, 6, 7, 8, 9}
	for i, exp := range expected {
		if flat[i].Value() != exp {
			t.Errorf("AsFlatSlice()[%d] = %d, want %d", i, flat[i].Value(), exp)
		}
	}
}

func TestFromBFieldSlice(t *testing.T) {
	tests := []struct {
		name    string
		slice   []field.Element
		wantErr bool
		want    XFieldElement
	}{
		{
			name:    "Valid slice",
			slice:   []field.Element{field.New(1), field.New(2), field.New(3)},
			wantErr: false,
			want:    New([3]field.Element{field.New(1), field.New(2), field.New(3)}),
		},
		{
			name:    "Too short",
			slice:   []field.Element{field.New(1), field.New(2)},
			wantErr: true,
		},
		{
			name:    "Too long",
			slice:   []field.Element{field.New(1), field.New(2), field.New(3), field.New(4)},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FromBFieldSlice(tt.slice)
			if tt.wantErr {
				if err == nil {
					t.Errorf("FromBFieldSlice() error = nil, wantErr %v", tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("FromBFieldSlice() unexpected error = %v", err)
			}

			if !got.Equal(tt.want) {
				t.Errorf("FromBFieldSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShahPolynomial(t *testing.T) {
	// The Shah polynomial is x³ - x + 1
	// Coefficients: [1, -1, 0, 1]
	poly := ShahPolynomial()

	coeffs := poly.Coefficients()
	if len(coeffs) != 4 {
		t.Errorf("Shah polynomial degree = %d, want 3 (length 4)", len(coeffs)-1)
	}

	// Check coefficients: [1, -1, 0, 1]
	expected := []uint64{1, field.P - 1, 0, 1} // -1 = P-1 in field
	for i, exp := range expected {
		if coeffs[i].Value() != exp {
			t.Errorf("Shah coefficient[%d] = %d, want %d", i, coeffs[i].Value(), exp)
		}
	}
}

// Benchmark tests
func BenchmarkXFieldElementAdd(b *testing.B) {
	x := New([3]field.Element{field.New(1), field.New(2), field.New(3)})
	y := New([3]field.Element{field.New(4), field.New(5), field.New(6)})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = x.Add(y)
	}
}

func BenchmarkXFieldElementMul(b *testing.B) {
	x := New([3]field.Element{field.New(1), field.New(2), field.New(3)})
	y := New([3]field.Element{field.New(4), field.New(5), field.New(6)})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = x.Mul(y)
	}
}

func BenchmarkXFieldElementInverse(b *testing.B) {
	x := New([3]field.Element{field.New(1), field.New(2), field.New(3)})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = x.Inverse()
	}
}
