// Package hash provides cryptographic hash functions optimized for STARKs.
//
// This package implements the Tip5 hash function, an arithmetization-oriented hash function
// designed for recursive STARKs. Tip5 uses a sponge construction with a permutation based on
// split-and-lookup operations, making it efficient to prove in zero-knowledge proof systems.
// Reference: https://eprint.iacr.org/2023/107.pdf
package hash

import (
	"github.com/vybium/vybium-crypto/pkg/vybium-crypto/field"
	"github.com/vybium/vybium-crypto/pkg/vybium-crypto/xfield"
)

// Tip5 constants defining the hash function parameters
const (
	StateSize         = 16
	NumSplitAndLookup = 4
	Log2StateSize     = 4
	Capacity          = 6
	Rate              = 10
	NumRounds         = 5
)

// Domain differentiates between modes of hashing.
type Domain int

const (
	// VariableLength is used for hashing objects that potentially serialize to more than RATE elements.
	VariableLength Domain = iota
	// FixedLength is used for hashing objects that always fit within RATE elements.
	FixedLength
)

// Tip5 represents the Tip5 hash function state.
// The state consists of 16 field elements that are permuted during hashing.
type Tip5 struct {
	state [StateSize]field.Element
}

// LookupTable is the lookup table used in Tip5 permutation.
// This table maps 8-bit values through a carefully chosen permutation for the split-and-lookup operation.
var LookupTable = [256]uint8{
	0, 7, 26, 63, 124, 215, 85, 254, 214, 228, 45, 185, 140, 173, 33, 240, 29, 177, 176, 32, 8,
	110, 87, 202, 204, 99, 150, 106, 230, 14, 235, 128, 213, 239, 212, 138, 23, 130, 208, 6, 44,
	71, 93, 116, 146, 189, 251, 81, 199, 97, 38, 28, 73, 179, 95, 84, 152, 48, 35, 119, 49, 88,
	242, 3, 148, 169, 72, 120, 62, 161, 166, 83, 175, 191, 137, 19, 100, 129, 112, 55, 221, 102,
	218, 61, 151, 237, 68, 164, 17, 147, 46, 234, 203, 216, 22, 141, 65, 57, 123, 12, 244, 54, 219,
	231, 96, 77, 180, 154, 5, 253, 133, 165, 98, 195, 205, 134, 245, 30, 9, 188, 59, 142, 186, 197,
	181, 144, 92, 31, 224, 163, 111, 74, 58, 69, 113, 196, 67, 246, 225, 10, 121, 50, 60, 157, 90,
	122, 2, 250, 101, 75, 178, 159, 24, 36, 201, 11, 243, 132, 198, 190, 114, 233, 39, 52, 21, 209,
	108, 238, 91, 187, 18, 104, 194, 37, 153, 34, 200, 143, 126, 155, 236, 118, 64, 80, 172, 89,
	94, 193, 135, 183, 86, 107, 252, 13, 167, 206, 136, 220, 207, 103, 171, 160, 76, 182, 227, 217,
	158, 56, 174, 4, 66, 109, 139, 162, 184, 211, 249, 47, 125, 232, 117, 43, 16, 42, 127, 20, 241,
	25, 149, 105, 156, 51, 53, 168, 145, 247, 223, 79, 78, 226, 15, 222, 82, 115, 70, 210, 27, 41,
	1, 170, 40, 131, 192, 229, 248, 255,
}

