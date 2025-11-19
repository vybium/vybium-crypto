// Field arithmetic examples demonstrating Goldilocks field operations
package main

import (
	"fmt"
	"math/big"

	"github.com/vybium/vybium-crypto/pkg/vybium-crypto/field"
)

func main() {
	fmt.Println("🔢 Vybium Crypto Field Operations Examples")
	fmt.Println("=================================")

	// Example 1: Basic Field Operations
	fmt.Println("\n1. Basic Field Operations:")
	demonstrateBasicOperations()

	// Example 2: Field Properties
	fmt.Println("\n2. Field Properties:")
	demonstrateFieldProperties()

	// Example 3: Large Number Operations
	fmt.Println("\n3. Large Number Operations:")
	demonstrateLargeNumbers()

	// Example 4: Montgomery Representation
	fmt.Println("\n4. Montgomery Representation:")
	demonstrateMontgomeryRepresentation()
}

func demonstrateBasicOperations() {
	a := field.New(12345)
	b := field.New(67890)

	fmt.Printf("   a = %v\n", a)
	fmt.Printf("   b = %v\n", b)
	fmt.Printf("   a + b = %v\n", a.Add(b))
	fmt.Printf("   a - b = %v\n", a.Sub(b))
	fmt.Printf("   a * b = %v\n", a.Mul(b))
	fmt.Printf("   a / b = %v\n", a.Div(b))
	fmt.Printf("   a^2 = %v\n", a.Square())
	fmt.Printf("   a^3 = %v\n", a.ModPow(3))
}

func demonstrateFieldProperties() {
	zero := field.Zero
	one := field.One
	a := field.New(42)

	fmt.Printf("   Zero: %v\n", zero)
	fmt.Printf("   One: %v\n", one)
	fmt.Printf("   a = %v\n", a)
	fmt.Printf("   a + 0 = %v (additive identity)\n", a.Add(zero))
	fmt.Printf("   a * 1 = %v (multiplicative identity)\n", a.Mul(one))
	fmt.Printf("   a + (-a) = %v (additive inverse)\n", a.Add(a.Neg()))
	fmt.Printf("   a * a^(-1) = %v (multiplicative inverse)\n", a.Mul(a.Inverse()))
}

func demonstrateLargeNumbers() {
	// Create a large number
	large := field.NewFromBigInt(big.NewInt(1234567890123456789))

	fmt.Printf("   Large number: %v\n", large)
	fmt.Printf("   Square: %v\n", large.Square())
	fmt.Printf("   Cube: %v\n", large.ModPow(3))
	fmt.Printf("   Inverse: %v\n", large.Inverse())

	// Verify inverse property
	inverse := large.Inverse()
	product := large.Mul(inverse)
	fmt.Printf("   a * a^(-1) = %v (should be 1)\n", product)
}

func demonstrateMontgomeryRepresentation() {
	a := field.New(12345)
	b := field.New(67890)

	fmt.Printf("   a = %v\n", a)
	fmt.Printf("   b = %v\n", b)

	// Show Montgomery representation
	fmt.Printf("   a in Montgomery: %v\n", a.Value())
	fmt.Printf("   b in Montgomery: %v\n", b.Value())

	// Demonstrate efficient multiplication
	product := a.Mul(b)
	fmt.Printf("   a * b = %v\n", product)
	fmt.Printf("   Product in Montgomery: %v\n", product.Value())
}
