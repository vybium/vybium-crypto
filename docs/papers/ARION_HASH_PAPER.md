Arion: Arithmetization-Oriented Permutation
and Hashing from Generalized Triangular
Dynamical Systems
Arnab Roy1 , Matthias Johann Steiner1 , and Stefano Trevisani2
1
Alpen-Adria-Universität Klagenfurt, Universitätsstraße 65-67, 9020 Klagenfurt am
Wörthersee, Austria
firstname.lastname@aau.at
2
Technische Universität Wien, Karlsplatz 13, 1040 Wien, Austria
firstname.lastname@tuwien.ac.at
Abstract. In this paper we propose the (keyed) permutation Arion and
the hash function ArionHash over Fp for odd and particularly large primes.
The design of Arion is based on the newly introduced Generalized Trian-
gular Dynamical System (GTDS), which provides a new algebraic frame-
work for constructing (keyed) permutation using polynomials over a ﬁnite
ﬁeld. At round level Arion is the ﬁrst design which is instantiated using
the new GTDS. We provide extensive security analysis of our construc-
tion including algebraic cryptanalysis (e.g. interpolation and Gröbner
basis attacks) that are particularly decisive in assessing the security of
permutations and hash functions over Fp . From an application perspec-
tive, ArionHash aims for eﬃcient implementation in zkSNARK protocols
and Zero-Knowledge proof systems. For this purpose, we exploit that
CCZ-equivalence of graphs can lead to a more eﬃcient implementation
of Arithmetization-Oriented primitives.
We compare the eﬃciency of ArionHash in R1CS and Plonk settings with
other hash functions such as Poseidon, Anemoi and Griffin. For demon-
strating the practical eﬃciency of ArionHash we implemented it with the
zkSNARK libraries libsnark and Dusk Network Plonk. Our result shows
that ArionHash is signiﬁcantly faster than Poseidon - a hash function
designed for zero-knowledge proof systems. We also found that an aggres-
sive version of ArionHash is considerably faster than Anemoi and Griffin
in a practical zkSNARK setting.
1
Introduction
With the advancement of Zero-Knowledge (ZK), Multi-Party Computation (MPC)
and Fully Homomorphic Encryption (FHE) in recent years new efficiency mea-
sures for symmetric-key primitives allowing efficient implementation in these
schemes, namely low multiplicative complexity and low multiplicative depth,
have been introduced. The block ciphers, permutations and hash functions with
low multiplicative complexity are also referred to as Arithmetization-Oriented
(AO) primitives. A significant number of these new types of AO primitives aredefined over large finite fields of prime order p ≫ 2 for target applications. Our
focus in this paper will be such a low multiplicative complexity construction
over Fp for large primes. Some generic definitions and results in this paper are
applicable to any odd prime, thus we describe these results and definitions ac-
cordingly. However, the security of the construction(s) will be analyzed only for
large primes.
To put this paper into context with previous AO constructions we give a short
overview of their developments. The AO primitives proposed in the literature
until now can be categorized into three generations.
Gen I: LowMC [3], MiMC [2]
Gen II: Hades [34], Poseidon [33], GMiMC [1], Rescue-Prime [4]
Gen III: Reinforced Concrete [32], Griffin [31], Anemoi [16], Arion (this
paper)
The first generation consists of constructions which demonstrated that one
can construct secure and efficient ciphers and hash functions with low degree
primitives at round level. In particular, LowMC introduced the partial Substitution
Permutation Network (SPN) strategy in AO.
In the second generation researchers tried to obtain further efficiency im-
provements from Feistel and (partial) SPNs to obtain new efficient primitives.
Moreover, more focus was given on constructions native in large prime fields
Fp rather than F2n . This resulted in Hades which combines full and partial
SPNs over Fp , and its derived sponge function Poseidon which is now a widely
deployed hash function for ZK applications.
The current third generation adopts new design principles which neither re-
duce to the Feistel nor the SPN that culminated in the Generalized Triangular
Dynamical System (GTDS) [48]. Moreover, this generation diverted from the
consensus that one needs low degree polynomials to instantiate a secure and
efficient AO primitive.
In this paper we propose new AO primitives - Arion (block cipher) and the
hash function derived from it ArionHash. At round level Arion (and ArionHash)
like Griffin, utilize(s) a polynomial of very high degree in one branch and low
degree polynomials in the remaining branches to significantly cut the number
of necessary rounds compared to the previous generations. Anemoi also utilizes
a high degree permutation, the so-called open Flystel, at round level, but to
limit the number of constraints in a prover circuit the authors proved that the
open Flystel is CCZ-equivalent (cf. [19] and [16, §4.2]) to a low degree function,
the so-called closed Flystel. Lastly, Reinforced Concrete is the first AO hash
function that utilizes look-up tables which significantly reduces the number of
necessary rounds of Reinforced Concrete and consequently also the number
of constraints in a prover circuit.
1.1
Our Results
In this paper we propose the block cipher Arion and the hash function ArionHash
(Section 2), using the Generalized Triangular Dynamical System [48]. The block
2cipher and hash function are constructed over Fp with the target to achieve low
multiplicative complexity in a prover circuit. Utilizing the structure of GTDS
enables us to provide a systematic security analysis of the newly proposed block
cipher and hash function. The GTDS structure also allows us to choose the best
suited parameters for the efficiency. We provide extensive security analysis of the
proposed block cipher and hash function against state-of-the-art cryptanalysis
techniques to justify their security (Section 3). Our construction validates the
soundness of the generic GTDS structure that uses polynomial dynamical system
for constructing cryptographic permutations over finite fields.
Although Arion and ArionHash are defined on arbitrary finite fields Fp , the
parameters of the block cipher and hash function are chosen in such way to
be compatible with the primes chosen for the target ZK application namely, for
BLS12 and BN254 curves. We propose aggressive versions of Arion and ArionHash
namely, α-Arion and α-ArionHash. The difference between Arion (and ArionHash)
and its aggressive version is that the former avoids a recently proposed proba-
bilistic Gröbner basis [27] attack (Section 3.2 and Appendix C in the full version
of the paper [49]).
To demonstrate and compare the efficiencies of our constructions (Section 4)
we implemented them using the zkSNARK libraries libsnark [50], a C++ library
used in the privacy protecting digital currency Zcash [35], and Dusk Network
Plonk [23], a Rust library used in the privacy-oriented blockchain protocol Dusk.
Our results show that ArionHash is significantly (2x) faster than Poseidon -
an efficient hash function designed for zkSNARK applications. The efficiency
of ArionHash is comparable to the recently proposed (but not yet published
at a peer-reviewed venue) hash functions Anemoi and Griffin. We find that α-
ArionHash for practical choices of parameters in a Merkle tree mode of hashing is
faster than Griffin. We also reveal that CCZ-equivalence between the graphs of
the ArionHash GTDS and another closely related GTDS leads to a more efficient
implementation of ArionHash in ZK applications compared to the naive circuit
for ArionHash (Section 4.1).
Our public GitHub repository
https://github.com/sca-research/Arion
contains reference implementations in SageMath, C++ and Rust, our OSCAR im-
plementation to perform Gröbner basis experiments, and our SageMath code to
estimate the security levels of Arion & ArionHash.
2The (Keyed) Permutation and Hash Function
2.1Overview of the Design Rationale
Before we define Arion and ArionHash, we quickly summarize the design rationale
behind our construction.
– By utilizing the GTDS to instantiate the permutation we aim to achieve
fast degree growth in each component like in SPNs and non-linear mixing
3between the components as in Feistel network. Our GTDS, see Definition 1,
incorporates the strength of both SPN and Feistel in a single primitive at
round level.
– It follows from the generic security analysis in [48, §5] that the univariate
permutations, the SPN part, of the GTDS determine worst case security
bounds against differential and linear cryptanalysis. Hence, we chose param-
eters that minimize these bounds.
– To thwart interpolation attacks we opted for a design that can achieve a
degree overflow in the input variables in the first round, see Lemma 2 and
Table 2. This is achieved in the SPN part of the GTDS by applying a low de-
gree univariate permutation polynomial p1 in all branches except the last one
and by applying a high-degree inverse permutation p−1
2 in the last branch.
– We opted for a linear layer that mixes all branches in every round. This is
achieved with a circulant matrix which has only non-zero entries.
– For the high degree inverse permutation, the naive circuit for p−1
2 (x) = y
introduces many multiplicative constraints, though one can always transform
such a circuit into a circuit for x = p2 (y) in constant time, see Section 4.1.
This trick plays a fundamental role in the efficiency of ArionHash circuits.
2.2
Keyed Permutation
We start with the definition of the generalized triangular dynamical system of
Arion.
Definition 1 (GTDS of Arion). Let p ∈ Z>4 be a prime, and let Fp be the
field with p elements. Let n, d1 , d2 , e ∈ Z>1 be integers such that
(i) d1 is the smallest positive integer such that gcd (d1 , p − 1) = 1,
(ii) d2 is an arbitrary integer such that gcd (d2 , p − 1) = 1, and
(iii) e · d2 ≡ 1 mod p − 1.
For 1 ≤ i ≤ n − 1 let αi,1 , αi,2 , βi ∈ Fp be such that α2i,1 − 4 · αi,2 is a quadratic
non-residue modulo p. The generalized triangular dynamical system FArion =
{f1 , . . . , fn } of Arion is defined as
fi (x1 , . . . , xn ) = xdi 1 · gi (σi+1,n ) + hi (σi+1,n ),
fn (x1 , . . . , xn ) = xen ,
where
gi (x) = x2 + αi,1 · x + αi,2 ,
hi (x) = x2 + βi · x,
n
X
xj + fj (x1 , . . . , xn ).
σi+1,n =
j=i+1
4
1 ≤ i ≤ n − 1,Note that the GTDS F = {f1 , . . . , fn } must be considered as ordered tuple of
polynomials since in general the order of the fi ’s cannot be interchanged. Since
α2i,1 − 4 · αi,2 is a non-residue modulo p for all 1 ≤ i ≤ n − 1 the polynomials
gi do not have any zeros over Fp , therefore we can invert the GTDS with the
procedure described in [48, Proposition 8, Corollary 9]. In Table 1 we propose
suitable exponents for d2 which can be evaluated with at most 9 multiplications.
All exponents are chosen so that Arion and ArionHash provide at least 128
bit security against Gröbner basis attacks while minimizing the number of mul-
tiplicative constraints in a prover circuit, see Sections 3.2 and 4.1.
Table 1. Eﬃcient evaluation of exponents d2 ∈ {121, 123, 125, 129, 161, 193, 195, 257}.
d2
y = x2
121
123
125
2
, z = y2 · y
y = x2 · x, z =
y = x2

129y=
161y = x2

2
x2
2
2
 2
2
· x, z =
y = x2 · x, z =
195y = x2 · x, z =
y=


x2

y
· x, z =
193
257
Number of
Multiplications
Evaluation chain
 2
2 2
2
 2
2 2

y2
,z=



, x121 = z 2
y2
y2
y2
,z=
2

·z·x9
·z·y9
, x125 = z 2 · z · y9
 28
, x123 = z
 2
y2
2 2
2
2

2
y2

2 2
, x129 = z · x
2
·x9
2·x9
·y9
, x257 = z · x9
, x161 = z 2

 2 2
2 2
2
, x193 = z 2
, x195 = z 2
 2
2 2
2
Let us compute the degrees of the polynomials in the GTDS.
Lemma 2. Let n, d1 , e ≥ 1 be integers, and let FArion = {f1 , . . . , fn } be an Arion
GTDS. Then
deg (fi ) = 2n−i · (d1 + e) − d1 .
Proof. We perform an upwards induction, for n and n − 1 the claim is clear.
Suppose the claim is true for indices greater than or equal to i, i.e. deg (fi ) =
deg (fi ) = 2n−i · (d1 + e) − d1 . By construction, the leading monomial of fi−1 is
1
· fi2 . Thus,
the leading term of the polynomial xdi−1



1
· fi2 = d1 + 2 · 2n−i · (d1 + e) − d1
deg (fi−1 ) = deg xdi−1
= 2n−(i−1) · (d1 + e) − d1 ,
which proves the claim.
⊓
⊔
5To introduce mixing between the blocks we chose a circulant matrix whose
product with a vector can be efficiently evaluated.
Definition 3 (Affine layer of Arion). Let p ∈ Z be a prime and let Fp be the
field with p elements. The affine layer of Arion is defined as
Lc : Fnp → Fnp ,
x 7→ circ (1, . . . , n) x + c,
where circ (1, . . . , n) ∈ Fpn×n is the circulant matrix3 with entries 1, . . . , n in the
first row and c ∈ Fnp is a constant vector.
Remark 4. For any prime number p ∈ Z with p > 130 and n = 2, 3, 4 the
matrix circ (1, . . . , n) is a MDS matrix over Fp .
The following algorithm provides an efficient way to evaluate the matrix-
vector product for circ (1, . . . , n).
Algorithm 1 Efficient evaluation of matrix-vector product
Input
v = (v1 , . . . , vn )⊺ ∈ Fn
p
Output
circ (1, . . . , n) v ∈ Fn
p
1: Initialize w = P
(0, . . . , 0) ∈ Fn
p.
Pn
n
2: Compute σ = i=1 vi and set w1 = σ + i=1 (i − 1) · vi .
3: Set i = 2.
4: while i ≤ n do
5:
Set wi = wi−1 − σ + n · vi−1 .
6:
i = i + 1.
7: return w
To define a keyed permutation we need a key addition which we denote as
Kk : Fnp × Fnp → Fnp ,
(x, k) 7→ x + k.
The keyed permutation Arion is now defined as follows.
Definition 5 (Arion). Let p ∈ Z be a prime and let Fp be the field with p
(i)
elements, and let n > 1 and r ≥ 1 be integers. For 1 ≤ i ≤ r let FArion : Fnp → Fnp
3
We shortly recall the deﬁnition of (right) circulant matrices: Let k be a ﬁeld and let
v = (v1 , . . . , vn ) ∈ kn , then the circulant matrix of v is deﬁned as


v1 v2 . . . vn−1 vn
vn v1 . . . vn−2 vn−1 
.
circ(v) = 
..


.
v2 v3 . . . vn v1
6be an Arion GTDS and for 1 ≤ i ≤ r let Lci : Fnp → Fnp be affine layers from
Definition 3. The ith round function of Arion is defined as
(i)
Rk : Fnp × Fnp → Fnp ,
(i)
(x, k) 7→ Kk ◦ Lci ◦ FArion (x) .
Then Arion is defined as the following composition
Arion : Fnp × Fn×(r+1)
→ Fnp ,
p