// RoundConstants are the round constants used in Tip5 permutation.
// Production implementation.
var RoundConstants = [NumRounds * StateSize]field.Element{
	field.New(13630775303355457758),
	field.New(16896927574093233874),
	field.New(10379449653650130495),
	field.New(1965408364413093495),
	field.New(15232538947090185111),
	field.New(15892634398091747074),
	field.New(3989134140024871768),
	field.New(2851411912127730865),
	field.New(8709136439293758776),
	field.New(3694858669662939734),
	field.New(12692440244315327141),
	field.New(10722316166358076749),
	field.New(12745429320441639448),
	field.New(17932424223723990421),
	field.New(7558102534867937463),
	field.New(15551047435855531404),
	field.New(17532528648579384106),
	field.New(5216785850422679555),
	field.New(15418071332095031847),
	field.New(11921929762955146258),
	field.New(9738718993677019874),
	field.New(3464580399432997147),
	field.New(13408434769117164050),
	field.New(264428218649616431),
	field.New(4436247869008081381),
	field.New(4063129435850804221),
	field.New(2865073155741120117),
	field.New(5749834437609765994),
	field.New(6804196764189408435),
	field.New(17060469201292988508),
	field.New(9475383556737206708),
	field.New(12876344085611465020),
	field.New(13835756199368269249),
	field.New(1648753455944344172),
	field.New(9836124473569258483),
	field.New(12867641597107932229),
	field.New(11254152636692960595),
	field.New(16550832737139861108),
	field.New(11861573970480733262),
	field.New(1256660473588673495),
	field.New(13879506000676455136),
	field.New(10564103842682358721),
	field.New(16142842524796397521),
	field.New(3287098591948630584),
	field.New(685911471061284805),
	field.New(5285298776918878023),
	field.New(18310953571768047354),
	field.New(3142266350630002035),
	field.New(549990724933663297),
	field.New(4901984846118077401),
	field.New(11458643033696775769),
	field.New(8706785264119212710),
	field.New(12521758138015724072),
	field.New(11877914062416978196),
	field.New(11333318251134523752),
	field.New(3933899631278608623),
	field.New(16635128972021157924),
	field.New(10291337173108950450),
	field.New(4142107155024199350),
	field.New(16973934533787743537),
	field.New(11068111539125175221),
	field.New(17546769694830203606),
	field.New(5315217744825068993),
	field.New(4609594252909613081),
	field.New(3350107164315270407),
	field.New(17715942834299349177),
	field.New(9600609149219873996),
	field.New(12894357635820003949),
	field.New(4597649658040514631),
	field.New(7735563950920491847),
	field.New(1663379455870887181),
	field.New(13889298103638829706),
	field.New(7375530351220884434),
	field.New(3502022433285269151),
	field.New(9231805330431056952),
	field.New(9252272755288523725),
	field.New(10014268662326746219),
	field.New(15565031632950843234),
	field.New(1209725273521819323),
	field.New(6024642864597845108),
}

// New creates a new Tip5 instance with the specified domain.
// Production implementation.
func New(domain Domain) *Tip5 {
	tip5 := &Tip5{}

	// Initialize state based on domain
	switch domain {
	case VariableLength:
		// Capacity is all zeros (default)
	case FixedLength:
		// Capacity is all ones
		for i := Rate; i < StateSize; i++ {
			tip5.state[i] = field.One
		}
	}

	return tip5
}

// Init creates a new Tip5 instance for variable-length hashing.
// This is the standard initialization for the Sponge trait.
func Init() *Tip5 {
	return New(VariableLength)
}

// Permutation applies the Tip5 permutation to the state.
// Production implementation.
func (t *Tip5) Permutation() {
	for i := 0; i < NumRounds; i++ {
		t.round(i)
	}
}

// round applies one round of the Tip5 permutation.
// Production implementation.
func (t *Tip5) round(roundIndex int) {
	t.sboxLayer()
	t.mdsGenerated()

	// Add round constants
	for i := 0; i < StateSize; i++ {
		t.state[i] = t.state[i].Add(RoundConstants[roundIndex*StateSize+i])
	}
}

// sboxLayer applies the S-box layer.
// Production implementation.
func (t *Tip5) sboxLayer() {
	// Split-and-lookup for first NUM_SPLIT_AND_LOOKUP elements
	for i := 0; i < NumSplitAndLookup; i++ {
		splitAndLookup(&t.state[i])
	}

	// Power map (x^7) for remaining elements
	for i := NumSplitAndLookup; i < StateSize; i++ {
		sq := t.state[i].Square()               // x^2
		qu := sq.Square()                       // x^4
		t.state[i] = t.state[i].Mul(sq).Mul(qu) // x * x^2 * x^4 = x^7
	}
}

// splitAndLookup applies the split-and-lookup operation.
// Production implementation.
func splitAndLookup(element *field.Element) {
	// Get raw bytes in Montgomery form
	bytes := element.ToBytes()

	// Apply lookup table to each byte
	for i := 0; i < 8; i++ {
		bytes[i] = LookupTable[bytes[i]]
	}

	// Create element from modified bytes (still in Montgomery form)
	*element = field.FromBytes(bytes)
}

