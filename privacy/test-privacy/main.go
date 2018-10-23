package main

import (
	"crypto/sha256"
	"fmt"
	"math/big"

	"github.com/ninjadotorg/cash/privacy"
)

func main() {

	spendingKey := privacy.GenSpendingKey(new(big.Int).SetInt64(123).Bytes())
	fmt.Printf("\nSpending key: %v\n", spendingKey)
	fmt.Println(len(spendingKey))

	address := privacy.GenAddress(spendingKey)
	fmt.Printf("\nAddress: %v\n", address)
	fmt.Println(len(address))

	receivingKey := privacy.GenReceivingKey(spendingKey)
	fmt.Printf("\nReceiving key: %v\n", receivingKey)
	fmt.Println(len(receivingKey))

	transmissionKey := privacy.GenTransmissionKey(receivingKey)
	fmt.Printf("\nTransmission key: %v\n", transmissionKey)
	fmt.Println(len(transmissionKey))

	msg := "hello, world"
	hash := sha256.Sum256([]byte(msg))

	r, s, err := privacy.Sign(hash[:], spendingKey)
	if err != nil {
		panic(err)
	}
	fmt.Printf("signature: (0x%x, 0x%x)\n", r, s)

	valid := privacy.Verify(r, s, hash[:], address)
	fmt.Println("signature verified:", valid)
}
