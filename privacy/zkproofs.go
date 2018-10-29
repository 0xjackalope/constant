package privacy

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/minio/blake2b-simd"
)

// ProtocolForPKCommittedValues is a protocol for Zero-knowledge Proof of Knowledge of committed values
// include witnesses
type ProtocolForPKCommittedValues struct {
	witnesses [][]byte
	// Proof     *ProofForPKCommittedValues
}

// ProofForPKCommittedValues contains proof's value
type ProofForPKCommittedValues struct {
	Alpha, Beta []byte
	Gammas      [][]byte
}

//func GetRandom() []byte {
//	rand.Seed(time.Now().UTC().Unix())
//	r := make([]byte, 32)
//	rand.Read(r)
//	return r
//}

func (pro *ProtocolForPKCommittedValues) SetWitness(witnesses [][]byte) {
	pro.witnesses = make([][]byte, len(witnesses))
	for i := 0; i < len(witnesses); i++ {
		copy(pro.witnesses[i], witnesses[i])
	}
}

// Prove create zero knowledge proof for an opening of a Pedersen commitment
func (pro *ProtocolForPKCommittedValues) Prove(commitmentValue []byte) *ProofForPKCommittedValues {
	proof := new(ProofForPKCommittedValues)
	//rand.Seed(time.Now().UTC().Unix())
	r := make([][]byte, CM_CAPACITY)
	for i := 0; i < 4; i++ {
		r[i] = RandBytes(32)
	}
	alpha := new(EllipticPoint)
	tmp := new(EllipticPoint)
	//
	alpha.X, alpha.Y = Curve.ScalarMult(Pcm.G[0].X, Pcm.G[0].Y, r[0])
	tmp.X, tmp.Y = Curve.ScalarMult(Pcm.G[1].X, Pcm.G[1].Y, r[1])
	alpha.X, alpha.Y = Curve.Add(alpha.X, alpha.Y, tmp.X, tmp.Y)
	tmp.X, tmp.Y = Curve.ScalarMult(Pcm.G[2].X, Pcm.G[2].Y, r[2])
	alpha.X, alpha.Y = Curve.Add(alpha.X, alpha.Y, tmp.X, tmp.Y)
	tmp.X, tmp.Y = Curve.ScalarMult(Pcm.G[3].X, Pcm.G[3].Y, r[3])
	alpha.X, alpha.Y = Curve.Add(alpha.X, alpha.Y, tmp.X, tmp.Y)
	proof.Alpha = make([]byte, 33)
	copy(proof.Alpha, CompressKey(*alpha))
	//fmt.Printf("Alpha: %+v\n", zkp.Alpha)
	//Compute commitment value

	// calculate beta
	hashFunc := blake2b.New256()
	appendStr := append(CompressKey(Pcm.G[0]), CompressKey(Pcm.G[1])...)
	appendStr = append(appendStr, CompressKey(Pcm.G[2])...)
	appendStr = append(appendStr, CompressKey(Pcm.G[3])...)
	appendStr = append(appendStr, commitmentValue...)
	appendStr = append(appendStr, CompressKey(*alpha)...)
	hashFunc.Write(appendStr)
	beta := hashFunc.Sum(nil)
	proof.Beta = make([]byte, 32)
	copy(proof.Beta, beta)

	b := new(big.Int)
	witness := new(big.Int)
	bMulWitness := new(big.Int)
	randTmp := new(big.Int)

	b.SetBytes(proof.Beta)
	fmt.Printf("Len witness: %v\n", len(pro.witnesses))
	proof.Gammas = make([][]byte, 4)

	for i := 0; i < len(pro.witnesses); i++ {
		witness.SetBytes(pro.witnesses[i])
		bMulWitness.Mul(b, witness)
		proof.Gammas[i] = make([]byte, 32)
		copy(proof.Gammas[i], bMulWitness.Add(bMulWitness, randTmp.SetBytes(r[i])).Bytes())
	}
	return proof
}

//ZkpPedersenCMVerify check the proof's value
func (pro *ProtocolForPKCommittedValues) Verify(proof ProofForPKCommittedValues, commitmentValue []byte) bool {
	// re-calculate beta and check whether it is equal to beta of proof or not
	hashFunc := blake2b.New256()
	appendStr := append(CompressKey(Pcm.G[0]), CompressKey(Pcm.G[1])...)
	appendStr = append(appendStr, CompressKey(Pcm.G[2])...)
	appendStr = append(appendStr, CompressKey(Pcm.G[3])...)
	appendStr = append(appendStr, commitmentValue...)
	appendStr = append(appendStr, proof.Alpha...)
	hashFunc.Write(appendStr)
	beta := hashFunc.Sum(nil)
	if !bytes.Equal(beta, proof.Beta) {
		fmt.Println("beta is not equal")
		return false
	}

	rightPoint := EllipticPoint{big.NewInt(0), big.NewInt(0)}
	tmpPoint := new(EllipticPoint)
	for i := 0; i < CM_CAPACITY; i++ {
		tmpPoint.X, tmpPoint.Y = Curve.ScalarMult(Pcm.G[i].X, Pcm.G[i].Y, proof.Gammas[i])
		rightPoint.X, rightPoint.Y = Curve.Add(rightPoint.X, rightPoint.Y, tmpPoint.X, tmpPoint.Y)
	}

	fmt.Printf("commitment value: %v\n", commitmentValue)
	commitmentPoint, error := DecompressCommitment(commitmentValue)

	if error != nil {
		fmt.Println("Cannot decompress commitments value to ECC point")
	}

	alphaPoint, error := DecompressKey(proof.Alpha)
	if error != nil {
		fmt.Println("Cannot decompress alpha to ECC point")
	}

	xY, yY := Curve.ScalarMult(commitmentPoint.X, commitmentPoint.Y, beta)
	LeftPoint := new(EllipticPoint)

	LeftPoint.X, LeftPoint.Y = Curve.Add(xY, yY, alphaPoint.X, alphaPoint.Y)

	if (rightPoint.X.CmpAbs(LeftPoint.X) == 0) && (rightPoint.Y.CmpAbs(LeftPoint.Y) == 0) {
		return false
	}
	return true
}