// mdsGenerated applies the MDS matrix using the optimized generated function.
// Production implementation.
func (t *Tip5) mdsGenerated() {
	var lo, hi [StateSize]uint64

	// Split into low and high 32-bit parts
	for i := 0; i < StateSize; i++ {
		b := t.state[i].RawValue()
		hi[i] = b >> 32
		lo[i] = b & 0xFFFFFFFF
	}

	// Apply generated function to each half
	lo = generatedFunction(lo)
	hi = generatedFunction(hi)

	// Recombine and reduce
	for r := 0; r < StateSize; r++ {
		s := (uint128(lo[r]) >> 4) + (uint128(hi[r]) << 28)

		sHi := uint64(s >> 32) // Fixed: shift by 32 instead of 64
		sLo := uint64(s)

		// Compute result with overflow handling
		res := sLo + sHi*0xFFFFFFFF
		over := res < sLo // overflow check

		if over {
			res += 0xFFFFFFFF
		}

		t.state[r] = field.NewFromRaw(res)
	}
}

// uint128 is a simple 128-bit unsigned integer for intermediate calculations
type uint128 uint64 // Simplified for shift operations

// generatedFunction is the optimized MDS matrix multiplication.
// This is a direct port of twenty-first's Tip5::generated_function()
//
// This function was automatically generated to optimize the MDS matrix multiplication.
// It uses a clever factorization to reduce the number of operations.
func generatedFunction(input [StateSize]uint64) [StateSize]uint64 {
	// First layer of additions
	node34 := input[0] + input[8]
	node38 := input[4] + input[12]
	node36 := input[2] + input[10]
	node40 := input[6] + input[14]
	node35 := input[1] + input[9]
	node39 := input[5] + input[13]
	node37 := input[3] + input[11]
	node41 := input[7] + input[15]

	// Second layer
	node50 := node34 + node38
	node52 := node36 + node40
	node51 := node35 + node39
	node53 := node37 + node41

	// Subtractions
	node160 := input[0] - input[8]
	node161 := input[1] - input[9]
	node165 := input[5] - input[13]
	node163 := input[3] - input[11]
	node167 := input[7] - input[15]
	node162 := input[2] - input[10]
	node166 := input[6] - input[14]
	node164 := input[4] - input[12]

	// Third layer
	node58 := node50 + node52
	node59 := node51 + node53
	node90 := node34 - node38
	node91 := node35 - node39
	node93 := node37 - node41
	node92 := node36 - node40

	// Multiplications and combinations
	node64 := (node58 + node59) * 524757
	node67 := (node58 - node59) * 52427
	node71 := node50 - node52
	node72 := node51 - node53

	node177 := node161 + node165
	node179 := node163 + node167
	node178 := node162 + node166
	node176 := node160 + node164

	node69 := node64 + node67
	node397 := node71*18446744073709525744 - node72*53918
	node1857 := node90 * 395512
	node99 := node91 + node93
	node1865 := node91 * 18446744073709254400
	node1869 := node93 * 179380
	node1873 := node92 * 18446744073709509368
	node1879 := node160 * 35608
	node185 := node161 + node163
	node1915 := node161 * 18446744073709340312
	node1921 := node163 * 18446744073709494992
	node1927 := node162 * 18446744073709450808
	node228 := node165 + node167
	node1939 := node165 * 18446744073709420056
	node1945 := node167 * 18446744073709505128
	node1951 := node166 * 216536
	node1957 := node164 * 18446744073709515080

	node70 := node64 - node67
	node702 := node71*53918 + node72*18446744073709525744
	node1961 := node90 * 18446744073709254400
	node1963 := node91 * 395512
	node1965 := node92 * 179380
	node1967 := node93 * 18446744073709509368
	node1970 := node160 * 18446744073709340312
	node1973 := node161 * 35608
	node1982 := node162 * 18446744073709494992
	node1985 := node163 * 18446744073709450808
	node1988 := node166 * 18446744073709505128
	node1991 := node167 * 216536
	node1994 := node164 * 18446744073709420056
	node1997 := node165 * 18446744073709515080

	node98 := node90 + node92
	node184 := node160 + node162
	node227 := node164 + node166

	node86 := node69 + node397
	node403 := node1857 - (node99*18446744073709433780 - node1865 - node1869 + node1873)
	node271 := node177 + node179
	node1891 := node177 * 18446744073709208752
	node1897 := node179 * 18446744073709448504
	node1903 := node178 * 115728
	node1909 := node185 * 18446744073709283688
	node1933 := node228 * 18446744073709373568

	node88 := node70 + node702
	node708 := node1961 + node1963 - (node1965 + node1967)
	node1976 := node178 * 18446744073709448504
	node1979 := node179 * 115728

	node87 := node69 - node397
	node897 := node1865 + node98*353264 - node1857 - node1873 - node1869
	node2007 := node184 * 18446744073709486416
	node2013 := node227 * 180000

	node89 := node70 - node702
	node1077 := node98*18446744073709433780 + node99*353264 - (node1961 + node1963) - (node1965 + node1967)
	node2020 := node184 * 18446744073709283688
	node2023 := node185 * 18446744073709486416
	node2026 := node227 * 18446744073709373568
	node2029 := node228 * 180000
	node2035 := node176 * 18446744073709550688
	node2038 := node176 * 18446744073709208752
	node2041 := node177 * 18446744073709550688

	node270 := node176 + node178

	node152 := node86 + node403
	node412 := node1879 + node185*18446744073709433780 - node1915 - node1921 - node1927
	node1237 := node2035 - node1891 - node1897 - node1903 - node1909

	node154 := node88 + node708
	node717 := node1921 + node2007 - node1970 - node1973 - node1982 - node1985
	node1375 := node1927 + node2013 - node1994 - node1997 - node1988 - node1991

	node156 := node87 + node897
	node906 := node1873 + node1909 + node2020 - node1879 - node1915 - node1921 - node1927
	node1492 := node1951 + node1933 + node2026 - node1939 - node1945 - node1957 - node1997

	node158 := node89 + node1077
	node1086 := node1961 + node1963 + node1979 + node2023 - node1973 - node1982 - node1985 - node1976
	node1657 := node1994 + node1997 + node1991 + node2029 - node1939 - node1945 - node1957 - node1988

	node153 := node270*114800 + node271*18446744073709433780 - node2038 - node2041 - node1976 - node1979 - (node2020 + node2023 - node1970 - node1973 - node1982 - node1985) - (node2026 + node2029 - node1994 - node1997 - node1988 - node1991)
	node155 := node270*18446744073709433780 + node271*114800 - node1891 - node1897 - node1903 - (node1879 + node1909 + node2020 - node1915 - node1921 - node1927) - (node1939 + node1933 + node2026 - node1951 - node1957 - node1988 - node1991)
	node157 := node1879 + node270*353264 - node2035 - node2038 - node2041 - node1976 - node1979 - (node1915 + node1909 + node2020 + node2023 - node1927 - node1982 - node1985 - node1973) - (node1939 + node1933 + node2026 + node2029 - node1951 - node1957 - node1988 - node1991)
	node159 := node1939 + node271*114800 - node2038 - node2041 - node1976 - node1979 - (node2020 + node2023 - node1970 - node1973 - node1982 - node1985) - (node2026 + node2029 - node1994 - node1997 - node1988 - node1991)

	return [StateSize]uint64{
		node152 + node412,
		node154 + node717,
		node156 + node906,
		node158 + node1086,
		node153 + node1237,
		node155 + node1375,
		node157 + node1492,
		node159 + node1657,
		node152 - node412,
		node154 - node717,
		node156 - node906,
		node158 - node1086,
		node153 - node1237,
		node155 - node1375,
		node157 - node1492,
		node159 - node1657,
	}
}

