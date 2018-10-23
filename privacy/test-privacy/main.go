package main

import (
	"crypto/sha256"
	"fmt"
	"math/big"

	"github.com/ninjadotorg/cash/privacy"
)

func main() {

	// Test transaction signing
	// privKey, _ := client.GenerateKey(rand.Reader)
	// tx := new(transaction.Tx)
	// tx.Version = 1
	// tx.Type = "Normal"
	// tx.LockTime = 123
	// tx.Fee = 1234
	// tx.SetTxId(tx.Hash())

	// fmt.Printf("Hash tx: %s\n", tx.GetTxId())
	// signed_tx, err := transaction.SignTx(tx, privKey)
	// if err != nil {
	// 	fmt.Printf("Error: %s", err)
	// }

	// pub := transaction.PubKeyToByteArray(&privKey.PublicKey)
	// signed_tx.JSPubKey = pub
	// fmt.Printf("Pub key: %v\n",signed_tx.JSPubKey)
	// fmt.Printf("Size of pub key: %d\n", len(signed_tx.JSPubKey))

	// fmt.Printf("Signature: %v\n", signed_tx.JSSig)
	// fmt.Printf("Size of signature: %d\n", len(signed_tx.JSSig))

	// fmt.Printf("Hash tx: %s\n", signed_tx.GetTxId())
	// res, _ := transaction.VerifySign(signed_tx)

	// fmt.Println(res)

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
