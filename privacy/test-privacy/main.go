package main

import (
	"crypto/sha256"
	"fmt"
	"math/big"

	"github.com/ninjadotorg/cash/privacy"
)

func main() {

	fmt.Printf("N: %X\n", privacy.Curve.Params().N)
	fmt.Printf("P: %X\n", privacy.Curve.Params().P)
	fmt.Printf("B: %X\n", privacy.Curve.Params().B)
	fmt.Printf("Gx: %x\n", privacy.Curve.Params().Gx)
	fmt.Printf("Gy: %X\n", privacy.Curve.Params().Gy)
	fmt.Printf("BitSize: %X\n", privacy.Curve.Params().BitSize)

	spendingKey := privacy.GenerateSpendingKey(new(big.Int).SetInt64(123).Bytes())
	fmt.Printf("\nSpending key: %v\n", spendingKey)
	fmt.Println(len(spendingKey))

	address := privacy.GenerateAddress(spendingKey)
	fmt.Printf("\nAddress: %v\n", address)
	fmt.Println(len(address))
	point, err := privacy.DecompressKey(address)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Address decom: %v\n", point)

	receivingKey := privacy.GenerateReceivingKey(spendingKey)
	fmt.Printf("\nReceiving key: %v\n", receivingKey)
	fmt.Println(len(receivingKey))

	transmissionKey := privacy.GenerateTransmissionKey(receivingKey)
	fmt.Printf("\nTransmission key: %v\n", transmissionKey)
	fmt.Println(len(transmissionKey))

	point, err = privacy.DecompressKey(transmissionKey)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Transmission key point decompress: %+v\n ", point)

	msg := "hello, world"
	hash := sha256.Sum256([]byte(msg))

	signature, err := privacy.Sign(hash[:], spendingKey)
	if err != nil {
		panic(err)
	}
	fmt.Printf("signature: %v\n", signature)

	valid := privacy.Verify(signature, hash[:], address)
	fmt.Println("\nsignature verified:", valid)

	curve1 := privacy.GetCurve()
	curve2 := privacy.GetCurve()
	fmt.Printf("Address curve 1: %v\n", &curve1)
	fmt.Printf("Address curve 2: %v\n", &curve2)
}
