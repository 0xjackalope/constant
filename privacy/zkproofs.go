package privacy

// type Proof interface{
// 	Params()

// 	CreateDDHProof()[]byte,

// }
// type Proofs struct {
// 	DDHProof []byte,
// 	CM01 []byte,

// }

// ZkpPedersenCMProof contains proof's value
type ZkpPedersenCMProof struct {
	Alpha, Beta, GammaAddr, GammaValue, GammaSN, GammaR []byte
}

// ZkpPedersenCMProve create zero knowledge proof for an opening of a Pedersen commitment
func ZkpPedersenCMProve(cm CommitmentParams, pubKey PublicKey, sn SerialNumber, value []byte) *ZkpPedersenCMProof {
	zkp := new(ZkpPedersenCMProof)
	return zkp
}

// ZkpPedersenCMVerify check the proof's value
func ZkpPedersenCMVerify(proofsvalue ZkpPedersenCMProof, commintmentsvalue []byte) bool {
	return true
}



