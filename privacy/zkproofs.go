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

// SetWitness sets witnesses
func (pro *ProtocolForPKCommittedValues) SetWitness(witnesses [][]byte) {
	pro.witnesses = make([][]byte, len(witnesses))
	for i := 0; i < len(witnesses); i++ {
		copy(pro.witnesses[i], witnesses[i])
	}
}

// Prove creates zero knowledge proof for an opening of a Pedersen commitment
func (pro *ProtocolForPKCommittedValues) Prove(commitmentValue []byte) (*ProofForPKCommittedValues, error) {
	if len(pro.witnesses) != 4 {
		return nil, Logger.log.Errorf("len of witnesses must be equal to 4")
	}

	proof := new(ProofForPKCommittedValues)

	// Calculate random numbers
	r := make([][]byte, CM_CAPACITY)
	for i := 0; i < CM_CAPACITY; i++ {
		r[i] = RandBytes(32)
	}

	// Calculate alpha
	alpha := new(EllipticPoint)
	tmp := new(EllipticPoint)
	alpha.X, alpha.Y = Curve.ScalarMult(Pcm.G[0].X, Pcm.G[0].Y, r[0])
	tmp.X, tmp.Y = Curve.ScalarMult(Pcm.G[1].X, Pcm.G[1].Y, r[1])
	alpha.X, alpha.Y = Curve.Add(alpha.X, alpha.Y, tmp.X, tmp.Y)
	tmp.X, tmp.Y = Curve.ScalarMult(Pcm.G[2].X, Pcm.G[2].Y, r[2])
	alpha.X, alpha.Y = Curve.Add(alpha.X, alpha.Y, tmp.X, tmp.Y)
	tmp.X, tmp.Y = Curve.ScalarMult(Pcm.G[3].X, Pcm.G[3].Y, r[3])
	alpha.X, alpha.Y = Curve.Add(alpha.X, alpha.Y, tmp.X, tmp.Y)
	proof.Alpha = make([]byte, 33)
	copy(proof.Alpha, CompressKey(*alpha))

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

	// Calculate gammas
	b := new(big.Int)
	witness := new(big.Int)
	bMulWitness := new(big.Int)
	randTmp := new(big.Int)

	b.SetBytes(proof.Beta)
	proof.Gammas = make([][]byte, CM_CAPACITY)

	for i := 0; i < CM_CAPACITY; i++ {
		witness.SetBytes(pro.witnesses[i])
		bMulWitness.Mul(b, witness)
		proof.Gammas[i] = make([]byte, 32)
		copy(proof.Gammas[i], bMulWitness.Add(bMulWitness, randTmp.SetBytes(r[i])).Bytes())
	}
	return proof, nil
}

// Verify check the proof's value
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

	// Calculate right point:
	rightPoint := EllipticPoint{big.NewInt(0), big.NewInt(0)}
	tmpPoint := new(EllipticPoint)
	for i := 0; i < CM_CAPACITY; i++ {
		tmpPoint.X, tmpPoint.Y = Curve.ScalarMult(Pcm.G[i].X, Pcm.G[i].Y, proof.Gammas[i])
		rightPoint.X, rightPoint.Y = Curve.Add(rightPoint.X, rightPoint.Y, tmpPoint.X, tmpPoint.Y)
	}

	Logger.log.Infof("commitment value: %v\n", commitmentValue)
	commitmentPoint, err := DecompressCommitment(commitmentValue)
	if err != nil {
		Logger.log.Errorf("Decompress commitment error: %v\n", err.Error())
	}

	alphaPoint, err := DecompressKey(proof.Alpha)
	if err != nil {
		Logger.log.Errorf("Decompress alpha error: %v\n", err.Error())
	}

	// Calculate left point:
	xY, yY := Curve.ScalarMult(commitmentPoint.X, commitmentPoint.Y, beta)
	LeftPoint := new(EllipticPoint)
	LeftPoint.X, LeftPoint.Y = Curve.Add(xY, yY, alphaPoint.X, alphaPoint.Y)

	// Check whether right point is equal left point or not
	if (rightPoint.X.CmpAbs(LeftPoint.X) == 0) && (rightPoint.Y.CmpAbs(LeftPoint.Y) == 0) {
		return false
	}
	return true
}
