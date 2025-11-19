// Extension field examples demonstrating F_p^3 operations
package main

import (
	"fmt"

	"github.com/vybium/vybium-crypto/pkg/vybium-crypto/field"
	"github.com/vybium/vybium-crypto/pkg/vybium-crypto/xfield"
)

func main() {
	fmt.Println("🔗 Vybium Crypto Extension Field Examples")
	fmt.Println("=======================================")

	// Example 1: Basic Extension Field Operations
	fmt.Println("\n1. Basic Extension Field Operations:")
	demonstrateBasicExtensionField()

	// Example 2: Extension Field Properties
	fmt.Println("\n2. Extension Field Properties:")
	demonstrateExtensionFieldProperties()

	// Example 3: Extension Field Arithmetic
	fmt.Println("\n3. Extension Field Arithmetic:")
	demonstrateExtensionFieldArithmetic()

	// Example 4: Extension Field Conversion
	fmt.Println("\n4. Extension Field Conversion:")
	demonstrateExtensionFieldConversion()
}

func demonstrateBasicExtensionField() {
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
}

func demonstrateExtensionFieldProperties() {
	// Create extension field elements
	zero := xfield.Zero
	one := xfield.One
	a := xfield.New([3]field.Element{
		field.New(1),
		field.New(2),
		field.New(3),
	})

	fmt.Printf("   Zero: %v\n", zero)
	fmt.Printf("   One: %v\n", one)
	fmt.Printf("   a = %v\n", a)
	fmt.Printf("   a + 0 = %v (additive identity)\n", a.Add(zero))
	fmt.Printf("   a * 1 = %v (multiplicative identity)\n", a.Mul(one))
	fmt.Printf("   a + (-a) = %v (additive inverse)\n", a.Add(a.Neg()))
	fmt.Printf("   a * a^(-1) = %v (multiplicative inverse)\n", a.Mul(a.Inverse()))
}

func demonstrateExtensionFieldArithmetic() {
	// Create extension field element: 1 + x + x²
	coeffs := [3]field.Element{
		field.New(1),
		field.New(1),
		field.New(1),
	}
	a := xfield.New(coeffs)

	fmt.Printf("   a = %v\n", a)
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

func demonstrateExtensionFieldConversion() {
	// Create extension field element
	coeffs := [3]field.Element{
		field.New(42),
		field.New(1337),
		field.New(2024),
	}
	a := xfield.New(coeffs)

	fmt.Printf("   Extension field element: %v\n", a)

	// Get constant term (unlift operation)
	unlifted := a.Unlift()
	fmt.Printf("   Unlifted (constant term): %v\n", unlifted)

	// Create constant extension field element from base field
	fromBase := xfield.NewConst(field.New(999))
	fmt.Printf("   From base field (999): %v\n", fromBase)

	// Create another constant extension field element
	lifted := xfield.NewConst(field.New(777))
	fmt.Printf("   Constant extension field element (777): %v\n", lifted)
}
