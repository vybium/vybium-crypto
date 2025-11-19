# vybium-crypto

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/vybium/vybium-crypto)](https://goreportcard.com/report/github.com/vybium/vybium-crypto)

A high-performance collection of cryptographic primitives in Go, optimized for zero-knowledge proof systems.

## Overview

**vybium-crypto** provides production-ready implementations of cryptographic primitives specifically designed for STARK-based proof systems.

## Features

### Hash Functions

- **Tip5** - STARK-optimized hash function with sponge construction
  - State size: 16 elements, 5 rounds, Rate: 10
  - Optimized for STARK proof systems

- **ARION** - Arithmetization-oriented hash based on Generalized Triangular Dynamical Systems
  - State size: 3 elements, 10 rounds, Rate: 2
  - GTDS-based non-linear layer

- **Poseidon** - Zero-knowledge friendly hash for constraint systems
  - Configurable security levels (128-bit, 256-bit)
  - Dynamic parameter generation with Grain LFSR

### Field Arithmetic

- **BFieldElement** (`field` package)
  - The prime-field type F_p where p = 2^64 - 2^32 + 1 (Goldilocks prime)
  - Montgomery representation for efficient modular arithmetic
  - Constant-time operations for timing attack resistance

- **XFieldElement** (`xfield` package)
  - Extension field F_p[x]/(x^3 - x + 1)
  - Degree-3 extension for advanced cryptographic protocols
  - Optimized multiplication and inversion

- **BFieldCodec** (`bfieldcodec` package)
  - Canonical encoding/decoding for proof system integration
  - Zero-copy serialization where possible
  - Type-safe codec trait implementation

### Number Theoretic Transform (NTT)

- **Fast Polynomial Multiplication** (`ntt` package)
  - Radix-2 Cooley-Tukey FFT algorithm
  - Efficient polynomial operations in frequency domain

### Polynomial Operations

- **Univariate Polynomials** (`polynomial` package)
  - Fast evaluation and interpolation
  - Division with remainder
  - Zerofier tree construction
  - Batch operations for proof generation

### Merkle Structures

- **Merkle Trees** (`merkle` package)
  - Efficient authentication with batch verification
  - Inclusion and exclusion proofs
  - Compatible with Tip5, ARION, and Poseidon hash functions

- **Merkle Mountain Ranges** (`mmr` package)
  - Append-only Merkle structures
  - Efficient proofs for growing datasets
  - Optimized for blockchain applications

## Installation

```bash
go get github.com/vybium/vybium-crypto
```

## Quick Start

### Field Arithmetic

```go
package main

import (
    "fmt"
    "github.com/vybium/vybium-crypto/pkg/vybium-crypto/field"
)

func main() {
    // Create field elements
    a := field.New(10)
    b := field.New(20)

    // Perform operations
    sum := a.Add(b)
    product := a.Mul(b)

    fmt.Printf("Sum: %v\n", sum)
    fmt.Printf("Product: %v\n", product)
}
```

### Hash Functions

```go
package main

import (
    "fmt"
    "github.com/vybium/vybium-crypto/pkg/vybium-crypto/field"
    "github.com/vybium/vybium-crypto/pkg/vybium-crypto/hash"
)

func main() {
    // Hash with Tip5
    input := []field.Element{field.New(1), field.New(2), field.New(3)}
    digest := hash.Tip5HashVarLen(input)
    fmt.Printf("Tip5 Hash: %v\n", digest)

    // Hash with ARION
    digest2 := hash.ArionHash(input)
    fmt.Printf("ARION Hash: %v\n", digest2)
}
```

### Merkle Trees

```go
package main

import (
    "fmt"
    "github.com/vybium/vybium-crypto/pkg/vybium-crypto/field"
    "github.com/vybium/vybium-crypto/pkg/vybium-crypto/hash"
    "github.com/vybium/vybium-crypto/pkg/vybium-crypto/merkle"
)

func main() {
    // Create leaves
    leaves := []hash.Digest{
        hash.Tip5HashVarLen([]field.Element{field.New(1)}),
        hash.Tip5HashVarLen([]field.Element{field.New(2)}),
        hash.Tip5HashVarLen([]field.Element{field.New(3)}),
        hash.Tip5HashVarLen([]field.Element{field.New(4)}),
    }

    // Build tree
    tree := merkle.New(leaves, hash.Tip5HashPair)
    root := tree.Root()

    // Generate proof
    proof := tree.AuthenticationStructure(1)

    // Verify proof
    valid := tree.Verify(leaves[1], 1, proof, root)
    fmt.Printf("Proof valid: %v\n", valid)
}
```

## Performance

Benchmarks run on Go 1.23, Linux (results from `go test -bench=.`):

### Field Operations

| Operation | Time       | Notes                 |
| --------- | ---------- | --------------------- |
| Add       | 0.26 ns/op | Constant-time         |
| Sub       | 0.26 ns/op | Constant-time         |
| Mul       | 2.3 ns/op  | Montgomery            |
| Inverse   | 464 ns/op  | Extended GCD          |
| ModPow    | 111 ns/op  | Binary exponentiation |

### Hash Functions

| Operation          | Time      | Notes               |
| ------------------ | --------- | ------------------- |
| Tip5 Permutation   | 1.9 µs/op | 16-state, 5 rounds  |
| Tip5 Hash10        | 1.9 µs/op | Fixed-size          |
| ARION Permutation  | 8.8 µs/op | 3-state, 10 rounds  |
| ARION Hash10       | 54 µs/op  | Fixed-size          |
| Poseidon (128-bit) | 63 µs/op  | Full+Partial rounds |

## Documentation

Full API documentation is available via `go doc`:

```bash
go doc github.com/vybium/vybium-crypto/pkg/vybium-crypto/field
go doc github.com/vybium/vybium-crypto/pkg/vybium-crypto/hash
go doc github.com/vybium/vybium-crypto/pkg/vybium-crypto/merkle
```

## Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run benchmarks
go test -bench=. ./...
```

## Project Structure

```
vybium-crypto/
├── pkg/vybium-crypto/
│   ├── field/          # Goldilocks field arithmetic
│   ├── xfield/         # Extension field (degree-3)
│   ├── hash/           # Tip5, ARION, Poseidon
│   ├── ntt/            # Number Theoretic Transform
│   ├── polynomial/     # Polynomial operations
│   ├── merkle/         # Merkle trees and proofs
│   ├── mmr/            # Merkle Mountain Ranges
│   ├── bfieldcodec/    # Canonical encoding
│   ├── sponge/         # Sponge construction
│   └── zerofier/       # Zerofier polynomials
├── examples/           # Example usage
└── docs/               # Additional documentation
```

## Security Considerations

- All field operations are implemented with constant-time algorithms to prevent timing attacks
- No unsafe code is used in the core cryptographic primitives
- The library has been designed for production use in zero-knowledge proof systems

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## Status

Version: 0.1.0
Status: Production Ready
Last Updated: November 2025

Built for zero-knowledge proof systems
