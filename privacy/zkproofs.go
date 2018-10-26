package privacy

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
<<<<<<< HEAD
	Alpha, Beta, GammaAddr, GammaValue, GammaSN, GammaR []byte
=======
	Gamma1, Gamma2, Gamma3, Gamma4 []byte
>>>>>>> aa3494236443b7891d5c2f5681f3035e08c77de1
}

//ZkpPedersenCMProve create zero knowledge proof for an opening of a Pedersen commitment
func ZkpPedersenCMProve(cm CommitmentParams, pubKey PublicKey, sn SerialNumber, value []byte) *ZkpPedersenCM {
	zkp := new(ZkpPedersenCM)
	return zkp
}

<<<<<<< HEAD
//ZkpPedersenCMVerify check the proof's value
func ZkpPedersenCMVerify(proofsvalue ZkpPedersenCM, commintmentsvalue []byte) bool {
	return true
}
=======
>>>>>>> aa3494236443b7891d5c2f5681f3035e08c77de1
