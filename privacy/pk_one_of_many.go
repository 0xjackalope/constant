package privacy

import (
	"fmt"
	"math/big"
	"strconv"
)

// PKOneOfManyProtocol is a protocol for Zero-knowledge Proof of Knowledge of one out of many commitments containing 0
// include witnesses: commitedValue, r []byte
type PKOneOfManyProtocol struct {
	witnesses [][]byte
}

// PKOneOfManyProof contains proof's value
type PKOneOfManyProof struct {
	ca, cb    []byte // 34 bytes
	f, za, zb []byte //32 bytes
}

// SetWitness sets witnesses
func (pro *PKOneOfManyProtocol) SetWitness(witnesses [][]byte) {
	pro.witnesses = make([][]byte, len(witnesses))
	for i := 0; i < len(witnesses); i++ {
		pro.witnesses[i] = make([]byte, len(witnesses[i]))
		copy(pro.witnesses[i], witnesses[i])
	}
}

// Prove creates proof for one out of many commitments containing 0
func (pro *PKOneOfManyProtocol) Prove(commitments [][]byte, indexIsZero int, commitmentValue []byte, index byte) (*PKOneOfManyProof, error) {
	n := len(commitments)
	// Check the number of commitment list's elements
	if !IsPowerOfTwo(n) {
		return nil, fmt.Errorf("the number of commitment list's elements must be power of two")
	}

	// Check indexIsZero
	if indexIsZero > n || index < 0 {
		return nil, fmt.Errorf("index is zero must be index in list of commitments")
	}

	// Check index
	if index < 0 || index > 2 {
		return nil, fmt.Errorf("index must be between 0 and 2")
	}

	// represent indexIsZero in binary
	indexIsZeroBinary := make([]byte, 32)
	str := strconv.FormatInt(int64(indexIsZero), 2)
	for i := 0; i < len(str); i++ {
		indexIsZeroBinary[i] = ConvertAsciiToInt(str[i])
	}
	fmt.Printf("inddex in binary: %v\n", indexIsZeroBinary)

	//
	r := make([][]byte, n+1)
	a := make([][]byte, n+1)
	s := make([][]byte, n+1)
	t := make([][]byte, n+1)
	u := make([][]byte, n+1)

	cl := make([][]byte, n)
	ca := make([][]byte, n)
	cb := make([][]byte, n)
	// cd := make([][]byte, n)

	for i := 1; i <= n; i++ {
		// Generate random numbers
		r[i] = make([]byte, 32)
		r[i] = RandBytes(32)
		a[i] = make([]byte, 32)
		a[i] = RandBytes(32)
		s[i] = make([]byte, 32)
		s[i] = RandBytes(32)
		t[i] = make([]byte, 32)
		t[i] = RandBytes(32)
		u[i-1] = make([]byte, 32)
		u[i-1] = RandBytes(32)

		// convert indexIsZeroBinary[i] to big.Int
		indexInt := big.NewInt(int64(indexIsZeroBinary[i]))

		// Calculate cl, ca, cb, cd
		// cl = Com(l, r)
		cl[i] = make([]byte, 34)
		cl[i] = Pcm.CommitSpecValue(indexInt.Bytes(), r[i], index)

		// ca = Com(a, s)
		ca[i] = make([]byte, 34)
		ca[i] = Pcm.CommitSpecValue(a[i], s[i], index)

		// cb = Com(la, t)
		la := new(big.Int)
		la.Mul(indexInt, new(big.Int).SetBytes(a[i]))
		cb[i] = make([]byte, 34)
		cb[i] = Pcm.CommitSpecValue(la.Bytes(), t[i], index)

		// cd =

	}

	return nil, nil
}

//TestPKOneOfMany test protocol for one of many commitment is commitment to zero
func TestPKOneOfMany() {
	Pcm.InitCommitment()
	pk := new(PKOneOfManyProtocol)

	indexIsZero := 23

	// list of commitments
	commitments := make([][]byte, 32)
	serialNumbers := make([][]byte, 32)
	randoms := make([][]byte, 32)

	for i := 0; i < 32; i++ {
		serialNumbers[i] = RandBytes(32)
		randoms[i] = RandBytes(32)
		commitments[i] = make([]byte, 34)
		commitments[i] = Pcm.CommitSpecValue(serialNumbers[i], randoms[i], SN_CM)
	}

	// create commitment to zero at indexIsZero

	serialNumbers[indexIsZero] = big.NewInt(0).Bytes()
	commitments[indexIsZero] = Pcm.CommitSpecValue(serialNumbers[indexIsZero], randoms[indexIsZero], SN_CM)
	res, err := pk.Prove(commitments, indexIsZero, commitments[indexIsZero], SN_CM)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
}
