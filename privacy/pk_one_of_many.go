package privacy

// PKOneOfManyProtocol is a protocol for Zero-knowledge Proof of Knowledge of one out of many commitments containing 0
// include witnesses: commitedValue, r []byte
type PKOneOfManyProtocol struct {
	witnesses [][]byte
}

// PKOneOfManyProof contains proof's value
type PKOneOfManyProof struct {
	ca, cb    []byte // 34 bytes
	f, za, zb []byte //32 bytes
}

// SetWitness sets witnesses
func (pro *PKOneOfManyProtocol) SetWitness(witnesses [][]byte) {
	pro.witnesses = make([][]byte, len(witnesses))
	for i := 0; i < len(witnesses); i++ {
		pro.witnesses[i] = make([]byte, len(witnesses[i]))
		copy(pro.witnesses[i], witnesses[i])
	}
}

// Prove creates proof for one out of many commitments containing 0
func (pro *PKOneOfManyProtocol) Prove(commitments [][]byte, l int, commitmentValue []byte, index byte) (*PKOneOfManyProof, error) {

	return nil, nil
}
