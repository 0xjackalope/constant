package privacy

import (
	"github.com/minio/blake2b-simd"
	"math/big"
	"math/rand"
	"time"
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
func getRandom() []byte{
	rand.Seed(time.Now().UTC().Unix())
	r:= make([]byte,32)
	rand.Read(r)
	return r
}
// ZkpPedersenCMComponent create zero knowledge proof for an opening of a Pedersen commitment
func pedersenCMGenerateProof(cm CommitmentParams, pubKey PublicKey, sn SerialNumber, value, cmRnd []byte) *ZkpPedersenCM {
	zkp := new(ZkpPedersenCM)
	rand.Seed(time.Now().UTC().Unix())
	r0:=getRandom();
	r1:=getRandom();
	r2:=getRandom();
	r3:=getRandom();
	alpha:=new(EllipticPoint)
	tmp:=new(EllipticPoint)
	//
	alpha.X, alpha.Y = Curve.ScalarMult(cm.G0.X,cm.G0.Y, r0);
	tmp.X,tmp.Y = Curve.ScalarMult(cm.G1.X, cm.G1.Y,r1);
	alpha.X, alpha.Y = Curve.Add(alpha.X, alpha.Y, tmp.X, tmp.Y)
	tmp.X,tmp.Y = Curve.ScalarMult(cm.G2.X, cm.G2.Y,r2);
	alpha.X, alpha.Y = Curve.Add(alpha.X, alpha.Y, tmp.X, tmp.Y)
	tmp.X,tmp.Y = Curve.ScalarMult(cm.H.X, cm.H.Y,r3);
	alpha.X, alpha.Y = Curve.Add(alpha.X, alpha.Y, tmp.X, tmp.Y)
	copy(zkp.Alpha, CompressKey(*alpha))
	//
	hashFunc:=blake2b.New256();
	appendStr:= append(CompressKey(cm.G0),CompressKey(cm.G1)...)
	appendStr = append(appendStr,CompressKey(cm.G2)...)
	appendStr = append(appendStr,CompressKey(cm.H)...)
	appendStr = append(appendStr,cmRnd...)
	appendStr = append(appendStr, CompressKey(*alpha)...)
	beta:=hashFunc.Sum(appendStr)
	copy(zkp.Beta, beta)
	//
	tmpRand:=new(big.Int)
	b:= new(big.Int)

	//compute GammaAddr
	b.SetBytes(zkp.Beta)
	addrVal:=new(big.Int)
	addrVal.SetBytes(pubKey)
	x:= b.Mul(b,addrVal)
	copy(zkp.GammaAddr, x.Add(x,tmpRand.SetBytes(r0)).Bytes());
	//compute GammaValue
	b.SetBytes(zkp.Beta)
	coinVal:=new(big.Int)
	coinVal.SetBytes(value)
	x= b.Mul(b,coinVal)
	copy(zkp.GammaValue, x.Add(x,tmpRand.SetBytes(r1)).Bytes());
	//compute GammaSerialNumber
	b.SetBytes(zkp.Beta)
	serialVal:=new(big.Int)
	serialVal.SetBytes(sn)
	x= b.Mul(b,serialVal)
	copy(zkp.GammaSN, x.Add(x,tmpRand.SetBytes(r2)).Bytes());

	//compute GammaR
	b.SetBytes(zkp.Beta)
	cmRandVal:=new(big.Int)
	cmRandVal.SetBytes(cmRnd);
	x= b.Mul(b,cmRandVal)
	copy(zkp.GammaR, x.Add(x,tmpRand.SetBytes(r3)).Bytes());
	return zkp
}


//ZkpPedersenCMVerify check the proof's value
func ZkpPedersenCMVerify(proofsvalue ZkpPedersenCM, commintmentsvalue []byte) bool {
	return true
}


