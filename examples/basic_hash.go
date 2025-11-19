// Basic hash examples demonstrating Tip5 and Poseidon hash functions
package main

import (
	"fmt"

	"github.com/vybium/vybium-crypto/pkg/vybium-crypto/field"
	"github.com/vybium/vybium-crypto/pkg/vybium-crypto/hash"
)

func main() {
	fmt.Println("🔐 Vybium Crypto Basic Hash Examples")
	fmt.Println("===================================")

	// Example 1: Tip5 Hash Function
	fmt.Println("\n1. Tip5 Hash Function:")
	demonstrateTip5Hash()

	// Example 2: Poseidon Hash Function
	fmt.Println("\n2. Poseidon Hash Function:")
	demonstratePoseidonHash()

	// Example 3: Hash Two Elements
	fmt.Println("\n3. Hash Two Elements:")
	demonstrateHashTwo()

	// Example 4: Variable Length Hashing
	fmt.Println("\n4. Variable Length Hashing:")
	demonstrateVarlenHash()
}

func demonstrateTip5Hash() {
	// Create some test data
	inputs := []field.Element{
		field.New(1),
		field.New(2),
		field.New(3),
		field.New(4),
		field.New(5),
	}

	// Hash the inputs using Tip5
	hashResult := hash.HashVarlen(inputs)

	fmt.Printf("   Input: %v\n", inputs)
	fmt.Printf("   Tip5 Hash: %v\n", hashResult)
}

func demonstratePoseidonHash() {
	// Create some test data
	inputs := []field.Element{
		field.New(42),
		field.New(1337),
		field.New(2024),
	}

	// Hash the inputs using Poseidon
	hashResult := hash.PoseidonHash(inputs)

	fmt.Printf("   Input: %v\n", inputs)
	fmt.Printf("   Poseidon Hash: %v\n", hashResult)
}

func demonstrateHashTwo() {
	left := field.New(12345)
	right := field.New(67890)

	// Hash two elements using Tip5
	leftDigest := [5]field.Element{left, field.Zero, field.Zero, field.Zero, field.Zero}
	rightDigest := [5]field.Element{right, field.Zero, field.Zero, field.Zero, field.Zero}
	tip5Hash := hash.HashPair(leftDigest, rightDigest)

	// Hash two elements using Poseidon
	poseidonHash := hash.PoseidonHashTwo(left, right)

	fmt.Printf("   Left: %v, Right: %v\n", left, right)
	fmt.Printf("   Tip5 Hash: %v\n", tip5Hash)
	fmt.Printf("   Poseidon Hash: %v\n", poseidonHash)
}

func demonstrateVarlenHash() {
	// Create a larger dataset
	inputs := make([]field.Element, 10)
	for i := 0; i < 10; i++ {
		inputs[i] = field.New(uint64(i * 1000))
	}

	// Hash using Tip5 varlen
	tip5Hash := hash.HashVarlen(inputs)

	// Hash using Poseidon
	poseidonHash := hash.PoseidonHash(inputs)

	fmt.Printf("   Input count: %d\n", len(inputs))
	fmt.Printf("   Tip5 Varlen Hash: %v\n", tip5Hash)
	fmt.Printf("   Poseidon Hash: %v\n", poseidonHash)
}
