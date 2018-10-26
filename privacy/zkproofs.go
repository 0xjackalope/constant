package privacy

import (
	"fmt"

	"github.com/minio/blake2b-simd"
)

// type Proof interface{
// 	Params()

// 	CreateDDHProof()[]byte,

// }
// type Proofs struct {
// 	DDHProof []byte,
// 	CM01 []byte,

// }

//ZkpPedersenCM contains proof's value
type ZkpPedersenCM struct {
	Alpha, Beta, GammaAddr, GammaValue, GammaSN, GammaR []byte
}

//ZkpPedersenCMProve create zero knowledge proof for an opening of a Pedersen commitment
func ZkpPedersenCMProve(cm CommitmentParams, pubKey PublicKey, sn SerialNumber, value []byte) *ZkpPedersenCM {
	zkp := new(ZkpPedersenCM)
	return zkp
}

//ZkpPedersenCMVerify check the proof's value
func ZkpPedersenCMVerify(cm CommitmentParams, proofsvalue ZkpPedersenCM, commitmentsvalue []byte) bool {

	plainBeta := append(CompressKey(cm.G0), CompressKey(cm.G1)...)
	plainBeta = append(plainBeta, CompressKey(cm.G2)...)
	plainBeta = append(plainBeta, CompressKey(cm.H)...)
	plainBeta = append(plainBeta, commitmentsvalue...)
	plainBeta = append(plainBeta, proofsvalue.Alpha...)

	hashMachine := blake2b.New256()
	hashMachine.Write(plainBeta)

	Beta := hashMachine.Sum(nil)

	xH, yH := Curve.ScalarMult(cm.H.X, cm.H.Y, proofsvalue.GammaR)
	xG0, yG0 := Curve.ScalarMult(cm.G0.X, cm.G0.Y, proofsvalue.GammaAddr)
	xG1, yG1 := Curve.ScalarMult(cm.G1.X, cm.G1.Y, proofsvalue.GammaValue)
	xG2, yG2 := Curve.ScalarMult(cm.G2.X, cm.G2.Y, proofsvalue.GammaSN)

	xRight, yRight := Curve.Add(xH, yH, xG0, yG0)
	xRight, yRight = Curve.Add(xRight, yRight, xG1, yG1)
	xRight, yRight = Curve.Add(xRight, yRight, xG2, yG2)

	commitmentsPoint, error := DecompressKey(commitmentsvalue)
	if error != nil {
		fmt.Println("Cannot decompress commitments value to ECC point")
	}

	alphaPoint, error := DecompressKey(proofsvalue.Alpha)
	if error != nil {
		fmt.Println("Cannot decompress alpha to ECC point")
	}

	xY, yY := Curve.ScalarMult(commitmentsPoint.X, commitmentsPoint.Y, Beta)

	xLeft, yLeft := Curve.Add(xY, yY, alphaPoint.X, alphaPoint.Y)

	if (xRight.CmpAbs(xLeft) == 0) && (yRight.CmpAbs(yLeft) == 0) {
		return false
	}
	return true
}
