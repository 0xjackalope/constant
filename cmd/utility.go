package main

import (
	"github.com/ninjadotorg/cash/wallet"
	"os"
	"strconv"
	"log"
)

func generateDummyPrivateKey() {
	argsWithoutProg := os.Args[1:]
	passPhrase := argsWithoutProg[0]
	numAccountStr := argsWithoutProg[1]
	numAccount, _ := strconv.Atoi(numAccountStr)

	walletObj := wallet.Wallet{}
	walletObj.Init(passPhrase, uint32(numAccount), "")
	log.Printf("Mnemonic: %s\n", walletObj.Mnemonic)
	log.Printf("Master priv key: %s\n", walletObj.MasterAccount.Key.Base58CheckSerialize(wallet.PriKeyType))
	for _, account := range walletObj.MasterAccount.Child {
		log.Printf("%s: %s\n", account.Name, account.Key.Base58CheckSerialize(wallet.PriKeyType))
	}
}
