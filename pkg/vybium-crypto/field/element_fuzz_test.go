package field

import (
	"testing"
)

func FuzzElementOperations(f *testing.F) {
	// Add seed values
	f.Add(uint64(0))
	f.Add(uint64(1))
	f.Add(uint64(42))
	f.Add(uint64(100))
	f.Add(uint64(1000))
	f.Add(uint64(10000))

	f.Fuzz(func(t *testing.T, value uint64) {
		// Create element from fuzzed value
		elem := New(value)

		// Test basic operations don't panic
		_ = elem.Add(elem)
		_ = elem.Sub(elem)
		_ = elem.Mul(elem)

		// Test inverse (should not panic for non-zero elements)
		if !elem.IsZero() {
			_ = elem.Inverse()
		}

		// Test serialization
		data, err := elem.MarshalBinary()
		if err != nil {
			t.Errorf("MarshalBinary failed: %v", err)
		}

		// Test deserialization
		restored := Element{}
		err = restored.UnmarshalBinary(data)
		if err != nil {
			t.Errorf("UnmarshalBinary failed: %v", err)
		}

		// Test round trip
		if !elem.Equal(restored) {
			t.Errorf("Round trip failed: %v != %v", elem, restored)
		}
	})
}

func FuzzElementArithmetic(f *testing.F) {
	// Add seed pairs
	f.Add(uint64(0), uint64(0))
	f.Add(uint64(1), uint64(1))
	f.Add(uint64(42), uint64(100))
	f.Add(uint64(1000), uint64(2000))

	f.Fuzz(func(t *testing.T, a, b uint64) {
		elemA := New(a)
		elemB := New(b)

		// Test addition
		_ = elemA.Add(elemB)
		// Note: Addition can result in zero element (e.g., a + (-a) = 0)
		// This is mathematically correct, so we don't test for non-zero

		// Test subtraction
		diff := elemA.Sub(elemB)
		if diff == (Element{}) && a != b {
			t.Error("Subtraction result is zero element when operands are different")
		}

		// Test multiplication
		prod := elemA.Mul(elemB)
		if prod == (Element{}) && a != 0 && b != 0 {
			t.Error("Multiplication result is zero element when operands are non-zero")
		}

		// Test division (if divisor is not zero)
		if !elemB.IsZero() {
			quot := elemA.Div(elemB)
			if quot == (Element{}) {
				t.Error("Division result is zero element")
			}
		}

		// Test commutativity of addition
		sum1 := elemA.Add(elemB)
		sum2 := elemB.Add(elemA)
		if !sum1.Equal(sum2) {
			t.Errorf("Addition commutativity failed: %v + %v != %v + %v", elemA, elemB, elemB, elemA)
		}

		// Test commutativity of multiplication
		prod1 := elemA.Mul(elemB)
		prod2 := elemB.Mul(elemA)
		if !prod1.Equal(prod2) {
			t.Errorf("Multiplication commutativity failed: %v * %v != %v * %v", elemA, elemB, elemB, elemA)
		}
	})
}

func FuzzElementComparison(f *testing.F) {
	// Add seed pairs
	f.Add(uint64(0), uint64(0))
	f.Add(uint64(1), uint64(1))
	f.Add(uint64(42), uint64(100))
	f.Add(uint64(1000), uint64(2000))

	f.Fuzz(func(t *testing.T, a, b uint64) {
		elemA := New(a)
		elemB := New(b)

		// Test equality
		equal := elemA.Equal(elemB)
		expectedEqual := (a == b)
		if equal != expectedEqual {
			t.Errorf("Equality test failed: %v.Equal(%v) = %v, expected %v", elemA, elemB, equal, expectedEqual)
		}

		// Test that equal elements are equal
		if !elemA.Equal(elemA) {
			t.Errorf("Element should be equal to itself: %v.Equal(%v) = false", elemA, elemA)
		}

		// Test that zero is equal to zero
		zero := Zero
		if !zero.Equal(zero) {
			t.Error("Zero should be equal to zero")
		}

		// Test that one is equal to one
		one := One
		if !one.Equal(one) {
			t.Error("One should be equal to one")
		}
	})
}

func FuzzElementSerialization(f *testing.F) {
	// Add seed values
	f.Add(uint64(0))
	f.Add(uint64(1))
	f.Add(uint64(42))
	f.Add(uint64(100))
	f.Add(uint64(1000))

	f.Fuzz(func(t *testing.T, value uint64) {
		elem := New(value)

		// Test serialization
		data, err := elem.MarshalBinary()
		if err != nil {
			t.Errorf("MarshalBinary failed: %v", err)
			return
		}

		if len(data) == 0 {
			t.Error("Serialized data is empty")
			return
		}

		// Test deserialization
		restored := Element{}
		err = restored.UnmarshalBinary(data)
		if err != nil {
			t.Errorf("UnmarshalBinary failed: %v", err)
			return
		}

		// Test round trip
		if !elem.Equal(restored) {
			t.Errorf("Round trip failed: %v != %v", elem, restored)
		}

		// Test that serialization is deterministic
		data2, err2 := elem.MarshalBinary()
		if err2 != nil {
			t.Errorf("Second MarshalBinary failed: %v", err2)
			return
		}

		if len(data) != len(data2) {
			t.Errorf("Serialization length mismatch: %d != %d", len(data), len(data2))
			return
		}

		for i := range data {
			if data[i] != data2[i] {
				t.Errorf("Serialization not deterministic at byte %d: %d != %d", i, data[i], data2[i])
			}
		}
	})
}

func FuzzElementMontgomery(f *testing.F) {
	// Add seed values
	f.Add(uint64(0))
	f.Add(uint64(1))
	f.Add(uint64(42))
	f.Add(uint64(100))
	f.Add(uint64(1000))

	f.Fuzz(func(t *testing.T, value uint64) {
		elem := New(value)

		// Test that Montgomery form is consistent
		canonical := elem.Value()
		expected := value % P
		if canonical != expected {
			t.Errorf("Montgomery conversion failed: New(%d).Value() = %d, expected %d", value, canonical, expected)
		}

		// Test that creating from canonical value gives same result
		elem2 := New(canonical)
		if !elem.Equal(elem2) {
			t.Errorf("Montgomery consistency failed: New(%d) != New(%d)", value, canonical)
		}

		// Test that raw value is in Montgomery form
		raw := elem.RawValue()
		if raw == 0 && value != 0 {
			t.Errorf("Raw value is zero for non-zero input: %d -> %d", value, raw)
		}
	})
}
