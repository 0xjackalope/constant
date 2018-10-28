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
	value := []byte("10")
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
	cm1 := coin.CommitPublicKey()
	fmt.Println(cm1)
	cm2 := coin.CommitValue()
	fmt.Println(cm2)
	cm3 := coin.CommitSerialNumber()
	fmt.Println(cm3)

	// fmt.Println(privacy.FULL_CM)
	// fmt.Println(privacy.PK_CM)
	// fmt.Println(privacy.VALUE_CM)
	// fmt.Println(privacy.SN_CM)

	//
	//proof := privacy.ZkpPedersenCMProve(pcm, coin.PublicKey, coin.SerialNumber,  coin.Value, coin.R, coin.CoinCommitment)
	//
	//fmt.Println(privacy.ZkpPedersenCMVerify(pcm, *proof, coin.CoinCommitment))

	//Gx, Gy :=privacy.Curve.Params().ScalarBaseMult(nil)
	//
	//c := privacy.EllipticPoint{big.NewInt(0), big.NewInt(0)}
	//Hx, Hy:=privacy.Curve.Params().ScalarBaseMult([]byte("10"))
	//res, _ := privacy.Curve.Add(Gx, Gy, Hx, Hy)
	//res1, _ := privacy.Curve.Add(c.X, c.Y, Hx, Hy)
	//fmt.Println(res)
	//fmt.Println(res1)

}