(1)
(r)
x, k0 , k1 , . . . , kr 7→ Rkr ◦ · · · ◦ Rk1 ◦ L0 ◦ Kk0 (x) .
Further, we denote with Arion-π the unkeyed permutation, i.e. Arion is instanti-
ated with the key k0 = . . . = kr = 0.
Since our final aim is to construct a hash function using the Arion-π, we
analyze Arion only for keys kj = k, where k ∈ Fnp , in every round. We do not
give any key scheduling algorithm for keys whose sizes are larger than the block
size. Instantiation of Arion with such a key is not the main topic of this paper.
2.3
Hash Function
For the hash function ArionHash over Fnp we instantiate Arion-π in sponge mode
[9, 10]. The state size n = r + c is split into the rate part r and the capac-
ity part c. In [10, Theorem 2] it has been proven that for a random permuta-
tion the sponge construction is indistinguishable from a random oracle up to
min pr , pc/2 queries. Therefore, to provide κ bits of security, pr ≥ 2κ and
pc/2 ≥ 2κ , we require that
r≥
κ
,
log2 (p)
and
c≥
2·κ
.
log2 (p)
(1)
Given an input message m we choose a similar padding rule as for Poseidon [33,
§4.2], we add the smallest number of zeros < r such that the size of m || 0∗ is
a multiple of r. If we have to pad the message, then we replace the initial value
IV ∈ Fcp with |m| || IV′ ∈ Fcp , where |m| ∈ Fp is the size of the input message m
and IV′ ∈ Fc−1
is an initial value.
p
2.4
Instantiations
Target primes for Arion & ArionHash are the 250 bit prime numbers BLS12 and
BN2544 . Since degree growth of the Arion GTDS is dominated by the power
permutation in the bottom component we list the smallest integer m ∈ Z such
that m · e ≥ p in Table 2.
4
BLS12=0x73eda753299d7d483339d80809a1d80553bda402f f f e5bf ef f f f f f f f 00000001 ,
BN254=0x30644e72e131a029b85045b68181585d2833e84879b9709143e1f 593f 0000001 .
7Table 2. Smallest positive integer so that m ∈ Z such that m · e ≥ p for BLS12 and
BN254.
⌈p/e⌉
d2BLS12BN254
121
123
125
129
161
193
195
257n.a.
n.a.
3
n.a.
3
13
n.a.
33
n.a.
2
n.a.
4
3
n.a.
2
Therefore, by Lemma 2 for n = 3 and all exponents except d2 = 193 the
Arion GTDS surpasses degree p in the first component.
In Table 3 we provide the parameters for Arion and ArionHash and their ag-
gressive versions α-Arion and α-ArionHash with d1 , d2 ∈ Z such that gcd (di , p − 1) =
1, d1 = 3, 5 and 121 ≤ d2 ≤ 257. The number of rounds for Arion and ArionHash
are chosen to provide 128 bit security against the most efficient probabilistic
algorithm (available till date) for polynomial system solving in a Gröbner basis
250
attack on ArionHash. Since 2 2 = 2125 we consider all possible rate-capacity
pairs n = c + r suitable for ArionHash over BLS12 and BN254. As hash output
of ArionHash over BLS12 and BN254 we recommend to use a single Fp element.
Table 3. Arion, ArionHash, α-Arion and α-ArionHash parameters over primes p ≥ 2250
with d1 , d2 ∈ Z such that gcd (di , p − 1) = 1 and 121 ≤ d2 ≤ 257 and 128 bit security.
Rounds
Blocks
Arion & ArionHash
α-Arion & α-ArionHash
d1 = 3
3
4
5
6
8
6
6
5
5
4
5
4
4
4
4
d1 = 5
3
4
5
6
8
6
5
5
5
4
4
4
4
4
4
8The number of rounds for α-Arion and α-ArionHash are chosen to provide 128
bit security against the most efficient deterministic algorithm for polynomial
system solving in a Gröbner basis attack on ArionHash. For more details on the
security with respect to Gröbner basis attacks we refer to Appendix C in the
full version of the paper [49].
3Security Analysis of Arion
3.1Statistical Cryptanalysis
Differential Cryptanalysis. In differential cryptanalysis [13] and its variants
the propagation of input differences through the rounds of a block cipher or
hash function is exploited to recover the secret key or to construct a collision.
For the Arion GTDS the probability that an input difference ∆x ∈ Fnq \{0} prop-
agates to the output difference is ∆y ∈ Fnq is bounded by (see [48, Theorem 18,
Corollary 19])
P [FArion : ∆x → ∆y] ≤

d2
p
wt(∆x)
≤
d2
,
p
(2)
where wt : Fnq → Z denotes Hamming weight. For p ≥ 2250 and d2 ≤ 29 this
probability is bounded by 2−241·wt(∆x) ≤ 2−241 . Under the assumption that the
rounds of Arion are statistically independent we can estimate the probability
of any non-trivial differential trail via 2−241·r . Moreover, even if an adversary
can search a restricted differential hull of size 2120 between the 2nd and the rth
round, then two rounds are already sufficient to provide 128 bit security against
differential cryptanalysis. For more details we refer to Appendix A.1 in the full
version of the paper [49].
Note that a small differential probability also rules out the boomerang attack
[37, 52] which exploits two complementary differential patterns that span the
cipher, of which one must at least cover two rounds.
Truncated Differential & Rebound Cryptanalysis. In a truncated differ-
ential attack [38] an attacker can only predict parts of the difference between
pairs of text. We expect that the Arion GTDS admits truncated differentials of
Hamming weight 1 with probability 1 for the first round. On the other hand, if
wt(v) = 1, then we have that wt circ(1, . . . , n)v = n. Therefore, such a trun-
cated differential activates all inputs in the second round of Arion. Hence, for
p ≥ 2250 and d2 ≤ 29 the differential probability for the second round is bounded
by 2−250·n . Even if an adversary can search restricted differential hulls of size
2120 after the first round, this probability and Equation (2) nullify truncated
differential attacks within the 128 bit security target.
In a rebound attack [41, 45] an adversary connects two (truncated) differen-
tials in the middle of a cipher or hash function. Probability 1 truncated differ-
entials can cover at most one round of Arion, so r − 2 rounds can be covered
9with an inside-out approach. By our previous analysis we do not expect that
a successful rebound attack can be mounted on 4 or more rounds on Arion &
ArionHash within the 128 bit security target.
For more details we refer to Appendix A.2 in the full version of the paper [49].
Linear Cryptanalysis. Linear cryptanalysis [5, 44] utilizes affine approxima-
tions of round functions for a sample of known plaintexts. For any additive
character χ : Fq → C and any affine approximation a, b ∈ Fnq \ {0}, the linear
probability of the Arion GTDS is bounded by (see [48, Theorem 24, Corollary 25])
LPFArion (χ, a, b) ≤
(d2 − 1)2
.
q
(3)
Therefore, for p ≥ 2250 and d2 ≤ 29 this probability is bounded by 2−232 , and
henceforth under the assumption of statistically independent rounds of Arion the
linear probability of any non-trivial linear trail is bounded by 2−232·r . Moreover,
even if an adversary can search a restricted linear hull of size 2120 between the
2nd and the rth round, then two rounds are already sufficient to provide 128 bit
security against linear cryptanalysis. For more details we refer to Appendix A.3
in the full version of the paper [49].
3.2
Algebraic Cryptanalysis
Interpolation & Integral Cryptanalysis. Interpolation attacks [36] con-
struct the polynomial vector representing a cipher without knowledge of the
secret key. If such an attack is successful against a cipher, then an adversary
can encrypt any plaintext without knowledge of the secret key. Recall that any
function F : Fnq → Fq can be represented by a polynomial f ∈ Fp [Xn ] =
Fq [x1 , . . . , xn ]/ (xq1 − x1 , . . . , xqn − xn ), thus at most q n monomials can be present
in f . After the first round of Arion-π we expect that the terms
!e
n
n
X
X
i · xi
(4)
i · xi +
i=1
i=1
are present in every branch. After another application of the round function we
expect to produce the terms
n
X
i=1
i · xi
!e

