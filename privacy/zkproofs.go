package privacy

import (
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

type ZkpPedersenCM struct {
	Gamma1, Gamma2, Gamma3, Gamma4 []byte
}
func getRandom() []byte{
	rand.Seed(time.Now().UTC().Unix())
	r:= make([]byte,32)
	rand.Read(r)
	return r
}
// ZkpPedersenCMComponent create zero knowledge proof for an opening of a Pedersen commitment
func pedersenCMGenerateProof(cm CommitmentParams, pubKey PublicKey, sn SerialNumber, value []byte) *ZkpPedersenCM{
	zkp := new(ZkpPedersenCM)
	rand.Seed(time.Now().UTC().Unix())
	r0:=getRandom();
	r1:=getRandom();
	r2:=getRandom();
	r3:=getRandom();
	alpha:=new(EllipticPoint)

	alpha.X, alpha.Y = Curve.ScalarMult(cm.G0.X,cm.G0.Y, r0);
	alpha.X, alpha.Y = Curve.Add(alpha.X, alpha.Y, Curve.ScalarMult(cm.G1.X, cm.G1.Y,r1))
	return zkp
}