// Hash10 hashes exactly 10 BFieldElements (one rate's worth).
// Production implementation.
func Hash10(input [Rate]field.Element) [DigestLen]field.Element {
	sponge := New(FixedLength)

	// Absorb once
	copy(sponge.state[:Rate], input[:])
	sponge.Permutation()

	// Squeeze once
	var digest [DigestLen]field.Element
	copy(digest[:], sponge.state[:DigestLen])
	return digest
}

// HashPair hashes two digests together.
// Production implementation.
func HashPair(left, right [DigestLen]field.Element) [DigestLen]field.Element {
	sponge := New(FixedLength)
	copy(sponge.state[:DigestLen], left[:])
	copy(sponge.state[DigestLen:2*DigestLen], right[:])

	sponge.Permutation()

	var digest [DigestLen]field.Element
	copy(digest[:], sponge.state[:DigestLen])
	return digest
}

// HashVarlen hashes a variable-length sequence of BFieldElements.
// Production implementation.
func HashVarlen(input []field.Element) [DigestLen]field.Element {
	sponge := Init()
	sponge.PadAndAbsorbAll(input)

	var digest [DigestLen]field.Element
	copy(digest[:], sponge.state[:DigestLen])
	return digest
}

// Tip5Permutation applies the Tip5 permutation to a 5-element state.
func Tip5Permutation(state [5]field.Element) [5]field.Element {
	tip5 := New(VariableLength)
	// Copy state to tip5 internal state
	for i := 0; i < 5; i++ {
		tip5.state[i] = state[i]
	}
	// Apply permutation
	tip5.Permutation()
	// Return the permuted state
	var result [5]field.Element
	for i := 0; i < 5; i++ {
		result[i] = tip5.state[i]
	}
	return result
}

