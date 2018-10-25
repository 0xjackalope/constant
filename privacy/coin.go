package privacy

// Coin represents a coin
type Coin struct {
	//PrdNumber
	Address,
	SerialNumber,
	Value,
	CoinCommitment,
	R,
	Info []byte
}

// Commit commits a coin
func (coin *Coin) Commit() {
	var cm Commitment
	coin.CoinCommitment = cm.Commit(coin.R, coin.Address, coin.Value, coin.SerialNumber)
}

// type Cryptosystem struct{
// 	privateKey
// 	publicKey

// }

// func genKey()

// func(self * Cryptosystem) Encrypt(){

// }
