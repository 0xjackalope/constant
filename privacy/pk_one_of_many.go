package privacy

import (
	"fmt"
	"math/big"
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
func (pro *PKOneOfManyProtocol) Prove(commitments [][]byte, indexIsZero int, rand []byte, commitmentValue []byte, index byte) (*PKOneOfManyProof, error) {
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
	indexIsZeroBinary := ConvertIntToBinany(indexIsZero)

	//
	r := make([][]byte, n+1)
	a := make([][]byte, n+1)
	s := make([][]byte, n+1)
	t := make([][]byte, n+1)
	u := make([][]byte, n)

	cl := make([][]byte, n+1)
	ca := make([][]byte, n+1)
	cb := make([][]byte, n+1)
	cd := make([][]byte, n)

	for j := 1; j <= n; j++ {
		// Generate random numbers
		r[j] = make([]byte, 32)
		r[j] = RandBytes(32)
		a[j] = make([]byte, 32)
		a[j] = RandBytes(32)
		s[j] = make([]byte, 32)
		s[j] = RandBytes(32)
		t[j] = make([]byte, 32)
		t[j] = RandBytes(32)
		u[j-1] = make([]byte, 32)
		u[j-1] = RandBytes(32)

		// convert indexIsZeroBinary[j] to big.Int
		indexInt := big.NewInt(int64(indexIsZeroBinary[j-1]))

		// Calculate cl, ca, cb, cd
		// cl = Com(l, r)
		cl[j] = make([]byte, 34)
		cl[j] = Pcm.CommitSpecValue(indexInt.Bytes(), r[j], index)

		// ca = Com(a, s)
		ca[j] = make([]byte, 34)
		ca[j] = Pcm.CommitSpecValue(a[j], s[j], index)

		// cb = Com(la, t)
		la := new(big.Int)
		la.Mul(indexInt, new(big.Int).SetBytes(a[j]))
		cb[j] = make([]byte, 34)
		cb[j] = Pcm.CommitSpecValue(la.Bytes(), t[j], index)

	}

	// cd_k =
	// Calculate: ci^pi,k
	for k:=0; k< n; k++{
		// Calculate pi,k which is coefficient of x^k in polynomial pi(x)
		res := big.NewInt(1)
		tmp := big.NewInt(0)
		for i:=0; i<n; i++{
			// represent i in binary
			iBinary := ConvertIntToBinany(i)
			pik := GetCoefficient(iBinary, k, n, a, indexIsZeroBinary)
			//pik := make([]byte, 32)
			tmp.Exp(new(big.Int).SetBytes(commitments[i]), pik, big.NewInt(0))
			res.Mul(res, tmp)
		}
		comZero := Pcm.CommitSpecValue(big.NewInt(0).Bytes(), u[k], index)
		res.Mul(res, new(big.Int).SetBytes(comZero))
		cd[k] = make([]byte, 32)
		copy(cd[k], res.Bytes())
	}


	// Calculate x
	x := big.NewInt(0)
	for j:=1; j<n; j++{
		x.SetBytes(Pcm.getHashOfValues([][]byte{x.Bytes(), cl[j], ca[j], cb[j], cd[j-1]}))
	}
	//x.Mod(x, Curve.Params().N)

	// Calculate za, zb zd
	//res := Poly{big.NewInt(1)}
	//var fji Poly
	f := make([][]byte, n+1)
	za := make([][]byte, n+1)
	zb := make([][]byte, n+1)
	zd := make([]byte, 32)

	for j:=1; j<=n; j++{
		// f = lx + a
		fInt:=  big.NewInt(0)
		fInt.Mul(big.NewInt(int64(indexIsZeroBinary[j])), x)
		fInt.Add(fInt, new(big.Int).SetBytes(a[j]))
		f[j] = fInt.Bytes()

		// za = s + rx
		zaInt := big.NewInt(0)
		zaInt.Mul(new(big.Int).SetBytes(r[j]), x)
		zaInt.Add(zaInt, new(big.Int).SetBytes(s[j]))
		za[j] = zaInt.Bytes()

		// zb = r(x - f) + t
		zbInt := big.NewInt(0)
		zbInt.Sub(x, fInt)
		zbInt.Mul(zbInt, new(big.Int).SetBytes(r[j]))
		zbInt.Add(zbInt, new(big.Int).SetBytes(t[j]))
		zb[j] = zbInt.Bytes()

	}

	zdInt := big.NewInt(0)
	zdInt.Exp(x, big.NewInt(int64(n)), nil )
	zdInt.Mul(zdInt, new(big.Int).SetBytes(rand))

	uxInt := big.NewInt(0)
	sumInt := big.NewInt(0)
	for k:=0; k<n; k++{
		uxInt.Exp(x, big.NewInt(int64(k)), nil )
		uxInt.Mul(uxInt, new(big.Int).SetBytes(u[k]))
		sumInt.Add(sumInt, uxInt)
	}

	zdInt.Sub(zdInt, sumInt)
	copy(zd, zdInt.Bytes())

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
	res, err := pk.Prove(commitments, indexIsZero, commitments[indexIsZero], randoms[indexIsZero], SN_CM)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
}

// Get coefficient of x^k in polynomial pi(x)
func GetCoefficient(iBinary []byte , k int, n int, a [][]byte, l []byte ) *big.Int{
	res := Poly{big.NewInt(1)}
	var fji Poly
	for j:=1; j<=n; j++{
		fj := Poly{new(big.Int).SetBytes(a[j]), big.NewInt(int64(l[j-1])) }
		if iBinary[j-1] == 0 {
			fji = Poly{big.NewInt(0), big.NewInt(1)}.Sub(fj, nil)

		} else{
			fji = fj
		}
		res = res.Mul(fji, nil)
	}
	return res[k+1]
}
