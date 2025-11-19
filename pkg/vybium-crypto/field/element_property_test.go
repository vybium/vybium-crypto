package field

import (
	"testing"
)

func TestElementProperties(t *testing.T) {
	t.Run("AdditiveIdentity", func(t *testing.T) {
		// a + 0 = a
		for i := 0; i < 100; i++ {
			a := New(uint64(i))
			zero := Zero
			result := a.Add(zero)

			if !result.Equal(a) {
				t.Errorf("Additive identity failed: %v + 0 != %v", a, a)
			}
		}
	})

	t.Run("MultiplicativeIdentity", func(t *testing.T) {
		// a * 1 = a
		for i := 1; i < 100; i++ {
			a := New(uint64(i))
			one := One
			result := a.Mul(one)

			if !result.Equal(a) {
				t.Errorf("Multiplicative identity failed: %v * 1 != %v", a, a)
			}
		}
	})

	t.Run("AdditiveInverse", func(t *testing.T) {
		// a + (-a) = 0
		for i := 1; i < 100; i++ {
			a := New(uint64(i))
			negA := a.Neg()
			result := a.Add(negA)

			if !result.IsZero() {
				t.Errorf("Additive inverse failed: %v + (-%v) != 0", a, a)
			}
		}
	})

	t.Run("MultiplicativeInverse", func(t *testing.T) {
		// a * a^(-1) = 1
		for i := 1; i < 100; i++ {
			a := New(uint64(i))
			invA := a.Inverse()
			result := a.Mul(invA)

			if !result.IsOne() {
				t.Errorf("Multiplicative inverse failed: %v * %v^(-1) != 1", a, a)
			}
		}
	})

	t.Run("Commutativity", func(t *testing.T) {
		// a + b = b + a
		// a * b = b * a
		for i := 1; i < 50; i++ {
			for j := 1; j < 50; j++ {
				a := New(uint64(i))
				b := New(uint64(j))

				// Test addition commutativity
				sum1 := a.Add(b)
				sum2 := b.Add(a)
				if !sum1.Equal(sum2) {
					t.Errorf("Addition commutativity failed: %v + %v != %v + %v", a, b, b, a)
				}

				// Test multiplication commutativity
				prod1 := a.Mul(b)
				prod2 := b.Mul(a)
				if !prod1.Equal(prod2) {
					t.Errorf("Multiplication commutativity failed: %v * %v != %v * %v", a, b, b, a)
				}
			}
		}
	})

	t.Run("Associativity", func(t *testing.T) {
		// (a + b) + c = a + (b + c)
		// (a * b) * c = a * (b * c)
		for i := 1; i < 20; i++ {
			for j := 1; j < 20; j++ {
				for k := 1; k < 20; k++ {
					a := New(uint64(i))
					b := New(uint64(j))
					c := New(uint64(k))

					// Test addition associativity
					sum1 := a.Add(b).Add(c)
					sum2 := a.Add(b.Add(c))
					if !sum1.Equal(sum2) {
						t.Errorf("Addition associativity failed: (%v + %v) + %v != %v + (%v + %v)", a, b, c, a, b, c)
					}

					// Test multiplication associativity
					prod1 := a.Mul(b).Mul(c)
					prod2 := a.Mul(b.Mul(c))
					if !prod1.Equal(prod2) {
						t.Errorf("Multiplication associativity failed: (%v * %v) * %v != %v * (%v * %v)", a, b, c, a, b, c)
					}
				}
			}
		}
	})

	t.Run("Distributivity", func(t *testing.T) {
		// a * (b + c) = a * b + a * c
		for i := 1; i < 20; i++ {
			for j := 1; j < 20; j++ {
				for k := 1; k < 20; k++ {
					a := New(uint64(i))
					b := New(uint64(j))
					c := New(uint64(k))

					left := a.Mul(b.Add(c))
					right := a.Mul(b).Add(a.Mul(c))

					if !left.Equal(right) {
						t.Errorf("Distributivity failed: %v * (%v + %v) != %v * %v + %v * %v", a, b, c, a, b, a, c)
					}
				}
			}
		}
	})
}

func TestElementMontgomeryProperties(t *testing.T) {
	t.Run("MontgomeryConversion", func(t *testing.T) {
		// Test that Montgomery conversion is consistent
		for i := 0; i < 100; i++ {
			value := uint64(i)
			elem := New(value)

			// Convert back to canonical form
			canonical := elem.Value()

			// Should be equal to original value (mod P)
			expected := value % P
			if canonical != expected {
				t.Errorf("Montgomery conversion failed: New(%d).Value() = %d, expected %d", value, canonical, expected)
			}
		}
	})

	t.Run("MontgomeryArithmetic", func(t *testing.T) {
		// Test that Montgomery arithmetic gives correct results
		for i := 1; i < 50; i++ {
			for j := 1; j < 50; j++ {
				a := New(uint64(i))
				b := New(uint64(j))

				// Test addition
				sum := a.Add(b)
				expectedSum := (uint64(i) + uint64(j)) % P
				if sum.Value() != expectedSum {
					t.Errorf("Montgomery addition failed: %d + %d = %d, expected %d", i, j, sum.Value(), expectedSum)
				}

				// Test multiplication
				prod := a.Mul(b)
				expectedProd := (uint64(i) * uint64(j)) % P
				if prod.Value() != expectedProd {
					t.Errorf("Montgomery multiplication failed: %d * %d = %d, expected %d", i, j, prod.Value(), expectedProd)
				}
			}
		}
	})
}

func TestElementSerializationProperties(t *testing.T) {
	t.Run("SerializationRoundTrip", func(t *testing.T) {
		// Marshal -> Unmarshal should preserve the element
		for i := 0; i < 100; i++ {
			original := New(uint64(i))

			data, err := original.MarshalBinary()
			if err != nil {
				t.Errorf("MarshalBinary failed: %v", err)
				continue
			}

			restored := Element{}
			err = restored.UnmarshalBinary(data)
			if err != nil {
				t.Errorf("UnmarshalBinary failed: %v", err)
				continue
			}

			if !original.Equal(restored) {
				t.Errorf("Serialization round trip failed: %v != %v", original, restored)
			}
		}
	})

	t.Run("SerializationConsistency", func(t *testing.T) {
		// Same element should serialize to same bytes
		for i := 0; i < 100; i++ {
			elem1 := New(uint64(i))
			elem2 := New(uint64(i))

			data1, err1 := elem1.MarshalBinary()
			data2, err2 := elem2.MarshalBinary()

			if err1 != nil || err2 != nil {
				t.Errorf("Serialization failed: %v, %v", err1, err2)
				continue
			}

			if len(data1) != len(data2) {
				t.Errorf("Serialization length mismatch: %d != %d", len(data1), len(data2))
				continue
			}

			for j := range data1 {
				if data1[j] != data2[j] {
					t.Errorf("Serialization data mismatch at byte %d: %d != %d", j, data1[j], data2[j])
				}
			}
		}
	})
}