// Absorb absorbs RATE elements into the sponge.
// This implements part of the Sponge trait.
func (t *Tip5) Absorb(input [Rate]field.Element) {
	copy(t.state[:Rate], input[:])
	t.Permutation()
}

// Squeeze squeezes RATE elements from the sponge.
// This implements part of the Sponge trait.
func (t *Tip5) Squeeze() [Rate]field.Element {
	var output [Rate]field.Element
	copy(output[:], t.state[:Rate])
	t.Permutation()
	return output
}

// PadAndAbsorbAll pads and absorbs all input elements.
// Production implementation.
func (t *Tip5) PadAndAbsorbAll(input []field.Element) {
	// Process full chunks
	for i := 0; i < len(input); i += Rate {
		end := i + Rate
		if end > len(input) {
			// Last chunk needs padding
			var lastChunk [Rate]field.Element
			remaining := len(input) - i
			copy(lastChunk[:remaining], input[i:])
			lastChunk[remaining] = field.One // Padding: [1, 0, 0, ...]
			t.Absorb(lastChunk)
		} else {
			// Full chunk
			var chunk [Rate]field.Element
			copy(chunk[:], input[i:end])
			t.Absorb(chunk)
		}
	}

	// If input was empty or a multiple of Rate, add padding chunk
	if len(input)%Rate == 0 {
		var paddingChunk [Rate]field.Element
		paddingChunk[0] = field.One
		t.Absorb(paddingChunk)
	}
}

// Trace returns the trace of applying the permutation.
// This is functionally equivalent to Permutation() but returns the state
// after each round, including the initial state.
// Production implementation.
func (t *Tip5) Trace() [1 + NumRounds][StateSize]field.Element {
	var trace [1 + NumRounds][StateSize]field.Element

	// Save initial state
	trace[0] = t.state

	// Apply rounds and save state after each
	for i := 0; i < NumRounds; i++ {
		t.round(i)
		trace[1+i] = t.state
	}

	return trace
}

// SampleIndices produces numIndices random integer values in the range [0, upperBound).
// The upperBound must be a power of 2.
//
// This method uses von Neumann rejection sampling.
// Specifically, if the top 32 bits of a BFieldElement are all ones, then the bottom 32 bits
// are not uniformly distributed, and so they are dropped. This method invokes squeeze until
// enough uniform u32s have been sampled.
//
// Production implementation.
func (t *Tip5) SampleIndices(upperBound uint32, numIndices int) []uint32 {
	// Verify upperBound is a power of 2
	if upperBound == 0 || (upperBound&(upperBound-1)) != 0 {
		panic("upperBound must be a power of 2")
	}

	indices := make([]uint32, 0, numIndices)
	var squeezedElements []field.Element

	for len(indices) < numIndices {
		if len(squeezedElements) == 0 {
			squeezed := t.Squeeze()
			// Reverse the order to match twenty-first's behavior
			squeezedElements = make([]field.Element, Rate)
			for i := 0; i < Rate; i++ {
				squeezedElements[Rate-1-i] = squeezed[i]
			}
		}

		element := squeezedElements[len(squeezedElements)-1]
		squeezedElements = squeezedElements[:len(squeezedElements)-1]

		// Reject if element is MAX (top 32 bits all ones)
		if element != field.Max {
			indices = append(indices, uint32(element.Value())%upperBound)
		}
	}

	return indices
}

// SampleScalars produces numElements random XFieldElement values.
//
// If numElements is not divisible by RATE, spill the remaining elements of the last squeeze.
//
// Production implementation.
func (t *Tip5) SampleScalars(numElements int) ([]xfield.XFieldElement, error) {
	// Import xfield here to avoid circular dependency
	// We'll need to import it at the top of the file
	const extensionDegree = 3

	numSqueezes := (numElements*extensionDegree + Rate - 1) / Rate // Ceiling division

	// Collect all squeezed elements
	allElements := make([]field.Element, 0, numSqueezes*Rate)
	for i := 0; i < numSqueezes; i++ {
		squeezed := t.Squeeze()
		allElements = append(allElements, squeezed[:]...)
	}

	// Group into XFieldElements (3 elements each)
	scalars := make([]xfield.XFieldElement, 0, numElements)
	for i := 0; i < numElements && i*extensionDegree+extensionDegree <= len(allElements); i++ {
		start := i * extensionDegree
		coeffs := [extensionDegree]field.Element{
			allElements[start],
			allElements[start+1],
			allElements[start+2],
		}
		scalars = append(scalars, xfield.New(coeffs))
	}

	return scalars, nil
}
