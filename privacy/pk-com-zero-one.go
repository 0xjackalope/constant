package privacy

import (
	"fmt"
	"math/big"
)

// PKComZeroOneProtocol is a protocol for Zero-knowledge Proof of Knowledge of committed zero or one
// include witnesses: pk, value, sn, r []byte
type PKComZeroOneProtocol struct {
	witnesses [][]byte
	// Proof     *ProofForPKCommittedValues
}

// PKComZeroOneProof contains proof's value
type PKComZeroOneProof struct {
	Alpha  []byte
	Gammas [][]byte
}

// SetWitness sets witnesses
func (pro *PKComZeroOneProtocol) SetWitness(witnesses [][]byte) {
	pro.witnesses = make([][]byte, len(witnesses))
	for i := 0; i < len(witnesses); i++ {
		pro.witnesses[i] = make([]byte, len(witnesses[i]))
		copy(pro.witnesses[i], witnesses[i])
	}
	// fmt.Printf("Witness: %+v\n", pro.witnesses)
}

func (pro *PKComZeroOneProtocol) Prove(commitmentValue []byte, index byte) (*PKComZeroOneProof, error) {
	// m := binary.BigEndian.Uint64(pro.witnesses[int(index)])
	indexInt := int(index)
	fmt.Printf("index int: %v\n", indexInt)
	if indexInt < 0 || indexInt > 3 {
		return nil, fmt.Errorf("index must be between 0 and 3")
	}

	witness := big.NewInt(0)
	witness.SetBytes(pro.witnesses[index])
	fmt.Printf("witness: %v\n", witness)

	if witness.Cmp(big.NewInt(0)) != 0 && witness.Cmp(big.NewInt(0)) != 1 {
		return nil, fmt.Errorf("witness must be zero or one")
	}

	// Generate random numbers
	a := RandBytes(32)
	s := RandBytes(32)
	t := RandBytes(32)

	// // Calculate ca, cb
	ca := Pcm.CommitSpecValue(a, s, index)
	am := make([]byte, 32)
	// if witness == 1 {
	// 	copy(am, a[:])
	// } else {
	// 	copy(am, pro.witnesses[indexInt])
	// }
	fmt.Printf("a: %v\n", a)
	fmt.Printf("am: %v\n", am)

	cb := Pcm.CommitSpecValue(am, t, index)
	fmt.Printf("ca: %v\n", ca)
	fmt.Printf("cb: %v\n", cb)

	// fmt.Printf("am: %v\n", am)
	// cb := Pcm.CommitSpecValue(pro.witnesses[0], pro.witnesses[1], index)
	return nil, nil
}
