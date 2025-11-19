# Changelog

All notable changes to Vybium Crypto will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.0] - 2025-11-10

### Initial Release

Core cryptographic primitives for zero-knowledge proof systems.

#### Added
- **Field Arithmetic**
  - Goldilocks field operations (2^64 - 2^32 + 1)
  - Extension field (degree 3) operations
  - Montgomery arithmetic optimization
  - Batch operations support

- **Hash Functions**
  - Tip5 hash with sponge construction
  - Poseidon hash function
  - Flexible sponge abstraction

- **Polynomial Operations**
  - Fast NTT (Number Theoretic Transform)
  - Polynomial evaluation and interpolation
  - Zerofier tree structures
  - Efficient coefficient operations

- **Merkle Structures**
  - Merkle tree implementation
  - MMR (Merkle Mountain Range) support
  - Proof generation and verification
  - Batch operations

- **Serialization**
  - BFieldCodec for canonical encoding
  - Zero-copy serialization
  - Proof system integration

- **Testing**
  - Comprehensive unit tests
  - Property-based testing
  - Fuzz testing for field operations
  - Benchmark suite

#### Performance
- Field operations: 1-2 ns per operation
- Tip5 hash: ~500 ns for 16 elements
- Poseidon hash: ~300 ns
- NTT (2^16): ~50 ms

#### Documentation
- Complete API documentation
- Usage examples for all modules
- Performance benchmarks
- Security considerations

[0.1.0]: https://github.com/vybium/vybium-crypto/releases/tag/v0.1.0
