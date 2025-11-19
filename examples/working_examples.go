// Working examples demonstrating core Vybium Crypto functionality
package main

import (
	"fmt"

	"github.com/vybium/vybium-crypto/pkg/vybium-crypto/field"
	"github.com/vybium/vybium-crypto/pkg/vybium-crypto/hash"
	"github.com/vybium/vybium-crypto/pkg/vybium-crypto/xfield"
)

func main() {
	fmt.Println("🔐 Vybium Crypto Working Examples")
	fmt.Println("=================================")

	// Example 1: Field Arithmetic
	fmt.Println("\n1. Field Arithmetic:")
	demonstrateFieldArithmetic()

	// Example 2: Hash Functions
	fmt.Println("\n2. Hash Functions:")
	demonstrateHashFunctions()

	// Example 3: Extension Field
	fmt.Println("\n3. Extension Field:")
	demonstrateExtensionField()

	// Example 4: Polynomial Operations
	fmt.Println("\n4. Polynomial Operations:")
	demonstratePolynomialOperations()
}

func demonstrateFieldArithmetic() {
	// Create field elements
	a := field.New(42)
	b := field.New(1337)

	fmt.Printf("   a = %v\n", a)
	fmt.Printf("   b = %v\n", b)
	fmt.Printf("   a + b = %v\n", a.Add(b))
	fmt.Printf("   a - b = %v\n", a.Sub(b))
	fmt.Printf("   a * b = %v\n", a.Mul(b))
	fmt.Printf("   a / b = %v\n", a.Div(b))
	fmt.Printf("   a^2 = %v\n", a.Square())
	fmt.Printf("   a^3 = %v\n", a.ModPow(3))
	fmt.Printf("   a^(-1) = %v\n", a.Inverse())

	// Verify inverse property
	inverse := a.Inverse()
	product := a.Mul(inverse)
	fmt.Printf("   a * a^(-1) = %v (should be 1)\n", product)
}

func demonstrateHashFunctions() {
	// Create test data
	inputs := []field.Element{
		field.New(1),
		field.New(2),
		field.New(3),
		field.New(4),
	}

	fmt.Printf("   Input: %v\n", inputs)

	// Hash using Tip5
	tip5Hash := hash.HashVarlen(inputs)
	fmt.Printf("   Tip5 Hash: %v\n", tip5Hash)

	// Hash using Poseidon
	poseidonHash := hash.PoseidonHash(inputs)
	fmt.Printf("   Poseidon Hash: %v\n", poseidonHash)

	// Hash two elements
	left := field.New(12345)
	right := field.New(67890)
	poseidonPair := hash.PoseidonHashTwo(left, right)
	fmt.Printf("   Poseidon Hash of (%v, %v): %v\n", left, right, poseidonPair)
}

func demonstrateExtensionField() {
	// Create extension field elements
	// a = 1 + 2x + 3x²
	coeffs1 := [3]field.Element{
		field.New(1),
		field.New(2),
		field.New(3),
	}

	// b = 4 + 5x + 6x²
	coeffs2 := [3]field.Element{
		field.New(4),
		field.New(5),
		field.New(6),
	}

	a := xfield.New(coeffs1)
	b := xfield.New(coeffs2)

	fmt.Printf("   a = %v\n", a)
	fmt.Printf("   b = %v\n", b)
	fmt.Printf("   a + b = %v\n", a.Add(b))
	fmt.Printf("   a - b = %v\n", a.Sub(b))
	fmt.Printf("   a * b = %v\n", a.Mul(b))
	fmt.Printf("   a / b = %v\n", a.Div(b))
	fmt.Printf("   a^2 = %v\n", a.Pow(2))
	fmt.Printf("   a^(-1) = %v\n", a.Inverse())
}

func demonstratePolynomialOperations() {
	// Create extension field element: 1 + x + x²
	coeffs := [3]field.Element{
		field.New(1),
		field.New(1),
		field.New(1),
	}
	a := xfield.New(coeffs)

	fmt.Printf("   Polynomial: %v\n", a)
	fmt.Printf("   a^2 = %v\n", a.Pow(2))
	fmt.Printf("   a^3 = %v\n", a.Pow(3))
	fmt.Printf("   a^4 = %v\n", a.Pow(4))

	// Demonstrate division
	b := a.Mul(a)        // b = a^2
	quotient := b.Div(a) // quotient = b/a = a
	fmt.Printf("   b = a^2 = %v\n", b)
	fmt.Printf("   b / a = %v\n", quotient)
	fmt.Printf("   Division verification: %t\n", quotient.Equal(a))
}
