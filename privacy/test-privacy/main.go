package main

import (
	"crypto/sha256"
	"fmt"
	"math/big"

	"github.com/ninjadotorg/cash/privacy"
)

func main() {

	spendingKey := privacy.GenerateSpendingKey(new(big.Int).SetInt64(123).Bytes())
	fmt.Printf("\nSpending key: %v\n", spendingKey)
	fmt.Println(len(spendingKey))

	address := privacy.GenerateAddress(spendingKey)
	fmt.Printf("\nAddress: %v\n", address)
	fmt.Println(len(address))

	receivingKey := privacy.GenerateReceivingKey(spendingKey)
	fmt.Printf("\nReceiving key: %v\n", receivingKey)
	fmt.Println(len(receivingKey))

	transmissionKey := privacy.GenerateTransmissionKey(receivingKey)
	fmt.Printf("\nTransmission key: %v\n", transmissionKey)
	fmt.Println(len(transmissionKey))

	msg := "hello, world"
	hash := sha256.Sum256([]byte(msg))

	signature, err := privacy.Sign(hash[:], spendingKey)
	if err != nil {
		panic(err)
	}
	fmt.Printf("signature: %v\n", signature)

	valid := privacy.Verify(signature, hash[:], address)
	fmt.Println("signature verified:", valid)
}
