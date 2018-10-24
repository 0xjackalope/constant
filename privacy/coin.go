package privacy

// Coin represents a coin
type Coin struct {
	Address,
	SerialNumber,
	Value,
	CoinCommitment,
	R []byte
}

// Commitment represents a commitment that includes 4 generators
type Commitment interface {
	// Params returns the parameters for the commitment
	Params() *CommitmentParams
	// InitCommitment initialize the parameters
	InitCommitment() *CommitmentParams
	// Commit commits
	Commit(address, value, serialNumber []byte) []byte
}

// CommitmentParams represents the parameters for the commitment
type CommitmentParams struct {
	G0, G1, G2, H EllipticPoint // generator
}

// Commit commits a coin
func (coin *Coin) Commit() {
	var cm Commitment
	coin.CoinCommitment = cm.Commit(coin.Address, coin.Value, coin.SerialNumber)
}

// hashGenerator create new generator from a generator using hash function
func hashGenerator(g EllipticPoint) EllipticPoint {
	// TODO: res.X = hash(g.X), res.Y = sqrt(res.X^3 - 3X + B)
	res := new(EllipticPoint)
	// res.X = new(big.Int).SetBytes(sha256.Sum256(g.X.Bytes()))
	return *res
}

// Params return the parameters of commitment
func (com *CommitmentParams) Params() *CommitmentParams {
	return com
}

// InitCommitment initial
func (com *CommitmentParams) InitCommitment() {
	// TODO: how to generate generators independently
	com.G0 = EllipticPoint{Curve.Params().Gx, Curve.Params().Gy}
	com.G1 = hashGenerator(com.G0)
	com.G2 = hashGenerator(com.G1)
	com.H = hashGenerator(com.G2)
}

// Commit commits 3 attributes of coin
func (com *CommitmentParams) Commit(address, value, serialNumber []byte) []byte {
	var res []byte
	// TODO: using Pedersen commitment
	// TODO: convert result from Elliptic to bytes array
	return res
}

// type Cryptosystem struct{
// 	privateKey
// 	publicKey

// }

// func genKey()

// func(self * Cryptosystem) Encrypt(){

// }
