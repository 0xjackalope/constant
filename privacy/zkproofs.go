package privacy

import (
	"math/big"
	"math/rand"
	"time"

	"github.com/minio/blake2b-simd"
)

//ZkpPedersenCMProof contains proof's value
type ZkpPedersenCMProof struct {
	Alpha, Beta, GammaAddr, GammaValue, GammaSN, GammaR []byte
}

func getRandom() []byte {
	rand.Seed(time.Now().UTC().Unix())
	r := make([]byte, 32)
	rand.Read(r)
	return r
}

// ZkpPedersenCMProve create zero knowledge proof for an opening of a Pedersen commitment
func ZkpPedersenCMProve(cm CommitmentParams, pubKey PublicKey, sn SerialNumber, value, cmRnd []byte) *ZkpPedersenCMProof {
	zkp := new(ZkpPedersenCMProof)
	rand.Seed(time.Now().UTC().Unix())
	r0 := getRandom()
	r1 := getRandom()
	r2 := getRandom()
	r3 := getRandom()
	alpha := new(EllipticPoint)
	tmp := new(EllipticPoint)
	//
	alpha.X, alpha.Y = Curve.ScalarMult(cm.G0.X, cm.G0.Y, r0)
	tmp.X, tmp.Y = Curve.ScalarMult(cm.G1.X, cm.G1.Y, r1)
	alpha.X, alpha.Y = Curve.Add(alpha.X, alpha.Y, tmp.X, tmp.Y)
	tmp.X, tmp.Y = Curve.ScalarMult(cm.G2.X, cm.G2.Y, r2)
	alpha.X, alpha.Y = Curve.Add(alpha.X, alpha.Y, tmp.X, tmp.Y)
	tmp.X, tmp.Y = Curve.ScalarMult(cm.H.X, cm.H.Y, r3)
	alpha.X, alpha.Y = Curve.Add(alpha.X, alpha.Y, tmp.X, tmp.Y)
	copy(zkp.Alpha, CompressKey(*alpha))
	//
	hashFunc := blake2b.New256()
	appendStr := append(CompressKey(cm.G0), CompressKey(cm.G1)...)
	appendStr = append(appendStr, CompressKey(cm.G2)...)
	appendStr = append(appendStr, CompressKey(cm.H)...)
	appendStr = append(appendStr, cmRnd...)
	appendStr = append(appendStr, CompressKey(*alpha)...)
	beta := hashFunc.Sum(appendStr)
	copy(zkp.Beta, beta)
	//
	tmpRand := new(big.Int)
	b := new(big.Int)

	//compute GammaAddr
	b.SetBytes(zkp.Beta)
	addrVal := new(big.Int)
	addrVal.SetBytes(pubKey)
	x := b.Mul(b, addrVal)
	copy(zkp.GammaAddr, x.Add(x, tmpRand.SetBytes(r0)).Bytes())
	//compute GammaValue
	b.SetBytes(zkp.Beta)
	coinVal := new(big.Int)
	coinVal.SetBytes(value)
	x = b.Mul(b, coinVal)
	copy(zkp.GammaValue, x.Add(x, tmpRand.SetBytes(r1)).Bytes())
	//compute GammaSerialNumber
	b.SetBytes(zkp.Beta)
	serialVal := new(big.Int)
	serialVal.SetBytes(sn)
	x = b.Mul(b, serialVal)
	copy(zkp.GammaSN, x.Add(x, tmpRand.SetBytes(r2)).Bytes())

	//compute GammaR
	b.SetBytes(zkp.Beta)
	cmRandVal := new(big.Int)
	cmRandVal.SetBytes(cmRnd)
	x = b.Mul(b, cmRandVal)
	copy(zkp.GammaR, x.Add(x, tmpRand.SetBytes(r3)).Bytes())
	return zkp
}

//ZkpPedersenCMVerify check the proof's value
func ZkpPedersenCMVerify(cm CommitmentParams, proofsvalue ZkpPedersenCMProof, commitmentsvalue []byte) bool {

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
		//fmt.Println("Cannot decompress commitments value to ECC point")
	}

	alphaPoint, error := DecompressKey(proofsvalue.Alpha)
	if error != nil {
		//fmt.Println("Cannot decompress alpha to ECC point")
	}

	xY, yY := Curve.ScalarMult(commitmentsPoint.X, commitmentsPoint.Y, Beta)

	xLeft, yLeft := Curve.Add(xY, yY, alphaPoint.X, alphaPoint.Y)

	if (xRight.CmpAbs(xLeft) == 0) && (yRight.CmpAbs(yLeft) == 0) {
		return false
	}
	return true
}
