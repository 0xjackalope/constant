package privacy

type SerialNumber []byte   //32 bytes
type CoinCommitment []byte //32 bytes

// Coin represents a coin
type Coin struct {
	Address        PublicKey
	SerialNumber   SerialNumber
	CoinCommitment CoinCommitment
	Value, R, Info []byte
}

// Commit commits a coin
func (coin *Coin) Commit() {
	var Cm CommitmentParams
	Cm.InitCommitment()
	coin.CoinCommitment = Cm.Commit(coin.R, coin.Address, coin.Value, coin.SerialNumber)
}
