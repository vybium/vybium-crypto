// Serialization examples demonstrating BFieldCodec
package main

import (
	"fmt"
	"log"

	"github.com/vybium/vybium-crypto/pkg/vybium-crypto/bfieldcodec"
	"github.com/vybium/vybium-crypto/pkg/vybium-crypto/field"
	"github.com/vybium/vybium-crypto/pkg/vybium-crypto/xfield"
)

func main() {
	fmt.Println("📦 Vybium Crypto Serialization Examples")
	fmt.Println("====================================")

	// Example 1: Basic Field Element Serialization
	fmt.Println("\n1. Basic Field Element Serialization:")
	demonstrateFieldElementSerialization()

	// Example 2: Extension Field Element Serialization
	fmt.Println("\n2. Extension Field Element Serialization:")
	demonstrateExtensionFieldSerialization()

	// Example 3: Slice Serialization
	fmt.Println("\n3. Slice Serialization:")
	demonstrateSliceSerialization()

	// Example 4: Error Handling
	fmt.Println("\n4. Error Handling:")
	demonstrateErrorHandling()
}

func demonstrateFieldElementSerialization() {
	// Create field elements
	elements := []field.Element{
		field.Zero,
		field.One,
		field.New(42),
		field.New(1337),
		field.New(2024),
	}

	fmt.Printf("   Original elements: %v\n", elements)

	// Encode elements
	encoded, err := bfieldcodec.EncodeSlice(elements)
	if err != nil {
		log.Fatalf("Error encoding elements: %v", err)
	}
	fmt.Printf("   Encoded length: %d bytes\n", len(encoded))
	fmt.Printf("   Encoded data: %v\n", encoded)

	// Decode elements
	decoded, err := bfieldcodec.DecodeSlice[field.Element](encoded, func() field.Element {
		return field.Zero
	})
	if err != nil {
		log.Fatalf("Error decoding elements: %v", err)
	}
	fmt.Printf("   Decoded elements: %v\n", decoded)

	// Verify round-trip
	fmt.Printf("   Round-trip verification: %t\n", elementsEqual(elements, decoded))
}

func demonstrateExtensionFieldSerialization() {
	// Create extension field elements
	coeffs1 := [3]field.Element{
		field.New(1),
		field.New(2),
		field.New(3),
	}
	coeffs2 := [3]field.Element{
		field.New(4),
		field.New(5),
		field.New(6),
	}

	elements := []xfield.XFieldElement{
		xfield.New(coeffs1),
		xfield.New(coeffs2),
	}

	fmt.Printf("   Original XField elements: %v\n", elements)

	// Encode elements
	encoded, err := bfieldcodec.EncodeSlice(elements)
	if err != nil {
		log.Fatalf("Error encoding XField elements: %v", err)
	}
	fmt.Printf("   Encoded length: %d bytes\n", len(encoded))

	// Decode elements
	decoded, err := bfieldcodec.DecodeSlice[xfield.XFieldElement](encoded, func() xfield.XFieldElement {
		return xfield.Zero
	})
	if err != nil {
		log.Fatalf("Error decoding XField elements: %v", err)
	}
	fmt.Printf("   Decoded elements: %v\n", decoded)

	// Verify round-trip
	fmt.Printf("   Round-trip verification: %t\n", xFieldElementsEqual(elements, decoded))
}

func demonstrateSliceSerialization() {
	// Create a large slice of field elements
	elements := make([]field.Element, 100)
	for i := 0; i < 100; i++ {
		elements[i] = field.New(uint64(i * 1000))
	}

	fmt.Printf("   Original slice length: %d\n", len(elements))

	// Encode slice
	encoded, err := bfieldcodec.EncodeSlice(elements)
	if err != nil {
		log.Fatalf("Error encoding large slice: %v", err)
	}
	fmt.Printf("   Encoded length: %d bytes\n", len(encoded))

	// Decode slice
	decoded, err := bfieldcodec.DecodeSlice[field.Element](encoded, func() field.Element {
		return field.Zero
	})
	if err != nil {
		log.Fatalf("Error decoding large slice: %v", err)
	}
	fmt.Printf("   Decoded slice length: %d\n", len(decoded))

	// Verify round-trip
	fmt.Printf("   Round-trip verification: %t\n", elementsEqual(elements, decoded))
}

func demonstrateErrorHandling() {
	// Test empty sequence
	empty := []field.Element{}
	encoded, err := bfieldcodec.EncodeSlice(empty)
	if err != nil {
		fmt.Printf("   Empty sequence error: %v\n", err)
	} else {
		fmt.Printf("   Empty sequence encoded: %v\n", encoded)
	}

	// Test invalid sequence
	invalid := []byte{0xFF, 0xFF, 0xFF, 0xFF} // Invalid sequence
	_, err = bfieldcodec.DecodeSlice[field.Element](invalid, func() field.Element {
		return field.Zero
	})
	if err != nil {
		fmt.Printf("   Invalid sequence error: %v\n", err)
	}
}

// Helper functions for comparison
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

func xFieldElementsEqual(a, b []xfield.XFieldElement) bool {
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
