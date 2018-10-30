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

	r1Int := big.NewInt(0)
	r2Int := big.NewInt(0)
	r1 := privacy.RandBytes(32)
	r2 := privacy.RandBytes(32)
	r1Int.SetBytes(r1)
	r2Int.SetBytes(r2)
	r1Int.Mod(r1Int, privacy.Curve.Params().P)
	r2Int.Mod(r2Int, privacy.Curve.Params().P)
	r1 = r1Int.Bytes()
	r2 = r2Int.Bytes()
	committemp1 := privacy.Pcm.CommitSpecValue(serialNumber, r1, 0)
	committemp2 := privacy.Pcm.CommitSpecValue(serialNumber, r2, 0)
	fmt.Println(committemp1)
	fmt.Println(committemp2)
	committemp1Point, err := privacy.DecompressKey(committemp1)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	committemp2Point, err := privacy.DecompressKey(committemp2)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	r1Int.Sub(r1Int, r2Int)
	zeroInt := big.NewInt(0)
	if r1Int.Cmp(zeroInt) < 0 {
		r1Int.Mod(r1Int, privacy.Curve.Params().P)
	}
	negcommittemp2Point := new(privacy.EllipticPoint)
	negcommittemp2Point.X = big.NewInt(0)
	negcommittemp2Point.Y = big.NewInt(0)
	negcommittemp2Point.X.SetBytes(committemp2Point.X.Bytes())
	negcommittemp2Point.Y.SetBytes(committemp2Point.Y.Bytes())
	negcommittemp2Point.Y.Sub(privacy.Curve.Params().P, committemp2Point.Y)

	//negcommittemp2Point.X, negcommittemp2Point.Y = privacy.Curve.Add(negcommittemp2Point.X, negcommittemp2Point.Y, committemp2Point.X, committemp2Point.Y)

	committemp1Point.X, committemp1Point.Y = privacy.Curve.Add(committemp1Point.X, committemp1Point.Y, negcommittemp2Point.X, negcommittemp2Point.Y)
	commitZero := privacy.CompressKey(*committemp1Point)
	proofZero, z := privacy.ProveIsZero(commitZero, r1Int.Bytes(), 0)
	boolValue := privacy.VerifyIsZero(commitZero, proofZero, 0, z)
	fmt.Println(boolValue)
	fmt.Println("Done")
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

	var zk privacy.PKComValProtocol

	// pk := zk.GetPKCommittedValues()
	zk.SetWitness(witness)
	proof, _ := zk.Prove(coin.CoinCommitment)

	fmt.Printf("Proof: %+v\n", proof)

	fmt.Println(zk.Verify(*proof, coin.CoinCommitment))
}
