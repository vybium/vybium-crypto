package hash

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/vybium/vybium-crypto/pkg/vybium-crypto/field"
)

// DigestLen is the number of field elements in a digest (equivalent to twenty-first's Digest::LEN).
const DigestLen = 5

// Digest represents the result of hashing a sequence of elements.
// It contains exactly 5 BFieldElements, matching twenty-first's Digest structure.
type Digest [DigestLen]field.Element

// NewDigest creates a new Digest from an array of field elements.
func NewDigest(elements [DigestLen]field.Element) Digest {
	return Digest(elements)
}

// Zero returns the all-zero digest.
func ZeroDigest() Digest {
	return Digest{field.Zero, field.Zero, field.Zero, field.Zero, field.Zero}
}

// Values returns the underlying field elements.
func (d Digest) Values() [DigestLen]field.Element {
	return [DigestLen]field.Element(d)
}

// Reversed returns a new digest with elements reversed.
// This function is an involutive endomorphism.
func (d Digest) Reversed() Digest {
	return Digest{d[4], d[3], d[2], d[1], d[0]}
}

// Equal returns true if two digests are equal.
func (d Digest) Equal(other Digest) bool {
	for i := 0; i < DigestLen; i++ {
		if !d[i].Equal(other[i]) {
			return false
		}
	}
	return true
}

// IsZero returns true if the digest is all zeros.
func (d Digest) IsZero() bool {
	for i := 0; i < DigestLen; i++ {
		if !d[i].IsZero() {
			return false
		}
	}
	return true
}

// String returns a human-readable string representation (comma-separated values).
func (d Digest) String() string {
	values := make([]string, DigestLen)
	for i := 0; i < DigestLen; i++ {
		values[i] = d[i].String()
	}
	return strings.Join(values, ",")
}

// Hex returns the hexadecimal representation of the digest.
func (d Digest) Hex() string {
	bytes := d.ToBytes()
	return hex.EncodeToString(bytes[:])
}

// ToBytes converts the digest to a byte array (40 bytes total: 5 elements × 8 bytes).
func (d Digest) ToBytes() [DigestLen * 8]byte {
	var result [DigestLen * 8]byte
	for i := 0; i < DigestLen; i++ {
		binary.LittleEndian.PutUint64(result[i*8:(i+1)*8], d[i].Value())
	}
	return result
}

// FromBytes creates a Digest from a byte array.
func DigestFromBytes(bytes [DigestLen * 8]byte) Digest {
	var result Digest
	for i := 0; i < DigestLen; i++ {
		value := binary.LittleEndian.Uint64(bytes[i*8 : (i+1)*8])
		result[i] = field.New(value)
	}
	return result
}

// FromHex creates a Digest from a hexadecimal string.
func DigestFromHex(s string) (Digest, error) {
	bytes, err := hex.DecodeString(s)
	if err != nil {
		return ZeroDigest(), fmt.Errorf("invalid hex string: %w", err)
	}
	if len(bytes) != DigestLen*8 {
		return ZeroDigest(), fmt.Errorf("invalid hex digest length: expected %d bytes, got %d", DigestLen*8, len(bytes))
	}
	var byteArray [DigestLen * 8]byte
	copy(byteArray[:], bytes)
	return DigestFromBytes(byteArray), nil
}

// Less returns true if this digest is less than the other (for ordering).
// Compares elements in reverse order (most significant first), matching twenty-first's Ord implementation.
func (d Digest) Less(other Digest) bool {
	for i := DigestLen - 1; i >= 0; i-- {
		if d[i].Less(other[i]) {
			return true
		}
		if d[i].Greater(other[i]) {
			return false
		}
	}
	return false
}

// Greater returns true if this digest is greater than the other.
func (d Digest) Greater(other Digest) bool {
	return other.Less(d)
}

// Clone creates a copy of the digest.
func (d Digest) Clone() Digest {
	return Digest{d[0], d[1], d[2], d[3], d[4]}
}
