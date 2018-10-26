package privacy

// type Proof interface{
// 	Params()

// 	CreateDDHProof()[]byte,

// }
// type Proofs struct {
// 	DDHProof []byte,
// 	CM01 []byte,

// }

type ZkpPedersenCM struct {
	Gamma1, Gamma2 []byte
}

// ZkpPedersenCMComponent create zero knowledge proof for an opening of a Pedersen commitment
func ZkpPedersenCMComponent(cm CommitmentParams, pubKey PublicKey, sn SerialNumber, value []byte) *ZkpPedersenCM{
	zkp := new(ZkpPedersenCM)
	return zkp
}
