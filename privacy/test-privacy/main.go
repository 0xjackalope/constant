package main

import (
	"fmt"
	"math/big"

	"github.com/ninjadotorg/cash/privacy"
)

func main() {

	// fmt.Printf("N: %X\n", privacy.Curve.Params().N)
	// fmt.Printf("P: %X\n", privacy.Curve.Params().P)
	// fmt.Printf("B: %X\n", privacy.Curve.Params().B)
	// fmt.Printf("Gx: %x\n", privacy.Curve.Params().Gx)
	// fmt.Printf("Gy: %X\n", privacy.Curve.Params().Gy)
	// fmt.Printf("BitSize: %X\n", privacy.Curve.Params().BitSize)

	//spendingKey := privacy.GenerateSpendingKey(new(big.Int).SetInt64(123).Bytes())
	//fmt.Printf("\nSpending key: %v\n", spendingKey)
	//fmt.Println(len(spendingKey))
	//
	//address := privacy.GeneratePublicKey(spendingKey)
	//fmt.Printf("\nAddress: %v\n", address)
	//fmt.Println(len(address))
	//point, err := privacy.DecompressKey(address)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Printf("Pk decom: %v\n", point)
	//
	//receivingKey := privacy.GenerateReceivingKey(spendingKey)
	//fmt.Printf("\nReceiving key: %v\n", receivingKey)
	//fmt.Println(len(receivingKey))
	//
	//transmissionKey := privacy.GenerateTransmissionKey(receivingKey)
	//fmt.Printf("\nTransmission key: %v\n", transmissionKey)
	//fmt.Println(len(transmissionKey))
	//
	//point, err = privacy.DecompressKey(transmissionKey)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Printf("Transmission key point decompress: %+v\n ", point)
	//
	//msg := "hello, world"
	//hash := sha256.Sum256([]byte(msg))
	//
	//signature, err := privacy.Sign(hash[:], spendingKey)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("signature: %v\n", signature)
	//
	//valid := privacy.Verify(signature, hash[:], address)
	//fmt.Println("\nsignature verified:", valid)
	//
	//tx, _ := transaction.CreateEmptyTxs()
	//fmt.Printf("Transaction: %+v\n", tx)

	//a := "aaaaaaaaaaaaaaaaa"

	// //var b privacy.PedersenCommitment
	//var xx privacy.PCParams
	//xx.InitCommitment()
	//
	//var sn privacy.SerialNumber
	//var v privacy.Value
	//sn = []byte("aaaaaaa")
	//v = []byte("bbbbbbb")
	//
	//m := make(map[string][]byte)
	//
	////m["sn"] = sn
	////m["v"] = v
	//m = map[string][]byte{
	//	"sn": sn,
	//	"v": v,
	//}
	//fmt.Printf("m['sn']: %+v\n", m["sn"])
	//fmt.Printf("m['v']: %+v\n", m["v"])
	//
	//fmt.Println(xx.CommitAll(m))

	privacy.Pcm.InitCommitment()

	spendingKey := privacy.GenerateSpendingKey(new(big.Int).SetInt64(123).Bytes())
	fmt.Printf("\nSpending key: %v\n", spendingKey)

	pubKey := privacy.GeneratePublicKey(spendingKey)
	serialNumber := privacy.RandBytes(32)

	// value := make([]byte, 32)
	c := big.NewInt(0)
	value := c.Bytes()
	// binary.LittleEndian.PutUint32(value, 1)
	fmt.Printf("Value: %v\n", value)
	r := privacy.RandBytes(32)
	coin := privacy.Coin{
		PublicKey:      pubKey,
		SerialNumber:   serialNumber,
		CoinCommitment: nil,
		R:              r,
		Value:          value,
	}
	coin.CommitAll()
	fmt.Println(coin.CoinCommitment)
	// cm1 := coin.CommitPublicKey()
	// fmt.Println(cm1)
	// cm2 := coin.CommitValue()
	// fmt.Println(cm2)
	// cm3 := coin.CommitSerialNumber()
	// fmt.Println(cm3)

	// witnesses := make([][]byte, privacy.CM_CAPACITY)

	witness := [][]byte{
		coin.PublicKey,
		coin.Value,
		coin.SerialNumber,
		coin.R,
	}

	var pk privacy.PKComZeroOneProtocol
	pk.SetWitness(witness)
	pk.Prove(coin.CoinCommitment, 1)

	// var zk privacy.ZKProtocols

	// pk := zk.GetPKCommittedValues()
	// pk.SetWitness(witness)
	// proof, _ := zk.GetPKCommittedValues().Prove(coin.CoinCommitment)

	// fmt.Printf("Proof: %+v\n", proof)

	// fmt.Println(zk.GetPKCommittedValues().Verify(*proof, coin.CoinCommitment))

}
