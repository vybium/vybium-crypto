// Simple hash examples demonstrating core functionality
package main

import (
	"fmt"

	"github.com/vybium/vybium-crypto/pkg/vybium-crypto/bfieldcodec"
	"github.com/vybium/vybium-crypto/pkg/vybium-crypto/field"
	"github.com/vybium/vybium-crypto/pkg/vybium-crypto/hash"
)

func main() {
	fmt.Println("🔐 Vybium Crypto Simple Hash Examples")
	fmt.Println("====================================")

	// Example 1: Basic Field Operations
	fmt.Println("\n1. Basic Field Operations:")
	demonstrateFieldOperations()

	// Example 2: Hash Functions
	fmt.Println("\n2. Hash Functions:")
	demonstrateHashFunctions()

	// Example 3: Serialization
	fmt.Println("\n3. Serialization:")
	demonstrateSerialization()
}

func demonstrateFieldOperations() {
	// Create field elements
	a := field.New(42)
	b := field.New(1337)

	fmt.Printf("   a = %v\n", a)
	fmt.Printf("   b = %v\n", b)
	fmt.Printf("   a + b = %v\n", a.Add(b))
	fmt.Printf("   a * b = %v\n", a.Mul(b))
	fmt.Printf("   a^2 = %v\n", a.Square())
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

func demonstrateSerialization() {
	// Create field elements
	elements := []field.Element{
		field.Zero,
		field.One,
		field.New(42),
		field.New(1337),
	}

	fmt.Printf("   Original elements: %v\n", elements)

	// Encode elements
	encoded, err := bfieldcodec.EncodeSlice(elements)
	if err != nil {
		fmt.Printf("   Error encoding: %v\n", err)
		return
	}
	fmt.Printf("   Encoded length: %d bytes\n", len(encoded))

	// Decode elements
	decoded, err := bfieldcodec.DecodeSlice[field.Element](encoded, func() field.Element {
		return field.Zero
	})
	if err != nil {
		fmt.Printf("   Error decoding: %v\n", err)
		return
	}
	fmt.Printf("   Decoded elements: %v\n", decoded)

	// Verify round-trip
	fmt.Printf("   Round-trip verification: %t\n", elementsEqual(elements, decoded))
}

// Helper function for comparison
func elementsEqual(a, b []field.Element) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !a[i].Equal(b[i]) {
			return false
		}
	}
	return true
}
