# Vybium Crypto Examples

This directory contains practical examples demonstrating the key features of the Vybium Crypto cryptographic library.

## Examples Overview

### 🔐 Basic Hash Examples (`basic_hash.go`)

Demonstrates the core hash functions:

- **Tip5 Hash**: STARK-optimized hash function with sponge construction
- **Poseidon Hash**: Advanced hash function with enhanced security
- **Variable Length Hashing**: Processing different input sizes
- **Hash Pairs**: Efficient hashing of two elements

**Usage:**

```bash
go run examples/basic_hash.go
```

### 🔢 Field Operations (`field_operations.go`)

Shows Goldilocks field arithmetic:

- **Basic Operations**: Addition, subtraction, multiplication, division
- **Field Properties**: Identity elements, inverses, and properties
- **Large Numbers**: Handling big integers and field operations
- **Montgomery Representation**: Efficient field arithmetic

**Usage:**

```bash
go run examples/field_operations.go
```

### 🌳 Merkle Tree Examples (`merkle_tree.go`)

Demonstrates authentication structures:

- **Basic Merkle Trees**: Building and root computation
- **Merkle Proofs**: Generating and verifying inclusion proofs
- **MMR (Merkle Mountain Range)**: Efficient append-only structures
- **Proof Verification**: Cryptographic authentication

**Usage:**

```bash
go run examples/merkle_tree.go
```

### 📐 Polynomial Operations (`polynomial_operations.go`)

Shows polynomial arithmetic and NTT:

- **Basic Polynomial Operations**: Addition, subtraction, multiplication
- **NTT (Number Theoretic Transform)**: Fast polynomial operations
- **Polynomial Multiplication**: Using NTT for efficiency
- **Polynomial Evaluation**: Computing values at specific points

**Usage:**

```bash
go run examples/polynomial_operations.go
```

### 🔗 Extension Field Examples (`extension_field.go`)

Demonstrates F_p^3 operations:

- **Basic Extension Field Operations**: Arithmetic in F_p^3
- **Extension Field Properties**: Identity elements and inverses
- **Extension Field Arithmetic**: Powers and division
- **Field Conversion**: Between base and extension fields

**Usage:**

```bash
go run examples/extension_field.go
```

### 📦 Serialization Examples (`serialization.go`)

Shows BFieldCodec usage:

- **Field Element Serialization**: Encoding/decoding field elements
- **Extension Field Serialization**: Handling XFieldElement
- **Slice Serialization**: Processing large arrays
- **Error Handling**: Robust error management

**Usage:**

```bash
go run examples/serialization.go
```

## Running All Examples

To run all examples at once:

```bash
# Run all examples
for example in examples/*.go; do
    echo "Running $example..."
    go run "$example"
    echo "---"
done
```

## Example Output

Each example provides detailed output showing:

- **Input values** and operations performed
- **Results** of cryptographic operations
- **Verification** of mathematical properties
- **Performance** characteristics where relevant

## Key Features Demonstrated

### Cryptographic Primitives

- **Field Arithmetic**: Goldilocks prime field operations
- **Hash Functions**: Tip5 and Poseidon implementations
- **Extension Fields**: Degree-3 extension field operations
- **Polynomial Operations**: NTT-accelerated polynomial arithmetic

### Data Structures

- **Merkle Trees**: Authentication and verification
- **MMR**: Efficient append-only data structures
- **Serialization**: BFieldCodec for proof system integration

### Mathematical Operations

- **Montgomery Representation**: Efficient field arithmetic
- **NTT**: Fast polynomial multiplication
- **Batch Operations**: Efficient cryptographic computations
- **Proof Systems**: Cryptographic authentication

## Integration Examples

These examples show how to integrate Vybium Crypto into:

- **Zero-Knowledge Proof Systems**: STARK proof generation
- **Blockchain Applications**: Cryptographic primitives
- **Cryptographic Protocols**: Hash functions and field operations
- **Data Authentication**: Merkle trees and proofs

## Performance Notes

- **Field Operations**: Optimized with Montgomery representation
- **Hash Functions**: STARK-optimized for proof systems
- **Polynomial Operations**: NTT acceleration for large polynomials
- **Batch Operations**: Efficient processing of multiple elements

## Security Considerations

- **Constant-Time Operations**: Timing-attack resistant implementations
- **Memory Safety**: Go's type system ensures safety
- **Input Validation**: Comprehensive sanitization
- **Side-Channel Resistance**: Protected against analysis attacks

---

For more information, see the main [README.md](../README.md) and [API documentation](../pkg/vybium-crypto/).
