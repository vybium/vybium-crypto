// Polynomial operations examples demonstrating NTT and polynomial arithmetic
package main

import (
	"fmt"
	"log"

	"github.com/vybium/vybium-crypto/pkg/vybium-crypto/field"
	"github.com/vybium/vybium-crypto/pkg/vybium-crypto/polynomial"
)

func main() {
	fmt.Println("📐 Vybium Crypto Polynomial Operations Examples")
	fmt.Println("=============================================")

	// Example 1: Basic Polynomial Operations
	fmt.Println("\n1. Basic Polynomial Operations:")
	demonstrateBasicPolynomialOperations()

	// Example 2: NTT (Number Theoretic Transform)
	fmt.Println("\n2. NTT Operations:")
	demonstrateNTT()

	// Example 3: Polynomial Multiplication
	fmt.Println("\n3. Polynomial Multiplication:")
	demonstratePolynomialMultiplication()

	// Example 4: Polynomial Evaluation
	fmt.Println("\n4. Polynomial Evaluation:")
	demonstratePolynomialEvaluation()
}

func demonstrateBasicPolynomialOperations() {
	// Create two polynomials
	coeffs1 := []field.Element{
		field.New(1),
		field.New(2),
		field.New(3),
	}
	coeffs2 := []field.Element{
		field.New(4),
		field.New(5),
		field.New(6),
	}

	poly1 := polynomial.New(coeffs1)
	poly2 := polynomial.New(coeffs2)

	fmt.Printf("   Polynomial 1: %v\n", poly1)
	fmt.Printf("   Polynomial 2: %v\n", poly2)
	fmt.Printf("   Degree of poly1: %d\n", poly1.Degree())
	fmt.Printf("   Degree of poly2: %d\n", poly2.Degree())

	// Basic operations
	sum := poly1.Add(poly2)
	diff := poly1.Sub(poly2)
	product := poly1.Mul(poly2)

	fmt.Printf("   Sum: %v\n", sum)
	fmt.Printf("   Difference: %v\n", diff)
	fmt.Printf("   Product: %v\n", product)
}

func demonstrateNTT() {
	// Create a polynomial
	coeffs := []field.Element{
		field.New(1),
		field.New(2),
		field.New(3),
		field.New(4),
	}

	poly := polynomial.New(coeffs)
	fmt.Printf("   Original polynomial: %v\n", poly)

	// Evaluate using NTT (this internally uses NTT for efficient evaluation)
	domainSize := 8 // Must be power of 2
	nttEvaluations := poly.EvaluateNTT(domainSize)
	fmt.Printf("   NTT evaluations (domain size %d): %v\n", domainSize, nttEvaluations)

	// Interpolate back (this uses inverse NTT)
	reconstructed := polynomial.InterpolateNTT(nttEvaluations)
	fmt.Printf("   Reconstructed polynomial: %v\n", reconstructed)

	// Verify round-trip
	fmt.Printf("   Round-trip verification: %t\n", poly.Equal(reconstructed))
}

func demonstratePolynomialMultiplication() {
	// Create two polynomials
	coeffs1 := []field.Element{
		field.New(1),
		field.New(2),
		field.New(3),
	}
	coeffs2 := []field.Element{
		field.New(4),
		field.New(5),
	}

	poly1 := polynomial.New(coeffs1)
	poly2 := polynomial.New(coeffs2)

	fmt.Printf("   Polynomial 1: %v\n", poly1)
	fmt.Printf("   Polynomial 2: %v\n", poly2)

	// Multiply using NTT
	product := poly1.MulNTT(poly2)
	fmt.Printf("   Product (NTT): %v\n", product)

	// Multiply using naive method
	naiveProduct := poly1.Mul(poly2)
	fmt.Printf("   Product (naive): %v\n", naiveProduct)

	// Verify they're the same
	fmt.Printf("   NTT vs naive verification: %t\n", product.Equal(naiveProduct))
}

func demonstratePolynomialEvaluation() {
	// Create a polynomial: 1 + 2x + 3x²
	coeffs := []field.Element{
		field.New(1),
		field.New(2),
		field.New(3),
	}

	poly := polynomial.New(coeffs)
	fmt.Printf("   Polynomial: %v\n", poly)

	// Evaluate at different points
	points := []field.Element{
		field.Zero,
		field.One,
		field.New(2),
		field.New(3),
	}

	for _, point := range points {
		result := poly.Evaluate(point)
		fmt.Printf("   P(%v) = %v\n", point, result)
	}
}
