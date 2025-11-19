package field

import (
	"math/big"
	"testing"
)

func TestElementBasicOperations(t *testing.T) {
	// Test basic field operations
	a := New(42)
	b := New(13)

	// Test addition
	sum := a.Add(b)
	expected := New(55)
	if !sum.Equal(expected) {
		t.Errorf("Addition failed: expected %v, got %v", expected, sum)
	}

	// Test subtraction
	diff := a.Sub(b)
	expected = New(29)
	if !diff.Equal(expected) {
		t.Errorf("Subtraction failed: expected %v, got %v", expected, diff)
	}

	// Test multiplication
	prod := a.Mul(b)
	expected = New(42 * 13)
	if !prod.Equal(expected) {
		t.Errorf("Multiplication failed: expected %v, got %v", expected, prod)
	}

	// Test division
	quot := a.Div(b)
	expected = a.Mul(b.Inverse())
	if !quot.Equal(expected) {
		t.Errorf("Division failed: expected %v, got %v", expected, quot)
	}
}

func TestElementInverse(t *testing.T) {
	// Test multiplicative inverse
	a := New(42)
	inv := a.Inverse()

	// a * a^(-1) should equal 1
	prod := a.Mul(inv)
	if !prod.Equal(One) {
		t.Errorf("Inverse test failed: %v * %v = %v, expected 1", a, inv, prod)
	}
}

func TestElementModPow(t *testing.T) {
	// Test modular exponentiation
	base := New(3)
	exp := uint64(5)
	result := base.ModPow(exp)
	expected := New(3 * 3 * 3 * 3 * 3)

	if !result.Equal(expected) {
		t.Errorf("ModPow failed: %v^%d = %v, expected %v", base, exp, result, expected)
	}

	// Test edge cases
	zero := New(0)
	if !zero.ModPow(0).Equal(One) {
		t.Error("0^0 should equal 1")
	}

	one := New(1)
	if !one.ModPow(100).Equal(One) {
		t.Error("1^100 should equal 1")
	}
}

func TestElementNegation(t *testing.T) {
	// Test additive inverse
	a := New(42)
	neg := a.Neg()

	// a + (-a) should equal 0
	sum := a.Add(neg)
	if !sum.Equal(Zero) {
		t.Errorf("Negation test failed: %v + %v = %v, expected 0", a, neg, sum)
	}
}

func TestElementFromInt64(t *testing.T) {
	// Test negative values
	neg := NewFromInt64(-5)
	expected := New(P - 5)
	if !neg.Equal(expected) {
		t.Errorf("Negative value test failed: expected %v, got %v", expected, neg)
	}

	// Test positive values
	pos := NewFromInt64(42)
	expected = New(42)
	if !pos.Equal(expected) {
		t.Errorf("Positive value test failed: expected %v, got %v", expected, pos)
	}
}

func TestElementFromBigInt(t *testing.T) {
	// Test large values
	bigVal := big.NewInt(0).SetUint64(0xFFFFFFFFFFFFFFFF)
	element := NewFromBigInt(bigVal)
	expected := New(0xFFFFFFFFFFFFFFFF)
	if !element.Equal(expected) {
		t.Errorf("BigInt test failed: expected %v, got %v", expected, element)
	}

	// Test values larger than P
	bigVal2 := big.NewInt(0).SetUint64(P + 100)
	element2 := NewFromBigInt(bigVal2)
	expected2 := New(100)
	if !element2.Equal(expected2) {
		t.Errorf("BigInt modulo test failed: expected %v, got %v", expected2, element2)
	}
}

func TestElementSerialization(t *testing.T) {
	// Test binary serialization
	original := New(0x123456789ABCDEF0)

	data, err := original.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary failed: %v", err)
	}

	if len(data) != 8 {
		t.Errorf("Expected 8 bytes, got %d", len(data))
	}

	var restored Element
	err = restored.UnmarshalBinary(data)
	if err != nil {
		t.Fatalf("UnmarshalBinary failed: %v", err)
	}

	if !restored.Equal(original) {
		t.Errorf("Serialization round-trip failed: original %v, restored %v", original, restored)
	}
}

func TestElementConstants(t *testing.T) {
	// Test constants
	if !Zero.IsZero() {
		t.Error("Zero constant is incorrect")
	}

	if !One.IsOne() {
		t.Error("One constant is incorrect")
	}

	if One.Value() != 1 {
		t.Errorf("One.Value() should be 1, got %d", One.Value())
	}

	if Max.Value() != P-1 {
		t.Errorf("Max.Value() should be %d, got %d", P-1, Max.Value())
	}

	// Test that One * One = One
	if !One.Mul(One).Equal(One) {
		t.Error("One * One should equal One")
	}
}

func TestElementComparison(t *testing.T) {
	a := New(42)
	b := New(13)
	c := New(42)

	// Test equality
	if !a.Equal(c) {
		t.Error("Equal test failed")
	}

	if a.Equal(b) {
		t.Error("Equal test failed for different values")
	}

	// Test less than
	if !b.Less(a) {
		t.Error("Less test failed")
	}

	if a.Less(b) {
		t.Error("Less test failed")
	}

	// Test greater than
	if !a.Greater(b) {
		t.Error("Greater test failed")
	}

	if b.Greater(a) {
		t.Error("Greater test failed")
	}
}

func TestElementEdgeCases(t *testing.T) {
	// Test division by zero
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for division by zero")
		}
	}()

	New(42).Div(Zero)
}

func TestElementModularReduction(t *testing.T) {
	// Test that values are properly reduced modulo P
	large := New(P + 100)
	expected := New(100)
	if !large.Equal(expected) {
		t.Errorf("Modular reduction failed: expected %v, got %v", expected, large)
	}
}