- n
  X
  i=1
  i · xi
  !e
  mod (xp1 − x1 , . . . , xpn − xn )
  (5)
  in every branch. By our specification e is the inverse exponent of a relatively
  low degree permutation, therefore we expect that after two rounds almost all
  monomials from Fp [Xn ] are present in every component of Arion. For more
  details we refer to Appendix B.1 in the full version of the paper [49].
  10For a polynomial f ∈ Fq [x1 , . . . , xn ] an integral distinguisher [11,39] exploits
  that for any affine subspace V ⊂ Fnq with deg (f ) < dim (V ) · (q − 1) one has
  that
  X
  f (x) = 0.
  (6)
  x∈V
  If almost all monomials are present in Arion-π, then deg (Arion-π) ≈ n · (q − 1)
  in every component, so only V = Fnq is a suitable subspace for an integral
  distinguisher. Therefore, we do not expect that non-trivial integral distinguishers
  exist for Arion-π. For more details we refer to Appendix B.2 in the full version
  of the paper [49].
  Gröbner Basis Analysis. In a Gröbner basis attack [17, 21] one models a
  cipher or hash function as fully determined polynomial system and then solves
  for the key or preimage. For Gröbner basis analysis of Arion & ArionHash we as-
  sume that a degree reverse lexicographic (DRL) Gröbner basis can be found in
  O(1). We base the security of Arion & ArionHash solely on the complexity of solv-
  ing their polynomial systems via state of the art deterministic and probabilistic
  Gröbner basis conversion algorithms [27–29] combined with the univariate poly-
  nomial solving algorithm of Bariant et al. [8, §3.1]. With this methods, solving
  a fully determined polynomial system over a finite field Fq with known DRL
  Gröbner basis via deterministic methods requires
  
  
  
  2
  O n · dω + d · log (q) · log (d) · log log (d) + d · log (d) · log log (d)
  (7)
  field operations, and with deterministic methods
  √
  
  
  n−1
  2
  O
  n · d2+ n + d · log (q) · log (d) · log log (d) + d · log (d) · log log (d)
  (8)
  field operations, where q is the size of the finite field, n is the number of variables,
  d is the Fq -vector space dimension of the polynomial ring modulo the polynomial
  system and 2 ≤ ω < 2.3727 is a linear algebra constant. We conjecture that the
  quotient space dimension of Arion grows or bounded by
  
  r
  n−1
  dimFp (FArion ) (n, r, d1 , d2 ) = d2 · (d1 + 2)
  ,
  (9)
  and for ArionHash we conjecture that the dimension grows or is bounded by
  r
  
  (10)
  dimFp (FArionHash) (n, r, d1 , d2 ) = 2n−1 · d2 · (d1 + 1) − d1 · d2 .
  Round numbers for Arion & ArionHash in Table 3 are chosen to resist determin-
  istic as well as probabilistic Gröbner basis attacks against an ideal adversary
  with ω = 2 within the 128 bit security target. Round numbers for α-Arion &
  α-ArionHash in Table 3 are chosen to resist only deterministic Gröbner basis
  attacks within the 128 bit security target.
  11For ArionHash one can set up a collision polynomial system by connecting two
  preimage polynomial systems. Note that this polynomial system is in general not
  fully determined, therefore an adversary has to randomly guess some variables
  before solving the system. If an adversary guesses output variables of the sponge
  until the collision polynomial system is fully determined, then we conjecture
  that the quotient space dimension of the collision polynomial system grows or is
  bounded by
  2
  dimFp (FArionHash,coll ) (n, r, d1 , d2 ) = dimFp (FArionHash ) (n, r, d1 , d2 ) . (11)
  Thus, we do not expect a collision Gröbner basis attack to be more performative
  than a preimage attack. For more details we refer to Appendix C in the full
  version of the paper [49].
  4
  Performance Evaluation
  In this section, we compare various instances of ArionHash, Anemoi, Griffin
  and Poseidon with respect to R1CS (Section 4.2) and Plonk (Section 4.3).
  For starters, we discuss the theoretical foundation of an efficient implementa-
  tion of an ArionHash circuit. In the Anemoi proposal it was revealed that CCZ-
  equivalence is a route to construct high degree permutations that can be verified
  with CCZ-equivalent low degree functions [16, §4.1]. In Section 4.1 we follow this
  approach to prove that a ArionHash circuit can be transformed into an efficient
  circuit that avoids the computation of xe via an affine transformation.
  4.1
  Reducing the Number of Constraints
  By definition of the Arion GTDS a prover circuit will have to verify that
  y = xe ,
  (12)
  though since e induces the inverse power permutation to d2 ∈ {121, 123, 125, 129,
  161, 193, 195, 257} the naive circuit for Equation (12) will introduce many con-
  straints. On the other hand, from a prover’s perspective Equation (12) is equiv-
  alent to
  d
  y d2 = (xe ) 2 = x,
  (13)
  for all x ∈ Fp . Thus, in an implementation to reduce the number of multiplicative
  constraints we are well advised to implement the equivalent circuit instead of
  the naive circuit. We also would like to note that the same trick was applied in
  Griffin [31] to reduce the number of constraints.
  In the design of Anemoi [16, §4] a new tool was introduced to reduce the num-
  ber constraints for an Anemoi circuit: CCZ-equivalence [19]. The authors have
  found a high degree permutation, the open Flystel, which is CCZ-equivalent to
  a low degree function, the closed Flystel. Consequently, this can be exploited
  to significantly reduce the number of constraints in a prover circuit (cf. [16,
  Corollary 2]). Let us now formalize the trick in Equation (13) in terms of CCZ-
  equivalence.
  12Definition 6. Let Fq be a finite field, and let F, G : Fnq → Fm
  q be functions.
  (1) The graph of F is defined as
  ΓF =
  n
  o
  
  x, F (x) | x ∈ Fnq .
  (2) F and G are said to be CCZ-equivalent if there exists an affine permutation
  A of Fnq × Fm
  q such that
  ΓF = A(ΓG ).
  Now let us describe a GTDS that is equivalent to the Arion GTDS.
  Proposition 7. Let Fp be a prime field, and let n, d1 , d2 , e ∈ Z>1 be integers
  such that
  (i) d1 is the smallest positive integer such that gcd (d1 , p − 1) = 1,
  (ii) d2 is an arbitrary integer such that gcd (d2 , p − 1) = 1, and
  (iii) e · d2 ≡ 1 mod p − 1.
  Let FArion = {f1 , . . . , fn } be the Arion GTDS, let gi , hi ∈ Fq [x] be the polynomials
  that define FArion , and let the GTDS G = {fˆ1 , . . . , fˆn } be defined as
  fˆi (x1 , . . . , xn ) = xdi 1 · gi (τi+1,n ) + hi (τi+1,n ),
  fˆn (x1 , . . . , xn ) = xd2 ,
  1 ≤ i ≤ n − 1,
  n
  where
  τi+1,n =
  n
  X
  xj + fˆj (x1 , . . . , xn ).
  j=i+1
  Then FArion is CCZ-equivalent to G.
  th
  2n
  Proof. We consider the affine permutation A : F2n
  q → Fq that swaps the n
  th
  element with the (2n) element, moreover we consider the substitution x =
  x̂1 , . . . , x̂n−1 , x̂dn2 . Now we apply the affine permutation to ΓFArion which yields
  
  
  {x̂i }1≤i≤n−1
  
  
  2
  x̂e·d
  
   n
  n
  o
  
  .
  
  A x, FArion (x) = 
  d2
  
   fi x̂i , . . . , x̂n
  
  1≤i≤n−1
  d2
  x̂n
  By construction of d2 and e we have that xe·d2 = x for every x ∈ Fp . Let’s now
  investigate what happens to the fi ’s. Starting with fn−1 , we have that
  
  2
  = x̂n + x̂dn2 = τn,n (x̂n ),
  σn,n x̂dn2 = x̂dn2 + x̂e·d
  n
  for all x̂n ∈ Fp , and therefore
  
  fn−1 x̂n−1 , x̂dn2 = fˆn−1 (x̂n−1 , x̂n ).
  
  Inductively, we now go through all the branches to conclude that fi x̂i , . . . , x̂dn2 =
  
  
  fˆi (x̂i , . . . , x̂n ) for all 1 ≤ i ≤ n − 1 which proves that A x, F (x) = x̂, G(x̂) .
  ⊓
  ⊔
  13Corollary 8. Verifying that (y1 , . . . , yn ) = F (x1 , . . . , xn ) is equivalent to veri-
  fying that (y1 , . . . , yn−1 , xn ) = G(x1 , . . . , xn−1 , yn ).
  Note that it follows from [48, Theorem 18, 24] that the Arion GTDS F and its
  CCZ-equivalent GTDS G from Proposition 7 are in the same security class with
  respect to differential and linear cryptanalysis. Unlike as for Anemoi the CCZ-
  equivalent GTDS G is not a low degree function, though when implementing it
  as prover circuit we never have use multiplications to compute τi+1,n .
  4.2
  R1CS Performance of ArionHash
  Estimating the number of multiplicative constraints in a R1CS circuit for Arion-
  Hash is straightforward.
  Lemma 9. Let Fp be a finite field, let r, n ≥ 2 be integers, and let ArionHash with
  r rounds be defined over Fnp . For i = 1, 2 denote with di,inc the minimal number of
  multiplications to compute the univariate power permutation xdi . Then a prover
  R1CS circuit for ArionHash needs
  
  NArionHash = r · (n − 1) · (d1,inc + 2) + d2,inc
  multiplicative constraints.
  Proof. By Corollary 8 one needs d2,inc constraints in the nth branch. In each of
  the remaining n − 1 branches one needs d1,inc constraints for the power permu-
  tation, 1 constraints for the computation of gi and hi and 1 multiplication for
  the product of the power permutation and gi .
  ⊓
  ⊔
  Analog the number of R1CS constraints for Anemoi, Griffin and Poseidon
  (cf. [16, §7], [31, §7.2] and [33]) are given by
  NGriffin = 2 · r · (dinc + n − 2) ,
  r·n
  NAnemoi =
  · (dinc + 2) ,
  2
  NPoseidon = dinc · (n · rf + rp ) .
  In Table 4 we compiled the round numbers of the hash functions.
  14
  (14)
  (15)
  (16)Table 4. Round numbers for Anemoi [16, §A.4], ArionHash (Table 3), Griffin [31,
  Table 2] and Poseidon [33, Table 1] for 256 bit prime ﬁelds and 128 bit security with
  d2 ∈ {121, 123, 125, 161, 193, 195, 257}.
  Rounds
  ArionHash
  α-ArionHash
  Griffin Anemoi
  n
  3
  4
  5
  6
  8
  6
  6
  4
  5
  4
  5
  416
  1412
  4
  41110
  10
  n
  3
  4
  5
  6
  8
  Poseidon
  d1 = 3
  rf = 8, rp = 84
  rf = 8, rp = 84
  rf = 8, rp = 84
  rf = 8, rp = 84
  rf = 8, rp = 84
  d1 = 5
  6
  5
  5
  5
  4
  4
  4
  4
  4
  4
  12
  1112
  910
  10
  rf = 8, rp = 56
  rf = 8, rp = 56
  rf = 8, rp = 56
  rf = 8, rp = 56
  rf = 8, rp = 56
  In Table 5 we compare the theoretical number of constraints for R1CS of
  various hash functions.
  Table 5. R1CS constraint comparison 256 bit prime ﬁelds and 128 bit security with d2 ∈
  {121, 123, 125, 161, 193, 195, 257}. Round numbers for Anemoi, Griffin and Poseidon
  are taken from [16, §A.4], [31, Table 1] and [33, Table 1].
  R1CS Constraints
  ArionHash
  α-ArionHash
  n
  3
  4
  5
  6
  8
  d1 = 3
  102
  126
  120
  145
  148
  85
  84
  100
  116
  148
  n
  3
  4
  5
  6
  8
  Griffin Anemoi Poseidon
  96
  11296
  176120
  160
  216
  232
  248
  264
  296
  d1 = 5
  114
  120
  125
  170
  176
  76
  96
  116
  136
  176
  96
  110120
  162150
  200
  240
  264
  288
  312
  360
  Moreover, in Appendix D.1 of the full version of the paper [49] we compare the
  performance of Arion, Griffin and Poseidon using the C++ library libsnark
  [50] that is used in the privacy-protecting digital currency Zcash [35].
  154.3
  Plonk Performance of ArionHash
  Plonk [30] is a zkSNARK proof system which does not utilize R1CS constraints.
  In Plonk a 2-(input)-wire constraint is of the form, see [30, §6],
  (a · b) · qM + a · qL + b · qR + c · qO + qC = 0,
  (17)
  a and b denote the left and right input variable, c denotes the output variable
  and qM , qL , qR , qO and qC denote the “selector coefficient” of the multiplication,
  the variables and the constant term. The 3-(input)-wire Plonk constraint has 3
  addition gates
  (a · b) · qM + a · qL + b · qR + c · qO + d · qF + qC = 0,
  (18)
  where d is the “fourth” variable and qF its selector coefficient.
  Counting the number of Plonk constraints is more subtle, since we now have
  to account for additions too.
  Lemma 10. Let Fp be a finite field, let r, n ≥ 2 be integers, and let ArionHash
  with r rounds be defined over Fnp . For i = 1, 2 denote with di,inc the minimal
  number of multiplications to compute the univariate power permutation xdi .
  (1) A prover circuit needs
  (n − 1) · (d1,inc + 6) + d2,inc − 1
  2-wire and
  (n − 1) · (d1,inc + 4) + d2,inc
  3-wire Plonk constraints for the ArionHash GTDS.
  (2) A prover circuit needs
  4 · (n − 1)
  2-wire and
  
  
  n,
  
  n + 2 +
  n = 2, 3,
  
  
  
  
  n−3
  n−4
