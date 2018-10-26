package privacy

type SerialNumber 	[]byte 	//32 bytes
type CoinCommitment []byte 	//32 bytes
type Random 		[]byte	//32 bytes
type Value 		[]byte	//32 bytes


// Coin represents a coin
type Coin struct {
	PublicKey      PublicKey
	SerialNumber   SerialNumber
	CoinCommitment CoinCommitment
	R              Random
	Value, Info    []byte
}

// Commit commits a coin
func (coin *Coin) Commit() {
	var Cm CommitmentParams
	Cm.InitCommitment()
	values := map[string][]byte{
		"pk": coin.PublicKey,
		"v": coin.Value,
		"sn": coin.SerialNumber,
		"r": coin.R,
	}
	coin.CoinCommitment, _ = Cm.Commit(values)
}
