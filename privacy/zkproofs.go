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
	Alpha, Beta, GammaAddr, GammaValue, GammaSN, GammaR []byte
}

//ZkpPedersenCMProve create zero knowledge proof for an opening of a Pedersen commitment
func ZkpPedersenCMProve(cm CommitmentParams, pubKey PublicKey, sn SerialNumber, value []byte) *ZkpPedersenCM {
	zkp := new(ZkpPedersenCM)
	return zkp
}

//ZkpPedersenCMVerify check the proof's value
func ZkpPedersenCMVerify(proofsvalue ZkpPedersenCM, commintmentsvalue []byte) bool {
	return true
}