- ,
  2
  2
  n ≥ 4,
  3-wire Plonk constraints for the affine layer of ArionHash.
  Then a prover circuit needs
  
  NArionHash,2 = r · (n− 1)·(d1,inc + 6)+ d2,inc − 1 + (r + 1)·
  2-wire and
  (
  n · (n − 1),
  4 · (n − 1),
  
  NArionHash,3 = r · (n − 1) · (d1,inc + 4) + d2,inc
  
  
  n,
  
   
  
- (r + 1) ·
  n−3
  n−4
  
- ,
  n + 2 +
  2
  2
  3-wire Plonk constraints for ArionHash.
  16
  n = 2, 3
  n ≥ 4,
  n = 2, 3,
  n ≥ 4,Proof. For (1), again we can use the CCZ-equivalent GTDS G, see Proposition 7,
  to build the circuit for the ArionHash GTDS. For xdi one needs di,inc , so we need
  (n − 1) · d1,inc + d2,inc constraints for the univariate permutation polynomials.
  For τn,n one needs one constraint, and for τi+1,n , where i < n − 1, one needs
  two 2-wire constraints, so in total one needs 1 + 2 · (n − 2) 2-wire constraints to
  compute the τi+1,n ’s. On the other, hand for 3-wire constraints one can compute
  all τi+1,n ’s with one constraint, so n − 1 in total. To compute
  gi = τ 2 + αi,1 · τi+1,n + αi,2 ,
  hi = τ 2 + βi · τi+1,n
  one needs two constraints since one can build any quadratic polynomial with the
  2-wire Plonk constraint, see Equation (17). To compute xdi 1 · gi + hi one needs
  two 2-wire constraints and one 3-wire constraint. We have to do this (n − 1)
  times, hence in total we need (n − 1) · d1,inc + d2,inc + 1 + 2 · (n − 2) + 4 · (n − 1)
  2-wire and (n − 1) · d1,inc + d2,inc + (n − 1) + (n − 1) · (2 + 1) 3-wire constraints.
  a sum of m
  For (2), we build the circuit with Algorithm 1. To
  
   compute
  3-wire
  constraints.
  elements one needs m − 1 2-wireP
  constraints and 1 + m−3
  2
  n
  We have to do this
  for
  σ
  and
  for
  (i
  −

1. ·
   v
   ,
   so
   we
   need
   (n
   −
1.

- (n − 2) + 1
  i
   n−3 
   n−4i=2
  
  2-wire
  and
  1
-
- 1
-
- 1
  3-wire
  constraints
  to
  compute
  w1 =
  2
  2
  Pn
  σ + i=2 (i − 1) · vi + c1 , where we do constant addition in the addition of the
  two sums. For the ith component, we have that wi = wi−1 − σ + n · vi − ci−1 + ci ,
  so we need two 2-wire and one 3-wire constraints. We have to do this n − 1
  times,
  total
   in
   we need (n − 1) + (n − 2) + 1 + 2 · (n − 1) 2-wire and
   hence
  n−4
-
- n − 1 3-wire constraints.
  ⊓
  ⊔
  3 + n−3
  2
  2
  Note that for n ≥ 4 Algorithm 1 yields more efficient 2-wire and 3-wire
  circuits thangeneric
   matrix multiplication which always needs n · (n − 1) 2-wire
  3-wire constraints.
  and n · 1 + n−3
  2
  For ArionHash’s main competitors Anemoi, Griffin and Poseidon we list
  the formulae to compute their Plonk constraints in Table 6.
  Table 6. Plonk constraints for Anemoi [16, §7.2], Griffin [31, §7.3] and Poseidon [33].
  Hash
  2-wire constraints
  Anemoi
  r·n
  · (dinc + 5) + (r + 1) ·
  2
  Griffin
  
  2,
  
  
  
  n · n2 − 1 ,
  
  10,
  
  
  16,
  
  5,
  
  
  
  8,
  r · (2 · dinc + 4 · n − 11) + (r + 1) ·
  
  24,
   8·n
  4
  Poseidon
  3-wire constraints
  
  n = 2,
  n = 4,
  n = 6,
  n=8
- 2 · n − 4,
  dinc · (n · rf + rp ) + (r + 1) · n · (n − 1)
  r·n
  · (dinc + 3) + (r + 1) ·
  2
  
  (r + 1) · n, n = 2, 4,
  
  6,
  12,
  n = 6,
  n=8
  
  3,
  n = 3,
  
  
  
  6,
  n = 4,
  r · (2 · dinc + 3 · n − 8) + (r + 1) ·
  20,
  n = 8,
  
  j
  k
  
  
   6·n + 4 · n4 −1 + n,
  n ≥ 12
  4
  2
  dinc · (n · rf + rp ) + (r + 1) · n ·
  
  n,
   n−3 
  2
  ,
  n = 3,
  n = 4,
  n = 8,
  n ≥ 12
  n = 2, 3,
  n≥4
  In Table 7 we compare the theoretical number of constraints for Plonk for
  various hash functions.
  17Table 7. Plonk constraint comparison 256 bit prime ﬁelds and 128 bit security with
  d2 ∈ {121, 123, 125, 161, 193, 195, 257}. Round numbers are the same as in Table 4.
  Plonk Constraints
  2-wire constraints
  3-wire constraints
  State size n
  3
  4
  5
  6
  8
  Hash
  ArionHash
  α-ArionHash
  Poseidon
  Griffin
  Anemoi
  200
  168
  768
  165276 296 360 396
  188 240 292 396
  1336 2088 3024 5448
  246
  563
  220
  320 456
  212
  144
  624
  173247 316 385 424
  200 256 312 424
  1032 1568 2232 3944
  275
  561
  244
  350 496
  Hash
  ArionHash
  α-ArionHash
  Poseidon
  Griffin
  Anemoi
  3
  45
  6
  8
  211
  143
  600
  202
  172219
  177
  708261 279
  211 279
  1368 2504
  460
  216 332
  192
  155
  520
  182
  196239
  193
  608286 307
  231 307
  1080 1896
  398
  246 372
  d1 = 3
  147
  123
  492
  131
  d1 = 5
  159
  107
  432
  123
  Moreover, in Appendix D.2 of the full version of the paper [49] we compare
  the performance of Arion and Poseidon using the Rust library Dusk Network
  Plonk [23].
  Acknowledgments. Matthias Steiner and Stefano Trevisani were supported
  by the KWF under project number KWF-3520|31870|45842.
  References

1. Albrecht, M.R., Grassi, L., Perrin, L., Ramacher, S., Rechberger, C., Rotaru, D.,
   Roy, A., Schofnegger, M.: Feistel structures for MPC, and more. In: Sako, K.,
   Schneider, S., Ryan, P.Y.A. (eds.) ESORICS 2019: 24th European Symposium
   on Research in Computer Security, Part II. Lecture Notes in Computer Science,
   vol. 11736, pp. 151–171. Springer, Heidelberg, Germany, Luxembourg (Sep 23–27,
   2019). https://doi.org/10.1007/978-3-030-29962-0_8
2. Albrecht, M.R., Grassi, L., Rechberger, C., Roy, A., Tiessen, T.: MiMC: Ef-
   ﬁcient encryption and cryptographic hashing with minimal multiplicative com-
   plexity. In: Cheon, J.H., Takagi, T. (eds.) Advances in Cryptology – ASI-
   ACRYPT 2016, Part I. Lecture Notes in Computer Science, vol. 10031, pp.
   191–219. Springer, Heidelberg, Germany, Hanoi, Vietnam (Dec 4–8, 2016).
   https://doi.org/10.1007/978-3-662-53887-6_7
3. Albrecht, M.R., Rechberger, C., Schneider, T., Tiessen, T., Zohner, M.: Ciphers
   for MPC and FHE. In: Oswald, E., Fischlin, M. (eds.) Advances in Cryptology
   – EUROCRYPT 2015, Part I. Lecture Notes in Computer Science, vol. 9056,
   pp. 430–454. Springer, Heidelberg, Germany, Soﬁa, Bulgaria (Apr 26–30, 2015).
   https://doi.org/10.1007/978-3-662-46800-5_17
4. Aly, A., Ashur, T., Ben-Sasson, E., Dhooghe, S., Szepieniec, A.: De-
   sign of symmetric-key primitives for advanced cryptographic protocols.
   IACR Transactions on Symmetric Cryptology 2020(3), 1–45 (2020).
   https://doi.org/10.13154/tosc.v2020.i3.1-45
5. Baignères, T., Stern, J., Vaudenay, S.: Linear cryptanalysis of non binary ciphers.
   In: Adams, C.M., Miri, A., Wiener, M.J. (eds.) SAC 2007: 14th Annual Interna-
   tional Workshop on Selected Areas in Cryptography. Lecture Notes in Computer
   Science, vol. 4876, pp. 184–211. Springer, Heidelberg, Germany, Ottawa, Canada
   (Aug 16–17, 2007). https://doi.org/10.1007/978-3-540-77360-3_13
6. Bardet, M., Faugère, J.C., Salvy, B.: Asymptotic behaviour of the in-
   dex of regularity of semi-regular quadratic polynomial systems. In: MEGA
   2005 - 8th International Symposium on Eﬀective Methods in Algebraic
   Geometry. pp. 1–17. Porto Conte, Alghero, Sardinia, Italy (05 2005),
   https://hal.archives-ouvertes.fr/hal-01486845
7. Bardet, M., Faugère, J.C., Salvy, B., Yang, B.Y.: Asymptotic behaviour of the
   degree of regularity of semi-regular polynomial systems. In: MEGA 2005 - 8th In-
   ternational Symposium on Eﬀective Methods in Algebraic Geometry. Porto Conte,
   Alghero, Sardinia, Italy (05 2005)
8. Bariant, A., Bouvier, C., Leurent, G., Perrin, L.: Algebraic attacks against some
   arithmetization-oriented primitives. IACR Transactions on Symmetric Cryptology
   2022(3), 73–101 (2022). https://doi.org/10.46586/tosc.v2022.i3.73-101
9. Bertoni, G., Daemen, J., Peeters, M., Van Assche, G.: Sponge functions. Ecrypt
   Hash Workshop (2007), https://keccak.team/ﬁles/SpongeFunctions.pdf
10. Bertoni, G., Daemen, J., Peeters, M., Van Assche, G.: On the indiﬀerentia-
    bility of the sponge construction. In: Smart, N.P. (ed.) Advances in Cryptol-
    ogy – EUROCRYPT 2008. Lecture Notes in Computer Science, vol. 4965, pp.
    181–197. Springer, Heidelberg, Germany, Istanbul, Turkey (Apr 13–17, 2008).
    https://doi.org/10.1007/978-3-540-78967-3_11
11. Beyne, T., Canteaut, A., Dinur, I., Eichlseder, M., Leander, G., Leurent, G.,
    Naya-Plasencia, M., Perrin, L., Sasaki, Y., Todo, Y., Wiemer, F.: Out of oddity -
    new cryptanalytic techniques against symmetric primitives optimized for integrity
    proof systems. In: Micciancio, D., Ristenpart, T. (eds.) Advances in Cryptology –
    CRYPTO 2020, Part III. Lecture Notes in Computer Science, vol. 12172, pp. 299–
12. Springer, Heidelberg, Germany, Santa Barbara, CA, USA (Aug 17–21, 2020).
    https://doi.org/10.1007/978-3-030-56877-1_11
13. Biham, E., Biryukov, A., Shamir, A.: Cryptanalysis of Skipjack reduced to 31
    rounds using impossible diﬀerentials. In: Stern, J. (ed.) Advances in Cryptol-
    ogy – EUROCRYPT’99. Lecture Notes in Computer Science, vol. 1592, pp. 12–
14. Springer, Heidelberg, Germany, Prague, Czech Republic (May 2–6, 1999).
    https://doi.org/10.1007/3-540-48910-X_2
15. Biham, E., Shamir, A.: Diﬀerential cryptanalysis of DES-like cryptosystems. Jour-
    nal of Cryptology 4(1), 3–72 (Jan 1991). https://doi.org/10.1007/BF00630563
16. Bogdanov, A., Rijmen, V.: Linear hulls with correlation zero and linear crypt-
    analysis of block ciphers. Des. Codes Cryptogr. 70(3), 369–383 (03 2014).
    https://doi.org/10.1007/s10623-012-9697-z
17. Bogdanov, A., Wang, M.: Zero correlation linear cryptanalysis with re-
    duced data complexity. In: Canteaut, A. (ed.) Fast Software Encryption
    – FSE 2012. Lecture Notes in Computer Science, vol. 7549, pp. 29–48.
    Springer, Heidelberg, Germany, Washington, DC, USA (Mar 19–21, 2012).
    https://doi.org/10.1007/978-3-642-34047-5_3
18. Bouvier, C., Briaud, P., Chaidos, P., Perrin, L., Salen, R., Velichkov, V., Willems,
    D.: New design techniques for eﬃcient arithmetization-oriented hash functions:
    Anemoi permutations and Jive compression mode. Cryptology ePrint Archive, Pa-
    per 2022/840 (2022), https://eprint.iacr.org/2022/840, Version: 20221017:092940
19. Buchberger, B.: Ein Algorithmus zum Auﬃnden der Basiselemente des Restklassen-
    ringes nach einem nulldimensionalen Polynomideal. Phd thesis, Universität Ins-
    bruck (1965)
20. Caminata, A., Gorla, E.: Solving multivariate polynomial systems and an in-
    variant from commutative algebra. In: Bajard, J.C., Topuzoğlu, A. (eds.) Arith-
    metic of Finite Fields. pp. 3–36. Springer International Publishing, Cham (2021).
    https://doi.org/10.1007/978-3-030-68869-1_1
21. Carlet, C., Charpin, P., Zinoviev, V.: Codes, bent functions and permutations suit-
    able for DES-like cryptosystems. Des. Codes Cryptogr. 15(2), 125–156 (11 1998).
    https://doi.org/10.1023/A:1008344232130
22. Cover, T.M., Joy, T.A.: Elements of Information Theory. John Wiley & Sons, Ltd,
    Hoboken, New Jersey, 2 edn. (2006). https://doi.org/10.1002/0471200611
23. Cox, D.A., Little, J., O’Shea, D.: Ideals, Varieties, and Algorithms: An Introduc-
    tion to Computational Algebraic Geometry and Commutative Algebra. Under-
    graduate Texts in Mathematics, Springer International Publishing, 4 edn. (2015).
    https://doi.org/10.1007/978-3-319-16721-3
24. Ding, J., Schmidt, D.: Solving degree and degree of regularity for polynomial sys-
    tems over a ﬁnite ﬁelds. In: Fischlin, M., Katzenbeisser, S. (eds.) Number Theory
    and Cryptography: Papers in Honor of Johannes Buchmann on the Occasion of
    His 60th Birthday. pp. 34–49. Springer Berlin Heidelberg, Berlin, Heidelberg (2013).
    https://doi.org/10.1007/978-3-642-42001-6_4
25. Dusk Network: PLONK (2022), https://github.com/dusk-network/plonk, Version:
    0.13.1
26. Dusk Network: Dusk-Poseidon (2023), https://github.com/dusk-network/Poseidon252,
    Version: 0.28.1
27. Faugère,
    J.C.:
    A
    new eﬃcient
    algorithm
    for
    computing
    Gröb-
    ner bases (F4). J. Pure Appl. Algebra 139(1), 61–88 (1999).
    https://doi.org/10.1016/S0022-4049(99)00005-5
28. Faugère, J.C.: A new eﬃcient algorithm for computing Gröbner bases without
    reduction to zero (F5). In: Proceedings of the 2002 International Symposium on
    Symbolic and Algebraic Computation. p. 75–83. ISSAC ’02, Association for Com-
    puting Machinery (2002). https://doi.org/10.1145/780506.780516
29. Faugère, J.C., Gaudry, P., Huot, L., Renault, G.: Sub-cubic change of ordering
    for Gröbner basis: A probabilistic approach. In: Proceedings of the 39th Inter-
    national Symposium on Symbolic and Algebraic Computation. p. 170–177. IS-
    SAC ’14, Association for Computing Machinery, New York, NY, USA (2014).
    https://doi.org/10.1145/2608628.2608669
30. Faugère, J.C., Gianni, P., Lazard, D., Mora, T.: Eﬃcient computation of zero-
    dimensional Gröbner bases by change of ordering. J. Symb. Comput. 16(4), 329–
    344 (1993). https://doi.org/10.1006/jsco.1993.1051
31. Faugère, J.C., Mou, C.: Sparse FGLM algorithms. J. Symb. Comput. 80, 538–569
    (2017). https://doi.org/https://doi.org/10.1016/j.jsc.2016.07.025
32. Gabizon, A., Williamson, Z.J., Ciobotaru, O.: PLONK: Permutations over
    lagrange-bases for oecumenical noninteractive arguments of knowledge. Cryptol-
    ogy ePrint Archive, Report 2019/953 (2019), https://eprint.iacr.org/2019/953
33. Grassi, L., Hao, Y., Rechberger, C., Schofnegger, M., Walch, R., Wang, Q.:
    Horst meets Fluid-SPN: Griffin for zero-knowledge applications. Cryptology
    ePrint Archive, Paper 2022/403 (2022), https://eprint.iacr.org/2022/403, Version:
    20230214:131048
34. Grassi, L., Khovratovich, D., Lüftenegger, R., Rechberger, C., Schofnegger, M.,
    Walch, R.: Reinforced Concrete: A fast hash function for veriﬁable computation.
    In: Proceedings of the 2022 ACM SIGSAC Conference on Computer and Commu-
    nications Security. p. 1323–1335. CCS ’22, Association for Computing Machinery,
    New York, NY, USA (2022). https://doi.org/10.1145/3548606.3560686
35. Grassi, L., Khovratovich, D., Rechberger, C., Roy, A., Schofnegger, M.: Poseidon:
    A new hash function for zero-knowledge proof systems. In: Bailey, M., Greenstadt,
    R. (eds.) USENIX Security 2021: 30th USENIX Security Symposium. pp. 519–535.
    USENIX Association (Aug 11–13, 2021)
36. Grassi, L., Lüftenegger, R., Rechberger, C., Rotaru, D., Schofnegger, M.: On a
    generalization of substitution-permutation networks: The HADES design strat-
    egy. In: Canteaut, A., Ishai, Y. (eds.) Advances in Cryptology – EURO-
    CRYPT 2020, Part II. Lecture Notes in Computer Science, vol. 12106, pp.
    674–704. Springer, Heidelberg, Germany, Zagreb, Croatia (May 10–14, 2020).
    https://doi.org/10.1007/978-3-030-45724-2_23
37. Hopwood, D., Bowe, S., Hornby, T., Wilcox, N.: Zcash protocol speciﬁcation
    (2022), https://github.com/zcash/zips/blob/main/protocol/protocol.pdf , Version:
    2022.3.8
38. Jakobsen, T., Knudsen, L.R.: The interpolation attack on block ciphers. In: Biham,
    E. (ed.) Fast Software Encryption – FSE’97. Lecture Notes in Computer Science,
    vol. 1267, pp. 28–40. Springer, Heidelberg, Germany, Haifa, Israel (Jan 20–22, 1997).
    https://doi.org/10.1007/BFb0052332
39. Joux, A., Peyrin, T.: Hash functions and the (ampliﬁed) boomerang attack. In:
    Menezes, A. (ed.) Advances in Cryptology – CRYPTO 2007. Lecture Notes in Com-
    puter Science, vol. 4622, pp. 244–263. Springer, Heidelberg, Germany, Santa Bar-
    bara, CA, USA (Aug 19–23, 2007). https://doi.org/10.1007/978-3-540-74143-5_14
40. Knudsen, L.R.: Truncated and higher order diﬀerentials. In: Preneel, B. (ed.) Fast
    Software Encryption – FSE’94. Lecture Notes in Computer Science, vol. 1008,
    pp. 196–211. Springer, Heidelberg, Germany, Leuven, Belgium (Dec 14–16, 1995).
    https://doi.org/10.1007/3-540-60590-8_16
41. Knudsen, L.R., Wagner, D.: Integral cryptanalysis. In: Daemen, J., Rijmen, V.
    (eds.) Fast Software Encryption – FSE 2002. Lecture Notes in Computer Science,
    vol. 2365, pp. 112–127. Springer, Heidelberg, Germany, Leuven, Belgium (Feb 4–6,
    2002). https://doi.org/10.1007/3-540-45661-9_9
42. Lai, X.: Higher order derivatives and diﬀerential cryptanalysis. In: Blahut, R.E.,
    Costello, D.J., Maurer, U., Mittelholzer, T. (eds.) Communications and Cryptog-
    raphy: Two Sides of One Tapestry. pp. 227–233. Springer US, Boston, MA (1994).
    https://doi.org/10.1007/978-1-4615-2694-0_23
43. Lamberger, M., Mendel, F., Rechberger, C., Rijmen, V., Schläﬀer, M.: Rebound
    distinguishers: Results on the full Whirlpool compression function. In: Matsui, M.
    (ed.) Advances in Cryptology – ASIACRYPT 2009. Lecture Notes in Computer Sci-
    ence, vol. 5912, pp. 126–143. Springer, Heidelberg, Germany, Tokyo, Japan (Dec 6–
    10, 2009). https://doi.org/10.1007/978-3-642-10366-7_8
44. Lazard, D.: Gröbner bases, Gaussian elimination and resolution of systems of al-
    gebraic equations. In: van Hulzen, J.A. (ed.) Computer Algebra. Lecture Notes
    in Computer Science, vol. 162, pp. 146–156. Springer Berlin Heidelberg (1983).
    https://doi.org/10.1007/3-540-12868-9_99
45. Lidl, R., Niederreiter, H.: Finite ﬁelds. Encyclopedia of mathematics and its appli-
    cations, Cambridge Univ. Press, Cambridge, 2 edn. (1997)
46. Matsui, M.: Linear cryptanalysis method for DES cipher. In: Helleseth, T. (ed.)
    Advances in Cryptology – EUROCRYPT’93. Lecture Notes in Computer Science,
    vol. 765, pp. 386–397. Springer, Heidelberg, Germany, Lofthus, Norway (May 23–
    27, 1994). https://doi.org/10.1007/3-540-48285-7_33
47. Mendel, F., Rechberger, C., Schläﬀer, M., Thomsen, S.S.: The rebound attack:
    Cryptanalysis of reduced Whirlpool and Grøstl. In: Dunkelman, O. (ed.) Fast
    Software Encryption – FSE 2009. Lecture Notes in Computer Science, vol. 5665,
    pp. 260–276. Springer, Heidelberg, Germany, Leuven, Belgium (Feb 22–25, 2009).
    https://doi.org/10.1007/978-3-642-03317-9_16
48. Nyberg, K.: Diﬀerentially uniform mappings for cryptography. In: Helleseth, T.
    (ed.) Advances in Cryptology – EUROCRYPT’93. Lecture Notes in Computer
    Science, vol. 765, pp. 55–64. Springer, Heidelberg, Germany, Lofthus, Norway
    (May 23–27, 1994). https://doi.org/10.1007/3-540-48285-7_6
49. OSCAR – open source computer algebra research system, Version: 0.12.0 (2023),
    https://oscar.computeralgebra.de
50. Roy, A., Steiner, M.: Generalized triangular dynamical system: An algebraic system
    for constructing cryptographic permutations over ﬁnite ﬁelds. arXiv: 2204.01802
    (2022). https://doi.org/10.48550/ARXIV.2204.01802, Version: 6
51. Roy, A., Steiner, M.J., Trevisani, S.: Arion: Arithmetization-oriented permutation
    and hashing from generalized triangular dynamical systems. arXiv: 2303.04639
    (2023). https://doi.org/10.48550/ARXIV.2303.04639, Version: 3
52. SCIPR Lab: libsnark: a c++ library for zkSNARK proofs (2017),
    https://github.com/scipr-lab/libsnark
53. The Sage Developers: SageMath, the Sage Mathematics Software System (Version
    9.3) (2022), https://www.sagemath.org
54. Wagner, D.: The boomerang attack. In: Knudsen, L.R. (ed.) Fast Soft-
    ware Encryption – FSE’99. Lecture Notes in Computer Science, vol. 1636,
    pp. 156–170. Springer, Heidelberg, Germany, Rome, Italy (Mar 24–26, 1999).
    https://doi.org/10.1007/3-540-48519-8_12
    AStatistical Attacks on Arion
    A.1Differential Cryptanalysis
    Differential cryptanalysis [13] and its variants are the most widely applied attack
    vectors against symmetric-key ciphers and hash functions. It is based on the
    propagation of input differences through the rounds of a block cipher. In its base
    form an attacker requests the ciphertexts for large numbers of chosen plaintexts.
    Then he assumes that for r − 1 rounds the input difference is ∆x ∈ Fnq \ {0} and
    the output difference is ∆y ∈ Fnq . Under the assumption that the differences in
    the last round are fixed the attacker can then deduce the possible keys. The key
    quantity to estimate the effectiveness of differential cryptanalysis is the so-called
    differential uniformity.
    Definition 11 (see [46]). Let Fq be a finite field, and let f : Fnq → Fm
    q be a
    function.
    22(1) The differential distribution table of f at a ∈ Fnq and b ∈ Fm
    q is defined as
    δf (a, b) = {x ∈ Fnq | f (x + a) − f (x) = f (b)} .
    (2) The differential uniformity of f is defined as
    δ(f ) =
    max
    a∈Fn
    q \{0},
    b∈Fm
    q
    δf (a, b).
    Given the differential uniformity of a function one can upper bound the
    success probability of differential cryptanalysis with input differences ∆x ∈
    Fnq \ {0} and ∆y ∈ Fm
    q by
    P [f : ∆x → ∆y] ≤
    δ(f )
    .
    qn
    (19)
    Naturally, the lower the differential uniformity the stronger is the resistance of
    a block cipher against differential cryptanalysis.
    By the choice of our parameters of Arion, see Definition 1 and thereafter, we
    have that d2 ≥ d1 . With [48, Theorem 18, Corollary 19] the maximal success
    probability of any differential for the Arion GTDS is bounded by
     wt(∆x)
    d2
    d2
    ≤ ,
    (20)
    P [FArion : ∆x → ∆y] ≤
    p
    p
    where ∆x ∈ Fnp \ {0}, ∆y ∈ Fnp , and 1 ≤ wt (∆x) ≤ n denotes the Hamming
    (i)
    (i)
    weight, i.e. the number of non-zero entries, of ∆x. Since Rk and FArion are affine
    (i)
    equivalent Equation (20) also applies to Rk . For Arion we use a differential in
    every round. If we assume that the differentials are independent among the
    rounds of Arion, then
    i  d wt(∆x1 )+...+wt(∆xr )
    h
    2
    (r)
    (1)
    P Rk ◦ · · · ◦ Rk : ∆x1 → . . . → ∆xr+1 ≤
    . (21)
    p
    Therefore, we can estimate the security level κ of Arion against any differen-
    tial trail via
     r
    
    
    d2
    ≤ 2−κ =⇒ κ ≤ r · log2 (p) − log2 (d2 ) .
    (22)
    p
    In Table 8 we list the security level of Arion against a differential characteristic
    for different prime sizes.
    Table 8. Security level of Arion against any diﬀerential characteristic for p ≥ 2N and
    d2 ≤ 29 .
    rNκ (bits)
    3
    2
    160
    120
    250153
    222
    241
    23With Equation (21) we can also estimate the probability of the differential
    hull of Arion.
    Theorem 12 (Differential hull of Arion). Let p ∈ Z be a prime and let Fp be
    (1)
    (r)
    the field with p elements, and let n > 1 and r ≥ 1 be integers. Let Rk , . . . , Rk :
    n
    n
    n
    n
    Fp × Fp → Fp be Arion round functions, and let ∆x1 , ∆xr+1 ∈ Fp be such
    that ∆x1 6= 0. Assume that the differentials among the rounds of Arion are
    independent. Then
    i  d wt(∆x1 ) 
    h
    r−1
    2
    (r)
    (1)
    n
    P Rk ◦ · · · ◦ Rk : ∆x1 → ∆xr+1 ≤
    · (d2 + 1) − 1
    .
    p
    Proof. Let us first do some elementary rearrangements and estimations
    i
    h
    (r)
    (1)
    P Rk ◦ · · · ◦ Rk : ∆x1 → ∆xr

#

" r
o
\ n (i)
X
Rk : ∆xi → ∆xi+1
P
=
i=1
∆x2 ,...,∆xr ∈Fn
p \{0}
(1)
=
X
∆x2 ,...,∆xr ∈Fn
p \{0}
(2)
=
X
∆x2 ,...,∆xr ∈Fn
p \{0}
 wt(∆x1 )
(3)
d2
·
≤
p

#

" r
o
\ n (i)
−1
FArion : ∆xi → circ(1, . . . , n) ∆xi+1
P
i=1

#

" r
o
\ n (i)
FArion : ∆xi → ∆xi+1
P
i=1
X
∆x2 ,...,∆xr ∈Fn
p \{0}

d2
p
wt(∆x2 )+...+wt(∆xr )
.
In (1) we expanded the definition of the round function and inverted the circulant
matrix, in (2) we exploited that we sum over all possible differentials and that
the affine Arion layer is invertible, moreover we implicitly substituted ∆xr+1 =
circ(1, . . . , n)−1 ∆xr+1 . Finally, in (3) we applied Equation (21). For ease of
writing we compute the sum only for one difference variable, then
X
∆x∈Fn
p \{0}
≤
n 
X
i=1

d2
p
wt(∆x)
n
(4) X
=
i=1
(p − 1)i ·
   i
n
d2
·
i
p

n
n
· di2 = (d2 + 1) − 1,
i
where in (4) we used that there are (p − 1)i ·
wt(∆x) = i. This proves the claim.

n
i
many vectors ∆x ∈ Fnp with
⊓
⊔
Remark 13. This estimation can be performed over any finite field Fq , any
invertible affine layer, and
 any primitive whose differential uniformity at round
level is in O p− wt(∆x) .
24With the theorem we can now estimate the security level of Arion with respect
to differential cryptanalysis and the full differential hull via
 wt(∆x1 ) 
r−1
d2
n
· (d2 + 1) − 1
≤ 2−κ
(23)
p




n
=⇒ κ ≤ wt(∆x1 ) · log2 (p) − log2 (d2 ) − (r − 1) · log2 (d2 + 1) − 1 .
(24)
In Table 9 report the security level for Arion against differential cryptanalysis
utilizing the full differential hull with the parameters from Table 3 and different
field sizes. Since our probability estimation from Theorem 12 could in principle
be > 1 for some parameters combinations, we always report the security level
with respect to the smallest wt(∆x1 ) such that the probability estimate is < 1.
Table 9. Security level of Arion against any diﬀerential characteristic for p ≥ 2N and
d2 ≤ 257 with the full diﬀerential hull.
N = 60
n
r
wt(∆x1 )
κ (bits)
3
6
3
35
4
5
3
27
4
6
4
47
5
5
4
47
N = 120
6
5
4
15
8
3
4
4
6
5
4
2
2
15 103 95
4
6
2
63
5
5
2
63
N = 250
6
5
1
31
8
3
4
4
4
6
5
6
2
1
1
1
31 121 113 81
5
5
1
81
6
5
1
49
8
4
1
49
Although, our target security of 128 bit is never met we should keep in mind
that the full differential hull is of size (p − 1)n·(r−1) . E.g., for our target primes
BLS12 and BN254 the smallest differential hull is of size ≈ 2250·15 , and for 60 bit
prime fields the smallest differential hull still would be of size 260·15 . Therefore,
we do not expect that differential cryptanalysis can break Arion & ArionHash
within the 128 bit security target.
Nevertheless, to convince skeptical readers we provide another hull estimation
for a computationally limited adversary, i.e. an adversary who can only search
a restricted number of differences within the differential hull.
Lemma 14 (Restricted differential hull of Arion). Let p ∈ Z be a prime
and let Fp be the field with p elements, and let n > 1 and r ≥ 1 be integers. Let
(1)
(r)
Rk , . . . , Rk : Fnp × Fnp → Fnp be Arion round functions, and let ∆x1 , ∆xr+1 ∈
n
Fp be such that ∆x1 6= 0. Assume that the differentials among the rounds of
Arion are independent and that in every round only up to M < (p − 1) · n many
differences ∆xi ∈ Fnp \ {0} can be utilized. Then
h
i  d wt(∆x1 )  M · d r−1
2
2
(r)
(1)
·
.
P Rk ◦ · · · ◦ Rk : ∆x1 → ∆xr+1 ≤
p
p
Proof. Recall that the bound for the differential uniformity of the Arion GTDS is
maximal if there is only one non-zero entry in ∆x ∈ Fnp \ {0}, see Equation (20).
25In every intermediate round there are (p − 1) · n many differences with one non-
zero entry, hence to maximize the estimation we first sum over those elements.
With Equation (21) we then obtain the claimed inequality.
⊓
⊔
Provided that M < dp2 , then the inequality is always less than 1. For the
restricted differential hull we can then estimate the security level of Arion via
 r
d2
M r−1 ·
≤ 2−κ
(25)
p


=⇒ κ ≤ r · log2 (p) − log2 (d2 ) − (r − 1) · log2 (M ) .
(26)
In Table 10 we report round number for various field sizes that achieve at least
√
128 bit security for a restricted hull of size p at round level.
Table 10. Security level of Arion against diﬀerential cryptanalysis for p ≥ 2N and
d2 ≤ 257 with a restricted diﬀerential hull of size 2N/2 at round level.
A.2
rNκ (bits)
5
4
360
120
250139
267
475
Truncated Differential and Rebound Attacks
In a truncated differential attack [38] an attacker can only predict parts of the dif-
ference between pairs of text. We expect that the Arion GTDS admits truncated
⊺ F
⊺
differentials with probability 1. For example (0, . . . , 0, α) −−Arion
−→ circ (1, . . . , n) (0, . . . , 0, β)
×
n
where α, β ∈ Fp . Note that for any vector v ∈ Fp with wt(v) one always has

that wt circ (1, . . . , n) v = n. Therefore, for any truncated differential of the
Arion GTDS where only one element is active in the input and output one has
that all components are active after application of the affine layer of Arion. For
truncated differentials we then estimate the security level of Arion via
 n
d2
≤ 2−κ =⇒ κ ≤ n · (log2 (p) − log2 (d2 )) − log2 (M ) ,
(27)
M·
p
where M denotes the size of the restricted differential hull available to a compu-
tationally limited adversary. In Table 11 we report the security level for various
parameters. For primes of size p ≥ 2120 one full round is already sufficient to
achieve 128 bit security for a restricted differential hull of size M ≥ 2120 . On the
other hand for p ≥ 260 the security target is not met. Therefore, if one would
like to instantiate Arion & ArionHash over such a prime one also has to consider
later rounds.
26Table 11. Security level of Arion against truncated diﬀerential cryptanalysis with
weight 1 truncated diﬀerentials in the ﬁrst round and for p ≥ 2N and d2 ≤ 257 with a
restricted diﬀerential hull of size 2M .
NnMκ (bits)
60
120
2504
3
3100
120
250107
215
475
Suppose an adversary can cover two rounds with a truncated differential of
probability 1, his best bet is to search the remaining rounds for input/output
differentials of weight 1. For such an adversary we can then estimate the security
level as
r−2

d2
≤ 2−κ =⇒ κ ≤ (r − 2) · (log2 (p) − log2 (d2 ) − log2 (M )) , (28)
M·
p
where M denotes the size of the restricted differential hull available to a compu-
tationally limited adversary at round level. In Table 12 we report the security
level for various parameters. For primes of size p ≥ 2250 four rounds of Arion
are sufficient to protect against two round truncated differentials of probability 1
within the 128 bit security target. On the other hand, for primes of size p ≥ 260
one needs at least 6 rounds to meet the 128 bit security target.
Table 12. Security level of Arion against two round truncated diﬀerentials for p ≥ 2N
and d2 ≤ 257 with a restricted diﬀerential hull of size 2N/2 at round level.
Nr−2κ (bits)
60
120
2504
3
287
155
233
In a rebound attack [41,45] one has to find two input/output pairs such that
the inputs satisfy a certain (truncated) input difference and the outputs satisfy
a certain (truncated) output difference. Such an attack can be split into two
phases: an inbound and an outbound phase. Let PArion : Fnq → Fnq be the target
permutation, then we split it into three sub-parts PArion = Pf w ◦ Pin ◦ Pout .
The inbound phase is placed in the middle preceded and followed by the two
outbound phases. Then, in the outbound phase two high-probability (truncated)
differential trails are constructed which are connected with the inbound phase.
Since a truncated differential with probability 1 can only cover a single round an
attacker can cover only r − 2 rounds with an inside-out approach. By Table 12,
for our target primes BLS12 and BN254 two inbound rounds are sufficient to
27achieve the 128 bit security target. So in total 4 rounds are sufficient to nullify
this attack vector.
A.3
Linear Cryptanalysis
In linear cryptanalysis [5,44] one tries to discover affine approximations of round
functions for a sample of known plaintexts. The key quantity to estimate the
effectiveness
Pn of linear cryptanalysis is the so-called correlation. We denote with
ha, bi = i=1 ai · bi the scalar product over Fnq .
Definition 15 (see [5, Definition 6, 15]). Let Fq be a finite field, let n ≥ 1,
let χ : Fq → C be a non-trivial additive character, let F : Fnq → Fnq be a function,
and let a, b ∈ Fnq .
(1) The correlation for the character χ of the linear approximation (a, b) of F
is defined as
CORRF (χ, a, b) =

1 X 
·
χ
a,
F
(x)

- hb,
  xi
  .
  qn
  n
  x∈Fq
  (2) The linear probability for the character χ of the linear approximation (a, b)
  of F is defined as
  LPF (χ, a, b) = |CORRF (χ, a, b)|2 .
  Remark 16. To be precise Baignères et al. [5] defined linear cryptanalysis over
  arbitrary abelian groups, in particular for maximal generality they defined the
  correlation with respect to two additive characters χ, ψ : Fq → C as
  CORRF (χ, ψ, a, b) =
   
  
  1 X 
  ·
  χ
  ·
  ψ
  .
  a,
  F
  (x)
  b,
  x
  qn
  n
  (29)
  x∈Fq
  Let Fq be a finite field of characteristic p, and let Tr : Fq → Fp be the absolute
  trace function, see [43, 2.22. Definition]. For all x ∈ Fq we define the function χ1
  as
  
  
  2·π·i
  · Tr(x) .
  χ1 (x) = exp
  p
  Then for every non-trivial additive character χ : Fq → C there exist a ∈ F×
  q such
  that χ(x) = χ1 (a · x), see [43, 5.7. Theorem]. Therefore, after an appropriate
  rescaling that we either absorb into a or b we can transform Equation (29) into
  Definition 15 (1).
  For Arion we approximate every round by affine functions, and we call the
  r+1
  a linear trail for Arion, where (ai−1 , ai ) is the
  tuple Ω = (a0 , . . . , ar ) ⊂ Fnp
  affine approximation of the ith round. Note that for A ∈ Fpn×n and c ∈ Fnp one has
  that LPAF +c (χ, a, b) = LPF (χ, A⊺ a, b). I.e., to bound the correlation of a single
  28round it suffices to bound the correlation of the Arion GTDS. By [48, Theorem 24,
  Corollary 25] we have that
  
  a, b = 0,
  
  1,
  2
  LPFArion (χ, a, b) ≤ (d2 − 1)
  (30)
  
  , else.
  
  q
  If we assume that the approximations among the rounds of Arion are indepen-
  dent, then
  !r
  (d2 − 1)2
  LPArion (χ, Ω) ≤
  .
  (31)
  q
  If a distinguisher is limited to D queries, then under heuristic assumptions
  Baignères et al. proved [5, Theorem 7] that the advantage of a distinguisher,
  which we call success probability, for a single linear trail is lower bounded by
  
  r
  2
  D
  psuccess  1 − e− 4 ·LPArion (χ,Ω) ≤ 1 − e
  −D
  4 ·
  (d2 −1)
  q
  .
  (32)
  For p ≥ 2N and d2 ≤ 2M we approximate exp(−x) ≈ 1 − x and estimate the
  security level κ of Arion against a linear trail Ω as
  !
  2 r
  D r·(2·M−N )
  (d2 − 1)
  D
  ≤
  ·
  ·2
  ≤ 2−κ
  (33)
  psuccess ≈
  4
  q
  4
  =⇒ κ ≤ 2 + r · (N − 2 · M ) − log2 (D) .
  (34)
  Now recall that Arion is supposed to be instantiated with a 64 bit prime
  number p & 264 and that d2 ∈ {121, 123, 125, 129, 161, 193, 195, 257}, then Equa-
  tion (33) implies the following security levels for Arion.
  Table 13. Security level of Arion against any linear trail for p ≥ 2N , d2 ≤ 29 and data
  amount 2M .
  N = 60
  r
  M
  κ (bits)
  5
  60
  152
  6
  120
  134
  8
  180
  158
  N = 120N = 250
  3
  120
  1882
  250
  216
  4
  240
  170
  3
  500
  198
  For a computationally limited adversary we next derive an estimation of the
  probability of a restricted linear hull.
  Lemma 17 (Restricted linear hull of Arion). Let p ∈ Z be a prime and let Fp
  be the field with p elements, let n > 1 and r ≥ 1 be integers, and let χ : Fp → C be
  (1)
  (r)
  a non-trivial additive character. Let Rk , . . . , Rk : Fnp ×Fnp → Fnp be Arion round
  292
  2
  functions, and let (a1 , ar ) ∈ Fnp \ {0} ∈ Fnp \ {0} . Assume that the linear
  approximations among the rounds of Arion are independent and that in every
  2
  round only up to M < (p − 1) · n many approximations (ai , ai+1 ) ∈ Fnp \ {0}
  can be utilized. Then
  !
  2 r
  (d2 − 1)
  r−1
  LPR(r) ◦···◦R(1) (χ, a1 , ar ) ≤ M
  ·
  .
  k
  k
  p
  Proof. By the independence of rounds we have that
  LPR(r) ◦···◦R(1) =
  k
  k
  r
  Y
  X
  LPR(i) (χ, ai , ai+1 )
  k
  a2 ,...,ar−1 i=1
  =
  r
  Y
  X
  LPF (i) (χ, circ(1, . . . , n)⊺ ai , ai+1 )
  a2 ,...,ar−1 i=1
  ≤
  X
  a2 ,...,ar−1
  2
  (d2 − 1)
  p
  !r
  2
  ≤M
  r−1
  ·
  (d2 − 1)
  p
  !r
  where the last two inequalities follow from Equation (30) and the assumption
  that only M approximations can be considered per round.
  ⊓
  ⊔
  p
  Provided that M < (d2 −1)
  2 , then the inequality is always less than 1. For the
  restricted differential hull we can then estimate the security level of Arion via
  !
  2 r
  (d2 − 1)
  r−1
  ≤ 2−κ
  (35)
  M
  ·
  p
  
  
  =⇒ κ ≤ r · log2 (p) − 2 · log2 (d2 − 1) − (r − 1) · log2 (M ) .
  (36)
  In Table 14 we report round number for various field sizes that achieve at least
  √
  128 bit security for a restricted hull of size p at round level.
  Table 14. Security level of Arion against linear cryptanalysis for p ≥ 2N and d2 ≤ 257
  with a restricted linear hull of size 2N/2 at round level.
  A.4
  rNκ (bits)
  5
  4
  360
  120
  250139
  267
  475
  Other Statistical Attacks
  By the definition of circ (1, . . . , n) a difference in one component can affect the
  whole state by a single round function call. Therefore, impossible differentials [12]
  and zero-correlation [14, 15] can hardly be mounted on 3 or more rounds.
  30Boomerang attacks [37, 52] search for quartets that satisfy two differential
  paths simultaneously. For our target primes BLS12 and BN254 no differentials
  with high probability exist for more than two rounds, see Equation (21). There-
  fore, a boomerang attack can hardly be mounted on 4 rounds of Arion & Arion-
  Hash.
  BAlgebraic Attacks on Arion
  B.1Higher-Order Differential & Interpolation Attacks
  Interpolation attacks [36] construct the polynomial vector representing a cipher
  without knowledge of the secret key. If such an attack is successful against a
  cipher, then an adversary can encrypt any plaintext without knowledge of the se-
  cret key. For a hash function the interpolated polynomial vector can be exploited
  to set up collision or forgery attacks. The cost of interpolating a polynomial de-
  pends on the number of monomials present in the polynomial vector representing
  the cipher function. Recall that any function F : Fnq → Fq can be represented
  by a unique polynomial f ∈ Fq [Xn ] = Fq [x1 , . . . , xn ]/ (xq1 − x1 , . . . , xqn − xn ),
  thus at most q n monomials can be present in f . Clearly, if f is dense, then an
  interpolation attack cannot be done faster than exhaustive search.
  Let R(i) ⊂ Fp [Xn ] denote the polynomial vector of the Arion-π round func-
  tion, we say that R(i+1) ◦ R(i) has a degree overflow if we have to reduce with at
  least one of the field equations to compute the unique representation in Fp [Xn ].
  As we saw in Table 2, for the fields BLS12 and BN254 some of our specified
  Arion parameters already achieve a degree overflow in the first round. Moreover,
  after the first round we expect that terms
  !e
  n
  n
  X
  X
  i · xi
  (37)
  i · xi +
  i=1
  i=1
  are present in every branch. After another application of the round function we
  expect to produce the terms
  !e
  !e
  n
  n
  X
  X
  i · xi
  mod (xp1 − x1 , . . . , xpn − xn )
  (38)
  i · xi +
  i=1
  i=1
  in every branch. By our specification e is the inverse exponent of a relatively
  low degree permutation, therefore we expect that the degrees of some of the
  aforementioned terms are close to p−2. In particular, this is the case if e2 ≥ p. In
  addition, we also expect that a big fraction of the monomials in Fp [Xn ] is present
  in the components of the polynomial vector of Arion-π after at least two iterations.
  Although an adversary can nullify the circulant matrix that is applied to the
  input vector in some attack scenarios, he cannot do so for later rounds. Therefore,
  we expect that at least after three rounds terms similar to Equation (38) are
  present in the polynomial vector of Arion-π. Further, to frustrate Meet-in-the-
  Middle (MitM) attacks we require that the number of rounds r ≥ 4.
  31We implemented Arion in SageMath [51] to compute the density of Arion-π
  for small primes. Our findings in Table 15 suggest that after two rounds Arion-π
  has already almost full density over Fp .
  Table 15. Observed minimum density for Arion-π for small primes, n = 3, 4, 5 and
  d1 , d2 = 3, 5.
  prMinimum density
  after 2 roundsDegree after 2
  roundsUnivariate degree
  after 2 rounds
  11
  13
  17
  19
  236
  6
  6
  6
  6≥ 82%
  ≥ 91%
  ≥ 91%
  ≥ 92%
  ≥ 90%n · (p − 1) − 1
  n · (p − 1) − 1
  n · (p − 1) − 1
  n · (p − 1) − 1
  n · (p − 1) − 1p−1
  p−1
  p−1
  p−1
  p−1
  If after two rounds the density of Arion-π is ≥ 0.8 · pn , then for p = 260
  and n = 3 already more than 2128 terms are present in Arion-π. Therefore, we
  expect that Arion resists interpolation attacks with the full key. An adversary
  can improve the capabilities of an interpolation attack by guessing parts of the
  key. If he can guess n − 1 parts correctly, then we expect that all univariate
  polynomials in Arion have degree close to p. To retrieve the remaining key an
  adversary has to factor at least one of the polynomials, due to the high degree
  his best choice to perform a greatest common
   divisor computation. We can then
  estimate the complexity with O p · log(p) . Therefore, for our target primes
  BLS12 and BN254 this complexity always exceeds the 128 bit security target
  after two rounds.
  For an interpolation attack on ArionHash we have a similar scenario as if
  we guessed parts of the key of Arion. At worst only one input of ArionHash is
  unknown to the adversary, he then still has to factor a polynomial of degree close
  to p, as mentioned before we do not expect that this defeats the 128 bit security
  target after two rounds.
  Higher-order differential attacks [11, 38, 40] exploit that higher differentials
  will vanish at some point. Since the density of Arion-π exceeds ≥ 0.8·pn after two
  rounds and its univariate degrees are close to p, we do not expect that higher-
  order differentials and distinguishers on Arion & ArionHash can undermine the
  128 bit security target for our target primes BLS12 and BN254 after two rounds.
  B.2
  Integral Attacks
  The notion of integral cryptanalysis was introduced by Knudsen & Wagner [39]
  and generalized to arbitrary finite fields by Beyne et al. [11]. It is based on the
  following properties of polynomial valued functions.
  Proposition 18 (see [11, Proposition 1, 2]). Let Fq be a finite field, and let
  F : Fnq → Fq be a polynomial valued function.
  32(1) Let V ⊂ Fnq be an affine subspace of dimension k. If deg (F ) < k · (q − 1),
  then
  X
  F (x) = 0.
  x∈V
  (2) Let G1 , . . . , Gn ⊂ F×
  q be multiplicative subgroups, and let G =
  product group. If degxi (F ) < |Gi | for all 1 ≤ i ≤ n, then
  X
  F (x) − |G| · F (0) = 0.
  ×ni=1 Gi their
  x∈G
  Proof. (1) was proven in [11, Corollary 1], for (2) we index the terms of F by j
  and then rearrange the sum
  X
  X
  X
  k
  F (x) =
  aj · x11,j · · · xknn,j
  x1 ∈G1 ,...,xn ∈Gn j≥0
  x∈G
  =
  X
  X
  j≥0 x1 ∈G1 ,...,xn ∈Gn
  k
  aj · x11,j · · · xknn,j .
  For the zero term this sums to |G| · F (0), now assume that at least one ki,j > 0,
  say k1,j , then we further rearrange
  !
  X
  X
  X
  k1,j
  k1,j
  kn,j
  kn,j
  aj · x1 · · · xn =
  aj · x1 · · · xn
  .
  x1 ∈G1 ,...,xn ∈Gn
  x2 ∈G1 ,...,xn ∈Gn
  x1 ∈G1
  For the inner sum on the right-hand side we can consider x2 , . . . , xn as non-zero
  field elements, by assumption k1,j < |G1 | so by [11, Proposition 2] the sum
  vanishes. This proves the claim.
  ⊓
  ⊔
  Let us investigate the capabilities of integral distinguishers on Arion-π. With-
  out loss of generality we can ignore the application circ(1, . . . , n) before the first
  Arion GTDS. To minimize the degree after the first Arion GTDS we fix all inputs
  except x1 , then only the first component is non-constant in x1 , but after appli-
  cation of the first affine layer all components contain the monomial xd11 . With
  our findings from interpolation attacks we expect that two additional rounds of
  Arion-π are sufficient to increase the degree in x1 of all components to at least
  p − 2.
  Therefore, for our target fields BLS12 and BN254 we do not expect that
  integral attacks can invalidate our 128 bit security target after two rounds of
  Arion & ArionHash. On the other hand, if Arion or ArionHash are instantiated
  over prime fields p < 2128 , then integral attacks could be a non-negligible threat
  since the degree of a univariate function can never exceed p.
  C
  Gröbner Basis Analysis
  In a Gröbner basis attack [17, 21] the adversary represents a cryptographic func-
  tion as fully determined system of polynomial equations and then solves for the
  33solutions of the system. Since the system is fully determined at least one solution
  of the polynomial system must contain the quantity of interest, e.g. the key of
  a block cipher or the preimage of a hash function. In general, a Gröbner basis
  attack proceeds in four steps:
  (1) Model the cryptographic function with a (iterated) system of polynomials.
  (2) Compute a Gröbner basis with respect to an efficient term order, e.g., the
  degree reverse lexicographic (DRL) order.
  (3) Perform a term order conversion to an elimination order, e.g., the lexico-
  graphic (LEX) order.
  (4) Solve the univariate equation.
  Let us for the moment assume that an adversary has already found a Gröbner
  basis and discuss the complexity of the remaining steps. Let k be a field, let
  I ⊂ k[x1 , . . . , xn ] be a zero-dimensional ideal modeling a cryptographic func-
  tion, and let d = dimk (k[x1 , . . . , xn ]/I) be the k-vector space dimension of the
  quotient space. With the original
  FGLM algorithm [28] the complexity of term
  
  order conversion is O n · d3 , but improved versions with probabilistic meth-
  ods achieve O (n · dω ) [27], where 2 ≤ 
  ω < 2.3727, and sparse linear algebra
  √
  n−1
  algorithms [29] achieve O
  n · d2+ n . To find the roots of the univariate
  polynomial from the elimination Gröbner basis5 f ∈ Fq [x] with d = deg (f )
  most efficiently we perform a greatest common divisor method that has recently
  been described in [8, §3.1].
  (1) Compute g = xq − x mod f .
  
  
  The computation of xq mod f requires O d · log (q) · log (d) · log log (d)
  field operations with a double-and-add algorithm.
  (2) Compute h = gcd (f, g).
  By construction h has the same roots as f in Fq since h = gcd (f, xq − x),
  but its degree is likely
   to be much lower.
  
  2
  field operations.
  This step requires O d · log (d) · log log (d)
  (3) Factor h.
  In general, the polynomial f coming from a 0-dimensional Gröbner basis has
  only a few roots in Fq .
  Thus, this step is negligible in complexity.
  Note that this method is only performative if deg (f ) < q else one has to exchange
  the roles of f and the field equation. Overall solving the polynomial system with
  probabilistic methods requires
  
  
  
  2
  (39)
  O n · dω + d · log (q) · log (d) · log log (d) + d · log (d) · log log (d)
  5
  If I is a radical ideal, then the degree of the univariate polynomial in the elimination
  Gröbner is indeed d, though for non-radical ideals the degree can be larger than d.
  For a simple example in the non-radical case consider (x2 ) ⊂ k[x].
  34field operations, and with deterministic methods
  √
  
  
  n−1
  2
  O
  n · d2+ n + d · log (q) · log (d) · log log (d) + d · log (d) · log log (d)
  (40)
  field operations.
  We must stress that we are unable to quantify the success probability of
  the probabilistic term order conversion. The probabilistic analysis of [29, §5.1]
  requires that the ideal is radical and that the homogenization of the DRL Gröb-
  ner basis is a regular sequence. We neither have a proof nor a disproof of these
  technical requirements for the Arion and ArionHash polynomial systems. There-
  fore, we cannot quantify how capable the probabilistic approach is aside from
  its complexity estimate.
  Let us now discuss the complexity of Gröbner basis computations. Today,
  the most efficient algorithms to compute Gröbner bases are Faugère’s linear
  algebra-based algorithms F4 [25] and F5 [26]. The main idea of linear algebra-
  based Gröbner basis algorithms can be traced back to Lazard [42]. Let F =
  {f1 , . . . , fm } ⊂ P = k[x1 , . . . , xn ] be a finite set of homogeneous polynomials
  over a field k. The homogeneous Macaulay matrix Md of degree d has columns
  indexed by monomials in s ∈ Pd and rows indexed by polynomials t · fj , where
  t ∈ P is a monomial such that deg (t · fj ) = d. The entry of row t · fj at column
  s is then the coefficient of the monomial s in t · fj . If F is an inhomogeneous
  system of polynomials, then one replaces Md by M≤d and the degree equality
  by an inequality. If one fixes a term order > on P , then by performing Gaussian
  elimination on Md respectively M≤d for a large enough value of d one produces
  a >-Gröbner basis of F . The least such d is called the solving degree sd> (F )
  of F . (The notion of solving degree was first introduced in [22], though we use
  the definition of [18].) The main improvement of F4/5 over Lazard’s method is
  the choice of efficient selection criteria. Conceptually, the Macaulay matrix will
  contain many redundant rows, if one is able to avoid many of these rows with
  selection criteria, then the running time of an algorithm will improve. Never-
  theless, with the notion of the solving degree it is possible to upper bound the
  maximal size of the computational universe of F4/5. It is well-known that the
  number of monomials in P of degree d is given by the binomial coefficient
  
  
  n+d−1
  N (n, d) =
  .
  (41)
  d
  We can now upper bound the size of the Macaulay matrix by m · d · N (n, d) × d ·
  N (n, d). Overall, we can bound the complexity of Gaussian elimination on the
  Macaulay matrix M≤d by
  
  ω !
  n+d−1
  O
  ,
  (42)
  d
  where we absorb m and d in the implied constant since in general N (n, d) ≫ m, d
  and ω ≥ 2 is a linear algebra constant.
  35In the cryptographic literature a generic approach to bound the solving de-
  gree is the so-called Macaulay bound. Assume that deg (f1 ) ≥ . . . ≥ deg (fm ) and
  let l = min{m, n + 1}, then the Macaulay bound of F is given by
  MBF =
  l
  X
  i=1
  deg (f1 ) + . . . + deg (fl ) − l + 1.
  (43)
  Up to date there are two known cases when one indeed has that sdDRL (F ) ≤
  MBF :
  (1) F is regular [6, 7], or
  (2) F is in generic coordinates [18, Theorem 9, 10].
  We stress that we were unable to prove that one of the two cases applies to Arion
  or ArionHash. Therefore, we can only hypothesize that the Macaulay bound
  upper bounds the solving degree.
  From a designer’s perspective, during all our small scale experiments the vec-
  tor space dimension of the quotient space behaved more stable with respect to
  the chosen primes, branch sizes and round numbers than the observed solving de-
  gree. Thus, all our extrapolated security claims of Arion and ArionHash
  with respect to Gröbner basis attacks are expressed in the complexity
  of solving for the solutions of the polynomial system after a Gröbner
  basis has been found.
  Moreover, in all our experiments we observed that the quotient space dimen-
  sion grows exponentially in r and the base only depends on n, d1 and d2 . We
  hypothesize that this behaviour is invariant under the chosen primes and parame-
  ters. Therefore, we will do an “educated guess” from our small scale experiments
  to obtain the base and extrapolate our findings for our security analysis.
  C.1
  Arion
  For the security of Arion against Gröbner basis attacks we consider the following
  polynomial model:
  (i) We do not consider a key schedule, i.e., in every round we add the same key
  k = (k1 , . . . , kn ) ∈ Fnp .
  (ii) We use a single plain/cipher pair p, c ∈ Fnp given by Arion to set up a fully
  determined polynomial system F .
  
  
  (i)
  (i)
  (iii) For 1 ≤ i ≤ r − 1 we denote with x(i) = x1 , . . . , xn the intermediate
  state variables, in addition we set x(0) = circ (1, . . . , n)· (p + k) and x(r) = c.
  Further, for 1 ≤ i ≤ r we denote with z (i) an auxiliary variable. For our
  polynomial model we choose a slight modification F̃ = {f˜1 , . . . , f˜n } of the
  Arion GTDS F = {f1 , . . . , fn }, see Definition 1, where we set f˜n (xn ) = xn .
  For 1 ≤ i ≤ n − 1 the f˜i follow the same iterative definition P
  as the original
  n−1
  polynomials in the GTDS, and we modify σ̃i+1,n = xn + zn + j=i+1
  xj + f˜j .
  36We consider the following polynomial model as the naive model Fnaive for
  Arion
  
  
  circ (1, . . . , n) F̃ x̂(i−1) + ci + k − x(r) = 0,
  e
  
  − z (i) = 0,
  x(i−1)
  n
  
  
  (i)
  (i)
  where x̂(i) = x1 , . . . , xn−1 , z (i) .
  (iv) Obviously the naive polynomial system Fnaive contains high degree equations
  given by the power permutation xe . Though, if we replace the auxiliary
  equations by
   d2
  (i)
  = 0,
  x(i)
  n − z
  then we obtain a polynomial system FArion whose polynomials are of small
  degree. Further, we expect the arithmetic of FArion to be independent of
  the chosen prime, i.e., for primes p, q ∈ P such that gcd (di , p − 1) = 1 =
  gcd (di , q − 1) we expect no notable difference for the complexity of a Gröb-
  ner basis attack.
  Lemma 19. Let Fp be a finite field, let Fnaive and FArion be the polynomial
  models from (iii) and (iv), and let F be the ideal of all field equations in the
  polynomial ring of Fnaive and FArion . Then
  (Fnaive ) + F = (FArion ) + F.
  
  e
  (i)
  Proof. By definition, we have that xn
  ≡ z (i) mod (Fnaive ) + F , by raising
  this congruence to the d2 th power yields
  e d2
   d2
  
  (i)
  ≡ x(i)
  x(i)
  n ≡ z
  n
  which proves the claim.
  mod (Fnaive ) + F
  ⊓
  ⊔
  I.e., on the solutions that solely come from the finite field Fp , which are
  the solutions of cryptographic interest, the varieties corresponding to Fnaive and
  FArion coincide, so
  V (Fnaive ) ∩ Fnp = V (FArion ) ∩ Fnp .
  (44)
  Thus, FArion is indeed a well-founded model for Arion.
  To compute the Macaulay bound of FArion we first apply circ(1, . . . , n)−1 in
  every round to cancel the mixing of the components, then we use Lemma 2 with
  e = 1 to compute the degree in the ith component. Further, we have to account
  for the auxiliary equation, overall we yield the Macaulay bound
  !
  n−1
  X
  
  n−i
- 1 − r · (n + 1)
  2
  · (d1 + 1) − d1
  MBn,d1 ,d2 (r) = r · d2 + 1 +
  (45)
  i=1
  
  
  = r · d2 + 2 · (d1 + 1) · 2n−1 − 1 − (n − 1) · d1 − n + 1
  37We stress that we were unable to prove or disprove that FArion is regular or
  in generic coordinates, thus we can only hypothesize the Macaulay bound as
  measure for the complexity of linear algebra-based Gröbner basis algorithms.
  We implemented FArion in the OSCAR computer algebra system [47] and com-
  puted the Gröbner basis with its F4 implementation. Unfortunately, the log
  function of F4 only prints the current working degree of the algorithm not its
  current solving degree. As remedy, we estimate the empirical solving degree as
  follows: We sum up all positive differences of the working degrees between consec-
  utive steps and add this quantity to the largest input polynomial degree. After
  computing the Gröbner basis we also computed the Fp -vector space dimension
  of the quotient ring. We conducted our experiments with the primes
  p1 = 1013,
  p2 = 10007,
  {z
  }
  |
  p3 = 1033,
  p4 = 15013.
  {z
  }
  |
  and
  gcd(3,pi −1)=1
  (46)
  gcd(5,pi −1)=1
  d1 = 3
  d1 = 5
  MB2,3,7 (r)
  MB2,5,7 (r)
  40
  dimFp (FArion )
  sdDRL (FArion )
  All computations were performed on an AMD EPYC-Rome (48) CPU with 94
  GB RAM.
  In Figure 1 we record our empirical results for n = 2 and in Figure 2 we
  record our empirical results for n = 3. From experiments, we hypothesize that
  for r > 1 the solving degree is indeed bounded by the Macaulay bound and that
  the quotient space dimension grows exponential in r.
  20
  0
  1
  2
  r
  3
  d1 = 3
  d1 = 5
  35r
  49r
  104
  102
  1
  2
  r
  3
  Fig. 1. Experimental solving degree and vector space dimension of the quotient ring
  for Arion with n = 2 and d2 = 7.
  38d1 = 3
  d1 = 5
  MB3,3,7 (r)
  MB3,5,7 (r)
  dimFp (FArion )
  sdDRL (FArion )
  80
  40
  0
  1
  105
  104
  103
  102
  2
  d1 = 3
  d1 = 5
  175r
  1
  2
  r
  r
  Fig. 2. Experimental solving degree and vector space dimension of the quotient ring
  for Arion with n = 3 and d2 = 7.
  To better understand the growth of the base of the quotient space dimension
  we computed the quotient space dimension for n ≤ 4 and r = 1, see Table 16.
  Table 16. Empirical growth of the vector space dimension of the quotient space of
  Arion with r = 1.
  nd1d2dimFp ,emp
  2
  2
  3
  4
  2
  2
  33
  3
  3
  3
  5
  5
  57
  257
  7
  7
  7
  257
  735
  1285
  175
  875
  49
  1799
  343
  From our experiments we conjecture that the quotient space dimension grows
  or is bounded via
  r
  
  n−1
  .
  (47)
  dimFp (FArion ) (n, r, d1 , d2 ) = d2 · (d1 + 2)
  We use two approaches to estimate the cost of Gröbner basis computations
  for FArion . First, note that if we exclude one polynomial from FArion , then we do
  not have a zero-dimensional polynomial system anymore. Therefore, the highest
  non-trivial lower bound for the solving degree of FArion has to be the highest
  degree in the Arion polynomial system, i.e.,
  
  (48)
  max d2 , 2n−1 · (d1 + 1) − d1 ≤ sdDRL (FArion ) .
  We combine this lower bound with Equation (42) to derive the “minimal pos-
  sible“ Gröbner basis complexity estimate. Second, under the assumption that
  the Macaulay bound always upper bounds the solving degree of Arion we can
  39combine Equations (42) and (43) to derive the “maximal possible” Gröbner ba-
  sis complexity estimate. For ease of computation we approximated the binomial
  coefficient with
    r
  n
  n
  ≈
  · 2n·H2 (k/n) ,
  (49)
  k
  π · k · (n − k)
  where H2 (p) = −p · log2 (p) − (1 − p) · log2 (1 − p) denotes the binary entropy (cf.
  [20, Lemma 17.5.1]). To estimate the cost of system solving we plug Equation (47)
  into Equations (39) and (40). For all estimates we assume that our adversary
  has an optimal Gaussian elimination algorithm with ω = 2.
  We base the security of Arion against Gröbner basis attacks solely
  on the complexity of solving the polynomial system. I.e., even if an
  adversary can compute a Gröbner basis in O(1) it must be computationally
  infeasible to find a solution of the polynomial system. Thus, we estimate the
  security level κ of Arion against a Gröbner basis attack via
  
  (50)
  κ ≤ log2 O (System solving) .
  Table 17. Empirical cost estimation of Gröbner basis attacks on Arion for primes
  p ≥ 260 and d2 ≥ 121. The column GB min contains the complexity of Gröbner basis
  computations estimated with the highest polynomial degree in the Arion polynomial
  system. The column GB MB contains the complexity Gröbner basis computations
  estimated via the Macaulay bound. We assume that the adversary has an optimal
  Gaussian elimination algorithm with ω = 2.
  n
  r
  GB min
  (bits)
  GB MB
  (bits)
  Deterministic
  Solving (bits)
  Probabilistic
  Solving (bits)
  d1 = 3
  3
  3
  4
  4
  5
  5
  6
  6
  8
  84
  6
  4
  5
  3
  4
  3
  4
  2
  3123
  162
  143
  166
  134
  162
  150
  181
  204
  279129
  170
  160
  186
  164
  201
  208
  257
  244
  338
  3
  3
  4
  4
  5
  5
  6
  8
  84
  5
  3
  5
  3
  4
  3
  2
  3123
  143
  118
  166
  134
  162
  172
  225
  309131
  153
  135
  195
  174
  215
  225
  263
  366
  137
  207
  165
  207
  145
  194
  166
  222
  138
  20896
  143
  115
  143
  101
  134
  115
  153
  96
  143
  149
  187
  136
  229
  162
  217
  187
  158
  238104
  129
  95
  158
  113
  149
  130
  110
  164
  d1 = 5
  40C.2
  ArionHash
  In a preimage attack on ArionHash we are given a given hash value α ∈ Fp and
  ′
  we have to find x ∈ Frp such that ArionHash(x) = α. In a second-preimage attack
  we assume that we are given a message that consists of two input blocks, i.e. y =
   ′ 2
  ′
  (y1 , y2 ) ∈ Frp . Now we again have to find x ∈ Frp such that ArionHash(y) =
  ArionHash(x). Consequently, both preimage attacks on ArionHash reduce to the
  same equation
    
  
  xin
  α
  Arion-π
  ,
  (51)
  =
  IV
  xout
  where α ∈ Fp is the output of ArionHash, IV ∈ Fcp is the initial value, and
  xin = (xin,1 , . . . , xin,r′ ) and xout = (xout,2 , . . . , xout,n ) are indeterminates. Analog
  to Arion we construct an iterated polynomial system where each round is modeled
  with a polynomial vector. Note that Equation (51) is not fully determined if n ≥ 3
  and r′ > 1, in such a case we have n − 1 + r′ > n many variables for xin and xout
  and (r − 1) · n many intermediate state variables, but we only have r · n many
  equations. In order to obtain a fully determined system the adversary either
  has to guess some values for xin or xout , but each guess has success probability
  1/p, so we can neglect this approach, or he has to add additional equations.
  If an adversary is unable to exploit additional algebraic structures of Arion-π
  (which are unknown to us), then he has just one generic choice: he has to add
  field equations until the system is fully determined. This approach introduces
  polynomials of very high degree, thus we do not expect this attack to be feasible
  below our 128 bit Gröbner basis security claim. Therefore, in the analysis of this
  section we always choose c = n − 1 to obtain a fully determined system. Note
  that the Macaulay bound for ArionHash is identical to the one for Arion, see
  Equation (45).
  We implemented ArionHash in the OSCAR [47] computer algebra system and
  computed the Gröbner basis of FArionHash with F4. As initial value we chose
  IV = 0c . Further, we computed the Fp -vector space dimension of the quotient
  space for ArionHash. We used the same primes as for Arion, see Equation (46). In
  Figure 3 we record our empirical results for n = 2 and in Figure 4 we record our
  empirical results for n = 3. From experiments, we hypothesize that the solving
  degree is indeed bounded by the Macaulay bound and that the quotient space
  dimension grows exponential in r.
  41dimFp (FArionHash )
  sdDRL (FArionHash )
  d1 = 3
  d1 = 5
  MB2,3,7 (r)
  MB2,5,7 (r)
  40
  20
  0
  1
  2
  r
  d1 = 3
  d1 = 5
  35r
  49r
  104
  102
  1
  3
  2
  r
  3
  80
  dimFp (FArionHash )
  sdDRL (FArionHash )
  Fig. 3. Experimental solving degree and vector space dimension of the quotient ring
  for ArionHash with n = 2 and d2 = 7.
  d1 = 3
  d1 = 5
  MB3,3 (r)
  MB3,5 (r)
  40
  0
  1
  d1 = 3
  d1 = 5
  91r
  133r
  104
  103
  102
  1
  2
  2
  r
  r
  Fig. 4. Experimental solving degree and vector space dimension of the quotient ring
  for ArionHash with n = 3 and d2 = 7.
  To better understand the growth of the base of the quotient space dimension
  we computed the quotient space dimension for n ≤ 5 and r = 1, see Table 18.
  42Table 18. Empirical growth of the vector space dimension of the quotient space of
  ArionHash.
  nd1d2dimFp ,emp
  2
  2
  3
  3
  4
  5
  2
  2
  3
  3
  4
  53
  3
  3
  3
  3
  3
  5
  5
  5
  5
  5
  57
  257
  7
  257
  7
  7
  7
  257
  7
  257
  7
  735
  1285
  91
  3341
  203
  427
  49
  1799
  133
  4833
  301
  637
  From our experiments we conjecture that the quotient space dimension grows
  or is bounded via
  r
  
  (52)
  dimFp (FArionHash) (n, r, d1 , d2 ) = 2n−1 · d2 · (d1 + 1) − d1 · d2 .
  As for Arion, we apply this equation in Equations (39) and (40) to estimate the
  cost of solving the ArionHash polynomial system.
  Since the polynomial degrees of ArionHash coincide with the ones for Arion
  also the Gröbner basis complexity estimates of ArionHash and Arion coincide.
  43Table 19. Empirical cost estimation of preimage Gröbner basis attacks on ArionHash
  for primes p ≥ 260 and d2 ≥ 121. The column GB min contains the complexity of
  Gröbner basis computations estimated with the highest polynomial degree in the Arion
  polynomial system. The column GB MB contains the complexity Gröbner basis com-
  putations estimated via the Macaulay bound. We assume that the adversary has an
  optimal Gaussian elimination algorithm with ω = 2.
  n
  r
  GB min
  (bits)
  GB MB
  (bits)
  Deterministic
  Solving (bits)
  Probabilistic
  Solving (bits)
  d1 = 3
  3
  3
  4
  4
  5
  5
  6
  6
  8
  8
  5
  6
  4
  6
  4
  5
  4
  5
  3
  4
  143
  162
  143
  187
  162
  187
  181
  208
  279
  345
  150
  170
  160
  210
  201
  235
  257
  301
  338
  423
  158
  190
  141
  212
  154
  193
  167
  208
  143
  191110
  132
  98
  146
  107
  133
  115
  143
  100
  132
  133
  200
  167
  195
  161
  201
  130
  217
  148
  19893
  138
  103
  128
  111
  139
  91
  149
  103
  137
  d1 = 5
  3
  3
  4
  4
  5
  5
  6
  6
  8
  8
  4
  6
  4
  5
  3
  5
  3
  5
  3
  4
  123
  162
  143
  179
  162
  187
  172
  244
  309
  385
  131
  173
  143
  166
  215
  251
  225
  328
  366
  461
  An adversary can also set up a collision attack via the equation
  ArionHash(x) = ArionHash(y),
  (53)
  where x, y ∈ F∗p are variables. We can construct the corresponding collision
  polynomial system by first preparing two preimage polynomial systems F1 and
  F2 , then we connect the last round of both systems via
  
  0 0 ... 0
  
  
  n 1 . . . n − 1 
  
  0
  
  
  (r−1)
  ,
- # c
  F̃
  x̂
  
  
  r
  ..
  xout
  
  
  .
  
  2 3 ...
  (54)
  1
  d2
  
  = 0,
  − z (r)
  x(r−1)
  n
  44
  (55)
  0 0 ... 0
  
  
  n 1 . . . n − 1 
  
  0
  
  
  (r−1)
  ,
- cr =
   F̃ ŷ
  
  ..
  yout
  
  
  .
  
  2 3 ...
  
  (56)
  1
  
  d2
  = 0,
  yn(r−1) − z (r)
  
  1 2 ... n
  0 0 . . . 0   
  
  
  
  
  
  = 0.
  F̃ x̂(r−1) − F̃ ŷ(r−1)
  
  
  .
  .. 
  
  0 0 ... 0
  (57)
  (58)
  Note that the collision polynomial system consists of 2 · r · (n + 1) − 1 many
  equations in 2 · r · (n + 1) many variables, hence it is not fully determined. There-
  fore, one has to guess at least one variable in the collision polynomial system. In
  all our experiments we guessed one value of yout and then computed the DRL
  Gröbner basis of the system. We observed that the position of our guess in yout
  did never affect the dimension of the quotient space dimension. In Table 20 we
  record the observed quotient space dimension in our experiments.
  Table 20. Empirical growth of the vector space dimension of the quotient space of the
  collision polynomial system of ArionHash.
  nrd1 d2
  2
  2
  2
  2
  3
  2
  3
  3
  3
  3
  4
  41
  2
  1
  1
  1
  1
  1
  1
  1
  1
  1
  13
  3
  5
  3
  3
  5
  3
  5
  3
  5
  3
  3
  3
  3
  5
  7
  7
  7
  3
  5
  7
  7
  3
  7
  dimFp ,emp
  225
  50625
  1225
  1225
  8281
  2401
  1521
  9025
  8281
  17689
  7569
  41209
  From our experiments we conjecture that the quotient space dimension grows
  or is bounded via
  2
  (59)
  dimFp (FArionHash,coll ) (n, r, d1 , d2 ) = dimFp (FArionHash) (n, r, d1 , d2 )
  Hence, the complexity estimate for the collision attack corresponds to the one
  for the preimage attack with doubled number of rounds. In particular, Table 19
  implies at least 128 bit security against a collision Gröbner basis attack.
  45DPerformance Evaluation
  D.1libsnark Implementation
  We implemented ArionHash, Griffin and Poseidon using the C++ library libsnark
  [50] that is used in the privacy-protecting digital currency Zcash [35]. The com-
  parative result of our experiment is given in Table 21. All experiments were run
  on a system with an Intel Core i7-11800H CPU and 32 GB RAM on a Clear
  Linux instance, using the g++-12 compiler with -O3 -march=native flags. Our
  result shows that ArionHash significantly outperforms Poseidon showing 2x ef-
  ficiency improvement, and the α-ArionHash is considerably faster than Griffin
  for practical parameter choices e.g. n = 3 or 4 in Merkle tree hashing mode.
  Table 21. Performance of various hash functions for proving the membership of a
  Merkle tree accumulator over BN254 with d2 = 257. Proving times are in ms. For all
  hash functions veriﬁcation requires less than 20 ms.
  Time (ms)
  Height
  ArionHash
  Griffin
  α-ArionHash
  Poseidon
  d=5n=3n=4n=8n=3n=4n=8n=3n=4n=8n=3n=4n=8
  4
  8
  16
  32101
  211
  392
  730103
  216
  401
  751142
  294
  554
  104673
  145
  278
  50987
  177
  334
  646143
  294
  553
  104788
  181
  338
  62299
  209
  387
  727133
  270
  505
  980186
  386
  745
  1422212
  417
  805
  1550274
  566
  1095
  2111
  D.2
  Dusk Network Plonk
  We implemented ArionHash in Rust using the Dusk Network Plonk [23] library
  which uses 3-wire constraints, and compared its performance to the Dusk Net-
  work Poseidon [24] reference implementation. We note that for ArionHash we
  only implemented a generic affine layer, i.e., ArionHash needs n · (n − 1) con-
  straints for matrix multiplication and n constraints for constant addition. The
  comparative result of our experiment is given in Table 22. All experiments were
  run on a system with a Intel Core i7-10700 and 64 GB RAM. Our result shows
  that ArionHash outperforms Poseidon showing 2x efficiency improvement in
  Merkle tree hashing mode.
  46Table 22. Performance of various hash functions for proving the membership of a
  Merkle tree accumulator over BLS12 for d1 = 5 and n = 5. For each hash function
  and each Merkle tree height 1000 samples were collected. Average times plus-minus the
  standard deviation are given in the table. All timings are in ms. For all hash functions
  veriﬁcation requires less than 10 ms.
  Time (ms)
  Native sponge
  Proof generation
  ArionHashPoseidon
  0.121 ± 0.001
  167 ± 10.101 ± 0.001
  313 ± 4
  Merkle Tree
  Height481632481632
  Proof generation323 ± 4533 ± 51011 ± 132034 ± 16524 ± 61008 ± 141975 ± 274015 ± 57
